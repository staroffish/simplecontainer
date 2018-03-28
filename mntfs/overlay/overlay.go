package overlay

import (
	"fmt"
	"os"
	"syscall"

	"github.com/staroffish/simplecontainer/config"
	"github.com/staroffish/simplecontainer/mntfs"

	"github.com/sirupsen/logrus"
)

type overlay struct{}

func init() {
	mntfs.SetMntInst(mntfs.OVERLAY, &overlay{})
}

func (o *overlay) InitMnt(name, imageName string) error {
	lowerPath := fmt.Sprintf("%s/%s/", config.ImagePath, imageName)
	mergePath := fmt.Sprintf("%s/%s/", config.MntPath, name)
	writePath := fmt.Sprintf("%s/%s", config.WirtelayPath, name)
	upperPath := fmt.Sprintf("%s/%s/", writePath, "upper")
	workPath := fmt.Sprintf("%s/%s/", writePath, "work")

	// 确认镜像目录是否存在
	_, err := os.Stat(lowerPath)
	if err != nil {
		logrus.Errorf("Load image incorrect:%v", err)
		return err
	}

	// 创建挂载目录
	if err = os.MkdirAll(mergePath, 755); err != nil {
		logrus.Errorf("make mnt dir error:%v", err)
		return err
	}

	// 创建可写层目录
	if err = os.MkdirAll(writePath, 755); err != nil {
		logrus.Errorf("make WriteLayer dir error:%v", err)
		return err
	}

	// 查看upper目录是否存在
	_, err = os.Stat(upperPath)
	if err != nil {
		if !os.IsNotExist(err) {
			logrus.Errorf("stat %s error:%v", upperPath, err)
			return err
		}
		// 如果不存在创建upper目录
		if err := os.MkdirAll(upperPath, 755); err != nil {
			logrus.Errorf("Mkdir %s error:%v", upperPath, err)
			return err
		}
	}

	// 查看work目录是否存在
	_, err = os.Stat(workPath)
	if err != nil {
		if !os.IsNotExist(err) {
			logrus.Errorf("stat %s error:%v", workPath, err)
			return err
		}
		// 如果不存在创建work目录
		if err := os.MkdirAll(workPath, 755); err != nil {
			logrus.Errorf("Mkdir %s error:%v", workPath, err)
			return err
		}
	}

	return nil
}

func (o *overlay) Mount(name, imageName string) error {
	lowerPath := fmt.Sprintf("%s/%s/", config.ImagePath, imageName)
	mergePath := fmt.Sprintf("%s/%s/", config.MntPath, name)
	writePath := fmt.Sprintf("%s/%s", config.WirtelayPath, name)
	upperPath := fmt.Sprintf("%s/%s/", writePath, "upper")
	workPath := fmt.Sprintf("%s/%s/", writePath, "work")

	// mount -t overlay overlay -olowerdir=./lower,upperdir=./upper,workdir=./work ./merged
	optStr := fmt.Sprintf("lowerdir=%s,upperdir=%s,workdir=%s", lowerPath, upperPath, workPath)
	err := syscall.Mount("overlay", mergePath, "overlay", 0, optStr)
	if err != nil {
		logrus.Errorf("Mount overlay file system error:%v", err)
		return err
	}
	return nil
}

func (o *overlay) Unmount(name string) error {
	return mntfs.Unmount(name)
}

func (o *overlay) Remove(name string) error {
	return mntfs.Remove(name)
}
