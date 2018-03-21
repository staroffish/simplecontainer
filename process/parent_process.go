package process

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"

	"github.com/sirupsen/logrus"
)

type ContainerInfo struct {
	Pid        string `json:"pid"`
	Name       string `json:"name"`
	CreateTime string `json:"createTime"`
	Status     string `json:"status"`
	ImageName  string `json:"imagename"`
	MemLimit   string `json:memlimit`
	Cpu        string `json:cpu`
}

var (
	RUNNING       = "running"
	STOP          = "stopped"
	EXIT          = "exited"
	InfoLocation  = "./config/container_info/%s/"
	MntUrl        = "./mnt/%s/"
	RootUrl       = "./"
	WriteLayerUrl = "./writeLayer/%s/"
	ConfigName    = "config.json"
)

// 启动容器父进程
func NewParentProcess(name, imageName string) (*exec.Cmd, error) {
	cmd := exec.Command("/proc/self/exe", "init", name)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID |
			syscall.CLONE_NEWNS | syscall.CLONE_NEWNET | syscall.CLONE_NEWIPC,
	}

	if err := initWorkspace(name, imageName); err != nil {
		return nil, err
	}

	return cmd, nil
}

// 创建容器需要用到的各个文件夹
func initWorkspace(name, imageName string) error {
	// 镜像目录
	readOnlyPath := fmt.Sprintf("%s/image/%s", RootUrl, imageName)
	_, err := os.Stat(readOnlyPath)
	if err != nil {
		logrus.Errorf("Image path incorrect:%v", err)
		return err
	}

	// 创建挂载目录
	mntPath := fmt.Sprintf(MntUrl, name)
	if err = os.MkdirAll(mntPath, 755); err != nil {
		logrus.Errorf("make mnt dir error:%v", err)
		return err
	}

	// 创建可写层目录
	writePath := fmt.Sprintf(WriteLayerUrl, name)
	if err = os.MkdirAll(writePath, 755); err != nil {
		logrus.Errorf("make WriteLayer dir error:%v", err)
		return err
	}

	return nil
}
