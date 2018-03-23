package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

var (
	MntPath      string
	WirtelayPath string
	ImagePath    string
	LogPath      string
	CInfopath    string
)

func init() {

	var scCfg = struct {
		MntPath      string `json:"mntpath"`
		WirtelayPath string `json:"writelaypath"`
		ImagePath    string `json:"imagepath"`
		LogPath      string `json:"logpath"`
		CInfopath    string `json:"cInfopath"`
	}{}
	data, err := ioutil.ReadFile("/etc/sc.json")
	if err != nil {
		data, err = ioutil.ReadFile("./sc.json")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Open config file error.\nThe configuration file sc.json must exist with /etc or the current directory.\n")
			os.Exit(-1)
		}
	}

	if err = json.Unmarshal(data, &scCfg); err != nil {
		fmt.Fprintf(os.Stderr, "Unmarshal config file error:%v\n", err)
		os.Exit(-1)
	}

	if scCfg.ImagePath == "" || scCfg.LogPath == "" || scCfg.WirtelayPath == "" || scCfg.MntPath == "" || scCfg.CInfopath == "" {
		fmt.Fprintf(os.Stderr, "Config file load error.\n")
		os.Exit(-1)
	}

	if !strings.HasSuffix(scCfg.ImagePath, "/") {
		scCfg.ImagePath = scCfg.ImagePath + "/"
	}
	ImagePath = scCfg.ImagePath

	if !strings.HasSuffix(scCfg.LogPath, "/") {
		scCfg.LogPath = scCfg.LogPath + "/"
	}
	LogPath = scCfg.LogPath

	if !strings.HasSuffix(scCfg.WirtelayPath, "/") {
		scCfg.WirtelayPath = scCfg.WirtelayPath + "/"
	}
	WirtelayPath = scCfg.WirtelayPath

	if !strings.HasSuffix(scCfg.MntPath, "/") {
		scCfg.MntPath = scCfg.MntPath + "/"
	}
	MntPath = scCfg.MntPath

	if !strings.HasSuffix(scCfg.CInfopath, "/") {
		scCfg.CInfopath = scCfg.CInfopath + "/"
	}
	CInfopath = scCfg.CInfopath
}
