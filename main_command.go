package main

import (
	"fmt"
	"os"

	container "github.com/staroffish/simplecontainer/container"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var runCommand = cli.Command{
	Name: "run",
	Usage: `Create a new container
	simplecontainer run [-m memroy_limit(mega)] [-cpu core_num] [-name container_name] [-net dhcp|static -ip ip -parent parent_dev -gateway gateway_ip] imagename`,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "name",
			Usage: "container name",
		},
		cli.StringFlag{
			Name:  "net",
			Usage: `container network setting`,
		},
		cli.StringFlag{
			Name:  "ip",
			Usage: `container ip address`,
		},
		cli.StringFlag{
			Name:  "parent",
			Usage: `parent device for macvlan`,
		},
		cli.StringFlag{
			Name:  "gateway",
			Usage: `container gateway ip`,
		},
	},
	Action: func(ctx *cli.Context) error {
		if len(ctx.Args()) < 1 {
			return fmt.Errorf("Missing imagename")
		}

		cInfo := &container.ContainerInfo{
			Name:      ctx.String("name"),
			ImageName: ctx.Args().Get(0),
		}

		cInfo.NetType = ctx.String("net")
		if cInfo.NetType != "" {
			if cInfo.NetType != "dhcp" && cInfo.NetType != "static" {
				return fmt.Errorf("net option must be dhcp or static")
			}
			cInfo.ParentNetwork = ctx.String("parent")
			if cInfo.ParentNetwork == "" {
				return fmt.Errorf("missing parent device name")
			}
			if cInfo.NetType == "static" {
				cInfo.Subnet = ctx.String("ip")
				if cInfo.Subnet == "" {
					return fmt.Errorf("Static network must specify IP")
				}
				cInfo.Gateway = ctx.String("gateway")
				if cInfo.Gateway == "" {
					return fmt.Errorf("missing Gateway")
				}
			}

		}

		return run(cInfo)
	},
}

var initCommand = cli.Command{
	Name:  "init",
	Usage: `Init container process run user's process in container. Do not call it outside`,
	Action: func(ctx *cli.Context) error {
		if len(ctx.Args()) < 1 {
			logrus.Errorf("missing container name")
			return fmt.Errorf("missing container name")
		}

		name := ctx.Args().Get(0)

		logrus.Infof("Enter container process:%s", name)
		container.ContianerProcess(name)

		return nil
	},
}

var execCommand = cli.Command{
	Name: "exec",
	Usage: `Execute the command in the container
	simplecontainer exec container_name command`,
	Action: func(ctx *cli.Context) error {
		_, ok := os.LookupEnv("mydocker_pid")
		if ok {
			return nil
		}
		if len(ctx.Args()) < 2 {
			logrus.Errorf("missing container name or command")
			return fmt.Errorf("missing container name or command")
		}

		containerName := ctx.Args().Get(0)
		cmd := ctx.Args().Get(1)

		if err := execCmd(containerName, cmd); err != nil {
			return err
		}

		return nil
	},
}

var startCommand = cli.Command{
	Name: "start",
	Usage: `Start container
	simplecontainer start container_name`,
	Action: func(ctx *cli.Context) error {
		if len(ctx.Args()) < 1 {
			logrus.Errorf("missing container name")
			return fmt.Errorf("missing container name")
		}

		if err := start(ctx.Args().Get(0)); err != nil {
			return err
		}

		return nil
	},
}

var stopCommand = cli.Command{
	Name: "stop",
	Usage: `Stop container
	simplecontainer stop container_name`,
	Action: func(ctx *cli.Context) error {
		if len(ctx.Args()) < 1 {
			logrus.Errorf("missing container name")
			return fmt.Errorf("missing container name")
		}

		if err := stop(ctx.Args().Get(0)); err != nil {
			return err
		}

		return nil
	},
}

var rmCommand = cli.Command{
	Name: "rm",
	Usage: `Remove container
	simplecontainer rm container_name`,
	Action: func(ctx *cli.Context) error {
		if len(ctx.Args()) < 1 {
			logrus.Errorf("missing container name")
			return fmt.Errorf("missing container name")
		}

		return remove(ctx.Args().Get(0))
	},
}

var psCommand = cli.Command{
	Name: "ps",
	Usage: `List up containers
	simplecontainer ps [-a]`,
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "a",
			Usage: "show all containers",
		},
	},
	Action: func(ctx *cli.Context) error {
		addFlg := ctx.Bool("a")

		return ps(addFlg)
	},
}

var imageCommand = cli.Command{
	Name: "image",
	Usage: `Image operation
	 simplecontainer image [command]`,
	Subcommands: []cli.Command{
		cli.Command{
			Name: "commit",
			Usage: `Create a new image from a container's changes
			simplecontainer image commit container_name new_image_name`,
			Action: func(ctx *cli.Context) error {
				if len(ctx.Args()) < 2 {
					logrus.Errorf("missing container_name or image_name")
					return fmt.Errorf("missing container_name or new_image_name")
				}

				return commit(ctx.Args().Get(0), ctx.Args().Get(1))
			},
		},
		cli.Command{
			Name: "import",
			Usage: `Import image from a tar file that compressed into gzip
			simplecontainer image import -name imagName file_name`,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "name",
					Usage: "image name",
				},
			},
			Action: func(ctx *cli.Context) error {
				if len(ctx.Args()) < 1 {
					logrus.Errorf("missing file_name")
					return fmt.Errorf("missing file_name")
				}

				imageName := ctx.String("name")

				return importImage(ctx.Args().Get(0), imageName)
			},
		},
		cli.Command{
			Name: "export",
			Usage: `Export image to a tar file that compressed into gzip
			simplecontainer image export image_name dst_path`,
			Action: func(ctx *cli.Context) error {
				if len(ctx.Args()) < 2 {
					logrus.Errorf("missing image_name or dst_path")
					return fmt.Errorf("missing image_name or dst_path")
				}

				return exportImage(ctx.Args().Get(0), ctx.Args().Get(1))
			},
		},
		cli.Command{
			Name: "rm",
			Usage: `Remove image
			simplecontainer image rm image_name`,
			Action: func(ctx *cli.Context) error {
				if len(ctx.Args()) < 1 {
					logrus.Errorf("missing image_name")
					return fmt.Errorf("missing image_name")
				}

				return removeImage(ctx.Args().Get(0))
			},
		},
	},
	Action: func(ctx *cli.Context) error {
		return imageList()
	},
}
