package container

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"math/rand"
	"mydocker/constant"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

func RecordContainerInfo(containerPID int, commandArray []string, containerName, containerId, volume, networkName, ip string, portMapping []string) (*Info, error) {
	// 如果未指定容器名，则使用随机生成的containerID
	if containerName == "" {
		containerName = containerId
	}
	command := strings.Join(commandArray, "")
	containerInfo := &Info{
		Pid:         strconv.Itoa(containerPID),
		Id:          containerId,
		Name:        containerName,
		Command:     command,
		CreatedTime: time.Now().Format("2006-01-02 15:04:05"),
		Status:      RUNNING,
		Volume:      volume,
		NetworkName: networkName,
		PortMapping: portMapping,
		IP:          ip,
	}

	jsonBytes, err := json.Marshal(containerInfo)
	if err != nil {
		return containerInfo, errors.WithMessage(err, "container info marshal failed")
	}
	jsonStr := string(jsonBytes)
	// 拼接出存储容器信息文件的路径，如果目录不存在则级联创建
	dirPath := fmt.Sprintf(InfoLocFormat, containerId)
	if err = os.MkdirAll(dirPath, constant.Perm0622); err != nil {
		return containerInfo, errors.WithMessagef(err, "mkdir %s failed", dirPath)
	}
	// 将容器信息写入文件
	fileName := path.Join(dirPath, ConfigName)
	file, err := os.Create(fileName)
	if err != nil {
		return containerInfo, errors.WithMessagef(err, "create file %s failed", fileName)
	}
	defer file.Close()
	if _, err = file.WriteString(jsonStr); err != nil {
		return containerInfo, errors.WithMessagef(err, "write container info to  file %s failed", fileName)
	}
	return containerInfo, nil
}

func DeleteContainerInfo(containerID string) error {
	dirPath := fmt.Sprintf(InfoLocFormat, containerID)
	if err := os.RemoveAll(dirPath); err != nil {
		return errors.WithMessagef(err, "remove dir %s failed", dirPath)
	}
	return nil
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

// GetLogfile build logfile name by containerId
func GetLogfile(containerId string) string {
	return fmt.Sprintf(LogFile, containerId)
}
