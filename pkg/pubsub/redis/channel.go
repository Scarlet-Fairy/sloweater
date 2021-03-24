package redis

import "fmt"

func ImageBuildChannel(jobId string) string {
	return fmt.Sprintf("/job/%s", jobId)
}
