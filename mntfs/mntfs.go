package mntfs

type MountFS interface {
	InitMnt(name, imageName string) error
	// Mount(name, imageName string) error
	// UnMount(name, imageName string) error
	Mount() error
	Unmount() error
}

var mntInst = make(map[string]MountFS)

func SetMntInst(fsName string, mntFs MountFS) {
	mntInst[fsName] = mntFs
}

func GetMountInst(fsName string) MountFS {
	return mntInst[fsName]
}
