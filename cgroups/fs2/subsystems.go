package fs2

import (
	"mydocker/cgroups/resource"
)

var Subsystems = []resource.Subsystem{
	&CpusetSubSystem{},
	&MemorySubSystem{},
	&CpuSubSystem{},
}
