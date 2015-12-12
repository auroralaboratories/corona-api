// Package exp11 allows displaying Go package source code with dot imports inlined, or merging entire Go package into a single file.
package exp11

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"os"
	"strings"

	"github.com/shurcooL/go/gists/gist5639599"
	"golang.org/x/tools/go/loader"
	"golang.org/x/tools/imports"
)

const parserMode = parser.ParseComments
const astMergeMode = 0*ast.FilterFuncDuplicates | ast.FilterUnassociatedComments | ast.FilterImportDuplicates

var dotImports []*loader.PackageInfo

func findDotImports(prog *loader.Program, pi *loader.PackageInfo) {
	for _, file := range pi.Files {
		for _, importSpec := range file.Imports {
			if importSpec.Name != nil && importSpec.Name.Name == "." {
				dotImportImportPath := strings.Trim(importSpec.Path.Value, `"`)
				dotImportPi := prog.Package(dotImportImportPath)
				dotImports = append(dotImports, dotImportPi)
				findDotImports(prog, dotImportPi)
			}
		}
	}
}

// InlineDotImports displays Go package source code with dot imports inlined.
func InlineDotImports(w io.Writer, importPath string) {
	/*imp2 := importer.New()
	imp2.Config.UseGcFallback = true
	cfg := types.Config{Import: imp2.Import}
	_ = cfg*/

	conf := loader.Config{
	//TypeChecker:   cfg,
	}

	conf.Import(importPath)

	prog, err := conf.Load()
	if err != nil {
		panic(err)
	}

	/*pi, err := imp.ImportPackage(importPath)
	if err != nil {
		panic(err)
	}
	_ = pi*/

	pi := prog.Imported[importPath]

	findDotImports(prog, pi)

	files := make(map[string]*ast.File)
	{
		// This package
		for _, file := range pi.Files {
			filename := prog.Fset.File(file.Package).Name()
			files[filename] = file
		}

		// All dot imports
		for _, pi := range dotImports {
			for _, file := range pi.Files {
				filename := prog.Fset.File(file.Package).Name()
				files[filename] = file
			}
		}
	}

	apkg := &ast.Package{Name: pi.Pkg.Name(), Files: files}

	merged := ast.MergePackageFiles(apkg, astMergeMode)

	WriteMergedPackage(w, prog.Fset, merged)
}

// WriteMergedPackage writes a merged package, typically coming from ast.MergePackageFiles, to w.
// It sorts and de-duplicates imports.
//
// TODO: Support comments.
func WriteMergedPackage(w io.Writer, fset *token.FileSet, merged *ast.File) {
	switch 3 {
	case 1:
		fmt.Fprintln(w, "package "+gist5639599.SprintAst(fset, merged.Name))
		fmt.Fprintln(w)
		fmt.Fprintln(w, `import (`)
		// TODO: SortImports (ala goimports).
		for _, importSpec := range merged.Imports {
			if importSpec.Name != nil && importSpec.Name.Name == "." {
				continue
			}
			fmt.Fprintln(w, "\t"+gist5639599.SprintAst(fset, importSpec))
		}
		fmt.Fprintln(w, `)`)
		fmt.Fprintln(w)

		for _, decl := range merged.Decls {
			if x, ok := decl.(*ast.GenDecl); ok && x.Tok == token.IMPORT {
				continue
			}

			fmt.Fprintln(w, gist5639599.SprintAst(fset, decl))
			fmt.Fprintln(w)
		}
	case 2:
		sortDecls(merged)

		//fmt.Fprintln(w, gist5639599.SprintAst(token.NewFileSet(), merged))

		//ast.SortImports(fset, merged)
		sortImports2(fset, merged)

		fmt.Fprintln(w, gist5639599.SprintAst(fset, merged))
	case 3:
		sortDecls(merged)

		// TODO: Clean up this mess...
		fset2, f2 := sortImports2(token.NewFileSet(), merged)

		fmt.Fprintln(w, "package "+gist5639599.SprintAst(fset, merged.Name))
		for _, decl := range f2.Decls {
			if x, ok := decl.(*ast.GenDecl); ok && x.Tok == token.IMPORT {
				fmt.Fprintln(w)
				fmt.Fprintln(w, gist5639599.SprintAst(fset2, decl))
			}
		}
		for _, decl := range merged.Decls {
			if x, ok := decl.(*ast.GenDecl); ok && (x.Tok == token.IMPORT || x.Tok == token.PACKAGE) {
				continue
			}

			fmt.Fprintln(w)
			fmt.Fprintln(w, gist5639599.SprintAst(fset, decl))
		}
	case 4:
		sortDecls(merged)

		src := []byte(gist5639599.SprintAst(fset, merged))

		out, err := imports.Process("", src, nil)
		if err != nil {
			panic(err)
		}
		os.Stdout.Write(out)
		fmt.Println()
	}
}

func sortDecls(merged *ast.File) {
	var sortedDecls []ast.Decl
	for _, decl := range merged.Decls {
		if x, ok := decl.(*ast.GenDecl); ok && x.Tok == token.PACKAGE {
			sortedDecls = append(sortedDecls, decl)
		}
	}
	/*for _, decl := range merged.Decls {
		if x, ok := decl.(*ast.GenDecl); ok && x.Tok == token.IMPORT {
			sortedDecls = append(sortedDecls, decl)
			goon.DumpExpr(decl)
		}
	}*/
	var specs []ast.Spec
	for _, importSpec := range merged.Imports {
		if importSpec.Name != nil && importSpec.Name.Name == "." {
			continue
		}
		importSpec.EndPos = 0
		specs = append(specs, importSpec)
	}
	sortedDecls = append(sortedDecls, &ast.GenDecl{
		Tok:    token.IMPORT,
		Lparen: (token.Pos)(1), // Needs to be non-zero to be considered as a group.
		Specs:  specs,
	})
	//goon.DumpExpr(sortedDecls[len(sortedDecls)-1])
	for _, decl := range merged.Decls {
		if x, ok := decl.(*ast.GenDecl); ok && (x.Tok == token.IMPORT || x.Tok == token.PACKAGE) {
			continue
		}
		sortedDecls = append(sortedDecls, decl)
	}
	merged.Decls = sortedDecls
}
