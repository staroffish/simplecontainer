package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	container "staroffish/simplecontainer/container"
	"strings"

	_ "staroffish/simplecontainer/nsenter"

	"github.com/sirupsen/logrus"
)

const ENV_EXEC_PID = "mydocker_pid"
const ENV_EXEC_CMD = "mydocker_cmd"

func execCmd(containerName, command string) error {
	logrus.Infof("exec command. containername=%s command=%s", containerName, command)
	cInfo := container.ReadContainerInfo(containerName)
	if cInfo == nil {
		return fmt.Errorf("Container info read error.")
	}

	envs, err := getEnvsByPid(cInfo.Pid)
	if err != nil {
		return err
	}
	cmd := exec.Command("/proc/self/exe", "exec")

	os.Setenv(ENV_EXEC_PID, cInfo.Pid)
	os.Setenv(ENV_EXEC_CMD, command)
	cmd.Env = append(os.Environ(), envs...)

	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr

	if err = cmd.Run(); err != nil {
		return fmt.Errorf("Exec command error:%s:%v", command, err)
	}

	return nil
}

func getEnvsByPid(pid string) ([]string, error) {
	path := fmt.Sprintf("/proc/%s/environ", pid)
	contentBytes, err := ioutil.ReadFile(path)
	if err != nil {
		logrus.Errorf("Read file %s error %v", path, err)
		return nil, fmt.Errorf("Read file %s error %v", path, err)
	}
	//env split by \u0000
	envs := strings.Split(string(contentBytes), "\u0000")
	return envs, nil
}
