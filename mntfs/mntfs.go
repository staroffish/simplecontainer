package mntfs

type MountFS interface {
	InitMnt(name, imageName string) error
	Mount(name, imageName string) error
	Unmount(name string) error
	Remove(name string) error
}

var mntInst = make(map[string]MountFS)

const (
	OVERLAY = "overlay"
)

func SetMntInst(fsName string, mntFs MountFS) {
	mntInst[fsName] = mntFs
}

func GetMountInst(fsName string) MountFS {
	return mntInst[fsName]
}
