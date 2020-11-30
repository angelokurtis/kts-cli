package kubectl

import (
	"github.com/enescakir/emoji"
	"time"
)

type (
	ContainerStateRunning struct {
		StartedAt time.Time `json:"startedAt"`
	}
	ContainerStateTerminated struct {
		ContainerID string    `json:"containerID"`
		ExitCode    int       `json:"exitCode"`
		FinishedAt  time.Time `json:"finishedAt"`
		Message     string    `json:"message"`
		Reason      string    `json:"reason"`
		Signal      int       `json:"signal"`
		StartedAt   time.Time `json:"startedAt"`
	}
	ContainerStateWaiting struct {
		Message string `json:"message"`
		Reason  string `json:"reason"`
	}
	ContainerState interface {
		Color() string
		GetStartTime() *time.Time
	}
)

func (c *ContainerStateRunning) Color() string {
	return emoji.GreenCircle.String()
}

func (c *ContainerStateRunning) GetStartTime() *time.Time {
	return &c.StartedAt
}

func (c *ContainerStateTerminated) Color() string {
	return emoji.BlackCircle.String()
}

func (c *ContainerStateTerminated) GetStartTime() *time.Time {
	return &c.StartedAt
}

func (c *ContainerStateWaiting) Color() string {
	return emoji.RedCircle.String()
}

func (c *ContainerStateWaiting) GetStartTime() *time.Time {
	return nil
}
