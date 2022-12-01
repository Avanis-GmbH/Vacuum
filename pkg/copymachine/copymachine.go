package copymachine

import (
	"fmt"
	"sync"
)

type CopyMachine struct {
	copyJobs          []*CopyJob
	copyJobStackMutes sync.Mutex
	running           bool
	Dry               bool
}

type CopyJob struct {
	FromPath       *string
	ToPath         *string
	ShredOnFinish  bool
	CopiedBytes    uint64
	CopyError      *error
	FinishCallBack func(cj *CopyJob)
}

var (
	copySync      sync.Once
	copySingleton *CopyMachine
)

func GetCopyMachine() *CopyMachine {
	copySync.Do(func() {
		copySingleton = &CopyMachine{
			copyJobs: make([]*CopyJob, 0),
		}
	})
	return copySingleton
}

func (qm *CopyMachine) EnqueueCopyJob(originalFullPath, copyFullPath string, shredOnFinish bool, finishCallback func(cj *CopyJob)) {

	// Create copy job
	cj := &CopyJob{
		FromPath:       &originalFullPath,
		ToPath:         &copyFullPath,
		ShredOnFinish:  shredOnFinish,
		CopyError:      nil,
		FinishCallBack: finishCallback,
	}

	// Append copy job to copy job stack
	qm.copyJobStackMutes.Lock()
	defer qm.copyJobStackMutes.Unlock()
	qm.copyJobs = append(qm.copyJobs, cj)
	fmt.Printf("Enqueued copy job %+v \n", cj)

	// Start the queuemaster if it's not running
	if !qm.running {
		go qm.copyQueueMasterRoutine()
	}
}
