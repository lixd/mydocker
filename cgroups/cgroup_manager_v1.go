package cgroups

import (
	"github.com/sirupsen/logrus"
	"mydocker/cgroups/fs"
	"mydocker/cgroups/resource"
)

type CgroupManagerV1 struct {
	// cgroup在hierarchy中的路径 相当于创建的cgroup目录相对于root cgroup目录的路径
	Path string
	// 资源配置
	Resource   *resource.ResourceConfig
	Subsystems []resource.Subsystem
}

func NewCgroupManagerV1(path string) *CgroupManagerV1 {
	return &CgroupManagerV1{
		Path:       path,
		Subsystems: fs.SubsystemsIns,
	}
}

// Apply 将进程pid加入到这个cgroup中
func (c *CgroupManagerV1) Apply(pid int) error {
	for _, subSysIns := range c.Subsystems {
		err := subSysIns.Apply(c.Path, pid)
		if err != nil {
			logrus.Errorf("apply subsystem:%s err:%s", subSysIns.Name(), err)
		}
	}
	return nil
}

// Set 设置cgroup资源限制
func (c *CgroupManagerV1) Set(res *resource.ResourceConfig) error {
	for _, subSysIns := range c.Subsystems {
		err := subSysIns.Set(c.Path, res)
		if err != nil {
			logrus.Errorf("apply subsystem:%s err:%s", subSysIns.Name(), err)
		}
	}
	return nil
}

// Destroy 释放cgroup
func (c *CgroupManagerV1) Destroy() error {
	for _, subSysIns := range c.Subsystems {
		if err := subSysIns.Remove(c.Path); err != nil {
			logrus.Warnf("remove cgroup fail %v", err)
		}
	}
	return nil
}
