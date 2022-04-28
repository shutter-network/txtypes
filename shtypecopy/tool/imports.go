// parts of the functionality is based on https://github.com/rogpeppe/govers
// for Copyright information see LICENSE file
package tool

import (
	"bufio"
	"fmt"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"regexp"
	"strconv"

	"github.com/pkg/errors"
)

func changeImports(match, newPackage string, editFiles []string) error {
	var (
		err           error
		oldPackagePat *regexp.Regexp
	)

	if match == "" {
		return errors.New("match is empty")

	}
	if newPackage == "" {
		return errors.New("new package name is empty")

	}

	oldPackagePat, err = regexp.Compile("^(" + match + ")")
	if err != nil {
		return errors.Wrap(err, "invalid match pattern")
	}
	ctxt := &context{
		newPackage:    newPackage,
		oldPackagePat: oldPackagePat,
	}

	for _, path := range editFiles {
		if ctxt.changeVersion(path) {
			logf("Import-replace: '%s' => '%s' in File %s\n", match, newPackage, path)
		}
	}
	return nil
}

type context struct {
	newPackage    string
	oldPackagePat *regexp.Regexp
	checked       map[string]bool
}

var printConfig = printer.Config{
	Mode:     printer.TabIndent | printer.UseSpaces,
	Tabwidth: 8,
}

func (ctxt *context) fixPath(p string) string {
	loc := ctxt.oldPackagePat.FindStringSubmatchIndex(p)
	if loc == nil {
		return p
	}
	i := loc[3]
	if p[0:i] != ctxt.newPackage {
		p = ctxt.newPackage + p[i:]
	}
	return p
}

// changeVersion changes the named go file to
// import the new version.
func (ctxt *context) changeVersion(path string) bool {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		logError(err, fmt.Sprintf("cannot parse %q", path))
	}
	changed := false
	for _, ispec := range f.Imports {
		impPath, err := strconv.Unquote(ispec.Path.Value)
		if err != nil {
			panic(err)
		}
		if p := ctxt.fixPath(impPath); p != impPath {
			ispec.Path.Value = strconv.Quote(p)
			changed = true
		}
	}
	if !changed {
		return changed
	}
	out, err := os.Create(path)
	if err != nil {
		logError(err, "cannot create file")
	}
	defer out.Close()
	w := bufio.NewWriter(out)
	if err := printConfig.Fprint(w, fset, f); err != nil {
		logError(err, "cannot write file")
	}
	if err := w.Flush(); err != nil {
		logError(err, "cannot write file")
	}
	return true
}
