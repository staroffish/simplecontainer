package main

import (
	"fmt"
	container "staroffish/simplecontainer/container"
	"staroffish/simplecontainer/mntfs"
	"strconv"
	"syscall"
)

func stop(name string) error {
	//取得容器信息
	cInfo := container.ReadContainerInfo(name)
	if cInfo == nil {
		return fmt.Errorf("Container %s not exists", name)
	}

	// 检查容器状态
	if cInfo.Status != container.RUNNING {
		return fmt.Errorf("Container %s already stopped", name)
	}

	pid, err := strconv.Atoi(cInfo.Pid)
	if err != nil {
		return fmt.Errorf("Pid error in container configuration")
	}

	// 停止容器
	syscall.Kill(pid, syscall.SIGTERM)
	mntFs := mntfs.GetMountInst(mntfs.OVERLAY)
	mntFs.Unmount(cInfo.Name)
	cInfo.Status = container.STOP
	cInfo.Pid = ""
	container.StoreContainerInfo(cInfo)

	return nil
}
