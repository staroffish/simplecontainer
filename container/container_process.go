package container

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/sirupsen/logrus"
)

// 容器进程
func ContianerProcess(name string) {
	var cInfo *ContainerInfo
	sigCh := make(chan os.Signal)

	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	logrus.Info("start ContianerProcess")
	select {
	case sig := <-sigCh:
		if sig == syscall.SIGINT {
			logrus.Info("Get SIGINT")
			// 初始化容器进程
			cInfo = ReadContainerInfo(name)
			if cInfo == nil {
				logrus.Errorf("ReadContainerInfo error")
				os.Exit(-1)
			}

			if cInfo.NetType == "dhcp" {
				// 获取DHCP
				cmd := exec.Command("/sbin/dhclient")
				output, err := cmd.CombinedOutput()
				if err != nil {
					logrus.Errorf("Execute dhcp cmd error %s:%v", output, err)
				}
			}

			// 切换rootfs 挂载proc文件系统
			err := setUpMount()
			if err != nil {
				logrus.Errorf("setUpMount error")
				os.Exit(-1)
			}
			signal.Reset(syscall.SIGINT)
		} else {
			logrus.Infof("get signal %s", sig.String())

			os.Exit(0)
		}
	}

	select {
	case sig := <-sigCh:
		logrus.Infof("get signal %s", sig.String())

		os.Exit(0)
	}
}

// 切换rootfs 挂载proc文件系统
func setUpMount() error {
	pwd, err := os.Getwd()
	if err != nil {
		logrus.Errorf("Get current location error %v", err)
		return err
	}

	logrus.Infof("Current location is %s", pwd)

	// 将根目录的挂载flag设成PRIVATE
	err = syscall.Mount("", "/", "", uintptr(syscall.MS_PRIVATE|syscall.MS_REC), "")
	if err != nil {
		logrus.Errorf(err.Error())
		return err
	}

	// 切换rootfs
	err = pivotRoot(pwd)
	if err != nil {
		logrus.Errorf("pivotRoot error %v", err)
		return err
	}

	// 挂载proc文件系统
	defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
	err = syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")
	if err != nil {
		logrus.Errorf(err.Error())
		return err
	}

	// 挂载tmpfs
	err = syscall.Mount("tmpfs", "/dev", "tmpfs", syscall.MS_NOSUID|syscall.MS_STRICTATIME, "mode=755")
	if err != nil {
		logrus.Errorf(err.Error())
		return err
	}
	return nil
}

// 切换rootfs
func pivotRoot(root string) error {
	// 由于pivot_root不能切换到自己分区里的目录,所以先将root bind一下
	if err := syscall.Mount(root, root, "bind", syscall.MS_BIND|syscall.MS_REC, ""); err != nil {
		return fmt.Errorf("Mount rootfs to itself error:%v", err)
	}

	pivotDir := filepath.Join(root, ".pivot_root")
	if err := os.Mkdir(pivotDir, 0777); err != nil && !os.IsExist(err) {
		return err
	}
	logrus.Infof("pivotDir=%s", pivotDir)

	// 切换rootfs
	if err := syscall.PivotRoot(root, pivotDir); err != nil {
		return fmt.Errorf("PivotRoot error:%v", err)
	}
	logrus.Infof("PivotRoot ok root=%s", root)

	// 将当前目录切换到 /
	if err := syscall.Chdir("/"); err != nil {
		return fmt.Errorf("Chdir /:%v", err)
	}

	pivotDir = "/.pivot_root"

	// 卸载原有rootfs
	if err := syscall.Unmount(pivotDir, syscall.MNT_DETACH); err != nil {
		return fmt.Errorf("umount /.pivot_root :%v", err)
	}

	return os.Remove(pivotDir)
}
