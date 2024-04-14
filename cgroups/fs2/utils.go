package fs2

import (
	"fmt"
	"mydocker/constant"
	"os"
	"path"
	"strconv"

	"github.com/pkg/errors"
)

// getCgroupPath 找到cgroup在文件系统中的绝对路径
/*
实际就是将根目录和cgroup名称拼接成一个路径。
如果指定了自动创建，就先检测一下是否存在，如果对应的目录不存在，则说明cgroup不存在，这里就给创建一个
*/
func getCgroupPath(cgroupPath string, autoCreate bool) (string, error) {
	// 不需要自动创建就直接返回
	cgroupRoot := UnifiedMountpoint
	absPath := path.Join(cgroupRoot, cgroupPath)
	if !autoCreate {
		return absPath, nil
	}
	// 指定自动创建时才判断是否存在
	_, err := os.Stat(absPath)
	// 只有不存在才创建
	if err != nil && os.IsNotExist(err) {
		err = os.Mkdir(absPath, constant.Perm0755)
		return absPath, err
	}
	return absPath, errors.Wrap(err, "create cgroup")
}

func applyCgroup(pid int, cgroupPath string) error {
	subCgroupPath, err := getCgroupPath(cgroupPath, true)
	if err != nil {
		return errors.Wrapf(err, "get cgroup %s", cgroupPath)
	}
	if err = os.WriteFile(path.Join(subCgroupPath, "cgroup.procs"), []byte(strconv.Itoa(pid)),
		constant.Perm0644); err != nil {
		return fmt.Errorf("set cgroup proc fail %v", err)
	}
	return nil
}
