package tool

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/pkg/errors"
)

func ExistsOrError(fs billy.Filesystem, path string) error {
	if _, err := fs.Stat(path); errors.Is(err, os.ErrNotExist) {
		return errors.Errorf("path '%s' does not exist\n", path)

	} else {
		if err != nil {
			return err
		}
	}
	return nil

}

type GitInfo struct {
	MainBranchOff  object.Commit
	Head           object.Commit
	URL            string
	LastVersionTag string
}

func GitCloneToFS(fs billy.Filesystem, url, mainBranch string) (*GitInfo, error) {
	gitSource := strings.Split(url, "@")
	if len(gitSource) != 2 {
		return nil, errors.Errorf("git url `%s` formatted incorrectly", gitSource)
	}
	branch := gitSource[1]
	url = fmt.Sprintf("https://%s", gitSource[0])

	// TODO use the git url from the "source" directive
	repo, err := git.Clone(memory.NewStorage(), fs, &git.CloneOptions{
		URL:        url,
		RemoteName: "origin",
	})
	revSh := plumbing.Revision(fmt.Sprintf("origin/%s", branch))
	shhash, err := repo.ResolveRevision(revSh)
	shHeadCommit, err := object.GetCommit(repo.Storer, *shhash)
	if err != nil {
		return nil, errors.Wrap(err, "resolving shutter head commit")
	}
	revMain := plumbing.Revision(fmt.Sprintf("origin/%s", mainBranch))
	mhash, err := repo.ResolveRevision(revMain)
	if err != nil {
		return nil, errors.Wrap(err, "resolving upstream head commit")
	}
	mHeadCommit, err := object.GetCommit(repo.Storer, *mhash)
	if err != nil {
		return nil, err
	}

	commonAncestorCommits, err := shHeadCommit.MergeBase(mHeadCommit)
	if err != nil {
		return nil, err
	}
	// TODO iterate over tags from the branch off on master
	// and find the last version tag (--> go-eth release)
	gi := &GitInfo{
		MainBranchOff: *commonAncestorCommits[0],
		Head:          *shHeadCommit,
		URL:           url,
	}
	return gi, nil

}

func Copy(srcfs, dstfs billy.Filesystem, src, dst string) error {
	in, err := srcfs.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := dstfs.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}

type Rules struct {
	source           string
	replaceMap       map[string]string
	importReplaceMap map[string]string
	newMap           map[string]string
}

func NewRules() *Rules {
	return &Rules{
		replaceMap:       make(map[string]string, 0),
		importReplaceMap: make(map[string]string, 0), newMap: make(map[string]string, 0)}
}

func readRules(fs billy.Filesystem, path string) (*Rules, error) {
	// TODO use regexes to parse the rules
	// instead of the string manipulations
	file, err := fs.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	rules := NewRules()
	lineNo := 0
	for scanner.Scan() {
		var src, target string
		ln := scanner.Text()
		splt := strings.Split(ln, ":")
		if len(splt) == 1 {
			return nil, errors.Errorf("syntax error in line %d: missing directive (\"new:\", \"replace:\", \"import-replace:\")", lineNo)
		}
		if splt[0] == "source" {
			url := strings.ReplaceAll(splt[1], " ", "")
			rules.source = url
			continue
		}
		rule := strings.ReplaceAll(splt[1], " ", "")
		spltInst := strings.Split(rule, "=>")
		if len(spltInst) == 2 {
			src = spltInst[0]
			target = spltInst[1]
		} else {
			src = spltInst[0]
			target = spltInst[0]
		}
		switch splt[0] {
		case "new":
			rules.newMap[src] = target
		case "replace":
			rules.replaceMap[src] = target
		case "import-replace":
			rules.importReplaceMap[src] = target
		default:
			err = errors.New("directive not known")
			log.Fatal()
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return rules, nil

}
