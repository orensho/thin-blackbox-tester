BLACKBOX_REPO=myrepo
BLACKBOX_IMAGENAME=thin-blackbox-tester
BLACKBOX_VERSION=1.0.0

build:
	go test .
	GOOS=linux GOARCH=amd64 go build -o thin-slack-bot
	docker build --platform=linux/amd64 -t ${BLACKBOX_REPO}/${BLACKBOX_IMAGENAME}:${BLACKBOX_VERSION} .