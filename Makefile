GOLANG := go

GOLANG_TEST_FLAGS ?= "-v"

.PHONY: test
test:
	${GOLANG} test ${GOLANG_TEST_FLAGS} ./... -timeout 30s

.PHONY: clean
clean:
	rm -f marble
	rm -f code-coverage.out

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
