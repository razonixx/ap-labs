package main

import (
	"fmt"
	"os"
	"path/filepath"
)

// scanDir stands for the directory scanning implementation
func scanDir(dir string) error {

	var files []os.FileInfo
	var numFiles = 0
	var dirs = 0
	var symLink = 0
	var otherFiles = 0

	root := os.Args[1]
	walk := filepath.Walk(root, func(fp string, fi os.FileInfo, err error) error {
		files = append(files, fi)
		return nil
	})
	if walk != nil {
		fmt.Printf("Error in walk function\n")
		return nil
	}
	for _, file := range files {
		numFiles++
		if file.IsDir() {
			dirs++
		}
		if file.Mode()&os.ModeSymlink != 0 {
			symLink++
		}
	}
	otherFiles = numFiles - dirs - symLink
	fmt.Printf("Path: %s\n", dir)
	fmt.Printf("Directories: %d\n", dirs)
	fmt.Printf("Symbolic Links: %d\n", symLink)
	fmt.Printf("Other Files: %d\n", otherFiles)
	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: ./dir-scan <directory>")
		os.Exit(1)
	}
	scanDir(os.Args[1])
}
