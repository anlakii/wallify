.PHONY: build install lint clean

export CGO_CFLAGS_ALLOW = -Xpreprocessor

build:
	go build -o wallify .

install:
	go install .

lint:
	golangci-lint run

clean:
	rm -f wallify
