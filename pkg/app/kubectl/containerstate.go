package kubectl

import (
	"github.com/enescakir/emoji"
	"time"
)

type (
	ContainerStateRunning struct {
		StartedAt *time.Time `json:"startedAt"`
	}
	ContainerStateTerminated struct {
		ContainerID string     `json:"containerID"`
		ExitCode    int        `json:"exitCode"`
		FinishedAt  *time.Time `json:"finishedAt"`
		Message     string     `json:"message"`
		Reason      string     `json:"reason"`
		Signal      int        `json:"signal"`
		StartedAt   *time.Time `json:"startedAt"`
	}
	ContainerStateWaiting struct {
		Message string `json:"message"`
		Reason  string `json:"reason"`
	}
	ContainerStateEvent interface {
		GetColor() string
		GetReason() string
		GetTime() *time.Time
	}
)

func (c *ContainerStateRunning) GetReason() string {
	return "Running"
}

func (c *ContainerStateRunning) GetColor() string {
	return emoji.GreenCircle.String()
}

func (c *ContainerStateRunning) GetTime() *time.Time {
	return c.StartedAt
}

func (c *ContainerStateTerminated) GetReason() string {
	return c.Reason
}

func (c *ContainerStateTerminated) GetColor() string {
	return emoji.BlackCircle.String()
}

func (c *ContainerStateTerminated) GetTime() *time.Time {
	if c.FinishedAt != nil {
		return c.FinishedAt
	}
	return c.StartedAt
}

func (c *ContainerStateWaiting) GetReason() string {
	return c.Reason
}

func (c *ContainerStateWaiting) GetColor() string {
	return emoji.RedCircle.String()
}

func (c *ContainerStateWaiting) GetTime() *time.Time {
	return nil
}
