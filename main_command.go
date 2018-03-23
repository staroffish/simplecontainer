package main

import (
	"fmt"
	"process"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var runCommand = cli.Command{
	Name: "run",
	Usage: `Create a new container
	sc run [-m memroy_limit(mega)] [-cpu core_num] [-name container_name] imagename`,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "name",
			Usage: "container name",
		},
	},
	Action: func(ctx *cli.Context) error {
		if len(ctx.Args()) < 1 {
			return fmt.Errorf("Missing imagename")
		}

		name := ctx.String("name")
		imageName := ctx.Args().Get(0)

		return Run(name, imageName)
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
		process.ContianerProcess(name)

		return nil
	},
}
