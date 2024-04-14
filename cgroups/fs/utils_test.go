package fs

import (
	"testing"
)

func TestFindCgroupMountpoint(t *testing.T) {
	t.Logf("cpu subsystem mount point %v\n", findCgroupMountpoint("cpu"))
	t.Logf("cpuset subsystem mount point %v\n", findCgroupMountpoint("cpuset"))
	t.Logf("memory subsystem mount point %v\n", findCgroupMountpoint("memory"))
}
