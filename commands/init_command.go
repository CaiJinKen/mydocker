package commands

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/CaiJinKen/mydocker/utils"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var initCommand = &cobra.Command{
	Use:    "init",
	Short:  "init container process",
	Hidden: true,
	Run: func(cmd *cobra.Command, args []string) {
		RunContainerInitProcess()
	},
}

//RunContainerInitProcess run container(child process)
func RunContainerInitProcess() error {
	cmdArray := readUserCommand()
	if len(cmdArray) == 0 {
		return fmt.Errorf("run container get user command error, cmd args is empty")
	}

	//pivot root & mount basic point
	setupMount()

	path, err := exec.LookPath(cmdArray[0])
	if err != nil {
		logrus.Errorf("exec %s look path error %v", cmdArray[0], err)
		return err
	}

	logrus.Infof("find path %s", path)
	if err := syscall.Exec(path, cmdArray, os.Environ()); err != nil {
		logrus.Errorf("syscall exec error %v", err)
	}
	return nil
}

//read command and args
func readUserCommand() []string {
	pipe := os.NewFile(uintptr(3), "pipe")
	args, err := ioutil.ReadAll(pipe)
	if err != nil {
		logrus.Errorf("init read pip error %v", err)
		return nil
	}

	return strings.Split(string(args), " ")
}

//setup mount
func setupMount() {
	pwd, err := os.Getwd()
	if err != nil {
		logrus.Errorf("get current pwd error %v", err)
		return
	}

	logrus.Infof("current location is %s", pwd)
	//pivot root
	if err = pivotRoot(pwd); err != nil {
		logrus.Errorf("pivot root error %v", err)
		return
	}

	//set basic info
	setenv()

	//mount basic point
	defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV

	if err := syscall.Mount("sysfs", "/sys", "sysfs", uintptr(defaultMountFlags), ""); err != nil {
		logrus.Errorf("mount /sys error %v", err)
	}

	if err := syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), ""); err != nil {
		logrus.Errorf("mount /proc error %v", err)
	}

	if err := syscall.Mount("tmpfs", "/dev", "tmpfs", syscall.MS_NOSUID|syscall.MS_STRICTATIME, "mode=755"); err != nil {
		logrus.Errorf("mount /dev error %v", err)
	}

}

func pivotRoot(root string) (err error) {

	//make sure current root mount point is private
	if err := syscall.Mount("", "/", "", syscall.MS_PRIVATE|syscall.MS_REC, ""); err != nil {
		return fmt.Errorf("make parent mount private error: %v", err)
	}

	//rebind mount point, make current mount point different the host before
	if err = syscall.Mount(root, root, "bind", syscall.MS_BIND|syscall.MS_REC, ""); err != nil {
		return fmt.Errorf("remount root error %v", err)
	}

	//make old mount point
	pivotDir := filepath.Join(root, ".pivot_root")
	if !utils.MustPathExist(pivotDir) {
		return fmt.Errorf("must path %s exist error", pivotDir)
	}

	//pivot root
	if err = syscall.PivotRoot(root, pivotDir); err != nil {
		return fmt.Errorf("pivot root error %v", err)
	}

	//change work dir
	if err = syscall.Chdir("/"); err != nil {
		return fmt.Errorf("chdir error %v", err)
	}

	pivotDir = filepath.Join("/", ".pivot_root")
	if err = syscall.Unmount(pivotDir, syscall.MNT_DETACH); err != nil {
		return fmt.Errorf("unmount pivot_root dir error %v", err)
	}

	return os.RemoveAll(pivotDir)

}

func setenv() {
	os.Setenv("HOME", "/")
	os.Setenv("PS1", "root@$(hostname):$(pwd)#")
	syscall.Sethostname([]byte("container"))
}
