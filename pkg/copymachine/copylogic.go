package copymachine

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func (cm *CopyMachine) copyQueueMasterRoutine() {
	cm.running = true

	// Get the amount of current queued jobs
	jobAmount := len(cm.copyJobs)

	for jobAmount > 0 {
		cm.copyJobStackMutes.Lock()
		jobAmount = len(cm.copyJobs)
		cj := cm.copyJobs[0]

		if jobAmount == 1 {
			cm.copyJobs = make([]*CopyJob, 0)
		} else {
			cm.copyJobs[0] = cm.copyJobs[jobAmount-1]
			cm.copyJobs[jobAmount-1] = nil
			cm.copyJobs = cm.copyJobs[:jobAmount-1]
		}

		jobAmount = len(cm.copyJobs)
		cm.copyJobStackMutes.Unlock()
		cm.performCopyJob(cj)
	}

	cm.running = false
}

func (cm *CopyMachine) performCopyJob(cj *CopyJob) {

	fmt.Printf("Performing copyjob %+v \n", cj)

	// Abort without any errors when doing a dry run
	if cm.Dry {
		cj.FinishCallBack(cj)
		return
	}

	// Create the target directory if it does not exist
	err := os.MkdirAll(filepath.Dir(*cj.ToPath), 0666)
	if err != nil {
		cj.CopyError = &err
		cj.FinishCallBack(cj)
		return
	}

	// Open the source file
	sourceF, err := os.Open(*cj.FromPath)
	if err != nil {
		cj.CopyError = &err
		cj.FinishCallBack(cj)
		return
	}
	defer sourceF.Close()

	// Create the destination file if not exist
	destF, err := os.Create(*cj.ToPath)
	if err != nil {
		cj.CopyError = &err
		cj.FinishCallBack(cj)
		return
	}
	defer destF.Close()

	// Copy the file content
	copBytes, err := io.Copy(destF, sourceF)
	if err != nil {
		cj.CopyError = &err
		cj.CopiedBytes = uint64(copBytes)
		cj.FinishCallBack(cj)
		return
	}

	// Finalize the copy job
	err = destF.Sync()
	if err != nil {
		cj.CopyError = &err
	}

	// Close the source file
	err = sourceF.Close()
	if err != nil {
		fmt.Printf("Error while closing file after copy: %v \n", err.Error())
	}

	cj.FinishCallBack(cj)
}
