
.PHONY: build push

all: build

# this one just build executables
build:
	cd reweave; ${MAKE}

clean:
	-cd reweave; ${MAKE} $@
	-docker plugin rm avhost/net-plugin:latest_release-amd64


# build also create sbom and push the images
push:
	-docker plugin rm avhost/net-plugin:latest_release-amd64
	cd reweave; ${MAKE} build publish

go-fmt:
	@gofmt

update-gomod:
	go get -u
	go mod tidy
