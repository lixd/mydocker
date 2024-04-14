package cgroups

import (
	log "github.com/sirupsen/logrus"
	"mydocker/cgroups/resource"
)

type CgroupManager interface {
	Apply(pid int) error
	Set(res *resource.ResourceConfig) error
	Destroy() error
}

func NewCgroupManager(path string) CgroupManager {
	if IsCgroup2UnifiedMode() {
		log.Infof("use cgroup v2")
		return NewCgroupManagerV2(path)
	}
	log.Infof("use cgroup v1")
	return NewCgroupManagerV1(path)
}
