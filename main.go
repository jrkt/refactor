package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	var dir string
	var oldString string
	var newString string
	var patterns string

	flag.StringVar(&dir, "dir", ".", "directory to perform refactor in")
	flag.StringVar(&oldString, "old", "", "old string to find")
	flag.StringVar(&newString, "new", "", "new string to replace old string with")
	flag.StringVar(&patterns, "patterns", "*", "comma-separated list of file patterns to perform refactor on")
	flag.Parse()

	if oldString == "" {
		fmt.Println("you must provide an 'old' string to search for")
		return
	}

	if newString == "" {
		fmt.Println("you must provide a 'new' string to replace with")
		return
	}

	patterns = strings.Replace(patterns, " ", "", -1)
	filePatterns := strings.Split(patterns, ",")

	err := refactor(dir, oldString, newString, filePatterns...)
	if err != nil {
		log.Fatalln("refactor failed:", err)
	}
}

func refactor(dir, old, new string, patterns ...string) error {
	return filepath.Walk(dir, refactorFunc(old, new, patterns))
}
func refactorFunc(old, new string, filePatterns []string) filepath.WalkFunc {
	return filepath.WalkFunc(func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !!fi.IsDir() {
			return nil
		}
		var matched bool
		for _, pattern := range filePatterns {
			var err error
			matched, err = filepath.Match(pattern, fi.Name())
			if err != nil {
				return err
			}
			if matched {
				read, err := ioutil.ReadFile(path)
				if err != nil {
					return err
				}

				fmt.Println("Refactoring:", path)
				newContents := strings.Replace(string(read), old, new, -1)
				err = ioutil.WriteFile(path, []byte(newContents), 0)
				if err != nil {
					return err
				}
			}
		}
		return nil
	})
}
