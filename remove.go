package main

import (
	"fmt"
	container "staroffish/simplecontainer/container"
	"staroffish/simplecontainer/mntfs"
)

func remove(name string) error {
	//取得容器信息
	cInfo := container.ReadContainerInfo(name)
	if cInfo == nil {
		return fmt.Errorf("Container %s not exists", name)
	}

	// 检查容器状态
	if cInfo.Status != container.STOP {
		return fmt.Errorf("Container %s was started", name)
	}

	// 删除容器信息
	if err := container.RemoveContainerInfo(name); err != nil {
		return fmt.Errorf("Remove container info error:%v", err)
	}

	mntFs := mntfs.GetMountInst(mntfs.OVERLAY)

	// 删除容器
	return mntFs.Remove(name)
}
