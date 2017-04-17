package protocol

import (
	"dkvgo/job"
)

type HeartBeatBag struct {
	Todo string
	Echo string
	Task *job.Task
	Report *job.TaskState
}