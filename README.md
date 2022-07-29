# Marble Language

An embeddable scripting language.

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
