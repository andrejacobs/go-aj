# Directory in which the temporary files are stored
BUILD_OUTPUT_DIR=build
REPORT_OUTPUT_DIR=${BUILD_OUTPUT_DIR}/reports

.PHONY: all
all: clean test

# Remove build output
.PHONY: clean
clean:
	@echo "Cleaning build output"
	go clean
	rm -rf ${BUILD_OUTPUT_DIR}

#------------------------------------------------------------------------------
# Code quality assurance
#------------------------------------------------------------------------------

# Run unit-testing with race detector and code coverage report
.PHONY: test
test:
	@echo "Running unit-tests"
	$(eval COVERAGE_REPORT := ${REPORT_OUTPUT_DIR}/codecoverage)
	@mkdir -p "${REPORT_OUTPUT_DIR}"
	@go test -v -count=1 -race ./... -coverprofile="${COVERAGE_REPORT}"

# Check if the last code coverage report met minimum coverage standard of 80%, if not make exit with error code
.PHONY: test-coverage-passed
test-coverage-passed:
	$(eval COVERAGE_REPORT := ${REPORT_OUTPUT_DIR}/codecoverage)
	@go tool cover -func "${COVERAGE_REPORT}" \
	| grep "total:" | awk '{code=((int($$3) > 80) != 1)} END{exit code}'

# Generate HTML from the last code coverage report
.PHONY: test-coverage-report
test-coverage-report:
	$(eval COVERAGE_REPORT := ${REPORT_OUTPUT_DIR}/codecoverage)
	@go tool cover -html="${COVERAGE_REPORT}" -o "${COVERAGE_REPORT}.html"
	@echo "Code coverage report: file://`realpath ${COVERAGE_REPORT}.html`"

# Check that the source code is formatted correctly according to the gofmt standards
.PHONY: check-formatting
check-formatting:
	@test -z $(shell gofmt -e -l ./ | tee /dev/stderr) || (echo "Please fix formatting first with gofmt" && exit 1)

# Check for other possible issues in the code
# NOTE: To install golangci-lint
# https://golangci-lint.run/welcome/install/#local-installation
# curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.64.5
.PHONY: check-lint
check-lint:
	@echo "Linting code"
	go vet ./...
ifneq (${CI}, true)
	golangci-lint run
	addlicense -check -c "Andre Jacobs" -l mit -ignore '.github/**' -ignore 'build/**' ./
endif

# Check code quality
.PHONY: check
check: check-formatting check-lint

#------------------------------------------------------------------------------
# Cross compilation testing
#------------------------------------------------------------------------------

# NOTE: I am using an Apple M2 at the moment and the following has not been tested on other systems

# Check the unit-tests (and thus the packages) compile on a 32bit x86 Linux machine
.PHONY: build-tests-32bit
build-tests-32bit:
	@echo "Building unit-tests to be run on 32 bit architecture"
	@mkdir -p ${BUILD_OUTPUT_DIR}/bin/tests
	@for pkg in $$(go list ./...); do \
    	GOOS=linux GOARCH=386 CGO_ENABLED=0 go test -c "$$pkg" -o "${BUILD_OUTPUT_DIR}/bin/tests/$$(basename $$pkg).test"; \
	done

# Run the unit-tests on a 32bit x86 Linux emulated machine using Docker
.PHONY: test-32bit
test-32bit:
	@command -v docker >/dev/null 2>&1 || { echo "docker not found. Please install docker first."; exit 1; }

	docker run --rm -v $$(pwd):/app -w /app --platform linux/386 golang:1.25.1 \
    	go test -v ./...

#------------------------------------------------------------------------------
# Miscellaneous
#------------------------------------------------------------------------------

# Fetch required go modules
.PHONY: go-deps
go-deps:
	go mod download

# Tidy up module references (also donwloads deps)
.PHONY: go-tidy
go-tidy:
	go mod tidy

# Add the copyright and license notice
.PHONY: addlic
addlic:
	@echo "Adding copyright and license notice"
	addlicense -v -c "Andre Jacobs" -l mit -ignore '.github/**' -ignore 'build/**' ./
