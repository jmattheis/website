TAGS=netgo osusergo
LD_FLAGS=-s -w -linkmode external -extldflags "-static" -X main.Mode=prod

.PHONY: build
build:
	rm -rf website/docs
	(cd website && hugo)
	GO111MODULE=off go run github.com/gobuffalo/packr/v2/packr2
	go build -ldflags="${LD_FLAGS}" -tags '${TAGS}' -o build/jmattheis.de main.go
