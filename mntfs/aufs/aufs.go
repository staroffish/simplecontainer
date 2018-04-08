package aufs

import (
	"fmt"
	"os"
	"syscall"

	"github.com/staroffish/simplecontainer/config"
	"github.com/staroffish/simplecontainer/mntfs"

	"github.com/sirupsen/logrus"
)

type aufs struct{}

func init() {
	mntfs.SetMntInst(mntfs.AUFS, &aufs{})
}

func (a *aufs) InitMnt(name, imageName string) error {
	rdOnlyPath := fmt.Sprintf("%s/%s/", config.ImagePath, imageName)
	mntPath := fmt.Sprintf("%s/%s/", config.MntPath, name)
	wirtelayPath := fmt.Sprintf("%s/%s", config.WirtelayPath, name)

	// 确认镜像目录是否存在
	_, err := os.Stat(rdOnlyPath)
	if err != nil {
		logrus.Errorf("Load image incorrect:%v", err)
		return err
	}

	// 创建挂载目录
	if err = os.MkdirAll(mntPath, 755); err != nil {
		logrus.Errorf("make mnt dir error:%v", err)
		return err
	}

	// 创建可写层目录
	if err = os.MkdirAll(wirtelayPath, 755); err != nil {
		logrus.Errorf("make WriteLayer dir error:%v", err)
		return err
	}

	return nil
}

func (a *aufs) Mount(name, imageName string) error {
	rdOnlyPath := fmt.Sprintf("%s/%s/", config.ImagePath, imageName)
	mntPath := fmt.Sprintf("%s/%s/", config.MntPath, name)
	wirtelayPath := fmt.Sprintf("%s/%s", config.WirtelayPath, name)

	// mount -t aufs none -o br=wirtelayPath:rdOnlyPath mntPath
	optStr := fmt.Sprintf("br=%s:%s", wirtelayPath, rdOnlyPath)
	err := syscall.Mount("none", mntPath, "aufs", 0, optStr)
	if err != nil {
		logrus.Errorf("Mount aufs file system error:%v", err)
		return err
	}
	return nil
}

func (a *aufs) Unmount(name string) error {
	return mntfs.Unmount(name)
}

func (a *aufs) Remove(name string) error {
	return mntfs.Remove(name)
}
