package helpers

import (
	"io/ioutil"
	"os"
	"strings"
)

func CreateProfiles(fullFilePath string, file []byte, perm os.FileMode) error {
	var (
		err      error
		f        = strings.Split(fullFilePath, "/")
		filepath = strings.Join(f[:(len(f)-1)], "/")
	)

	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		err = os.MkdirAll(filepath, 0755)
	}
	_ = ioutil.WriteFile(fullFilePath, file, perm)
	return err
}
