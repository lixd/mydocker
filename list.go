package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"text/tabwriter"

	"mydocker/container"

	log "github.com/sirupsen/logrus"
)

func ListContainers() {
	// 读取存放容器信息目录下的所有文件
	files, err := os.ReadDir(container.InfoLoc)
	if err != nil {
		log.Errorf("read dir %s error %v", container.InfoLoc, err)
		return
	}
	containers := make([]*container.Info, 0, len(files))
	for _, file := range files {
		tmpContainer, err := getContainerInfo(file)
		if err != nil {
			log.Errorf("get container info error %v", err)
			continue
		}
		containers = append(containers, tmpContainer)
	}
	// 使用tabwriter.NewWriter在控制台打印出容器信息
	// tabwriter 是引用的text/tabwriter类库，用于在控制台打印对齐的表格
	w := tabwriter.NewWriter(os.Stdout, 12, 1, 3, ' ', 0)
	_, err = fmt.Fprint(w, "ID\tNAME\tPID\tIP\tSTATUS\tCOMMAND\tCREATED\n")
	if err != nil {
		log.Errorf("Fprint error %v", err)
	}
	for _, item := range containers {
		_, err = fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\t%s\n",
			item.Id,
			item.Name,
			item.Pid,
			item.IP,
			item.Status,
			item.Command,
			item.CreatedTime)
		if err != nil {
			log.Errorf("Fprint error %v", err)
		}
	}
	if err = w.Flush(); err != nil {
		log.Errorf("Flush error %v", err)
	}
}

func getContainerInfo(file os.DirEntry) (*container.Info, error) {
	// 根据文件名拼接出完整路径
	configFileDir := fmt.Sprintf(container.InfoLocFormat, file.Name())
	configFileDir = path.Join(configFileDir, container.ConfigName)
	// 读取容器配置文件
	content, err := os.ReadFile(configFileDir)
	if err != nil {
		log.Errorf("read file %s error %v", configFileDir, err)
		return nil, err
	}
	info := new(container.Info)
	if err = json.Unmarshal(content, info); err != nil {
		log.Errorf("json unmarshal error %v", err)
		return nil, err
	}

	return info, nil
}
