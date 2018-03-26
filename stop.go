package main

import (
	"fmt"
	"strconv"
	"syscall"
	"time"

	"github.com/staroffish/simplecontainer/cgroups"
	container "github.com/staroffish/simplecontainer/container"
	"github.com/staroffish/simplecontainer/mntfs"
)

func stop(name string) error {
	//取得容器信息
	cInfo := container.ReadContainerInfo(name)
	if cInfo == nil {
		return fmt.Errorf("Container %s not exists", name)
	}

	// 检查容器状态
	if cInfo.Status != container.RUNNING && cInfo.Status != container.STARTING {
		return fmt.Errorf("Container %s already stopped", name)
	}

	if cInfo.Status != container.STARTING {
		pid, err := strconv.Atoi(cInfo.Pid)
		if err != nil {
			return fmt.Errorf("Pid error in container configuration")
		}
		// 停止容器
		syscall.Kill(pid, syscall.SIGTERM)
	}

	time.Sleep(1 * time.Second)
	mntFs := mntfs.GetMountInst(mntfs.OVERLAY)
	mntFs.Unmount(cInfo.Name)

	if len(cInfo.CPU) != 0 {
		cgroups.UnsetCPULimit(cInfo.Name)
	}
	if len(cInfo.MemLimit) != 0 {
		cgroups.UnsetMemroyLimit(cInfo.Name)
	}

	cInfo.Status = container.STOP
	cInfo.Pid = ""
	container.StoreContainerInfo(cInfo)

	return nil
}
