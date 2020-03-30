package container

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/CaiJinKen/mydocker/utils"
)

const (
	RootURL       = "/root/"
	WriteLayerURL = "/var/lib/mydocker/%s/writeLayer"
	MntURL        = "/run/mydocker/%s/mnt"
	_rootfsFile   = "ubuntu-base-16.04.1-base-amd64.tar.gz"
	_rootfsURL    = "http://cdimage.ubuntu.com/ubuntu-base/releases/16.04/release/ubuntu-base-16.04.1-base-amd64.tar.gz"
)

//NewWorkSpace create container workspace
func NewWorkSpace(containerID string, volumeURLs []string) {
	writeLayerURL, mntURL := getContainerURL(containerID)

	//CreateReadOnlyLayer(RootURL)
	CreateWriteLayer(writeLayerURL)
	CreateMountPoint(RootURL, writeLayerURL, mntURL, containerID)
	CreateMountVolumePoint(mntURL, volumeURLs)
}

func getContainerURL(containerID string) (string, string) {
	return fmt.Sprintf(WriteLayerURL, containerID), fmt.Sprintf(MntURL, containerID)
}

func CreateReadOnlyLayer(rootURL string) {
	rootfs := fmt.Sprintf("%s%s", rootURL, "rootfs")
	utils.MustPathExist(rootfs)
	files, err := ioutil.ReadDir(rootfs)
	if err == nil && len(files) > 0 {
		logrus.Infof("%s is not empty", rootfs)
		return
	}

	rootZipFile := fmt.Sprintf("%s%s", rootURL, _rootfsFile)
	if _, err := os.Stat(rootZipFile); os.IsNotExist(err) {
		utils.Exec("curl", "-o", _rootfsFile, "-L", _rootfsURL)
		utils.Exec("move", _rootfsFile, rootZipFile)
	}

	utils.Exec("tar", "-xzf", rootZipFile, "-C", rootfs)

}

func CreateWriteLayer(writeLayerURL string) {
	utils.MustPathExist(writeLayerURL)
}

func CreateMountPoint(rootURL, writeLayerURL, mntURL, containerID string) {
	utils.MustPathExist(mntURL)

	dirs := fmt.Sprintf("dirs=%s:%srootfs", writeLayerURL, rootURL)

	utils.Exec("mount", "-t", "aufs", "-o", dirs, containerID, mntURL)

}

func CreateMountVolumePoint(mntURL string, volumeURLs []string) {
	volumes := getCorrectVolumes(volumeURLs)

	for i := 0; i < len(volumes)/2; i++ {
		mountVolume(volumes[i*2], path.Join(mntURL, volumes[i*2+1]))
	}
}

func getCorrectVolumes(volumeURLs []string) []string {
	result := make([]string, 0)

	for _, v := range volumeURLs {
		urls := strings.Split(v, ":")
		if len(urls) != 2 || urls[0] == "" || urls[1] == "" {
			continue
		}
		result = append(result, urls...)
	}

	return result
}

func mountVolume(source, dest string) (err error) {
	utils.MustPathExist(source)
	utils.MustPathExist(dest)
	_, err = utils.Exec("mount", "-t", "aufs", "-o", fmt.Sprintf("dirs=%s", source), "none", dest)
	return
}

//DeleteWorkSpace delete container workspace
func DeleteWorkSpace(writeLayerURL, mntURL string, volumeURLs []string) {
	volumes := getCorrectVolumes(volumeURLs)
	for i := 0; i < len(volumes)/2; i++ {
		UnmountVolume(path.Join(mntURL, volumes[2*i+1]))
	}

	DeleteMountPoint(mntURL)
	DeleteWriteLayer(writeLayerURL)
}

func DeleteMountPoint(mntURL string) {
	utils.Exec("umount", mntURL)
	utils.RemoveAll(path.Dir(mntURL))
}

func DeleteWriteLayer(writeLayerURL string) {
	utils.RemoveAll(path.Dir(writeLayerURL))
}

func UnmountVolume(volumeMntURL string) {
	utils.Exec("umount", volumeMntURL)
	//utils.RemoveAll(volumeMntURL)
}

func CleanUpWorkspace(containerID string) error {
	info, err := getContainerInfoByID(containerID)
	if err != nil {
		logrus.Errorf("get container %s info error %v", containerID, err)
	}
	writeLayerURL, mntURL := getContainerURL(containerID)

	DeleteWorkSpace(writeLayerURL, mntURL, info.Volumes)
	if info.Rm {
		DeleteContainerInfo(containerID)
	}
	return nil
}
