package container

import (
	"fmt"
	"path"

	"github.com/CaiJinKen/mydocker/ufs"

	"github.com/CaiJinKen/mydocker/utils"
)

//NewWorkSpace create container workspace
func NewWorkSpace(rootURL, mntURL string) {
	CreateReadOnlyLayer(rootURL)
	CreateWriteLayer(rootURL)
	CreateMountPoint(rootURL, mntURL)
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

	utils.Exec("mount", "-t", ufs.GetSupportedUFS(), "-o", dirs, "container", mntURL)

}

//DeleteWorkSpace delete container workspace
func DeleteWorkSpace(rootURL, mntURL string) {
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
