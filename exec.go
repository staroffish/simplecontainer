package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/staroffish/simplecontainer/cgroups"
	container "github.com/staroffish/simplecontainer/container"

	_ "github.com/staroffish/simplecontainer/nsenter"

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

	if err = cmd.Start(); err != nil {
		return fmt.Errorf("Exec command error:%s:%v", command, err)
	}

	pid := fmt.Sprintf("%d", cmd.Process.Pid)

	// Cgroup设定
	if len(cInfo.MemLimit) != 0 {
		if err = cgroups.SetMemroyLimit(cInfo.Name, cInfo.MemLimit, pid); err != nil {
			logrus.Errorf("Set Memroy Limit error")
			return err
		}
	}
	if len(cInfo.CPU) != 0 {
		if err = cgroups.SetCPULimit(cInfo.Name, cInfo.CPU, pid); err != nil {
			logrus.Errorf("Set CPU Limit error")
			return err
		}
	}

	cmd.Wait()

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
