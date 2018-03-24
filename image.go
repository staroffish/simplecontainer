package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"staroffish/simplecontainer/config"
	container "staroffish/simplecontainer/container"

	"github.com/sirupsen/logrus"
)

func commit(name, imageName string) error {
	cInfo := container.ReadContainerInfo(name)
	if cInfo == nil {
		logrus.Errorf("container %s does not exists", name)
		return fmt.Errorf("container %s does not exists", name)
	}

	if cInfo.Status != container.RUNNING {
		logrus.Errorf("container %s must be started", name)
		return fmt.Errorf("container %s must be started", name)
	}

	//取得容器挂载目录
	mntPath := fmt.Sprintf("%s/%s", config.MntPath, name)

	newImgPath := fmt.Sprintf("%s/%s", config.ImagePath, imageName)

	cmd := exec.Command("cp", "-rp", mntPath, newImgPath)

	output, err := cmd.CombinedOutput()
	if err != nil {
		logrus.Errorf("cp image error:%v", err)
		return fmt.Errorf("cp image error:%v", err)
	}
	if len(output) != 0 {
		logrus.Errorf("cp image error:%s", output)
		return fmt.Errorf("cp image error:%s", output)
	}

	return nil
}

func imageList() error {
	//取得镜像信息
	dirs, err := ioutil.ReadDir(config.ImagePath)
	if err != nil {
		logrus.Errorf("read dir %s error:%v", config.ImagePath, err)
		return fmt.Errorf("read dir %s error:%v", config.ImagePath, err)
	}

	fmt.Println("IMAGE")
	for _, dir := range dirs {
		if dir.IsDir() {
			fmt.Println(dir.Name())
		}
	}

	return nil
}

func importImage(gzFile string) error {

	if _, err := exec.Command("tar", "-xzvf", gzFile, "-C", config.ImagePath).CombinedOutput(); err != nil {
		logrus.Errorf("unTar file %s error %v", gzFile, err)
		return fmt.Errorf("unTar file %s error %v", gzFile, err)
	}

	return nil
}

func exportImage(imageName, gzPath string) error {
	//取得镜像信息
	if err := checkImageExists(imageName); err != nil {
		return err
	}

	gzFile := fmt.Sprintf("%s/%s.tgz", gzPath, imageName)
	output, err := exec.Command("tar", "-C", config.ImagePath, "-czvf", gzFile, imageName).CombinedOutput()
	if err != nil {
		logrus.Errorf("Compress file %s error %s", gzFile, output)
		return fmt.Errorf("Compress file %s error %s", gzFile, output)
	}

	return nil
}

func removeImage(name string) error {

	if err := checkImageExists(name); err != nil {
		return err
	}

	imagePath := fmt.Sprintf("%s/%s", config.ImagePath, name)
	if err := os.RemoveAll(imagePath); err != nil {
		logrus.Errorf("remove dir %s error:%v", imagePath, err)
		return fmt.Errorf("remove dir %s error:%v", imagePath, err)
	}
	return nil
}

func checkImageExists(name string) error {
	//取得镜像信息
	dirs, err := ioutil.ReadDir(config.ImagePath)
	if err != nil {
		logrus.Errorf("read dir %s error:%v", config.ImagePath, err)
		return fmt.Errorf("read dir %s error:%v", config.ImagePath, err)
	}

	exists := false
	for _, dir := range dirs {
		if dir.IsDir() && dir.Name() == name {
			exists = true
			break
		}
	}

	if !exists {
		logrus.Errorf("Image %s does not exists", name)
		return fmt.Errorf("Image %s does not exists", name)
	}

	return nil
}
