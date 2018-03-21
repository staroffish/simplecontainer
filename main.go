package main

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "sc"
	app.Usage = "A very simple container runtime implatemention."
	app.Commands = []cli.Command{runCommand}
	app.Version = "0.0.1"

	app.Before = func(ctx *cli.Context) error {
		logrus.SetFormatter(&logrus.TextFormatter{})
		file, err := os.OpenFile("./sc.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			fmt.Printf("Error to open ./sc.log:%v\n", err)
			os.Exit(-1)
		}
		logrus.SetOutput(file)
		return nil
	}

	if err := app.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}
}
