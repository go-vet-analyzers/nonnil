<p align="center"><img src="https://raw.githubusercontent.com/go-vet-analyzers/brand/main/social/go-vet-analyzers.png" alt="go-vet-analyzers/nonnil" width="720"></p>

# go-vet-analyzers/nonnil

[![ci](https://github.com/go-vet-analyzers/nonnil/actions/workflows/ci.yml/badge.svg)](https://github.com/go-vet-analyzers/nonnil/actions/workflows/ci.yml)

A `go vet` analyzer that enforces the [go-composites](https://github.com/go-composites)
**Null-Object invariant**: a function whose result is a *Null-Object interface*
(any interface declaring `IsNull() bool`) must never return a bare `nil` — it must
hand back a real Null object (`Null.New()`, `NullError.New()`, …).

This turns "always return a Null object, never `nil`" from a matter of human
discipline into a check that **fails CI** — closing the one gap that lets the
nil-dereference the pattern is designed to prevent slip back in.

## Why

The composition style avoids `nil` dereferences because constructors initialise
every field to a Null object and callers send messages instead of testing
`!= nil`. But Go's type system still lets any function return a bare `nil` for an
interface result; one stray `return nil` and a caller panics. `nonnil` makes that
mistake un-mergeable.

## Install & use

```sh
go install github.com/go-vet-analyzers/nonnil/cmd/nonnil@latest
go vet -vettool=$(which nonnil) ./...
```

It reports, for example:

```
src/error.go:43:46: return a Null object (e.g. Null.New()) instead of nil for
  github.com/go-composites/error/src.Interface: it is a Null-Object interface
  (IsNull() bool), and nil reintroduces the nil-dereference the pattern prevents
```

## What it flags

A bare `nil` targeting a Null-Object interface (one with `IsNull() bool`,
directly or embedded), at every site where it would reintroduce a nil to
dereference:

- **`return nil`** — any result position (single and multi-value), incl. function literals.
- **`x = nil`** — assignment to a Null-Object-typed variable/field.
- **`var x T = nil`** — initialisation with an explicit Null-Object type.
- **composite literals** — `T{Field: nil}`, `[]T{nil}`, `[N]T{nil}`, `map[K]T{k: nil}`.

It deliberately does **not** flag interfaces without `IsNull()` (e.g. the builtin
`error`, or `Result.Interface`), so false positives are minimal. A deliberate
typed conversion `return (Thing)(nil)`, a nil interface variable, an omitted
struct field (its zero value), or a nil map *key* are left alone — those are
explicit, structural, or uncommon.

## CI

Add to each repo's workflow:

```yaml
  - run: go install github.com/go-vet-analyzers/nonnil/cmd/nonnil@latest
  - run: go vet -vettool=$(which nonnil) ./...
```

## License

BSD-3-Clause © the go-vet-analyzers/nonnil authors.
