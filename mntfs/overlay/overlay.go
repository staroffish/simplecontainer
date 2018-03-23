package overlay

import (
	"config"
	"fmt"
	"mntfs"
	"os"
	"syscall"

	"github.com/sirupsen/logrus"
)

type overlay struct{}

var (
	upperPath string // up层路径
	lowerPath string // lower层路径
	mergePath string // 挂载路径
	workPath  string // work层路径
)

func init() {
	mntfs.SetMntInst("overlay", &overlay{})
}

func (o *overlay) InitMnt(name, imageName string) error {
	lowerPath = fmt.Sprintf("%s/%s/", config.ImagePath, imageName)
	mergePath = fmt.Sprintf("%s/%s/", config.MntPath, name)
	writePath := fmt.Sprintf("%s/%s", config.WirtelayPath, name)
	upperPath = fmt.Sprintf("%s/%s/", writePath, "upper")
	workPath = fmt.Sprintf("%s/%s/", writePath, "work")

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

func (o *overlay) Mount() error {

	// mount -t overlay overlay -olowerdir=./lower,upperdir=./upper,workdir=./work ./merged
	optStr := fmt.Sprintf("lowerdir=%s,upperdir=%s,workdir=%s", lowerPath, upperPath, workPath)
	err := syscall.Mount("overlay", mergePath, "overlay", 0, optStr)
	if err != nil {
		logrus.Errorf("Mount overlay file system error:%v", err)
		return err
	}
	return nil
}

func (o *overlay) Unmount() error {

	err := syscall.Unmount(mergePath, 0)
	if err != nil {
		logrus.Errorf("Unmount overlay file system error:%v", err)
		return err
	}
	return nil
}
