package container

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"mydocker/constant"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

func RecordContainerInfo(containerPID int, commandArray []string, containerName, containerId string) error {
	// 如果未指定容器名，则使用随机生成的containerID
	if containerName == "" {
		containerName = containerId
	}
	command := strings.Join(commandArray, "")
	containerInfo := &Info{
		Id:          containerId,
		Pid:         strconv.Itoa(containerPID),
		Command:     command,
		CreatedTime: time.Now().Format("2006-01-02 15:04:05"),
		Status:      RUNNING,
		Name:        containerName,
	}

	jsonBytes, err := json.Marshal(containerInfo)
	if err != nil {
		return errors.WithMessage(err, "container info marshal failed")
	}
	jsonStr := string(jsonBytes)
	// 拼接出存储容器信息文件的路径，如果目录不存在则级联创建
	dirPath := fmt.Sprintf(InfoLocFormat, containerInfo)
	if err := os.MkdirAll(dirPath, constant.Perm0622); err != nil {
		return errors.WithMessagef(err, "mkdir %s failed", dirPath)
	}
	// 将容器信息写入文件
	fileName := path.Join(dirPath, ConfigName)
	file, err := os.Create(fileName)
	defer file.Close()
	if err != nil {
		return errors.WithMessagef(err, "create file %s failed", fileName)
	}
	if _, err = file.WriteString(jsonStr); err != nil {
		return errors.WithMessagef(err, "write container info to  file %s failed", fileName)
	}
	return nil
}

func DeleteContainerInfo(containerID string) {
	dirPath := fmt.Sprintf(InfoLocFormat, containerID)
	if err := os.RemoveAll(dirPath); err != nil {
		log.Errorf("Remove dir %s error %v", dirPath, err)
	}
}

func GenerateContainerID() string {
	return randStringBytes(IDLength)
}

func randStringBytes(n int) string {
	letterBytes := "1234567890"
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
