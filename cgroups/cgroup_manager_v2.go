package cgroups

import (
	"github.com/sirupsen/logrus"
	"mydocker/cgroups/fs2"
	"mydocker/cgroups/resource"
)

type CgroupManagerV2 struct {
	Path       string
	Resource   *resource.ResourceConfig
	Subsystems []resource.Subsystem
}

func NewCgroupManagerV2(path string) *CgroupManagerV2 {
	return &CgroupManagerV2{
		Path:       path,
		Subsystems: fs2.Subsystems,
	}
}

// Apply 将进程pid加入到这个cgroup中
func (c *CgroupManagerV2) Apply(pid int) error {
	for _, subSysIns := range c.Subsystems {
		err := subSysIns.Apply(c.Path, pid)
		if err != nil {
			logrus.Errorf("apply subsystem:%s err:%s", subSysIns.Name(), err)
		}
	}
	return nil
}

// Set 设置cgroup资源限制
func (c *CgroupManagerV2) Set(res *resource.ResourceConfig) error {
	for _, subSysIns := range c.Subsystems {
		err := subSysIns.Set(c.Path, res)
		if err != nil {
			logrus.Errorf("apply subsystem:%s err:%s", subSysIns.Name(), err)
		}
	}
	return nil
}

// Destroy 释放cgroup
func (c *CgroupManagerV2) Destroy() error {
	for _, subSysIns := range c.Subsystems {
		if err := subSysIns.Remove(c.Path); err != nil {
			logrus.Warnf("remove cgroup fail %v", err)
		}
	}
	return nil
}
