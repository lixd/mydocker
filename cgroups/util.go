package cgroups

import (
	"os"
	"sync"

	"golang.org/x/sys/unix"
)

const (
	unifiedMountpoint = "/sys/fs/cgroup"
)

var (
	isUnifiedOnce sync.Once
	isUnified     bool
)

// IsCgroup2UnifiedMode returns whether we are running in cgroup v2 unified mode.
func IsCgroup2UnifiedMode() bool {
	isUnifiedOnce.Do(func() {
		var st unix.Statfs_t
		err := unix.Statfs(unifiedMountpoint, &st)
		if err != nil && os.IsNotExist(err) {
			// For rootless containers, sweep it under the rug.
			isUnified = false
			return
		}
		isUnified = st.Type == unix.CGROUP2_SUPER_MAGIC
	})
	return isUnified
}
