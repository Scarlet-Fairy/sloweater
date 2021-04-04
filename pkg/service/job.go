package service

import (
	"encoding/json"
	"fmt"
)

type Status byte

const (
	StatusError     Status = 1
	StatusLoading   Status = 2
	StatusCompleted Status = 3
)

func (s Status) IsValid() bool {
	return s == StatusError || s == StatusLoading || s == StatusCompleted
}

type Step byte

const (
	StepInit  Step = 0
	StepClone Step = 1
	StepBuild Step = 2
	StepPush  Step = 3
)

func (s Step) IsValid() bool {
	return s == StepInit || s == StepClone || s == StepBuild || s == StepPush
}

type Job struct {
	Id     string  `json:"id"`
	Status Status  `json:"status"`
	Steps  []Steps `json:"steps"`
}

type Steps struct {
	Step  Step   `json:"step"`
	Error string `json:"error"`
}

func (job Job) MarshalBinary() ([]byte, error) {
	return json.Marshal(job)
}

type JobId string

const (
	JobTypeImageBuild = "ImageBuild"
	NameImageBuilder  = "cobold"
)

func (id JobId) NameImageBuild() string {
	return fmt.Sprintf("%s.%s", JobTypeImageBuild, id)
}

func (id JobId) ImageName(registry string) string {
	return fmt.Sprintf("%s/%s/%s", registry, NameImageBuilder, id)
}
