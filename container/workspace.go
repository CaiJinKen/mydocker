package container

import (
	"fmt"
	"path"
	"strings"

	"github.com/CaiJinKen/mydocker/utils"
)

//NewWorkSpace create container workspace
func NewWorkSpace(rootURL, mntURL string, volumeURLs []string) {
	CreateReadOnlyLayer(rootURL)
	CreateWriteLayer(rootURL)
	CreateMountPoint(rootURL, mntURL)
	CreateMountVolumePoint(mntURL, volumeURLs)
}

func CreateReadOnlyLayer(rootURL string) {
	//use /roo/rootfs as readonly layer

}

func CreateWriteLayer(rootURL string) {
	writeURL := path.Join(rootURL, "writeLayer")
	utils.MustPathExist(writeURL)
}

func CreateMountPoint(rootURL, mntURL string) {
	utils.MustPathExist(mntURL)

	dirs := fmt.Sprintf("dirs=%swriteLayer:%srootfs", rootURL, rootURL)

	utils.Exec("mount", "-t", "aufs", "-o", dirs, "container", mntURL)

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
func DeleteWorkSpace(rootURL, mntURL string, volumeURLs []string) {
	volumes := getCorrectVolumes(volumeURLs)
	for i := 0; i < len(volumes)/2; i++ {
		UnmountVolume(path.Join(mntURL, volumes[2*i+1]))
	}

	DeleteMountPoint(mntURL)
	DeleteWriteLayer(rootURL)
}

func DeleteMountPoint(mntURL string) {
	utils.Exec("umount", mntURL)
	utils.RemoveAll(mntURL)
}

func DeleteWriteLayer(rootURL string) {
	writeURL := path.Join(rootURL, "writeLayer")
	utils.RemoveAll(writeURL)
}

func UnmountVolume(volumeMntURL string) {
	utils.Exec("umount", volumeMntURL)
	//utils.RemoveAll(volumeMntURL)
}
