package main

import (
	"fmt"
	"go/token"
	"go/types"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/go/ssa"
	"golang.org/x/tools/go/ssa/ssautil"
	"golang.org/x/tools/go/types/typeutil"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	if len(os.Args) < 2 {
		return fmt.Errorf("usage: ssa <PACKAGE> [FUNCTIONS...]")
	}

	packagePath := os.Args[1]
	functionNames := os.Args[2:]

	// Determine the absolute path of the package, along
	// with the load pattern.
	packageAbsPath, err := filepath.Abs(packagePath)
	if err != nil {
		return err
	}
	var dotdotdot bool
	packageAbsPath, dotdotdot = strings.CutSuffix(packageAbsPath, "...")
	if s, err := os.Stat(packageAbsPath); err != nil {
		return err
	} else if !s.IsDir() {
		// Make sure we're loading whole packages.
		packageAbsPath = filepath.Dir(packageAbsPath)
	}
	var pattern string
	if dotdotdot {
		pattern = "./..."
	} else {
		pattern = "."
	}

	// Load the package(s), parse syntax, and check types.
	conf := &packages.Config{
		Mode: packages.NeedName | packages.NeedModule |
			packages.NeedImports | packages.NeedDeps |
			packages.NeedFiles | packages.NeedSyntax |
			packages.NeedTypes | packages.NeedTypesInfo,
		Fset: token.NewFileSet(),
		Dir:  packageAbsPath,
		Env:  os.Environ(),
	}
	pkgs, err := packages.Load(conf, pattern)
	if err != nil {
		return fmt.Errorf("packages.Load %q: %w", packagePath, err)
	}

	// Check that the package(s) are valid.
	var moduleDir string
	for _, p := range pkgs {
		if p.Module == nil {
			return fmt.Errorf("package %s is not part of a module", p.PkgPath)
		}
		if moduleDir == "" {
			moduleDir = p.Module.Dir
		} else if moduleDir != p.Module.Dir {
			return fmt.Errorf("pattern more than one module (%s + %s)", moduleDir, p.Module.Dir)
		}
	}
	err = nil
	packages.Visit(pkgs, func(p *packages.Package) bool {
		for _, e := range p.Errors {
			err = e
			break
		}
		return err == nil
	}, nil)
	if err != nil {
		return err
	}

	// Build the SSA program.
	prog, packages := ssautil.Packages(pkgs, ssa.InstantiateGenerics|ssa.GlobalDebug)
	prog.Build()

	// Print a summary of members in the package(s) if no
	// function names were specified.
	if len(functionNames) == 0 {
		for _, p := range packages {
			if _, err := p.WriteTo(os.Stdout); err != nil {
				return err
			}
		}
		return nil
	}

	// Otherwise, print the specified functions.
	functionSet := map[string]struct{}{}
	for _, name := range functionNames {
		functionSet[name] = struct{}{}
	}
	for _, p := range packages {
		for name, member := range p.Members {
			switch m := member.(type) {
			case *ssa.Function:
				if _, ok := functionSet[name]; ok {
					if _, err := m.WriteTo(os.Stdout); err != nil {
						return err
					}
				}
			case *ssa.Type:
				mset := typeutil.IntuitiveMethodSet(m.Type(), &prog.MethodSets)
				for _, selection := range mset {
					tfn := selection.Obj().(*types.Func)
					name := tfn.Name()
					if _, ok := functionSet[name]; ok {
						fn := prog.FuncValue(tfn)
						if _, err := fn.WriteTo(os.Stdout); err != nil {
							return err
						}
					}
				}
			}
		}
	}
	return nil
}
