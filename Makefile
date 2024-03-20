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
.PHONY: check-lint
check-lint:
	@echo "Linting code"
	go vet ./...
ifneq (${CI}, true)
	golangci-lint run
endif

# Check code quality
.PHONY: check
check: check-formatting check-lint

#------------------------------------------------------------------------------
# Go miscellaneous
#------------------------------------------------------------------------------

# Fetch required go modules
.PHONY: go-deps
go-deps:
	go mod download

# Tidy up module references (also donwloads deps)
.PHONY: go-tidy
go-tidy:
	go mod tidy
