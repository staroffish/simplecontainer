package macvlan

import (
	"fmt"
	"net"
	"os"
	"runtime"
	"strconv"

	"github.com/staroffish/simplecontainer/container"
	"github.com/staroffish/simplecontainer/network"

	"github.com/vishvananda/netns"

	"github.com/sirupsen/logrus"

	"github.com/vishvananda/netlink"
)

type macvlan struct{}

func init() {
	mac := &macvlan{}
	network.SetNetworkInterface(network.MACVLAN, mac)
}

func (m *macvlan) Name() string {
	return network.MACVLAN
}

// 设置macvlan网络
func (m *macvlan) SetupNetwork(cInfo *container.ContainerInfo) error {
	la := netlink.NewLinkAttrs()
	if cInfo.NetDeviceName == "" {
		cInfo.NetDeviceName = network.GetInterfaceName(cInfo.Name)
	}

	la.Name = cInfo.NetDeviceName
	var subnet *net.IPNet
	var ip net.IP
	var gwIP net.IP
	if cInfo.Subnet != "" {
		netIp, sub, err := net.ParseCIDR(cInfo.Subnet)
		if err != nil {
			logrus.Errorf("ParseSubnet error:%s:%v", cInfo.Subnet, err)
			return fmt.Errorf("ParseSubnet:%s:%v", cInfo.Subnet, err)
		}
		subnet = sub
		ip = netIp

		gwIP = net.ParseIP(cInfo.Gateway)
		if gwIP == nil {
			logrus.Errorf("ParseIP error:%s", cInfo.Gateway)
			return fmt.Errorf("ParseIP:%s", cInfo.Gateway)
		}
	}

	parent, err := netlink.LinkByName(cInfo.ParentNetwork)
	if err != nil {
		logrus.Errorf("get parent interface error:%s:%v", cInfo.ParentNetwork, err)
		return fmt.Errorf("get parent interface error:%s:%v", cInfo.ParentNetwork, err)
	}
	la.ParentIndex = parent.Attrs().Index

	// 添加网络接口
	mac := &netlink.Macvlan{LinkAttrs: la, Mode: netlink.MACVLAN_MODE_BRIDGE}
	if err := netlink.LinkAdd(mac); err != nil {
		logrus.Errorf("Add link error:%v", err)
		return fmt.Errorf("Add link error:%v", err)
	}

	pid, err := strconv.Atoi(cInfo.Pid)
	if err != nil {
		logrus.Errorf("Container Pid error %s:%v", cInfo.Name, err)
		return fmt.Errorf("Container Pid error %s:%v", cInfo.Name, err)
	}

	// 取得容器进程的net namespace
	nsPath := fmt.Sprintf("/proc/%s/ns/net", cInfo.Pid)
	file, err := os.Open(nsPath)
	if err != nil {
		logrus.Errorf("Get container net namespace error %s:%v", nsPath, err)
		return fmt.Errorf("Get container net namespace error %s:%v", nsPath, err)
	}
	fd := file.Fd()

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	// 将PID 设到 namespace里
	err = netlink.LinkSetNsFd(mac, int(fd))
	if err != nil {
		logrus.Errorf("Set link to ns error %d:%v", pid, err)
		return fmt.Errorf("Container Pid error %d:%v", pid, err)
	}

	// 取得当前的命名空间
	curNs, err := netns.Get()
	if err != nil {
		logrus.Errorf("Get current net namespace error:%v", err)
		return fmt.Errorf("Get current net namespace error:%v", err)
	}

	// 在函数结束时切换回原来的命名空间
	defer func() {
		// 切换到容器的命名空间
		err = netns.Set(netns.NsHandle(curNs))
		if err != nil {
			logrus.Errorf("Change net namespace error:%v", err)
		}
	}()

	// 切换到容器的命名空间
	err = netns.Set(netns.NsHandle(fd))
	if err != nil {
		logrus.Errorf("Change net namespace error:%v", err)
		return fmt.Errorf("Change net namespace error:%v", err)
	}

	if cInfo.Subnet != "" {
		// 设置静态IP
		addr := &netlink.Addr{}
		addr.IPNet = subnet
		addr.IP = ip

		if err := netlink.AddrAdd(mac, addr); err != nil {
			logrus.Errorf("Set link up error %d:%v", pid, err)
			return fmt.Errorf("Set link up error %d:%v", pid, err)
		}
	} else {
		// 获取DHCP
		// cmd := exec.Command("/sbin/dhclient")
		// output, err := cmd.CombinedOutput()
		// if err != nil {
		// 	logrus.Errorf("Execute dhcp cmd error %s:%v", output, err)
		// 	return fmt.Errorf("Execute dhcp cmd error %s:%v", output, err)
		// }
	}

	// 启动网络接口
	err = setInterfaceUP(mac.Attrs().Name)
	if err != nil {
		logrus.Errorf("Set link up error %s:%v", mac.Attrs().Name, err)
		return fmt.Errorf("Set link up error %s:%v", mac.Attrs().Name, err)
	}

	err = setInterfaceUP("lo")
	if err != nil {
		logrus.Errorf("Set link up error %s:%v", "lo", err)
		return fmt.Errorf("Set link up error %s:%v", "lo", err)
	}

	if cInfo.Subnet != "" {
		// 设定路由信息
		route := &netlink.Route{LinkIndex: mac.Attrs().Index, Table: 254, Gw: gwIP}
		if err = netlink.RouteAdd(route); err != nil {
			logrus.Errorf("set route error:%s:%v", cInfo.Gateway, err)
			return fmt.Errorf("set route error:%s:%v", cInfo.Gateway, err)
		}
	}

	return nil
}

func (m *macvlan) ShutdownNetwork(devName string) error {
	link, err := netlink.LinkByName(devName)
	if err != nil {
		logrus.Errorf("shutdown network error:%s:%v", devName, err)
		return fmt.Errorf("shutdown network error:%s:%v", devName, err)
	}
	return netlink.LinkDel(link)
}

func setInterfaceUP(interfaceName string) error {
	iface, err := netlink.LinkByName(interfaceName)
	if err != nil {
		logrus.Errorf("Error retrieving a link named [ %s ]: %v", iface.Attrs().Name, err)
		return fmt.Errorf("Error retrieving a link named [ %s ]: %v", iface.Attrs().Name, err)
	}

	if err := netlink.LinkSetUp(iface); err != nil {
		logrus.Errorf("Error enabling interface for %s: %v", interfaceName, err)
		return fmt.Errorf("Error enabling interface for %s: %v", interfaceName, err)
	}
	return nil
}
