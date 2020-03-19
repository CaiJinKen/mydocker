package subsystems

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/CaiJinKen/mydocker/utils"

	"github.com/sirupsen/logrus"
)

func findCgroupMountPoint(subsystem string) string {
	f, err := os.Open("/proc/self/mountinfo")
	if err != nil {
		logrus.Errorf("open /proc/self/mountinfo error %v", err)
		return ""
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		txt := scanner.Text()
		fields := strings.Split(txt, " ")
		for _, opt := range strings.Split(fields[len(fields)-1], ",") {
			if opt == subsystem {
				return fields[4]
			}
		}
	}
	return ""
}

func getCgroupPath(subsystem, cgroupPath string, autoCreate bool) (string, error) {
	cgroupRoot := findCgroupMountPoint(subsystem)
	resultPath := path.Join(cgroupRoot, cgroupPath)
	if utils.PathExist(resultPath, autoCreate) {
		return resultPath, nil
	}
	return "", fmt.Errorf("cgroup path error")
}

//set cgroup
func cgroupSet(cgroupName, cgroupPath, fileName, value string) error {
	if value == "" {
		return nil
	}
	subsysCgroupPath, err := getCgroupPath(cgroupName, cgroupPath, true)
	if err != nil {
		return err
	}

	if value != "" {
		if err = ioutil.WriteFile(path.Join(subsysCgroupPath, fileName), []byte(value), 0644); err != nil {
			return fmt.Errorf("set cgroup %s error %v", fileName, err)
		}
	}

	return nil
}

//apply pid to cgroup
func cgroupApply(cgroupName, cgroupPath string, pid int) error {
	subsysCgroupPath, err := getCgroupPath(cgroupName, cgroupPath, false)
	if err != nil {
		return err
	}

	if err = ioutil.WriteFile(path.Join(subsysCgroupPath, "tasks"), []byte(strconv.Itoa(pid)), 0644); err != nil {
		return fmt.Errorf("set cgroup proc error %v", err)
	}
	return nil
}

//remove cgroup
func cgroupRemove(cgroupName, cgroupPath string) error {
	subsysCgroupPath, err := getCgroupPath(cgroupName, cgroupPath, false)
	if err != nil {
		return err
	}

	return os.Remove(subsysCgroupPath)
}
