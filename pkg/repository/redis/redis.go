package redis

import (
	"context"
	"encoding/json"
	"github.com/Scarlet-Fairy/sloweater/pkg/service"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type redisRepository struct {
	client *redis.Client
}

const (
	KeyJobIdList = "jobs"
)

var (
	ErrInvalidJobFormat    = errors.New("invalid job format")
	ErrInvalidJobsIdFormat = errors.New("invalid jobs id format")
)

type jobList struct {
	JobsId []string `json:"jobs_id"`
}

func New(client *redis.Client) service.Repository {
	return redisRepository{
		client: client,
	}
}

func (r redisRepository) CreateJob(ctx context.Context) (string, error) {
	job := service.Job{
		Id:     uuid.NewString(),
		Status: service.StatusLoading,
		Step:   service.StepInit,
	}

	if err := r.client.Set(ctx, job.Id, job, 0).Err(); err != nil {
		return "", errors.Wrap(err, "creating job")
	}

	return job.Id, nil
}

func (r redisRepository) GetJob(ctx context.Context, jobId string) (*service.Job, error) {
	cmd := r.client.Get(ctx, jobId)
	if err := cmd.Err(); err != nil {
		return nil, errors.Wrap(err, "getting job")
	}

	job := &service.Job{}
	if err := json.Unmarshal([]byte(cmd.String()), job); err != nil {
		return nil, ErrInvalidJobFormat
	}

	return job, nil
}

func (r redisRepository) ListJobs(ctx context.Context) ([]string, error) {
	cmd := r.client.Get(ctx, KeyJobIdList)
	if err := cmd.Err(); err != nil {
		return nil, errors.Wrap(err, "listing jobs")
	}

	jobsId := &jobList{}
	if err := json.Unmarshal([]byte(cmd.String()), jobsId); err != nil {
		return nil, ErrInvalidJobsIdFormat
	}

	return jobsId.JobsId, nil
}

func (r redisRepository) SetJobStep(ctx context.Context, jobId string, step service.Step) error {
	job, err := r.GetJob(ctx, jobId)
	if err != nil {
		return err
	}

	job.Step = step

	if err := r.client.Set(ctx, jobId, job, 0).Err(); err != nil {
		return errors.Wrap(err, "settings job")
	}

	return nil
}

func (r redisRepository) SetJobStatus(ctx context.Context, jobId string, status service.Status) error {
	job, err := r.GetJob(ctx, jobId)
	if err != nil {
		return err
	}

	job.Status = status

	if err := r.client.Set(ctx, jobId, job, 0).Err(); err != nil {
		return errors.Wrap(err, "setting job")
	}

	return nil
}
