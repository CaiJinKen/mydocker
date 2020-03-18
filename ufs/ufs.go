package ufs

import (
	"bytes"
	"sort"

	"github.com/CaiJinKen/mydocker/utils"
	"github.com/sirupsen/logrus"
)

var ufs []string

var knownUfs = []string{"aufs", "btrfs", "overlay", "overlay2", "vfs", "zfs"}

func GetSupportedUFS() string {
	if len(ufs) > 0 {
		return ufs[0]
	}

	fss := getSupportedFS()

	for _, v := range knownUfs {
		for _, f := range fss {
			if v == f {
				ufs = append(ufs, v)
			}
		}
	}

	sort.Strings(ufs)

	if len(ufs) > 0 {
		return ufs[0]
	}

	return ""

}

func getSupportedFS() (fss []string) {
	output, err := utils.Exec("cat", "/proc/filesystems")
	if err != nil {
		logrus.Errorf("get supported file system error %v", err)
		return
	}

	lines := bytes.Split(output, []byte{'\n'})
	if len(lines) == 0 {
		return
	}

	for _, line := range lines {
		if len(line) == 0 {
			continue
		}

		if fields := bytes.Split(line, []byte{'\t'}); len(fields) > 1 {
			fss = append(fss, string(fields[len(fields)-1]))
		}

	}

	return
}
