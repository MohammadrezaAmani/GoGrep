package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sync"
)

var wg2 sync.WaitGroup

func main2() {
	searchPattern := flag.String("pattern", "", "search pattern in regex format")
	rootDir := flag.String("dir", ".", "root directory to start search from")
	recursive := flag.Bool("r", false, "search recursively through subdirectories")
	caseInsensitive := flag.Bool("i", false, "perform case-insensitive search")

	flag.Parse()

	if *searchPattern == "" {
		fmt.Println("Please provide a search pattern")
		os.Exit(1)
	}

	var re *regexp.Regexp
	var err error

	if *caseInsensitive {
		re, err = regexp.Compile("(?i)" + *searchPattern)
	} else {
		re, err = regexp.Compile(*searchPattern)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error compiling regex: %v\n", err)
		os.Exit(1)
	}

	count := 0

	filepath.Walk(*rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() && !*recursive && path != *rootDir {
			return filepath.SkipDir
		}

		if !info.Mode().IsRegular() {
			return nil
		}

		wg2.Add(1)
		go func() {
			defer wg2.Done()

			dat, err := os.ReadFile(path)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error reading file %s: %v\n", path, err)
				return
			}

			if re.MatchString(string(dat)) {
				fmt.Println(path)
				count++
			}
		}()

		return nil
	})

	wg2.Wait()

	fmt.Printf("Found %d matches\n", count)
}
