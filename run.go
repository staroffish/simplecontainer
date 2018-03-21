package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"process"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
)

func Run(name, imageName string) (err error) {
	err = nil
	cInfo := &process.ContainerInfo{}
	cInfo.Name = name
	cInfo.ImageName = imageName

	// 创建容器进程的命令
	cmd, err := process.NewParentProcess(name, imageName)
	if err != nil {
		return
	}

	// 启动容器进程
	if err = cmd.Start(); err != nil {
		logrus.Errorf("start command error:%v", err)
		return
	}

	defer func() {
		if err != nil {
			cmd.Process.Kill()
			process.Cleanup(name)
		}
	}()

	cInfo.Pid = fmt.Sprintf("%d", cmd.Process.Pid)
	// Cgroup设定

	// 网络设定

	// 保存容器设定
	cInfo.CreateTime = time.Now().Format("2006-01-02 15:04:05")
	err = StoreContainerInfo(cInfo)
	if err != nil {
		return
	}

	// 给容器发送Cgroup和网络设定完成的信号
	err = cmd.Process.Signal(syscall.SIGINT)
	if err != nil {
		logrus.Errorf("send sighup to child error:%v", err)
		return
	}

	return
}

// 保存容器信息
func StoreContainerInfo(cInfo *process.ContainerInfo) error {
	// 序列化容器信息
	jstr, err := json.Marshal(cInfo)
	if err != nil {
		logrus.Errorf("marshal container info error:%v", err)
		return err
	}

	// 创建容器信息目录
	confDir := fmt.Sprintf(process.InfoLocation, cInfo.Name)
	if err := os.MkdirAll(confDir, 0755); err != nil {
		logrus.Errorf("make container info dir error:%s:%v", confDir, err)
		return err
	}

	// 写入容器信息
	confPath := fmt.Sprintf("%s/%s", confDir, process.ConfigName)
	if err := ioutil.WriteFile(confPath, jstr, 0644); err != nil {
		logrus.Errorf("OpenFile container info file error:%s:%v", confPath, err)
		return err
	}

	return nil
}
