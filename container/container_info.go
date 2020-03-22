package container

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/sirupsen/logrus"

	"github.com/CaiJinKen/mydocker/utils"

	uuid "github.com/satori/go.uuid"
)

//container info
type Info struct {
	Pid       string `json:"pid"`
	ID        string `json:"id"`
	Name      string `json:"name"`
	Command   string `json:"command"`
	CreatedAt string `json:"created_at"`
	Status    string `json:"status"`
}

const (
	Running = "RUNNING"
	Stop    = "STOP"
	Exit    = "EXIT"
)

var (
	DefaultInfoLocation = "/var/run/mydocker/"
	ConfigName          = "config.json"
)

func generateUUID() string {
	return strings.ReplaceAll(fmt.Sprintf("%s", uuid.NewV4()), "-", "")
}

func RecordContainerInfo(containerPID int, containerName string, cmdArgs []string) (string, error) {
	containerInfo := &Info{
		Pid:       strconv.Itoa(containerPID),
		ID:        generateUUID(),
		Name:      containerName,
		Command:   strings.Join(cmdArgs, " "),
		CreatedAt: utils.TimeNowString(),
		Status:    Running,
	}

	infoBytes, err := json.Marshal(containerInfo)
	if err != nil {
		logrus.Errorf("marshal container info error %v", err)
		return "", nil
	}

	infoStr := string(infoBytes)
	infoFilePath := containerInfoPath(containerInfo.ID)
	utils.MustPathExist(infoFilePath)

	fileName := fmt.Sprintf("%s/%s", infoFilePath, ConfigName)
	file, err := os.Create(fileName)
	defer file.Close()

	if err != nil {
		logrus.Errorf("create file %s error %v", fileName, err)
		return "", err
	}

	if _, err = file.WriteString(infoStr); err != nil {
		logrus.Errorf("write file %s error %v", fileName, err)
		return "", err
	}

	return containerInfo.ID, nil

}

func containerInfoPath(containerID string) string {
	return path.Join(DefaultInfoLocation, containerID)
}

func DeleteContainerInfo(containerID string) {
	if err := utils.RemoveAll(containerInfoPath(containerID)); err != nil {
		logrus.Errorf("delete container info error %v", err)
	}
}

func ListContainer(all bool) {
	files, err := ioutil.ReadDir(DefaultInfoLocation)
	if err != nil {
		logrus.Errorf("read container base path %s error %v", DefaultInfoLocation, err)
		return
	}

	var infos []*Info
	for _, file := range files {
		info, err := getContainerInfo(file.Name(), all)
		if info == nil || err != nil {
			continue
		}
		infos = append(infos, info)
	}

	w := tabwriter.NewWriter(os.Stdout, 12, 1, 3, ' ', 0)
	fmt.Fprintf(w, "ID\tNAME\tPID\tSTATUS\tCOMMAND\tCREATED\n")
	for _, info := range infos {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\n",
			info.ID,
			info.Name,
			info.Pid,
			info.Status,
			info.Command,
			info.CreatedAt,
		)
	}

	if err = w.Flush(); err != nil {
		logrus.Errorf("flush io error %v", err)
	}

}

func getContainerInfo(containerID string, all bool) (*Info, error) {
	//avoid . or .. file
	if len(containerID) < 3 {
		return nil, nil
	}

	fileName := fmt.Sprintf("%s/%s/%s", DefaultInfoLocation, containerID, ConfigName)
	bytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		logrus.Errorf("read file %s error %v", fileName, err)
		return nil, err
	}

	var info Info
	if err = json.Unmarshal(bytes, &info); err != nil {
		logrus.Errorf("unmarshal info error %v", err)
		return nil, err
	}
	if !all && info.Status != Running {
		return nil, nil
	}

	return &info, nil
}
