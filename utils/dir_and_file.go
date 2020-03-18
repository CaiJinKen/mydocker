package utils

import "os"

//FileExist return file whether exist
func FileExist(path string) (ok bool) {
	if _, err := os.Stat(path); err == nil {
		return true
	}
	return false
}

//PathExist path whether exist, if not and authCreate is bool then return create result
func PathExist(path string, autoCreate bool) bool {
	if autoCreate {
		return MustPathExist(path)
	}
	return FileExist(path)
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
