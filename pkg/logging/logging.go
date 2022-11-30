package logging

import "os"

type Logger interface {
	LogOldFile(f *os.FileInfo, maxAgeInYears uint) error
	LogCopiedFile(f *os.FileInfo, originPath, copyPath string) error
	LogShreddedFile(f *os.FileInfo, originPath string, err error) error
}
