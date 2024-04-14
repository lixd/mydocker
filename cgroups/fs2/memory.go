package fs2

import (
	"fmt"
	"mydocker/cgroups/resource"
	"mydocker/constant"
	"os"
	"path"
)

type MemorySubSystem struct {
}

// Name 返回cgroup名字
func (s *MemorySubSystem) Name() string {
	return "memory"
}

// Set 设置cgroupPath对应的cgroup的内存资源限制
func (s *MemorySubSystem) Set(cgroupPath string, res *resource.ResourceConfig) error {
	if res.MemoryLimit == "" {
		return nil
	}
	subCgroupPath, err := getCgroupPath(cgroupPath, true)
	if err != nil {
		return err
	}
	// 设置这个cgroup的内存限制，即将限制写入到cgroup对应目录的memory.limit_in_bytes 文件中。
	if err := os.WriteFile(path.Join(subCgroupPath, "memory.max"), []byte(res.MemoryLimit), constant.Perm0644); err != nil {
		return fmt.Errorf("set cgroup memory fail %v", err)
	}
	return nil
}

// Apply 将pid加入到cgroupPath对应的cgroup中
func (s *MemorySubSystem) Apply(cgroupPath string, pid int) error {
	return applyCgroup(pid, cgroupPath)
}

// Remove 删除cgroupPath对应的cgroup
func (s *MemorySubSystem) Remove(cgroupPath string) error {
	subCgroupPath, err := getCgroupPath(cgroupPath, false)
	if err != nil {
		return err
	}
	return os.RemoveAll(subCgroupPath)
}
