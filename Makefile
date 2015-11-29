GODEP=$(GOPATH)/bin/godep
PREFIX?=/build

GOFILES = $(shell find . -type f -name '*.go')
httpbeat: $(GOFILES)
	# first make sure we have godep
	go get github.com/tools/godep
	$(GODEP) go build

.PHONY: getdeps
getdeps:
	go get -t -u -f

.PHONY: test
test:
	$(GODEP) go test ./...

.PHONY: updatedeps
updatedeps:
	$(GODEP) update ...

.PHONY: install_cfg
install_cfg:
	cp etc/httpbeat.yml $(PREFIX)/httpbeat-linux.yml
	cp etc/httpbeat.template.json $(PREFIX)/httpbeat.template.json
	# darwin
	cp etc/httpbeat.yml $(PREFIX)/httpbeat-darwin.yml
	# win
	cp etc/httpbeat.yml $(PREFIX)/httpbeat-win.yml

.PHONY: gofmt
gofmt:
	go fmt ./...

.PHONY: cover
cover:
	# gotestcover is needed to fetch coverage for multiple packages
	go get github.com/pierrre/gotestcover
	GOPATH=$(shell $(GODEP) path):$(GOPATH) $(GOPATH)/bin/gotestcover -coverprofile=profile.cov -covermode=count github.com/christiangalsterer/httpbeat/...
	mkdir -p cover
	$(GODEP) go tool cover -html=profile.cov -o cover/coverage.html

.PHONY: clean
clean:
	rm -r cover || true
	rm profile.cov || true
	rm httpbeat || true
