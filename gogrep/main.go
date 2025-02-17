package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

type grepOptions struct {
	caseInsensitive bool
	invertMatch     bool
	lineNumbers     bool
	onlyMatch       bool
	countMatches    bool
	matchWholeWords bool
	recursive       bool
	concurrency     int
}

func searchFile(filePath string, re *regexp.Regexp, options grepOptions, wg *sync.WaitGroup, results chan<- string) {
	defer wg.Done()

	file, err := os.Open(filePath)
	if err != nil {
		results <- fmt.Sprintf("Error opening file %s: %v", filePath, err)
		return
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	lineNumber := 0
	matchCount := 0

	for {
		line, err := reader.ReadString('\n')
		if err != nil && err.Error() != "EOF" {
			results <- fmt.Sprintf("Error reading file %s: %v", filePath, err)
			break
		}
		if err == io.EOF {
			break
		}

		lineNumber++
		matches := re.FindAllString(line, -1)

		if options.invertMatch && len(matches) == 0 || !options.invertMatch && len(matches) > 0 {
			if options.countMatches {
				matchCount++
				continue
			}

			var output string
			if options.lineNumbers {
				output = fmt.Sprintf("%s:%d:%s", filePath, lineNumber, line)
			} else if options.onlyMatch {
				output = fmt.Sprintf("%s:%s", filePath, strings.Join(matches, " "))
			} else {
				output = fmt.Sprintf("%s:%s", filePath, line)
			}
			results <- output
		}
	}

	if options.countMatches {
		results <- fmt.Sprintf("%s:%d", filePath, matchCount)
	}
}

func searchDirectory(directory string, re *regexp.Regexp, options grepOptions, wg *sync.WaitGroup, results chan<- string) {
	defer wg.Done()

	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			wg.Add(1)
			go searchFile(path, re, options, wg, results)
		}
		return nil
	})

	if err != nil {
		results <- fmt.Sprintf("Error walking directory %s: %v", directory, err)
	}
}

func main() {
	options := grepOptions{}
	flag.BoolVar(&options.caseInsensitive, "i", false, "Perform case-insensitive matching")
	flag.BoolVar(&options.invertMatch, "v", false, "Invert match, show non-matching lines")
	flag.BoolVar(&options.lineNumbers, "n", false, "Show line numbers with output")
	flag.BoolVar(&options.onlyMatch, "o", false, "Only show the matched parts of the line")
	flag.BoolVar(&options.countMatches, "c", false, "Count the number of matching lines")
	flag.BoolVar(&options.matchWholeWords, "w", false, "Match whole words only")
	flag.BoolVar(&options.recursive, "r", false, "Search recursively in directories")
	flag.IntVar(&options.concurrency, "j", 10, "Maximum number of concurrent tasks")
	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Println("Usage: grep [options] pattern [file|directory...]")
		os.Exit(1)
	}

	pattern := flag.Arg(0)
	if options.caseInsensitive {
		pattern = "(?i)" + pattern
	}
	if options.matchWholeWords {
		pattern = `\b` + pattern + `\b`
	}

	re, err := regexp.Compile(pattern)
	if err != nil {
		fmt.Printf("Error compiling regex: %v\n", err)
		os.Exit(1)
	}

	files := flag.Args()[1:]
	if len(files) == 0 {
		files = append(files, ".")
	}

	results := make(chan string, options.concurrency)
	var wg sync.WaitGroup

	for _, file := range files {
		fileInfo, err := os.Stat(file)
		if err != nil {
			fmt.Printf("Error accessing path %s: %v\n", file, err)
			continue
		}

		if fileInfo.IsDir() && options.recursive {
			wg.Add(1)
			go searchDirectory(file, re, options, &wg, results)
		} else if !fileInfo.IsDir() {
			wg.Add(1)
			go searchFile(file, re, options, &wg, results)
		}
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	for result := range results {
		fmt.Println(result)
	}
}
