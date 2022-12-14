package steps

import (
	"bytes"
	"context"
	"image"
	_ "image/png" // png decode support
	"math/bits"
	"time"

	"golang.org/x/image/draw"

	"github.com/chromedp/chromedp"
	"github.com/go-playground/validator/v10"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

const validateStepType = "validate-step"
const hashThresholdBits = 5

const pollFunction = `() => {
		const unloaded_images = [...document.querySelectorAll("img")].filter((image) => { return image.height == 0 })
		return unloaded_images.length == 0
	}
`

type validateStepConf struct {
	Hash     uint64 `validate:"required"`
	Isolated bool
}

type validateStep struct {
	name string
	conf validateStepConf
}

func (s *validateStep) GetType() string {
	return validateStepType
}

func (s *validateStep) GetName() string {
	return s.name
}

func (s *validateStep) Init(name string, input map[string]interface{}) error {
	var conf validateStepConf
	err := mapstructure.Decode(input, &conf)
	if err != nil {
		return errors.Wrapf(err, "failed parsing step '%s' configuration", s.GetType())
	}

	// validate conf using validate tags
	err = validator.New().Struct(conf)
	if err != nil {
		return errors.Wrapf(err, "failed validating step '%s' configuration", s.GetType())
	}

	s.name = name
	s.conf = conf

	return nil
}

func (s *validateStep) Run(logger *log.Entry, proxy string) chromedp.Tasks {
	logger.Infof("Running validate step with conf %+v", s.conf)

	validationTasks := chromedp.Tasks{}

	validationTasks = append(validationTasks, chromedp.ActionFunc(func(ctx context.Context) error {
		var buf []byte

		const timeoutMS = 100
		// wait for the image to load, this is required because isolated pages are ""loaded"" immediately
		// and then JS pops in elements in the background
		chromedp.PollFunction(pollFunction, nil, chromedp.WithPollingTimeout(timeoutMS*time.Millisecond))
		err := chromedp.FullScreenshot(&buf, 100).Do(ctx)
		if err != nil {
			logger.Error("failed to take screenshot")
		}

		return s.validateScreenshot(logger, buf)
	}))

	return validationTasks
}

// isolated tabs have a bar at the top of the screen
// this is testing to see if the bar is present in the screenshot, to check if it's isolated
// bar and its colours are defined in the WI system rule
func pixelTest(logger *log.Entry, img image.Image) bool {
	logger.Info("Pixel testing for isolation indicator")

	// converting between rgb types in golang (0-65535), and WI (0-255)
	const isolatedRed = 49 * 257
	const isolatedGreen = 122 * 257
	const isolatedBlue = 53 * 257

	// the upper left-hand-most corner pixel
	r, g, b, _ := img.At(0, 0).RGBA()

	return r == isolatedRed && g == isolatedGreen && b == isolatedBlue
}

func perceptualHash(logger *log.Entry, img image.Image) uint64 {
	// https://www.hackerfactor.com/blog/index.php?/archives/432-Looks-Like-It.html
	logger.Info("Finding perceptual hash")

	// resize the image to 8x8
	resized := image.NewRGBA(image.Rect(0, 0, 8, 8))
	draw.BiLinear.Scale(resized, resized.Rect, img, img.Bounds(), draw.Over, nil)

	// convert each pixel to grayscale (luminance), while taking an average across all pixels
	var luminanceTotal float64
	pixels := make([]float64, 64)

	for x := 0; x < 8; x++ {
		for y := 0; y < 8; y++ {
			// current pixel values
			r, g, b, _ := resized.At(x, y).RGBA()

			// https://e2eml.school/convert_rgb_to_grayscale.html
			// linear approximation magic, how we convert doesn't really matter, just needs to be consistent
			luminance := 0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b)

			pixels[(x*8)+y] = luminance
			luminanceTotal += luminance
		}
	}
	luminanceMean := luminanceTotal / 64

	// per pixel, is the luminance above average
	// store the bool result as a bit in a uint64 (8x8 pixels), finished uint64 is the hash
	var hash uint64

	for i := range pixels {
		if pixels[i] > luminanceMean {
			hash |= 1 << uint(len(pixels)-1-i)
		}
	}

	return hash
}

func (s *validateStep) validateScreenshot(logger *log.Entry, buf []byte) error {
	logger.Infof("validating test screenshot")

	// convert screenshot bytes to an image
	img, _, err := image.Decode(bytes.NewReader(buf))

	if err != nil {
		return errors.Errorf("Failed to decode screenshot bytes to image: %v", err)
	}

	// does the isolation indicator appear, and do we expect it
	isolated := pixelTest(logger, img)

	if isolated != s.conf.Isolated {
		return errors.Errorf("Tab isolated state expected %v, got %v.", s.conf.Isolated, isolated)
	}

	hash := perceptualHash(logger, img)
	logger.Infof("Got hash: %d", hash)

	// 1 if the bits are different in a given position
	hashXor := hash ^ s.conf.Hash

	// number of 1s (number of bits that were different)
	differentBits := bits.OnesCount64(hashXor)

	if differentBits > hashThresholdBits {
		return errors.Errorf("Too many different hash bits. Expected <= %d got %d.", hashThresholdBits, differentBits)
	}

	logger.Info("Done validating image hash")
	return nil
}
