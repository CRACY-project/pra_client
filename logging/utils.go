package logging

import (
	"os"
)

func MoveOldLogFiles(oldPath, newPath string) {
	// Check if old path exists
	if _, err := os.Stat(oldPath); os.IsNotExist(err) {
		return
	}

	// Move old log files to new path
	err := os.Rename(oldPath, newPath)
	if err != nil {
		Log(WARNING, "Failed to move old log files: ", err)
	} else {
		Log(DEBUG, "Moved old log files successfully!")
	}
}
