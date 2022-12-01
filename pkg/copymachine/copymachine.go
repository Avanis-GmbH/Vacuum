package copymachine

import (
	"fmt"
	"sync"
)

type CopyMachine struct {
	copyJobs      []*CopyJob
	jobStackMutes sync.Mutex
	running       bool
	Dry           bool
}

type CopyJob struct {
	FromPath       *string
	ToPath         *string
	ShredOnFinish  bool
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

	cj := &CopyJob{
		FromPath:       &originalFullPath,
		ToPath:         &copyFullPath,
		ShredOnFinish:  shredOnFinish,
		CopyError:      nil,
		FinishCallBack: finishCallback,
	}

	qm.jobStackMutes.Lock()
	defer qm.jobStackMutes.Unlock()
	qm.copyJobs = append(qm.copyJobs, cj)
	fmt.Printf("Enqueued copy job %+v \n", cj)

	if !qm.running {
		go qm.copyQueueMasterRoutine()
	}
}
