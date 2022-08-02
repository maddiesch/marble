# Marble Language

An embedded scripting language written in Go.

## CLI

### REPL

There is a REPL available!

```bash
make marble

./marble repl
```

## Pipeline

### Lexing & Tokenization

Print example tokens to stdout

```bash
make marble

./marble tokenize ./example.marble
```

### Parsing

Parse the source file and print the resulting AST to stdout

```bash
make marble

./marble parse ./example.marble
```

### Evaluating

Run the source file through the evaluator

```bash
make marble

./marble run ./example.marble
```
