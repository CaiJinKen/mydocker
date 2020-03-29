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
	Pid       string   `json:"pid"`
	ID        string   `json:"id"`
	Name      string   `json:"name"`
	Command   string   `json:"command"`
	CreatedAt string   `json:"created_at"`
	Status    string   `json:"status"`
	Rm        bool     `json:"remove"`
	Tty       bool     `json:"tty"`
	Volumes   []string `json:"volumes"`
	Resource  string   `json:"resource"`
}

const (
	Running = "RUNNING"
	Stop    = "STOPPED"
	Exit    = "EXITED"
)

var (
	DefaultInfoLocation = "/var/lib/mydocker/infos"
	ConfigName          = "config.json"
)

func (containerInfo *Info) Stop() error {
	containerInfo.Pid = ""
	containerInfo.Status = Stop
	return containerInfo.Save()
}

func (containerInfo *Info) Save() error {
	infoBytes, err := json.Marshal(containerInfo)
	if err != nil {
		logrus.Errorf("marshal container info error %v", err)
		return nil
	}

	infoStr := string(infoBytes)
	infoFilePath := containerInfoPath(containerInfo.ID)
	utils.MustPathExist(infoFilePath)

	fileName := fmt.Sprintf("%s/%s", infoFilePath, ConfigName)
	file, err := os.Create(fileName)
	defer file.Close()

	if err != nil {
		logrus.Errorf("create file %s error %v", fileName, err)
		return err
	}

	if _, err = file.WriteString(infoStr); err != nil {
		logrus.Errorf("write file %s error %v", fileName, err)
		return err
	}

	return nil
}

func (containerInfo *Info) Remove() error {
	infoFilePath := containerInfoPath(containerInfo.ID)
	utils.RemoveAll(infoFilePath)
	return nil
}

//GenerateUUID container UUID
func GenerateUUID() string {
	return strings.ReplaceAll(fmt.Sprintf("%s", uuid.NewV4()), "-", "")
}

//RecordContainerInfo record container info
func RecordContainerInfo(containerPID int, containerID, containerName string, cmdArgs []string) (string, error) {
	containerInfo := &Info{
		Pid:       strconv.Itoa(containerPID),
		ID:        containerID,
		Name:      containerName,
		Command:   strings.Join(cmdArgs, " "),
		CreatedAt: utils.TimeNowString(),
		Status:    Running,
	}

	err := containerInfo.Save()

	return containerInfo.ID, err

}

func containerInfoPath(containerID string) string {
	return path.Join(DefaultInfoLocation, containerID)
}

func DeleteContainerInfo(containerID string) {
	if err := utils.RemoveAll(containerInfoPath(containerID)); err != nil {
		logrus.Errorf("delete container info error %v", err)
	}
}

func getContainerInfos() (infos []*Info) {
	files, err := ioutil.ReadDir(DefaultInfoLocation)
	if err != nil {
		logrus.Errorf("read container base path %s error %v", DefaultInfoLocation, err)
		return
	}

	for _, file := range files {
		info, err := getContainerInfoByID(file.Name())
		if info == nil || err != nil {
			continue
		}

		infos = append(infos, info)

	}
	return
}

func ListContainer(all bool) {

	infos := getContainerInfos()

	w := tabwriter.NewWriter(os.Stdout, 12, 1, 3, ' ', 0)
	fmt.Fprintf(w, "ID\tNAME\tPID\tSTATUS\tCOMMAND\tCREATED\n")
	for _, info := range infos {
		if !all && info.Status != Running {
			continue
		}
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\n",
			info.ID,
			info.Name,
			info.Pid,
			info.Status,
			info.Command,
			info.CreatedAt,
		)
	}

	if err := w.Flush(); err != nil {
		logrus.Errorf("flush io error %v", err)
	}

}

func getContainerInfoByID(containerID string) (*Info, error) {
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

	return &info, nil
}

func getContainerInfoByName(containerName string) (*Info, error) {

	infos := getContainerInfos()
	for _, info := range infos {
		if info.Name == containerName {
			return info, nil
		}
	}

	return nil, fmt.Errorf("container %s not found", containerName)
}

func GetContainerInfoByIdentification(containerNameOrID string) (*Info, error) {
	infos := getContainerInfos()
	for _, info := range infos {
		if info.Name == containerNameOrID || info.ID == containerNameOrID {
			return info, nil
		}
	}

	return nil, fmt.Errorf("container %s not found", containerNameOrID)
}

func GetContainerPidByID(containerId string) (string, error) {
	info, err := getContainerInfoByID(containerId)
	if err != nil {
		return "", err
	}
	return info.Pid, nil
}

func GetContainerInfoByName(containerName string) (*Info, error) {
	return getContainerInfoByName(containerName)
}

func GetContainerPidByName(containerName string) (string, error) {
	info, err := getContainerInfoByName(containerName)
	if err != nil {
		return "", err
	}
	return info.Pid, nil
}

func GetContainerPidByIdentification(containerNameOrID string) (string, error) {
	info, err := GetContainerInfoByIdentification(containerNameOrID)
	if err != nil {
		return "", err
	}
	return info.Pid, nil
}
