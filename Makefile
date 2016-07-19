httpbeat:
	go build

.PHONY: getdeps
getdeps:
	glide up --strip-vcs --update-vendored

.PHONY: test
test:
	go test . ./beat/...

.PHONY: updatedeps
updatedeps:
	glide up --strip-vcs --update-vendored

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
	echo 'mode: atomic' > coverage.txt && go list . ./beat | xargs -n1 -I{} sh -c 'go test -covermode=atomic -coverprofile=coverage.tmp {} && tail -n +2 coverage.tmp >> coverage.txt' && rm coverage.tmp

.PHONY: clean
clean:
	rm -r cover || true
	rm profile.cov || true
	rm httpbeat || true
	rm coverage.txt || true
