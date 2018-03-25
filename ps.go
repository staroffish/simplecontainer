package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"text/tabwriter"

	"github.com/staroffish/simplecontainer/config"
	container "github.com/staroffish/simplecontainer/container"

	"github.com/sirupsen/logrus"
)

func ps(allFlg bool) error {
	//取得容器信息
	files, err := ioutil.ReadDir(config.CInfopath)
	if err != nil {
		logrus.Errorf("Read dir error:%s:%v", config.CInfopath, err)
		return fmt.Errorf("Read dir error:%s:%v", config.CInfopath, err)
	}

	cInfos := make([]*container.ContainerInfo, 0)
	for _, files := range files {
		if files.IsDir() {
			cInfo := container.ReadContainerInfo(files.Name())
			if cInfo == nil ||
				(cInfo.Status != container.RUNNING && !allFlg) {
				continue
			}
			cInfos = append(cInfos, cInfo)
		}
	}

	w := tabwriter.NewWriter(os.Stdout, 12, 1, 3, ' ', 0)
	fmt.Fprint(w, "NAME\tPID\tSTATUS\tIMAGE\tCREATED\n")
	for _, item := range cInfos {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
			item.Name,
			item.Pid,
			item.Status,
			item.ImageName,
			item.CreateTime)
	}
	if err := w.Flush(); err != nil {
		logrus.Errorf("Flush error %v", err)
		return err
	}

	return nil
}
