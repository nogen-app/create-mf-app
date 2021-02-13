#
# Makefile
#
VERSION = snapshot
GHRFLAGS =
.PHONY: build release

default: build

build:
	goxc -d=pkg -pv=$(VERSION) -bc="linux darwin,amd64"

release:
	ghr  -u rareinator  $(GHRFLAGS) v$(VERSION) pkg/$(VERSION)
