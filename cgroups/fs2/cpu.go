package fs2

import (
	"fmt"
	"mydocker/cgroups/resource"
	"os"
	"path"
	"strconv"

	"mydocker/constant"
)

type CpuSubSystem struct {
}

const (
	PeriodDefault = 100000
	Percent       = 100
)

func (s *CpuSubSystem) Name() string {
	return "cpu"
}

func (s *CpuSubSystem) Set(cgroupPath string, res *resource.ResourceConfig) error {
	if res.CpuCfsQuota == 0 {
		return nil
	}
	subCgroupPath, err := getCgroupPath(cgroupPath, true)
	if err != nil {
		return err
	}

	// cpu.cfs_period_us & cpu.cfs_quota_us 控制的是CPU使用时间，单位是微秒，比如每1秒钟，这个进程只能使用200ms，相当于只能用20%的CPU
	// v2 中直接将 cpu.cfs_period_us & cpu.cfs_quota_us 统一记录到 cpu.max 中，比如 5000 10000 这样就是限制使用 50% cpu
	if res.CpuCfsQuota != 0 {
		// cpu.cfs_quota_us 则根据用户传递的参数来控制，比如参数为20，就是限制为20%CPU，所以把cpu.cfs_quota_us设置为cpu.cfs_period_us的20%就行
		// 这里只是简单的计算了下，并没有处理一些特殊情况，比如负数什么的
		if err = os.WriteFile(path.Join(subCgroupPath, "cpu.max"), []byte(fmt.Sprintf("%s %s", strconv.Itoa(PeriodDefault/Percent*res.CpuCfsQuota), PeriodDefault)), constant.Perm0644); err != nil {
			return fmt.Errorf("set cgroup cpu share fail %v", err)
		}
	}
	return nil
}

func (s *CpuSubSystem) Apply(cgroupPath string, pid int) error {
	return applyCgroup(pid, cgroupPath)
}

func (s *CpuSubSystem) Remove(cgroupPath string) error {
	subCgroupPath, err := getCgroupPath(cgroupPath, false)
	if err != nil {
		return err
	}
	return os.RemoveAll(subCgroupPath)
}
