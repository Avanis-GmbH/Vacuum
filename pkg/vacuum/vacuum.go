package vacuum

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

var RECURSIVE, DRY_RUN, SHRED_ORIGINAL, NO_PROTOCOL bool
var MIN_AGE_IN_YEARS = 11
var TARGET_DIR string

func Clean(rootDir, branchDir string) (*OperationStats, error) {
	entries, err := os.ReadDir(filepath.Join(rootDir, branchDir))
	if err != nil {
		return nil, err
	}

	stats := &OperationStats{}

	for _, f := range entries {

		if f.IsDir() {

			fmt.Printf("Found directory %v in %v \n", f.Name(), filepath.Join(rootDir, branchDir))

			if RECURSIVE {
				stat, err := Clean(rootDir, filepath.Join(branchDir, f.Name()))
				if err != nil {
					//TODO LOG
					stats.Errors++
					continue
				}

				stats.Add(stat)
			}

			continue
		}

		fInfo, err := f.Info()
		if err != nil {
			//TODO log
			stats.Errors++
			continue
		}

		fmt.Printf("Found file %v in %v \n", fInfo.Name(), filepath.Join(rootDir, branchDir))

		if fInfo.ModTime().Year() > time.Now().Year()-MIN_AGE_IN_YEARS {
			continue
		}

		fmt.Printf("File %v in %v is older than %v years: %v", fInfo.Name(), filepath.Join(rootDir, branchDir), fmt.Sprint(MIN_AGE_IN_YEARS), fInfo.ModTime())

		//TODO rest
	}

	return stats, nil
}
