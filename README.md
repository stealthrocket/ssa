# SSA

A tool that prints the [SSA](https://pkg.go.dev/golang.org/x/tools/go/ssa)
representation of packages and/or functions.

## Usage

Install the tool:

```console
go install github.com/stealthrocket/ssa@latest
```

Print a summary of members in a package

```console
ssa path/to/go/package
````

Print the SSA form of one or more functions in a package:

```console
ssa path/to/go/package Function SomeOtherFunction
````

## Notes

This prints the `golang.org/x/tools/go/ssa` representation,
which is not the same as the SSA form that the Go compiler
uses (`cmd/compile/internal/ssa`).

See:
* https://pkg.go.dev/golang.org/x/tools/go/ssa
* https://github.com/golang/go/blob/master/src/cmd/compile/README.md
