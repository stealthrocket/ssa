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

Print the SSA form of one or more functions/methods in a package:

```console
ssa path/to/go/package Function SomeOtherFunction SomeOtherMethod
````

## Notes

This prints the `golang.org/x/tools/go/ssa` representation,
which is not the same as the SSA form that the Go compiler
uses (`cmd/compile/internal/ssa`).

See:
* https://pkg.go.dev/golang.org/x/tools/go/ssa
* https://github.com/golang/go/blob/master/src/cmd/compile/README.md

A related tool is [x/tools/cmd/callgraph](https://pkg.go.dev/golang.org/x/tools@v0.16.1/cmd/callgraph),
which prints the callgraph for a set of packages. It uses the same
SSA form for its analysis.

```console
go install golang.org/x/tools/cmd/callgraph@latest
```
