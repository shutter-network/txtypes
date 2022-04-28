package tool

import (
	"fmt"
	"path/filepath"

	"github.com/pkg/errors"
)

func Run(inDir, outDir, ruleFile string) error {
	rules, err := readRules(ruleFile)
	if err != nil {
		return err
	}
	ExistsOrError(outDir)

	outFiles := make([]string, 0)
	for k, v := range rules.replaceMap {
		in := filepath.Join(inDir, k)
		out := filepath.Join(outDir, v)
		err = ExistsOrError(out)
		if err != nil {
			logError(err, "Replace failed")
			continue
		}
		err := Copy(in, out)
		if err != nil {
			logError(err, "Replace failed")
			continue

		} else {
			logf("Replace '%s' => '%s' \n", in, out)
			outFiles = append(outFiles, out)
		}
	}
	for k, v := range rules.newMap {
		in := filepath.Join(inDir, k)
		out := filepath.Join(outDir, v)
		err := Copy(in, out)
		if err != nil {
			logError(err, "New failed")
			continue
		} else {
			logf("New '%s' => '%s' \n", in, out)
			outFiles = append(outFiles, out)
		}
	}
	for k, v := range rules.importReplaceMap {
		err := changeImports(k, v, outFiles)
		if err != nil {
			logError(err, "Replace failed")
			continue
		}
	}
	return nil
}

func logf(format string, a ...any) (n int, err error) {
	return fmt.Printf(format, a...)
}

func logError(err error, message string) {
	err = errors.Wrap(err, "Replace failed")
	fmt.Println(err)
}
