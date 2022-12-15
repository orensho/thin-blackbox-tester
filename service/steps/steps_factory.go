package steps

import (
	"github.com/pkg/errors"
)

type StepFactoryInterface interface {
	NewStep(stepType string) (StepInterface, error)
}

type stepFactoryImpl struct{}

func NewStepFactory() StepFactoryInterface {
	return &stepFactoryImpl{}
}

func (sf *stepFactoryImpl) NewStep(stepType string) (StepInterface, error) {
	switch stepType {
	case navigateStepType:
		return &navigateStep{}, nil
	case waitStepType:
		return &waitStep{}, nil
	case validateStepType:
		return &validateStep{}, nil
	default:
		return nil, errors.Errorf("Undefined step '%s'", stepType)
	}
}
