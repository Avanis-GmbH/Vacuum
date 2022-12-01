package vacuum

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/Avanis_GmbH/Go-Dust-Vacuum/pkg/copymachine"
	"github.com/Avanis_GmbH/Go-Dust-Vacuum/pkg/logging"
)

var RECURSIVE, DRY_RUN, SHRED_ORIGINAL, NO_PROTOCOL bool
var MIN_AGE_IN_YEARS = 11
var TARGET_DIR string

var copyJobsPlanned = 0
var copyJobCountMutex sync.RWMutex

var stats *OperationStats
var statsMutex sync.Mutex

func PerformCleaning(rootDir string, log logging.Logger) *OperationStats {
	fmt.Printf("Cleaning directory %v", rootDir)

	// Abort if copy jobs are still running
	if copyJobsPlanned > 0 {
		errs := make([]*error, 1)
		err := fmt.Errorf("can't perform cleaning while copy jobs are still running")
		errs[0] = &err

		return &OperationStats{
			Errors: errs,
		}
	}

	// Lock the statistics object and recreate it
	statsMutex.Lock()
	stats = &OperationStats{
		Errors: make([]*error, 0),
	}
	statsMutex.Unlock()

	// Perform cleanDirectory in root directory
	cleanDirectory(rootDir, "", log)

	// Wait for every planned copy job to complete
	for copyJobsPlanned > 0 {
		time.Sleep(time.Second)
	}

	// Return the statistics
	return stats
}

func cleanDirectory(rootDir, branchDir string, log logging.Logger) {
	// Get all files from current directory
	entries, err := os.ReadDir(filepath.Join(rootDir, branchDir))
	if err != nil {
		appendErrorToStatistics(&err)
	}

	// Iterate over all retreived files
	for _, f := range entries {

		if f.IsDir() {
			fmt.Printf("Found directory %v in %v \n", f.Name(), filepath.Join(rootDir, branchDir))

			// Perform another cleanDirectory call if the file is a directory and the vacuum is running in recursive mode
			if RECURSIVE {
				cleanDirectory(rootDir, filepath.Join(branchDir, f.Name()), log)
				if err != nil {
					log.LogGenericError(err)
					appendErrorToStatistics(&err)
					continue
				}
			}

			continue
		}

		// Retreive the fileinfo, log and continue if it fails
		fInfo, err := f.Info()
		if err != nil {
			log.LogGenericError(err)
			appendErrorToStatistics(&err)
			continue
		}

		fmt.Printf("Found file %v in %v \n", fInfo.Name(), filepath.Join(rootDir, branchDir))

		// Continue if the file is not old enough
		if fInfo.ModTime().Year() > time.Now().Year()-MIN_AGE_IN_YEARS {
			continue
		}

		fmt.Printf("File %v in %v is older than %v years: %v", fInfo.Name(), filepath.Join(rootDir, branchDir), fmt.Sprint(MIN_AGE_IN_YEARS), fInfo.ModTime())

		//TODO rest

	}

}

// TODO implement
func copyJobFinishCallback(cj *copymachine.CopyJob) {

	// Update statistics
	statsMutex.Lock()
	stats.CopiedBytes += cj.CopiedBytes
	stats.CopiedFiles++
	if cj.CopyError != nil {
		stats.Errors = append(stats.Errors, cj.CopyError)
	}
	statsMutex.Unlock()

	// Shred original if enabled
	if SHRED_ORIGINAL {
		// TODO shred original
	}

	// Update copy jobs planned counter
	copyJobCountMutex.Lock()
	copyJobsPlanned--
	copyJobCountMutex.Unlock()
}

func appendErrorToStatistics(err *error) {
	statsMutex.Lock()
	stats.Errors = append(stats.Errors, err)
	statsMutex.Unlock()
}
