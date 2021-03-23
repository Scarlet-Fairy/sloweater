package service

type Status byte

const (
	Error     Status = 1
	Loading   Status = 2
	Completed Status = 3
)

type Step byte

const (
	Clone Step = 1
	Build Step = 2
	Push  Step = 3
)

type Job struct {
	Id     string
	Status Status
	Step   Step
}
