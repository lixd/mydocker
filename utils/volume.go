package utils

import "strings"

// VolumeUrlExtract 通过冒号分割解析volume目录，比如 -v /tmp:/tmp
func VolumeUrlExtract(volume string) []string {
	return strings.Split(volume, ":")
}
