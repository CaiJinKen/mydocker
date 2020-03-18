package utils

import "os"

//FileExist return file whether exist
func FileExist(path string) (ok bool) {
	if _, err := os.Stat(path); err == nil {
		return true
	}
	return false
}

//MustPathExist make sure path is exist
func MustPathExist(path string) (ok bool) {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(path, 077)
			return err == nil
		} else {
			return false
		}
	}

	return true
}
