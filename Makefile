TAGS=netgo osusergo
LD_FLAGS=-s -w -X main.Mode=prod

.PHONY: build
build:
	chromium --headless --no-pdf-header-footer --print-to-pdf=assets/jmattheis-resume.pdf --print-to-pdf-no-header ./assets/resume.html
	[ $$(pdfinfo assets/jmattheis-resume.pdf | awk '/^Pages:/ {print $$2}') = 1 ] || exit 1
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -ldflags="${LD_FLAGS}" -tags '${TAGS}' -o build/jmattheis.de main.go
