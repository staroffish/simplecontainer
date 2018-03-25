package main

import (
	"fmt"

	container "github.com/staroffish/simplecontainer/container"
)

func start(name string) error {
	//取得容器信息
	cInfo := container.ReadContainerInfo(name)
	if cInfo == nil {
		return fmt.Errorf("Container %s not exists", name)
	}

	// 检查容器状态
	if cInfo.Status != container.STOP {
		return fmt.Errorf("Container %s already started", name)
	}

	// 启动容器
	return startContainer(cInfo)
}
