package commands

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"

	"github.com/spf13/cobra"

	"github.com/CaiJinKen/mydocker/utils"

	"github.com/CaiJinKen/mydocker/cgroups"

	"github.com/CaiJinKen/mydocker/cgroups/subsystems"

	"github.com/sirupsen/logrus"

	"github.com/CaiJinKen/mydocker/container"
)

var runCommand = &cobra.Command{
	Use:   "run [flags] [command]",
	Short: "run container",
	Long:  "create a container with namespace and cgroup limit",
	Run: func(cmd *cobra.Command, args []string) {

		if tty && detach {
			logrus.Errorf("cannot use ti and d at same time")
			return
		}

		resConf := &subsystems.ResourceConfig{
			MemoryLimit: memory,
			CpuSet:      cpuset,
			CpuShare:    cpushare,
		}

		logrus.Infof("volumeMappings: %v", volumeMappings)

		logrus.Infof("containerName: %s", containerName)

		Run(tty, rm, args, resConf, volumeMappings, envs, containerName)

	},
}

//Run fork a new process to start container
func Run(tty, rm bool, cmdArgs []string, res *subsystems.ResourceConfig, volumeURLs, envs []string, containerName string) {
	logrus.Infof("Run tty %b, args: %v", tty, cmdArgs)
	containerID := container.GenerateUUID()
	parent, writePipe := NewParentProcess(tty, volumeURLs, envs, containerID)
	if parent == nil {
		logrus.Errorf("new parent process error")
		return
	}

	if err := parent.Start(); err != nil {
		logrus.Errorf("parent process start error %v", err)
	}

	//log container info
	logrus.Infof("container pid is %d", parent.Process.Pid)
	containerInfo := &container.Info{
		Pid:       strconv.Itoa(parent.Process.Pid),
		ID:        containerID,
		Name:      containerName,
		Command:   strings.Join(cmdArgs, " "),
		Status:    container.Running,
		Rm:        rm,
		Tty:       tty,
		Volumes:   volumeURLs,
		Resource:  res.String(),
		Envs:      envs,
		CreatedAt: utils.TimeNowString(),
	}

	if err := containerInfo.Save(); err != nil {
		logrus.Errorf("record container info error %v", err)
		return
	}

	//set cgroup limit
	cgroupManager := cgroups.NewCgroupManager(containerID)
	defer cgroupManager.Destroy()
	cgroupManager.Set(res)
	cgroupManager.Apply(parent.Process.Pid)

	//send command with args to init process
	sendInitCommand(cmdArgs, writePipe)

	//tty waiting exit: tty container exit / mydocker stop
	if tty {
		parent.Wait()

		//change container status
		info, err := container.GetContainerInfoByIdentification(containerID)
		if err != nil {
			logrus.Errorf("get container info by identification error %v", err)
		}
		if info != nil {
			info.Status = container.Exit
			info.Pid = ""
			info.Save()
		}

		container.CleanUpWorkspace(containerID)
	}
	//parent exit

}

func sendInitCommand(cmdArgs []string, writePipe *os.File) {
	command := strings.Join(cmdArgs, " ")
	logrus.Infof("command is %s", command)
	writePipe.WriteString(command)
	writePipe.Close()
}

//NewParentProcess fork a new process and pass command and args to new process
func NewParentProcess(tty bool, volumeURLs, envs []string, containerID string) (*exec.Cmd, *os.File) {
	readPipe, writePipe, err := newPipe()
	if err != nil {
		return nil, nil
	}

	//fork a child process, setup namespace & io
	cmd := exec.Command("/proc/self/exe", "init")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWIPC | syscall.CLONE_NEWNET | syscall.CLONE_NEWNS,
		//Credential: &syscall.Credential{Uid: 0, Gid: 0},
	}

	if tty {
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	//extra read pipe to child process
	cmd.ExtraFiles = []*os.File{readPipe}

	//set environment
	cmd.Env = append(os.Environ(), envs...)

	//create container workspace
	mntUrl := fmt.Sprintf(container.MntURL, containerID)
	container.NewWorkSpace(containerID, volumeURLs)

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
