GOLANG := go

GOLANG_TEST_FLAGS ?= -json -v
GOLANG_TEST_TIMEOUT ?= 30s
TEST_RUN ?= .

.PHONY: test
test:
	${GOLANG} test ${GOLANG_TEST_FLAGS} -run ${TEST_RUN} ./... -timeout ${GOLANG_TEST_TIMEOUT} 2>&1 | tee /tmp/gotest.log | gotestfmt

.PHONY: clean
clean:
	rm -f marble
	rm -f code-coverage.out
	go clean --cache

GO_SOURCE_FILES := $(shell find . -name "*.go" ! -name "*_test.go")
GO_SOURCE_TEST_FILES := $(shell find . -name "*_test.go")

marble: ${GO_SOURCE_FILES}
	${GOLANG} build -o $@ ./cmd/marble

.PHONY: run-repl
run-repl: marble
	./marble repl

.PHONY: debugger
debugger: marble
	dlv exec ./marble parse ./example.marble

code-coverage.out: ${GO_SOURCE_FILES} ${GO_SOURCE_TEST_FILES}
	GOLANG_TEST_FLAGS="-covermode count -coverprofile $@" ${MAKE} test

.PHONY: view-coverage
view-coverage: code-coverage.out
	${GOLANG} tool cover -html code-coverage.out
