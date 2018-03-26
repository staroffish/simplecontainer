package cgroups

import (
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"strconv"
	"syscall"

	"github.com/sirupsen/logrus"
)

const (
	cgroupMemPath       = "/sys/fs/cgroup/memory"
	cgroupCpuPath       = "/sys/fs/cgroup/cpu"
	memLimitFileName    = "memory.limit_in_bytes"
	cfsPeriodUsFileName = "cpu.cfs_period_us"
	cfsQuotaUsName      = "cpu.cfs_quota_us"
	cgroupProcName      = "cgroup.procs"
	tasksName           = "tasks"
)

func SetMemroyLimit(name, limit, pid string) error {
	memCgPath := fmt.Sprintf("%s/%s", cgroupMemPath, name)
	if err := os.Mkdir(memCgPath, 0755); err != nil && !os.IsExist(err) {
		logrus.Errorf("Mkdir %s error:%v", memCgPath, err)
		return fmt.Errorf("Mkdir %s error:%v", memCgPath, err)
	}

	memLimit, err := strconv.Atoi(limit)
	if err != nil {
		logrus.Errorf("convert MemLimit %s to int error:%v", limit, err)
		return fmt.Errorf("convert MemLimit %s to int error:%v", limit, err)
	}
	memLimit = memLimit * 1024 * 1024

	memLimitFile := fmt.Sprintf("%s/%s", memCgPath, memLimitFileName)

	if err := ioutil.WriteFile(memLimitFile, []byte(fmt.Sprintf("%d", memLimit)), 0644); err != nil {
		logrus.Errorf("set mem limit error %s:%v", memLimitFile, err)
		return fmt.Errorf("set mem limit error %s:%v", memLimitFile, err)
	}

	procFile := fmt.Sprintf("%s/%s", memCgPath, cgroupProcName)
	if err := ioutil.WriteFile(procFile, []byte(pid), 0644); err != nil {
		logrus.Errorf("set pid to procfile error %s:%v", procFile, err)
		return fmt.Errorf("set pid to procfile error %s:%v", procFile, err)
	}
	return nil
}

func SetCPULimit(name, limit, pid string) error {
	cpuCgPath := fmt.Sprintf("%s/%s", cgroupCpuPath, name)
	if err := os.Mkdir(cpuCgPath, 0755); err != nil && !os.IsExist(err) {
		logrus.Errorf("Mkdir %s error:%v", cpuCgPath, err)
		return fmt.Errorf("Mkdir %s error:%v", cpuCgPath, err)
	}

	sysinfo := &syscall.Sysinfo_t{}
	if err := syscall.Sysinfo(sysinfo); err != nil {
		logrus.Errorf("Get sysinfo error:%v", err)
		return fmt.Errorf("Get sysinfo error:%v", err)
	}

	cpuNum, err := strconv.Atoi(limit)
	if err != nil {
		logrus.Errorf("Convert cpu num to int error:%s:%v", limit, err)
		return fmt.Errorf("Convert cpu num to int error:%s:%v", limit, err)
	}

	maxCPU := runtime.NumCPU()

	if cpuNum > maxCPU {
		logrus.Errorf("cpu must be less then %d", maxCPU)
		return fmt.Errorf("cpu must be less then %d", maxCPU)
	}

	periodPath := fmt.Sprintf("%s/%s", cpuCgPath, cfsPeriodUsFileName)
	data, err := ioutil.ReadFile(periodPath)
	if err != nil {
		logrus.Errorf("Read %s error:%v", cfsPeriodUsFileName, err)
		return fmt.Errorf("Read %s error:%v", cfsPeriodUsFileName, err)
	}

	if data[len(data)-1] == '\n' {
		data = data[:len(data)-1]
	}
	periodUs, err := strconv.Atoi(string(data))
	if err != nil {
		logrus.Errorf("Convert period us to int error:%s:%v", data, err)
		return fmt.Errorf("Convert period us to int error:%s:%v", data, err)
	}

	quotaUs := fmt.Sprintf("%d", periodUs*cpuNum)

	quotaPath := fmt.Sprintf("%s/%s", cpuCgPath, cfsQuotaUsName)
	err = ioutil.WriteFile(quotaPath, []byte(quotaUs), 0644)
	if err != nil {
		logrus.Errorf("write %s error:%v", quotaPath, err)
		return fmt.Errorf("write %s error:%v", quotaPath, err)
	}

	procFile := fmt.Sprintf("%s/%s", cpuCgPath, cgroupProcName)
	if err := ioutil.WriteFile(procFile, []byte(pid), 0644); err != nil {
		logrus.Errorf("set pid to procfile error %s:%v", procFile, err)
		return fmt.Errorf("set pid to procfile error %s:%v", procFile, err)
	}
	return nil
}

func UnsetMemroyLimit(name string) error {
	memCgPath := fmt.Sprintf("%s/%s", cgroupMemPath, name)

	if err := os.Remove(memCgPath); err != nil {
		logrus.Errorf("remove %s error:%v", memCgPath, err)
		return fmt.Errorf("remove %s error:%v", memCgPath, err)
	}

	return nil
}

func UnsetCPULimit(name string) error {
	cpuCgPath := fmt.Sprintf("%s/%s", cgroupCpuPath, name)

	if err := os.Remove(cpuCgPath); err != nil {
		logrus.Errorf("remove %s error:%v", cpuCgPath, err)
		return fmt.Errorf("remove %s error:%v", cpuCgPath, err)
	}
	return nil
}
