.SUFFIXES: .go

VERSION = 1.0.1

DEFAULT_BUILDVARS := CGO_ENABLED=0 GOOS=${GOOS} GOARCH=amd64
BUILDVARS  := $(DEFAULT_BUILDVARS)
COMMIT_HASH := $(shell git describe --match=NeVeRmAtCh --always --abbrev=8 --dirty)
BUILDFLAGS := -a -tags netgo -ldflags "-w -extldflags \"-static\" -X main.GitCommit=$(COMMIT_HASH) -X main.Version=$(VERSION)"

default: tarsplitter

clean:
	@rm -f tarsplitter

tarsplitter: tarsplitter.go cmd/tarsplitter/main.go
	$(BUILDVARS) go build $(BUILDFLAGS) github.com/messiaen/tarsplitter/cmd/tarsplitter

install: tarsplitter.go cmd/tarsplitter/main.go
	$(BUILDVARS) go install $(BUILDFLAGS) github.com/messiaen/tarsplitter/cmd/tarsplitter
