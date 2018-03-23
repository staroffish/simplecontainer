package main

import (
	"fmt"
	"mntfs"
	"process"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
)

func Run(name, imageName string) (err error) {
	err = nil
	// 创建容器进程的命令
	cmd, err := process.NewParentProcess(&name, imageName)
	if err != nil {
		return
	}

	cInfo := &process.ContainerInfo{}
	cInfo.Name = name
	cInfo.ImageName = imageName

	// 挂载容器目录
	mntFs := mntfs.GetMountInst("overlay")

	if err := mntFs.InitMnt(name, imageName); err != nil {
		return err
	}

	if err := mntFs.Mount(); err != nil {
		return err
	}

	// 启动容器进程
	if err = cmd.Start(); err != nil {
		logrus.Errorf("start command error:%v", err)
		return
	}

	defer func() {
		if err != nil {
			cmd.Process.Kill()
			mntFs.Unmount()
		}
	}()

	cInfo.Pid = fmt.Sprintf("%d", cmd.Process.Pid)
	// Cgroup设定

	// 网络设定

	// 保存容器设定
	cInfo.CreateTime = time.Now().Format("2006-01-02 15:04:05")
	err = process.StoreContainerInfo(cInfo)
	if err != nil {
		return
	}

	time.Sleep(5 * time.Second)
	// 给容器发送Cgroup和网络设定完成的信号
	err = cmd.Process.Signal(syscall.SIGINT)
	if err != nil {
		logrus.Errorf("send sighup to child error:%v", err)
		return
	}

	cmd.Wait()
	mntFs.Unmount()
	return
}
