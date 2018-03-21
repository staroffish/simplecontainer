package main

import (
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var runCommand = cli.Command{
	Name: "run",
	Usage: `Create a new container
	sc run [-m memroy_limit(mega)] [-cpu core_num] [-name container_name] imagename`,
	Action: func(ctx *cli.Context) error {
		if len(ctx.Args()) < 1 {
			logrus.Errorf("Missing imagename")
		}

		return nil
	},
}
