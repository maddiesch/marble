GOLANG := go

.PHONY: test
test:
	${GOLANG} test -v ./...

GO_SOURCE_FILES := $(shell find . -name "*.go" ! -name "*_test.go")

marble: ${GO_SOURCE_FILES}
	${GOLANG} build -o $@ ./cmd/marble

.PHONY: run-repl
run-repl: marble
	./marble repl
