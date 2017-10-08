VERSION ?= $(shell git describe --tags)

IMAGE = dkr.yldr.io/vulcand
PKG = github.com/yieldr/vulcand

OS ?= darwin
ARCH ?= amd64

GOBUILDFLAGS = -a -tags netgo -ldflags '-w'

build:
	GOOS=$(OS) GOARCH=$(ARCH) go build -o bin/vulcand $(GOBUILDFLAGS)
	GOOS=$(OS) GOARCH=$(ARCH) go build -o bin/vctl $(GOBUILDFLAGS) ./vctl

install:
	go install .
	go install ./vctl

test:
	go test

configure:
	vctl backend upsert -id nginx
	vctl server upsert -b nginx -id nginx-1 -url http://nginx:80
	vctl frontend upsert -id nginx -b nginx -route 'PathRegexp("/.*")'
	vctl oauth2 upsert -f nginx -id nginx-oauth \
		-domain yieldr.eu.auth0.com \
		-clientId JklNORC4LOPSotjX25sZVcam6ZWpM53f \
		-clientSecret Pu0kl30h4ut7pFh5baczOhLlyCpBv-pm9iQOKFsVEsVdeUgEGlh4RY0zeknl4oUx \
		-redirectUrl http://localhost:8181/callback

docker-all: docker-build docker-image docker-push

docker-build:
	@docker run -i --rm -v "$(PWD):/go/src/$(PKG)" $(IMAGE):build make build

docker-test:
	@docker run -i --rm -v "$(PWD):/go/src/$(PKG)" $(IMAGE):build make test

docker-image:
	@docker build -t $(IMAGE):$(VERSION) .
	@docker tag -f $(IMAGE):$(VERSION) $(IMAGE):latest
	@echo " ---> $(IMAGE):$(VERSION)\n ---> $(IMAGE):latest"

docker-push:
	@docker push $(IMAGE):$(VERSION)
	@docker push $(IMAGE):latest

docker-builder-image:
	@docker build -t $(IMAGE):build -f Dockerfile.build .
