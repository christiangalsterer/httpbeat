BEATNAME=httpbeat
BEAT_DESCRIPTION=Queries arbitrary HTTP endpoints in a regular interval and ships the result to the configured output channel.
BEAT_DIR=github.com/christiangalsterer/httpbeat
SYSTEM_TESTS=false
TEST_ENVIRONMENT=false
ES_BEATS?=./vendor/github.com/elastic/beats
# GOPACKAGES=$(shell glide novendor)
GOPACKAGES=$(shell go list ${BEAT_DIR}/... | grep -v /vendor/)
PREFIX?=.

# Path to the libbeat Makefile
-include $(ES_BEATS)/libbeat/scripts/Makefile

# Update dependencies
.PHONY: getdeps
getdeps:
	glide up --strip-vcs --update-vendored

# Initial beat setup
.PHONY: setup
setup: copy-vendor
	make update

# Copy beats into vendor directory
.PHONY: copy-vendor
copy-vendor:
	mkdir -p vendor/github.com/elastic/
	cp -R ${GOPATH}/src/github.com/elastic/beats vendor/github.com/elastic/
	rm -rf vendor/github.com/elastic/beats/.git

# This is called by the beats packer before building starts
.PHONY: before-build
before-build:

.PHONY: cover
cover:
	echo 'mode: atomic' > coverage.txt && go list . ./beater | xargs -n1 -I{} sh -c 'go test -covermode=atomic -coverprofile=coverage.tmp {} && tail -n +2 coverage.tmp >> coverage.txt' && rm coverage.tmp

# This is called by the beats packer before building starts
.PHONY: before-build
before-build:

# Collects all dependencies and then calls update
.PHONY: collect
collect:
