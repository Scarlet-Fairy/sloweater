package service

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
	Id     string
	Status Status
	Step   Step
}
