package fulllogger

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type FullLogger struct {
	Initialized bool
	Finished    bool

	oldFileLog      *os.File
	copiedFileLog   *os.File
	copyErrorLog    *os.File
	shredFileLog    *os.File
	shredErrorLog   *os.File
	genericErrorLog *os.File
}

func (fl *FullLogger) Init() error {
	if fl.Finished {
		return fmt.Errorf("logger already finished and became unusable")
	}

	if fl.Initialized {
		return nil
	}

	// Create the log directory
	dirName := filepath.Join("./", time.Now().String())
	err := os.Mkdir(dirName, 0755)
	if err != nil {
		return fmt.Errorf("could not create log directory: %v", err.Error())
	}

	// Create old files log file
	fl.oldFileLog, err = os.OpenFile(filepath.Join(dirName, "old_files_found"), os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0664)
	if err != nil {
		return fmt.Errorf("could not create old files log file: %v", err.Error())
	}
	// Create copied files log file
	fl.copiedFileLog, err = os.OpenFile(filepath.Join(dirName, "copied_files"), os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0664)
	if err != nil {
		return fmt.Errorf("could not create copied files log file: %v", err.Error())
	}
	// Create copy errors log file
	fl.copyErrorLog, err = os.OpenFile(filepath.Join(dirName, "copy_errors"), os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0664)
	if err != nil {
		return fmt.Errorf("could not create copy errors log file: %v", err.Error())
	}
	// Create shredded files log file
	fl.shredFileLog, err = os.OpenFile(filepath.Join(dirName, "shredded_files"), os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0664)
	if err != nil {
		return fmt.Errorf("could not create shredded files log file: %v", err.Error())
	}
	// Create shred errors log file
	fl.shredErrorLog, err = os.OpenFile(filepath.Join(dirName, "shred_errors"), os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0664)
	if err != nil {
		return fmt.Errorf("could not create shred errors log file: %v", err.Error())
	}
	// Create generic errors log file
	fl.genericErrorLog, err = os.OpenFile(filepath.Join(dirName, "generic_errors"), os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0664)
	if err != nil {
		return fmt.Errorf("could not create generic errors log file: %v", err.Error())
	}

	fl.Initialized = true
	return nil
}

func (fl *FullLogger) LogOldFile(f os.FileInfo, maxAgeInYears uint) {
	if !fl.Initialized || fl.Finished {
		fmt.Println("Logger not initialized")
		return
	}

	_, err := fl.oldFileLog.WriteString(fmt.Sprintf("Found file %v being older than %v years \n", f.Name(), maxAgeInYears))
	if err != nil {
		fmt.Println("Could not write to old file log: " + err.Error())
	}
}

func (fl *FullLogger) LogCopiedFile(originPath, copyPath string, copiedBytes uint64) {
	if !fl.Initialized || fl.Finished {
		fmt.Println("Logger not initialized")
		return
	}

	_, err := fl.copiedFileLog.WriteString(fmt.Sprintf("Copied %v bytes from %v to %v \n", copiedBytes, originPath, copyPath))
	if err != nil {
		fmt.Println("Could not write to copied file log: " + err.Error())
	}
}

func (fl *FullLogger) LogFailedCopy(originPath, copyPath string, err error) {
	if !fl.Initialized || fl.Finished {
		fmt.Println("Logger not initialized")
		return
	}

	_, err = fl.copyErrorLog.WriteString("Could not copy file " + originPath + " to " + copyPath + ": " + err.Error() + "\n")
	if err != nil {
		fmt.Println("Could not write to copy file error log: " + err.Error())
	}
}

func (fl *FullLogger) LogShreddedFile(originPath string) {
	if !fl.Initialized || fl.Finished {
		fmt.Println("Logger not initialized")
		return
	}

	_, err := fl.shredFileLog.WriteString("Shredded file  " + originPath + "\n")
	if err != nil {
		fmt.Println("Could not write to shredded file log: " + err.Error())
	}
}

func (fl *FullLogger) LogFailedShred(originPath string, err error) {
	if !fl.Initialized || fl.Finished {
		fmt.Println("Logger not initialized")
		return
	}

	_, err = fl.shredErrorLog.WriteString("Could not shred file  " + originPath + ": " + err.Error() + "\n")
	if err != nil {
		fmt.Println("Could not write to shred error log: " + err.Error())
	}
}

func (fl *FullLogger) LogGenericError(err error) {
	if !fl.Initialized || fl.Finished {
		fmt.Println("Logger not initialized")
		return
	}

	_, err = fl.genericErrorLog.WriteString(err.Error() + "\n")
	if err != nil {
		fmt.Println("Could not write to generic error log: " + err.Error())
	}
}

func (fl *FullLogger) Finish() {
	err := fl.oldFileLog.Close()
	if err != nil {
		fmt.Printf("Error closing old file log file: %v", err.Error())
	}

	err = fl.copiedFileLog.Close()
	if err != nil {
		fmt.Printf("Error closing copied file log file: %v", err.Error())
	}

	err = fl.copyErrorLog.Close()
	if err != nil {
		fmt.Printf("Error closing copy error log file: %v", err.Error())
	}

	err = fl.shredFileLog.Close()
	if err != nil {
		fmt.Printf("Error closing shred file log file: %v", err.Error())
	}

	err = fl.shredErrorLog.Close()
	if err != nil {
		fmt.Printf("Error closing shred error log file: %v", err.Error())
	}

	err = fl.genericErrorLog.Close()
	if err != nil {
		fmt.Printf("Error closing generic error log file: %v", err.Error())
	}
}
