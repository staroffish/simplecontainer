package main

import (
	"fmt"
	container "staroffish/simplecontainer/container"
	"staroffish/simplecontainer/mntfs"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
)

func run(cInfo *container.ContainerInfo) (err error) {
	if len(cInfo.Name) > 0 && container.ReadContainerInfo(cInfo.Name) != nil {
		return fmt.Errorf("Container %s exists", cInfo.Name)
	}
	return startContainer(cInfo)
}

func startContainer(cInfo *container.ContainerInfo) (err error) {
	err = nil
	// 创建容器进程的命令
	cmd, err := container.NewParentProcess(&cInfo.Name, cInfo.ImageName)
	if err != nil {
		return
	}

	// 挂载容器目录
	mntFs := mntfs.GetMountInst(mntfs.OVERLAY)

	if err := mntFs.InitMnt(cInfo.Name, cInfo.ImageName); err != nil {
		return err
	}

	if err := mntFs.Mount(cInfo.Name, cInfo.ImageName); err != nil {
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
			mntFs.Unmount(cInfo.Name)
			cInfo.Status = container.STOP
			container.StoreContainerInfo(cInfo)
		}
	}()

	cInfo.Pid = fmt.Sprintf("%d", cmd.Process.Pid)
	// Cgroup设定

	// 网络设定

	// 保存容器设定
	cInfo.CreateTime = time.Now().Format("2006-01-02 15:04:05")
	cInfo.Status = container.STARTING
	err = container.StoreContainerInfo(cInfo)
	if err != nil {
		return
	}

	time.Sleep(1 * time.Second)
	// 给容器发送Cgroup和网络设定完成的信号
	err = cmd.Process.Signal(syscall.SIGINT)
	if err != nil {
		logrus.Errorf("send sighup to child error:%v", err)
		return
	}
	cInfo.Status = container.RUNNING
	container.StoreContainerInfo(cInfo)
	return
}
