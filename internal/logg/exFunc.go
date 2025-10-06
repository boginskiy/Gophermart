package logg

import (
	"os"
)

func CreateFile(pathToFile string) (*os.File, error) {
	return os.OpenFile(pathToFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
}
