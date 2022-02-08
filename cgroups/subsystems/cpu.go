package subsystems

import (
	"fmt"
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

func (s *CpuSubSystem) Set(cgroupPath string, res *ResourceConfig) error {
	if res.CpuCfsQuota == 0 && res.CpuShare == "" {
		return nil
	}
	subsysCgroupPath, err := getCgroupPath(s.Name(), cgroupPath, true)
	if err != nil {
		return err
	}
	// cpu.shares 控制的是CPU使用比例，不是绝对值
	if res.CpuShare != "" {
		if err = os.WriteFile(path.Join(subsysCgroupPath, "cpu.shares"), []byte(res.CpuShare), constant.Perm0644); err != nil {
			return fmt.Errorf("set cgroup cpu share fail %v", err)
		}
	}

	// cpu.cfs_period_us & cpu.cfs_quota_us 控制的是CPU使用时间，单位是微秒，比如每1秒钟，这个进程只能使用200ms，相当于只能用20%的CPU
	if res.CpuCfsQuota != 0 {
		// cpu.cfs_period_us 默认为100000，即100ms
		if err = os.WriteFile(path.Join(subsysCgroupPath, "cpu.cfs_period_us"), []byte(strconv.Itoa(PeriodDefault)), constant.Perm0644); err != nil {
			return fmt.Errorf("set cgroup cpu share fail %v", err)
		}
		// cpu.cfs_quota_us 则根据用户传递的参数来控制，比如参数为20，就是限制为20%CPU，所以把cpu.cfs_quota_us设置为cpu.cfs_period_us的20%就行
		// 这里只是简单的计算了下，并没有处理一些特殊情况，比如负数什么的
		if err = os.WriteFile(path.Join(subsysCgroupPath, "cpu.cfs_quota_us"), []byte(strconv.Itoa(PeriodDefault/Percent*res.CpuCfsQuota)), constant.Perm0644); err != nil {
			return fmt.Errorf("set cgroup cpu share fail %v", err)
		}
	}
	return nil
}

func (s *CpuSubSystem) Apply(cgroupPath string, pid int, res *ResourceConfig) error {
	if res.CpuCfsQuota == 0 && res.CpuShare == "" {
		return nil
	}

	subsysCgroupPath, err := getCgroupPath(s.Name(), cgroupPath, false)
	if err != nil {
		return fmt.Errorf("get cgroup %s error: %v", cgroupPath, err)
	}
	if err = os.WriteFile(path.Join(subsysCgroupPath, "tasks"), []byte(strconv.Itoa(pid)), constant.Perm0644); err != nil {
		return fmt.Errorf("set cgroup proc fail %v", err)
	}
	return nil
}

func (s *CpuSubSystem) Remove(cgroupPath string) error {
	subsysCgroupPath, err := getCgroupPath(s.Name(), cgroupPath, false)
	if err != nil {
		return err
	}
	return os.RemoveAll(subsysCgroupPath)
}
