package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Avanis-GmbH/Go-Dust-Vacuum/pkg/fulllogger"
	"github.com/Avanis-GmbH/Go-Dust-Vacuum/pkg/logging"
	"github.com/Avanis-GmbH/Go-Dust-Vacuum/pkg/nologger"
	"github.com/Avanis-GmbH/Go-Dust-Vacuum/pkg/vacuum"
)

const VERSION = "1.0.0-beta"

var runHelp bool
var rootDir string
var noLog bool

func main() {
	showStartupBanner()
	parseFlags()

	if runHelp {
		showHelp()
		os.Exit(0)
	}

	// Set up logging
	var l logging.Logger
	if !noLog {
		l = &fulllogger.FullLogger{}
		err := l.Init()
		if err != nil {
			panic(err)
		}
	} else {
		l = &nologger.NoLogger{}
	}

	if vacuum.MinAgeInYears < 0 {
		l.LogGenericError(fmt.Errorf("minimum age cannot be lower than 0 years! Parsed value: %v", vacuum.MinAgeInYears))
		os.Exit(1)
	}

	if vacuum.MinAgeInYears == 0 {
		fmt.Println("=====================================================================")
		fmt.Println(">>>ATTENTION! SETTING THE MINIMUM FILE AGE TO 0 WOULD AFFECT EVERY FILE IN THE ROOT DIRECTORY!<<<")
		fmt.Printf("===================================================================== \n\n")
		fmt.Println("Do you want to continue? (Type 'yes' to continue. Any other input will exit the software)")
		var answer string
		_, err := fmt.Scanln(&answer)
		if err != nil {
			panic(err)
		}

		if answer != "yes" {
			fmt.Println("Aborted!")
		}

		fmt.Println("Continuing despite warning...")
	}

	s := vacuum.PerformCleaning(rootDir, l)

	l.Finish()

	fmt.Printf("Finished process with following statistics: \n")
	fmt.Printf("Copied files: %v \n", s.CopiedFiles)
	fmt.Printf("Copied bytes: %v \n", s.CopiedBytes)
	fmt.Printf("Deleted files: %v \n", s.DeletedFiles)
	fmt.Printf("Number of errors: %v \n", len(s.Errors))
}

func showStartupBanner() {
	fmt.Println("=============================================================")
	fmt.Printf("Go Dust Vacuum v%v - Created by Simon Nils Rach \n", VERSION)
	fmt.Printf("=============================================================\n\n")
}

func parseFlags() {
	flag.Usage = func() {
		showHelp()
	}
	flag.BoolVar(&vacuum.Recursive, "r", true, "Should all subdirectories be included.")
	flag.BoolVar(&vacuum.DryRun, "dry", false, "Should the application perform a dry run without any io operations?")
	flag.BoolVar(&vacuum.ShredOriginal, "shred", false, "Should the original file get deleted after copy?")
	flag.IntVar(&vacuum.MinAgeInYears, "older-than", 11, "How old the last edit of a file should be (in years) to consider it for archiving.")
	flag.BoolVar(&runHelp, "help", false, "Shows usage information about this software")
	flag.BoolVar(&noLog, "nolog", false, "If no log files should be written for the process. Use at own risk only!")
	flag.StringVar(&rootDir, "root-dir", ">INVALID<", "The root directory which should be scanned for old files. [REQUIRED]")
	flag.StringVar(&vacuum.TargetDir, "target-dir", ">INVALID<", "The target directory where the old files should be copied to. [REQUIRED]")
	flag.Parse()
}

func showHelp() {
	fmt.Printf("Scans a chosen root directory for old files and copies them to a target directory for archiving. \n \n")

	fmt.Printf("Usage: \ngo-dust-vacuum -root-dir=<root directory> -target-dir=<target directory> [additional flags]\n \n")
	fmt.Printf("Available flags: \n")
	flag.PrintDefaults()
}
