
.PHONY: build push

all: build

# this one just build executables
build:
	cd reweave; ${MAKE}

clean:
	-cd reweave; ${MAKE} $@

# build also create sbom and push the images
push:
	cd reweave; ${MAKE} build

go-fmt:
	@gofmt

update-gomod:
	go get -u
	go mod tidy
