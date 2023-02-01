package vacuum

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/Avanis-GmbH/Go-Dust-Vacuum/pkg/copymachine"
	"github.com/Avanis-GmbH/Go-Dust-Vacuum/pkg/logging"
)

var Recursive, DryRun, ShredOriginal bool
var MinAgeInYears = 11
var TargetDir string

var copyJobsEnqueued = 0
var copyJobCountMutex sync.RWMutex

var stats *OperationStats
var statsMutex sync.Mutex

var logger logging.Logger

func PerformCleaning(rootDir string, log logging.Logger) *OperationStats {
	fmt.Printf("Cleaning directory %v\n", rootDir)

	// Abort if copy jobs are still running
	if copyJobsEnqueued > 0 {
		errs := make([]*error, 1)
		err := fmt.Errorf("can't perform cleaning while copy jobs are still running")
		errs[0] = &err

		return &OperationStats{
			Errors: errs,
		}
	}

	logger = log

	copymachine.GetCopyMachine().Dry = DryRun

	// Lock the statistics object and recreate it
	statsMutex.Lock()
	stats = &OperationStats{
		Errors: make([]*error, 0),
	}
	statsMutex.Unlock()

	// Perform cleanDirectory in root directory
	cleanDirectory(rootDir, "", log)

	// Wait for every planned copy job to complete
	for copyJobsEnqueued > 0 {
		time.Sleep(time.Second)
	}

	// Return the statistics
	return stats
}

func cleanDirectory(rootDir, branchDir string, log logging.Logger) {
	// Get all files from current directory
	entries, err := os.ReadDir(filepath.Join(rootDir, branchDir))
	if err != nil {
		err = fmt.Errorf("could not read directory %v: %v", filepath.Join(rootDir, branchDir), err.Error())
		logger.LogGenericError(err)
		appendErrorToStatistics(&err)
	}

	// Iterate over all retreived files
	for _, f := range entries {

		if f.IsDir() {
			fmt.Printf("Found directory %v in %v \n", f.Name(), filepath.Join(rootDir, branchDir))
			// Perform another cleanDirectory call if the file is a directory and the vacuum is running in recursive mode
			if Recursive {
				cleanDirectory(rootDir, filepath.Join(branchDir, f.Name()), log)
			}
			continue
		}

		// Retreive the fileinfo, log and continue if it fails
		fInfo, err := f.Info()
		if err != nil {
			err = fmt.Errorf("could not obtain fileinfo for file %+v: %v", f.Name(), err.Error())
			logger.LogGenericError(err)
			appendErrorToStatistics(&err)
			continue
		}

		fmt.Printf("Found file %v in %v \n", fInfo.Name(), filepath.Join(rootDir, branchDir))

		// Continue if the file is not old enough
		if fInfo.ModTime().Year() > time.Now().Year()-MinAgeInYears {
			continue
		}

		// Enqueue the copy job
		fmt.Printf("File %v in %v is older than %v years: %v \n", fInfo.Name(), filepath.Join(rootDir, branchDir), MinAgeInYears, fInfo.ModTime())
		logger.LogOldFile(fInfo, uint(MinAgeInYears)-1)
		copymachine.GetCopyMachine().EnqueueCopyJob(filepath.Join(rootDir, branchDir, fInfo.Name()), filepath.Join(TargetDir, branchDir, fInfo.Name()), ShredOriginal, copyJobFinishCallback)
		copyJobCountMutex.Lock()
		copyJobsEnqueued++

		copyJobCountMutex.Unlock()
	}

}

func copyJobFinishCallback(cj *copymachine.CopyJob) {

	// Update copy jobs planned counter
	copyJobCountMutex.Lock()
	copyJobsEnqueued--
	copyJobCountMutex.Unlock()

	// Update statistics
	statsMutex.Lock()

	if cj.CopyError != nil {
		stats.Errors = append(stats.Errors, cj.CopyError)
		fmt.Printf("Failed copyjob from %v to %v\n", *cj.FromPath, *cj.ToPath)
		logger.LogFailedCopy(*cj.FromPath, *cj.ToPath, *cj.CopyError)
	} else {
		fmt.Printf("Finished copyjob from %v to %v | Copied %v bytes\n", *cj.FromPath, *cj.ToPath, cj.CopiedBytes)
		stats.CopiedBytes += cj.CopiedBytes
		stats.CopiedFiles++
		logger.LogCopiedFile(*cj.FromPath, *cj.ToPath, cj.CopiedBytes)
	}
	statsMutex.Unlock()

	// Shred original if enabled
	if cj.ShredOnFinish && cj.CopyError == nil {

		fmt.Printf("Shredding file %v \n", *cj.FromPath)
		// Only shred the file if this is not a dry run
		var err error
		if !DryRun {
			err = os.Remove(*cj.FromPath)
			checkAndDeleteEmptyDirectoryTree(filepath.Dir(*cj.FromPath))
		}

		// Update statistics
		statsMutex.Lock()
		if err != nil {
			stats.Errors = append(stats.Errors, &err)
			fmt.Printf("Failed to shred file %v \n", *cj.FromPath)
			logger.LogFailedShred(*cj.FromPath, err)
		} else {
			stats.DeletedFiles++
			fmt.Printf("Shredded file %v \n", *cj.FromPath)
			logger.LogShreddedFile(*cj.FromPath)
		}

		statsMutex.Unlock()
	}

}

func appendErrorToStatistics(err *error) {
	statsMutex.Lock()
	stats.Errors = append(stats.Errors, err)
	statsMutex.Unlock()
}

func checkAndDeleteEmptyDirectoryTree(treeLeafPath string) {
	// Check if the directory is empty after the deletion and delete it if it's empty
	entries, err := os.ReadDir(treeLeafPath)
	if err != nil {
		logger.LogGenericError(err)
	}

	if len(entries) > 1 {
		return
	}

	// Check if directory is empty or only has thumbs db (and delete thumbs db)
	if len(entries) == 1 {
		if entries[0].Name() != "Thumbs.db" {
			return
		}

		err = os.Remove(filepath.Join(treeLeafPath, "Thumbs.db"))
		if err != nil {
			logger.LogFailedShred(treeLeafPath, err)
			return
		}
	}

	// Delete empty directory
	err = os.Remove(treeLeafPath)
	if err != nil {
		logger.LogFailedShred(treeLeafPath, err)
	}

	checkAndDeleteEmptyDirectoryTree(filepath.Dir(treeLeafPath))
}
