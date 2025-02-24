.PHONY: build
build: versioninfo
	go build -ldflags "${GO_LDFLAGS}" -o buildinfo main.go

.PHONY: versioninfo
versioninfo:
	$(eval GIT_COMMIT_HASH := $(shell git rev-parse HEAD))
	$(eval GIT_TAG := $(shell git describe --tags --dirty --always))

	$(eval GO_LDFLAGS := -s -w -X github.com/andrejacobs/go-micropkg/buildinfo.AppName=example -X github.com/andrejacobs/go-micropkg/buildinfo.Version=${GIT_TAG} -X github.com/andrejacobs/go-micropkg/buildinfo.GitCommitHash=${GIT_COMMIT_HASH})
