package container

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"mydocker/constant"
	"os"
	"os/exec"
	"path"
	"strings"
)

// volume 相关工具

// mountVolume 使用 bind mount 挂载 volume
func mountVolume(mntPath, hostPath, containerPath string) {
	// 创建宿主机目录
	if err := os.Mkdir(hostPath, constant.Perm0777); err != nil {
		log.Infof("mkdir parent dir %s error. %v", hostPath, err)
	}
	// 拼接出对应的容器目录在宿主机上的的位置，并创建对应目录
	containerPathInHost := path.Join(mntPath, containerPath)
	if err := os.Mkdir(containerPathInHost, constant.Perm0777); err != nil {
		log.Infof("mkdir container dir %s error. %v", containerPathInHost, err)
	}
	// 通过bind mount 将宿主机目录挂载到容器目录
	// mount -o bind /hostPath /containerPath
	cmd := exec.Command("mount", "-o", "bind", hostPath, containerPathInHost)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Errorf("mount volume failed. %v", err)
	}
}

func umountVolume(mntPath, containerPath string) {
	// mntPath 为容器在宿主机上的挂载点，例如 /root/merged
	// containerPath 为 volume 在容器中对应的目录，例如 /root/tmp
	// containerPathInHost 则是容器中目录在宿主机上的具体位置，例如 /root/merged/root/tmp
	containerPathInHost := path.Join(mntPath, containerPath)
	cmd := exec.Command("umount", containerPathInHost)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Errorf("Umount volume failed. %v", err)
	}
}

// volumeExtract 通过冒号分割解析volume目录，比如 -v /tmp:/tmp
func volumeExtract(volume string) (sourcePath, destinationPath string, err error) {
	parts := strings.Split(volume, ":")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid volume [%s], must split by `:`", volume)
	}

	sourcePath, destinationPath = parts[0], parts[1]
	if sourcePath == "" || destinationPath == "" {
		return "", "", fmt.Errorf("invalid volume [%s], path can't be empty", volume)
	}

	return sourcePath, destinationPath, nil
}
