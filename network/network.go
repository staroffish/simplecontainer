package network

import container "github.com/staroffish/simplecontainer/container"

type NetworkInterface interface {
	Name() string
	SetupNetwork(cInfo *container.ContainerInfo) error
	ShutdownNetwork(devName string) error
	// SetupMasterNet() error
	// SetupContainerNet(containerName string) error
}

var networks = make(map[string]NetworkInterface)

var (
	MACVLAN = "macvlan"
)

func SetNetworkInterface(name string, network NetworkInterface) {
	networks[name] = network
}

func GetNetworkInterface(name string) NetworkInterface {
	return networks[name]
}

func GetInterfaceName(name string) string {
	return name[0:3] + "-" + container.RandStringBytes(6)
}
