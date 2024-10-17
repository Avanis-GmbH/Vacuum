package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
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

type stats struct {
	archived uint
	errors   uint
	deleted  uint
}

var (
	source  string
	target  string
	age     int
	dry     bool
	shred   bool
	recurse bool
	result  stats
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
	err := processFlags()
	if err != nil {
		log.Fatal(err)
	}
	err = filepath.Walk(source, processEntry)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("archived:", result.archived)
	fmt.Println("errors:", result.errors)
	fmt.Println("deleted:", result.deleted)
}

func processFlags() error {
	flag.Parse()
	if source == "" {
		flag.Usage()
		return fmt.Errorf("source directory is required")
	}
	if target == "" {
		flag.Usage()
		return fmt.Errorf("target directory is required")
	}
	stat, err := os.Stat(source)
	if err != nil {
		return err
	}
	if !stat.IsDir() {
		return fmt.Errorf("source is not a directory")
	}
	stat, err = os.Stat(target)
	if err != nil {
		return err
	}
	if !stat.IsDir() {
		return fmt.Errorf("target is not a directory")
	}
	if age <= 0 {
		return fmt.Errorf("age must be a positive integer")
	}
	if dry {
		fmt.Println("dry run enabled no files will be moved")
	}
	return nil
}

func processEntry(path string, info os.FileInfo, err error) error {
	if err != nil {
		fmt.Printf("error accessing path %s: %v\n", path, err)
		result.errors++
		return nil
	}
	if info.IsDir() && !recurse && path != source {
		return filepath.SkipDir
	}
	if info.ModTime().Before(time.Now().AddDate(-age, 0, 0)) {
		processFile(path)
	}
	return nil
}

func processFile(path string) {
	relPath, err := filepath.Rel(source, path)
	if err != nil {
		fmt.Printf("error determining relative path for %s: %v\n", path, err)
		result.errors++
		return
	}
	targetPath := filepath.Join(target, relPath)
	if dry {
		fmt.Printf("dry: would archive file: %s to %s\n", path, targetPath)
		result.archived++
		return
	}
	err = os.MkdirAll(filepath.Dir(targetPath), os.ModePerm)
	if err != nil {
		fmt.Printf("error creating target directory for %s: %v\n", targetPath, err)
		result.errors++
		return
	}
	err = copyFile(path, targetPath)
	if err != nil {
		fmt.Printf("error copying file %s to %s: %v\n", path, targetPath, err)
		result.errors++
		return
	}
	fmt.Printf("archived file: %s to %s\n", path, targetPath)
	result.archived++
	if !shred {
		return
	}
	err = os.Remove(path)
	if err != nil {
		fmt.Printf("error deleting original file %s: %v\n", path, err)
		result.errors++
		return
	}
	fmt.Printf("deleted original file: %s\n", path)
	result.deleted++
}

func copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()
	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()
	_, err = srcFile.Stat()
	if err != nil {
		return err
	}
	_, err = dstFile.ReadFrom(srcFile)
	return err
}
