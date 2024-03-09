TAGS=netgo osusergo
LD_FLAGS=-s -w -linkmode external -extldflags "-static" -X main.Mode=prod

.PHONY: build
build:
	rm -rf assets/jmattheis-resume.pdf
	chromium --headless --print-to-pdf=assets/jmattheis-resume.pdf --print-to-pdf-no-header ./assets/resume.html
	go build -ldflags="${LD_FLAGS}" -tags '${TAGS}' -o build/jmattheis.de main.go
