package copymachine

func (cm *CopyMachine) copyQueueMasterRoutine() {
	cm.running = true
	jobAmount := len(cm.copyJobs)

	for jobAmount > 0 {
		cm.jobStackMutes.Lock()
		cj := cm.copyJobs[0]

		if jobAmount == 1 {
			cm.copyJobs = make([]*CopyJob, 0)
			cm.jobStackMutes.Unlock()
		} else {
			cm.copyJobs[0] = cm.copyJobs[jobAmount-1]
			cm.copyJobs[jobAmount-1] = nil
			cm.copyJobs = cm.copyJobs[:jobAmount-1]
		}

		cm.jobStackMutes.Unlock()
		cm.performCopyJob(cj)
	}

	cm.running = false
}

func (cm *CopyMachine) performCopyJob(cj *CopyJob) {

}
