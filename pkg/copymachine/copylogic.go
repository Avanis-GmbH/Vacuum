package copymachine

import "fmt"

func (cm *CopyMachine) copyQueueMasterRoutine() {
	cm.running = true

	// Get the amount of current queued jobs
	jobAmount := len(cm.copyJobs)

	for jobAmount > 0 {
		cm.copyJobStackMutes.Lock()
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

//TODO implement
func (cm *CopyMachine) performCopyJob(cj *CopyJob) {
	err := fmt.Errorf("could not copy file from %+v to %+v: not implemented", *cj.FromPath, *cj.ToPath)
	cj.CopyError = &err

	cj.FinishCallBack(cj)
}
