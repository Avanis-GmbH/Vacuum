package logging

import "os"

type Logger interface {
	Init() error
	LogOldFile(f *os.FileInfo, maxAgeInYears uint)
	LogCopiedFile(originPath, copyPath string, copiedBytes uint64)
	LogShreddedFile(originPath string)
	LogGenericError(err error)
	LogFailedCopy(originPath, copyPath string, err error)
	LogFailedShred(originPath string, err error)
	Finish() error
}
