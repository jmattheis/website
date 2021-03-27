TAGS=netgo osusergo
LD_FLAGS=-s -w -linkmode external -extldflags "-static" -X main.Mode=prod

.PHONY: build
build:
	rm -rf website/docs
	packr2
	go build -ldflags="${LD_FLAGS}" -tags '${TAGS}' -o build/jmattheis.de main.go
