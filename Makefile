SUFFIX=$(GOOS)_$(GOARCH)

.PHONY: build

basenine:
	cd server/ && go build -gcflags="-e" -o ../basenine *.go

install: basenine
	mv basenine /usr/local/bin/

clean:
	rm basenine || true

test: clean basenine test-client-go
	cd server/ && go test *.go -v -covermode=atomic -coverprofile=coverage.out
	$(MAKE) coverage

test-client-go:
	test/client_go.sh

build:
	cd server/ && CGO_ENABLED=0 go build -gcflags="-e" -o ../build/basenine_$(SUFFIX) *.go && \
	cd ../build/ && shasum -a 256 basenine_${SUFFIX} > basenine_${SUFFIX}.sha256

build-all:
	rm -rf build/ && \
	$(MAKE) build GOOS=linux GOARCH=amd64
	$(MAKE) build GOOS=linux GOARCH=arm64
	$(MAKE) build GOOS=linux GOARCH=386
	$(MAKE) build GOOS=darwin GOARCH=amd64
	$(MAKE) build GOOS=darwin GOARCH=arm64

coverage:
	cp server/coverage.out coverage.out && sed 1,1d client/go/coverage.out >> coverage.out

install-init-systemd:
	cp scripts/init/systemd/basenine.service /etc/systemd/system/ && \
	systemctl daemon-reload && \
	systemctl restart basenine && \
	systemctl status basenine
