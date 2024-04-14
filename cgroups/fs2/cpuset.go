package fs2

import (
	"fmt"
	"mydocker/cgroups/resource"
	"mydocker/constant"
	"os"
	"path"
)

type CpusetSubSystem struct {
}

func (s *CpusetSubSystem) Name() string {
	return "cpuset"
}

func (s *CpusetSubSystem) Set(cgroupPath string, res *resource.ResourceConfig) error {
	if res.CpuSet == "" {
		return nil
	}
	subCgroupPath, err := getCgroupPath(cgroupPath, true)
	if err != nil {
		return err
	}
	if err := os.WriteFile(path.Join(subCgroupPath, "cpuset.cpus"), []byte(res.CpuSet), constant.Perm0644); err != nil {
		return fmt.Errorf("set cgroup cpuset fail %v", err)
	}
	return nil
}

func (s *CpusetSubSystem) Apply(cgroupPath string, pid int) error {
	return applyCgroup(pid, cgroupPath)
}

func (s *CpusetSubSystem) Remove(cgroupPath string) error {
	subsysCgroupPath, err := getCgroupPath(cgroupPath, false)
	if err != nil {
		return err
	}
	return os.RemoveAll(subsysCgroupPath)
}
