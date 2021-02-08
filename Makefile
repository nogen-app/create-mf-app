#
# Makefile
#
VERSION = snapshot
GHRFLAGS =
.PHONY: build release

default: build

build:
	gox

release:
	ghr  -u rareinator  $(GHRFLAGS) v$(VERSION) pkg/$(VERSION)
