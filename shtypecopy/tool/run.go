package tool

import (
	"fmt"

	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-billy/v5/osfs"
	"github.com/pkg/errors"
)

func Run(inDir, outDir, ruleFile string) error {
	outFs := osfs.New("./")
	outDir = outFs.Join(outDir)
	err := ExistsOrError(outFs, outDir)
	if err != nil {
		return err
	}

	rules, err := readRules(outFs, ruleFile)
	if err != nil {
		return err
	}

	outFiles := make([]string, 0)
	if rules.source == "" {
		return errors.New("'source' directive is mandatory")
	}

	// Read to memory from git
	inFs := memfs.New()
	mainBranch := "master"
	gitInfo, err := GitCloneToFS(
		inFs,
		rules.source,
		mainBranch,
	)
	logf("Common ancestor commit with branch '%s': %s\n", mainBranch, gitInfo.MainBranchOff.Hash)
	file, err := outFs.Create("./GOETHEREUM_BRANCH_OFF")
	defer file.Close()
	if err != nil {
		return err
	}
	file.Write([]byte(fmt.Sprintf("%s\n", gitInfo.MainBranchOff.Hash)))

	file, err = outFs.Create("./GOETHEREUM_SOURCE")
	defer file.Close()
	if err != nil {
		return err
	}
	file.Write([]byte(fmt.Sprintf("%s\n", rules.source)))

	file, err = outFs.Create("./GOETHEREUM_COMMIT")
	defer file.Close()
	if err != nil {
		return err
	}
	file.Write([]byte(fmt.Sprintf("%s\n", gitInfo.Head.Hash)))

	if err != nil {
		return err
	}
	for k, v := range rules.replaceMap {
		in := inFs.Join(inDir, k)
		out := outFs.Join(outDir, v)
		err = ExistsOrError(outFs, out)
		if err != nil {
			logError(err, "Replace failed")
			continue
		}
		err := Copy(inFs, outFs, in, out)
		if err != nil {
			logError(err, "Replace failed")
			continue

		} else {
			logf("Replace '%s' => '%s' \n", in, out)
			outFiles = append(outFiles, out)
		}
	}
	for k, v := range rules.newMap {
		in := inFs.Join(inDir, k)
		out := outFs.Join(outDir, v)
		err := Copy(inFs, outFs, in, out)
		if err != nil {
			logError(err, "New failed")
			continue
		} else {
			logf("New '%s' => '%s' \n", in, out)
			outFiles = append(outFiles, out)
		}
	}
	for k, v := range rules.importReplaceMap {
		err := changeImports(outFs, k, v, outFiles)
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
