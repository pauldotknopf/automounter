
.PHONY: ci_deps vendor build

default: build

build:
	@echo "building..."
vendor:
	vndr
ci_deps:
	@echo "fetching vndr"
	@go get -u github.com/LK4D4/vndr