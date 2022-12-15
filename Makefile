BLACKBOX_REPO=myrepo
BLACKBOX_IMAGENAME=thin-blackbox-tester
BLACKBOX_VERSION=1.0.0

build_run_linux:
	GOOS=linux GOARCH=amd64 go build -o ./bin/service/${BLACKBOX_IMAGENAME} ./service/cmd/service/main.go
	chmod +x ./bin/service/${BLACKBOX_IMAGENAME}
	./bin/service/${BLACKBOX_IMAGENAME}

build_run_darwin:
	GOOS=darwin GOARCH=amd64 go build -o ./bin/service/${BLACKBOX_IMAGENAME} ./service/cmd/service/main.go
	chmod +x ./bin/service/${BLACKBOX_IMAGENAME}
	./bin/service/${BLACKBOX_IMAGENAME}