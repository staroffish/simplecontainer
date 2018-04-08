package mntfs

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"syscall"

	"github.com/sirupsen/logrus"
	"github.com/staroffish/simplecontainer/config"
)

type MountFS interface {
	InitMnt(name, imageName string) error
	Mount(name, imageName string) error
	Unmount(name string) error
	Remove(name string) error
}

var mntInst = make(map[string]MountFS)

const (
	fileSystemPath    = "/proc/filesystems"
	fileSystemModPath = "/lib/modules/%s/kernel/fs"
	OVERLAY           = "overlay"
	AUFS              = "aufs"
)

var SupportedFileSystem = []string{OVERLAY, AUFS}

func SetMntInst(fsName string, mntFs MountFS) {
	mntInst[fsName] = mntFs
}

func GetMountInst(fsName string) MountFS {
	return mntInst[fsName]
}

func Unmount(name string) error {
	mergePath := fmt.Sprintf("%s%s", config.MntPath, name)

	err := syscall.Unmount(mergePath, 0)
	if err != nil {
		logrus.Errorf("Unmount file system error:%s:%v", mergePath, err)
		return err
	}
	return nil
}

func Remove(name string) error {
	mergePath := fmt.Sprintf("%s/%s", config.MntPath, name)
	writePath := fmt.Sprintf("%s/%s", config.WirtelayPath, name)

	if err := os.RemoveAll(mergePath); err != nil {
		logrus.Errorf("Remove %s error:%v", mergePath, err)
		return fmt.Errorf("Remove %s error:%v", mergePath, err)
	}

	if err := os.RemoveAll(writePath); err != nil {
		logrus.Errorf("Remove %s error:%v", writePath, err)
		return fmt.Errorf("Remove %s error:%v", writePath, err)
	}

	return nil
}

// 返回系统支持的文件系统
func GetSupportedFs() string {

	data, err := ioutil.ReadFile(fileSystemPath)
	if err != nil {
		logrus.Errorf("Read filesystem error %s:%v", fileSystemPath, err)
		return ""
	}

	for _, fs := range SupportedFileSystem {
		if bytes.Contains(data, []byte("\t"+fs)) {
			return fs
		}
	}

	uts := &syscall.Utsname{}
	err = syscall.Uname(uts)
	if err != nil {
		logrus.Errorf("Get linux kernel version error %s:%v", fileSystemPath, err)
		return ""
	}

	buf := bytes.Buffer{}
	for _, char := range uts.Release {
		if char == 0 {
			break
		}
		buf.WriteByte(uint8(char))
	}

	modulePath := fmt.Sprintf(fileSystemModPath, buf.String())

	dirs, err := ioutil.ReadDir(modulePath)
	if err != nil {
		logrus.Errorf("Read filesystem error %s:%v", fileSystemPath, err)
		return ""
	}

	for _, fs := range SupportedFileSystem {
		for _, dir := range dirs {
			if dir.IsDir() {
				if strings.Contains(dir.Name(), fs) {
					return fs
				}
			}
		}
	}

	return ""
}
