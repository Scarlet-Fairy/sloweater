package service

import (
	"fmt"
)

type WorkloadId string

const (
	JobTypeImageBuild = "imagebuild"
	JobTypeWorkload   = "workload"
	NameImageBuilder  = "cobold"
)

func (id WorkloadId) NameImageBuild() string {
	return fmt.Sprintf("%s.%s", JobTypeImageBuild, id)
}

func (id WorkloadId) NameWorkload() string {
	return fmt.Sprintf("%s.%s", JobTypeWorkload, id)
}

func (id WorkloadId) ImageName(registry string) string {
	return fmt.Sprintf("%s/%s/%s", registry, NameImageBuilder, id)
}
