package container

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"syscall"
	"time"

	"github.com/staroffish/simplecontainer/config"

	"github.com/sirupsen/logrus"
)

type ContainerInfo struct {
	Pid           string `json:"pid"`
	Name          string `json:"name"`
	CreateTime    string `json:"createTime"`
	Status        string `json:"status"`
	ImageName     string `json:"imagename"`
	MemLimit      string `json:"memlimit"`
	CPU           string `json:"cpu"`
	NetType       string `json:"nettype"`
	Subnet        string `json:"subnet"`
	ParentNetwork string `json:"parentnetwork"`
	NetDeviceName string `json:"devicename"`
	Gateway       string `json:"gateway"`
}

var (
	STARTING   = "starting"
	RUNNING    = "running"
	STOP       = "stopped"
	EXIT       = "exited"
	ConfigName = "config.json"
)

// 启动容器父进程
func NewParentProcess(name *string, imageName string) (*exec.Cmd, error) {
	if *name == "" {
		*name = RandStringBytes(16)
	}
	cmd := exec.Command("/proc/self/exe", "init", *name)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID |
			syscall.CLONE_NEWNS | syscall.CLONE_NEWNET | syscall.CLONE_NEWIPC,
	}

	cmd.Dir = fmt.Sprintf("%s/%s/", config.MntPath, *name)

	return cmd, nil
}

func RandStringBytes(n int) string {
	letterBytes := "1234567890"
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}

	return string(b)
}

// 保存容器信息
func StoreContainerInfo(cInfo *ContainerInfo) error {
	// 序列化容器信息
	jstr, err := json.Marshal(cInfo)
	if err != nil {
		logrus.Errorf("marshal container info error:%v", err)
		return err
	}

	// 创建容器信息目录
	confDir := fmt.Sprintf("%s/%s", config.CInfopath, cInfo.Name)
	if err := os.MkdirAll(confDir, 0755); err != nil {
		logrus.Errorf("make container info dir error:%s:%v", confDir, err)
		return err
	}

	// 写入容器信息
	confPath := fmt.Sprintf("%s/%s", confDir, ConfigName)
	logrus.Debug(confPath)
	logrus.Debug(string(jstr))
	if err := ioutil.WriteFile(confPath, jstr, 0644); err != nil {
		logrus.Errorf("Write container info to file error:%s:%v", confPath, err)
		return err
	}

	return nil
}

// 取得容器信息
func ReadContainerInfo(name string) *ContainerInfo {
	cInfo := &ContainerInfo{Name: name}
	confPath := fmt.Sprintf("%s/%s", fmt.Sprintf("%s/%s", config.CInfopath, cInfo.Name), ConfigName)

	// 读取容器信息
	data, err := ioutil.ReadFile(confPath)
	if err != nil {
		logrus.Errorf("Read container info from file error:%s:%v", confPath, err)
		return nil
	}

	// 序列化容器信息
	if err := json.Unmarshal(data, cInfo); err != nil {
		logrus.Errorf("marshal container info error:%v", err)
		return nil
	}

	return cInfo
}

// 删除容器信息
func RemoveContainerInfo(name string) error {
	confPath := fmt.Sprintf("%s/%s", config.CInfopath, name)
	return os.RemoveAll(confPath)
}
