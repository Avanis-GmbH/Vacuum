package nologger

import "os"

type NoLogger struct{}

func (nl *NoLogger) LogOldFile(f *os.FileInfo, maxAgeInYears uint) error {
	return nil
}

func (nl *NoLogger) LogCopiedFile(f *os.FileInfo, originPath, copyPath string) error {
	return nil
}

func (nl *NoLogger) LogShreddedFile(f *os.FileInfo, originPath string, err error) error {
	return nil
}

func (nl *NoLogger) LogGenericError(err error) error {
	return nil
}

func (nl *NoLogger) LogFailedCopy(originPath, copyPath string, err error) error {
	return nil
}
