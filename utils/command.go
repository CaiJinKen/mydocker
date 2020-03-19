package utils

import (
	"bytes"
	"os"
	"os/exec"

	"github.com/sirupsen/logrus"
)

//Exec execute command with args
func Exec(command string, arg ...string) (output []byte, err error) {
	logrus.Infof("exec %s %v", command, arg)

	var b bytes.Buffer

	cmd := exec.Command(command, arg...)
	cmd.Stdout = &b
	cmd.Stderr = &b
	if err = cmd.Run(); err != nil {
		logrus.Errorf("exec %s %v error %v", command, arg, err)
		return nil, err
	}
	return b.Bytes(), err
}

//RemoveAll remove path
func RemoveAll(path string) error {
	logrus.Infof("remove dir %s", path)
	if err := os.RemoveAll(path); err != nil {
		logrus.Errorf("remove dir %s error %v", path, err)
		return err
	}
	return nil
}
