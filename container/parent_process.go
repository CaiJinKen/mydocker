package container

import (
	"os"
	"os/exec"
	"syscall"

	"github.com/sirupsen/logrus"
)

//NewParentProcess fork a new process and pass command and args to new process
func NewParentProcess(tty bool) (*exec.Cmd, *os.File) {
	readPipe, writePipe, err := newPipe()
	if err != nil {
		return nil, nil
	}

	//fork a child process, setup namespace & io
	cmd := exec.Command("/proc/self/exe", "init")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWIPC | syscall.CLONE_NEWNET | syscall.CLONE_NEWNS,
		Credential: &syscall.Credential{Uid: 0, Gid: 0},
	}

	if tty {
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	//extra read pipe to child process
	cmd.ExtraFiles = []*os.File{readPipe}

	//create container workspace
	rootUrl := "/root/"
	mntUrl := "/root/mnt"
	NewWorkSpace(rootUrl, mntUrl)

	//setup work dir
	cmd.Dir = mntUrl

	return cmd, writePipe
}

func newPipe() (*os.File, *os.File, error) {
	read, write, err := os.Pipe()
	if err != nil {
		logrus.Errorf("create new pipe error %v", err)
		return nil, nil, err
	}
	return read, write, nil
}
