package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	ExplainSource  = "Path to the source directory from which files will be moved."
	ExplainTarget  = "Path to the target directory where files will be archived."
	ExplainDry     = "Perform a dry run without executing any file operations."
	ExplainShred   = "Delete the original file after copying it to the target directory."
	ExplainRecurse = "Recursively include all subdirectories."
	ExplainAge     = "Minimum file age in years to consider for archiving."
)

var (
	dry, shred, recurse bool
	source, target      string
	age                 int
)

func init() {
	flag.StringVar(&source, "source", "", ExplainSource)
	flag.StringVar(&target, "target", "", ExplainTarget)
	flag.BoolVar(&dry, "dry", false, ExplainDry)
	flag.BoolVar(&shred, "shred", false, ExplainShred)
	flag.BoolVar(&recurse, "recurse", false, ExplainRecurse)
	flag.IntVar(&age, "age", 3, ExplainAge)
}

func main() {
	err := validateFlags()
	if err != nil {
		log.Fatal(err)
		return
	}
	process(source)
}

func validateFlags() error {
	flag.Parse()
	if source == "" || target == "" {
		flag.Usage()
		return fmt.Errorf("missing required parameters")
	}
	err := checkDirectory(source, "source")
	if err != nil {
		return err
	}
	err = checkDirectory(target, "target")
	if err != nil {
		return err
	}
	return nil
}

func checkDirectory(path string, name string) error {
	dirInfo, err := os.Stat(path)
	if os.IsNotExist(err) || !dirInfo.IsDir() {
		return fmt.Errorf("invalid %s directory", name)
	}
	testFile := filepath.Join(path, "test")
	err = os.WriteFile(testFile, []byte("test"), 0644)
	if err != nil {
		return fmt.Errorf("unwritable %s directory", name)
	}
	os.Remove(testFile)
	return nil
}

func process(dir string) {
	files, err := os.ReadDir(dir)
	if err != nil {
		log.Printf("error reading directory: %s\n", dir)
	}
	for _, file := range files {
		path := filepath.Join(dir, file.Name())
		if file.IsDir() {
			if !recurse {
				continue
			}
			go process(path)
		}
		if !shouldProcessFile(file, path) {
			continue
		}
		err := handleFile(path)
		if err != nil {
			log.Printf("error handling file: %s\n", path)
		}
	}
}

func shouldProcessFile(file os.DirEntry, path string) bool {
	fileInfo, err := file.Info()
	if err != nil {
		log.Printf("error getting file info: %s\n", path)
		return false
	}
	return time.Since(fileInfo.ModTime()) >= time.Duration(age*365*24)*time.Hour
}

func handleFile(path string) error {
	targetPath := filepath.Join(target, strings.TrimPrefix(path, source))
	if dry {
		fmt.Printf("moving file: %s to %s\n", path, targetPath)
		return nil
	}
	err := moveFile(path, targetPath)
	if err != nil {
		return fmt.Errorf("error moving file: %w", err)
	}
	if !shred {
		return nil
	}
	err = shredFile(path)
	if err != nil {
		return fmt.Errorf("error shredding file: %w", err)
	}
	return nil
}

func moveFile(src, dst string) error {
	err := os.MkdirAll(filepath.Dir(dst), os.ModePerm)
	if err != nil {
		return err
	}
	input, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	err = os.WriteFile(dst, input, 0644)
	if err != nil {
		return err
	}
	log.Printf("file moved from %s to %s\n", src, dst)
	return nil
}

func shredFile(path string) error {
	err := os.Remove(path)
	if err != nil {
		return err
	}
	log.Printf("shredded %s\n", path)
	return nil
}
