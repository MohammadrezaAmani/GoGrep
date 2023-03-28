package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var wg sync.WaitGroup
var recursive bool
var caseInsensitive bool
var rootDir string

func check(e error) {
	if e != nil {
		panic(e)
	}
}

var searchText string

func ReadFile(path string) {
	defer wg.Done()
	dat, err := os.ReadFile(path)
	check(err)

	f, err := os.Open(path)
	check(err)
	defer f.Close()
	if strings.Contains(string(dat), searchText) {
		fmt.Println(path)
		count += 1
	}
	f.Close()
}

var count int

func WalkDir(path string) {
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && !recursive && path != rootDir {
			return filepath.SkipDir
		}
		{
			if info.Mode().IsRegular() {
				wg.Add(1)
				go ReadFile(path)
			}
		}
		return nil
	})
	if err != nil {
		fmt.Println(err)
	}
}

func main() {

	searchPattern := flag.String("p", "", "search pattern in regex format")
	rootDir = string(*flag.String("d", ".", "root directory to start search from"))
	recursiveFlag := flag.Bool("r", false, "search recursively through subdirectories")
	caseInsensitiveFlag := flag.Bool("i", false, "perform case-insensitive search")
	flag.Parse()
	recursive = *recursiveFlag
	caseInsensitive = *caseInsensitiveFlag
	searchText = *searchPattern
	if caseInsensitive {
		searchText = strings.ToLower(searchText)
	}
	fmt.Println("*"+searchText+"*")
	WalkDir(rootDir)
	fmt.Println(count)

}
