package tool

import (
	"bufio"
	"io"
	"log"
	"os"
	"strings"

	"github.com/pkg/errors"
)

func ExistsOrError(path string) error {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return errors.Errorf("path '%s' does not exist\n", path)

	} else {
		if err != nil {
			return err
		}
	}
	return nil

}

func Copy(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
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
	replaceMap       map[string]string
	importReplaceMap map[string]string
	newMap           map[string]string
}

func NewRules() *Rules {
	return &Rules{
		replaceMap:       make(map[string]string, 0),
		importReplaceMap: make(map[string]string, 0), newMap: make(map[string]string, 0)}
}

func readRules(path string) (*Rules, error) {
	file, err := os.Open(path)
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
