#
# Makefile
#
VERSION = snapshot
GHRFLAGS =
.PHONY: build release

default: build

build:
	cd ./src
	goxc -d=pkg -pv=$(VERSION) -bc="linux darwin,amd64"

release:
	cd ./src
	ghr  -u nogen-app  $(GHRFLAGS) v$(VERSION) pkg/$(VERSION)

clean:
	rm -rf ./pkg
