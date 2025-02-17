# GoGrep - A Fast, Concurrent Search Tool

`GoGrep` is a command-line tool for searching files and directories for patterns using regular expressions. It offers features similar to the classic `grep` tool, but with added concurrency, flexibility, and performance enhancements. It supports options such as recursive searching, case-insensitive matching, whole-word matching, line-number display, and more.

## Features

- **Case-insensitive Search**: Use the `-i` flag for case-insensitive matching.
- **Invert Match**: Use the `-v` flag to show lines that do not match the pattern.
- **Line Numbers**: Use the `-n` flag to display the line numbers where matches occur.
- **Show Only Matches**: Use the `-o` flag to show only the matching part of the line.
- **Count Matches**: Use the `-c` flag to count the number of matching lines.
- **Whole Word Matching**: Use the `-w` flag to match only whole words.
- **Recursive Search**: Use the `-r` flag to search directories recursively.
- **Concurrency**: Use the `-j` flag to specify the maximum number of concurrent tasks for faster searching in large files and directories.

## Installation

### Prerequisites

Ensure that you have Go installed on your system. You can download it from [https://golang.org/dl/](https://golang.org/dl/).

### Clone the repository

```bash
git clone https://github.com/MohammadrezaAmani/GoGrep.git
cd GoGrep
```

### Build the application

```bash
go build -o gogrep gogrep/main.go
```

## Usage

```bash
./gogrep [options] pattern [file|directory...]
```

### Options

- `-i`: Perform case-insensitive matching.
- `-v`: Invert the match, showing non-matching lines.
- `-n`: Show line numbers with output.
- `-o`: Show only the matched parts of the line.
- `-c`: Count the number of matching lines.
- `-w`: Match whole words only.
- `-r`: Search recursively in directories.
- `-j`: Set the maximum number of concurrent tasks (default is 10).

### Example Usages

1. **Search for a pattern in a file with case-insensitive matching:**

   ```bash
   ./gogrep -i "pattern" file.txt
   ```

2. **Search recursively in a directory for a pattern:**

   ```bash
   ./gogrep -r "pattern" /path/to/directory
   ```

3. **Show line numbers along with matching lines:**

   ```bash
   ./gogrep -n "pattern" file.txt
   ```

4. **Count the number of matching lines in a file:**

   ```bash
   ./gogrep -c "pattern" file.txt
   ```

5. **Search for a pattern and show only the matched parts of the line:**

   ```bash
   ./gogrep -o "pattern" file.txt
   ```

6. **Search recursively with a limit on concurrency (e.g., 5 concurrent tasks):**

   ```bash
   ./gogrep -r -j 5 "pattern" /path/to/directory
   ```

## How it Works

`GoGrep` uses regular expressions to find patterns in files. It opens each file, scans through its contents line by line, and applies the specified regular expression. If a match is found, the tool processes the result according to the selected options (e.g., showing line numbers, counting matches, etc.).

For large files or directories, `GoGrep` leverages Go’s concurrency model to process multiple files in parallel, significantly speeding up the search process.

## Concurrency

The `-j` flag allows you to control the maximum number of concurrent search tasks. This is particularly useful when you are searching large directories with many files. The default value is 10, but you can adjust it to suit your needs.

## Performance

`GoGrep` performs efficiently with large files and directories due to its concurrency model, minimizing the time spent on searching by splitting the workload across multiple goroutines. It uses `bufio.Reader` to read large files without token size limitations.
