package nologger

import "os"

type NoLogger struct{}

func (nl *NoLogger) Init() error {
	return nil
}

func (nl *NoLogger) LogOldFile(f os.FileInfo, maxAgeInYears uint) {

}

func (nl *NoLogger) LogCopiedFile(originPath, copyPath string, copiedBytes uint64) {

}

func (nl *NoLogger) LogShreddedFile(originPath string) {

}

func (nl *NoLogger) LogGenericError(err error) {

}

func (nl *NoLogger) LogFailedCopy(originPath, copyPath string, err error) {

}

func (nl *NoLogger) LogFailedShred(originPath string, err error) {

}

func (nl *NoLogger) Finish() {
}
