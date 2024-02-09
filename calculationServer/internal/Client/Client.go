package Client

import "time"

type Task struct {
	Id         int
	Expression string
	ShowAlive  time.Duration
}

type Result struct {
	Id         int
	Expression string
	Logs       string
	Answer     int
}

type Client struct {
}

func (c *Client) CheckTask() (bool, *Task) {
	return false, &Task{}
}

func (c *Client) SendAlive() {

}

func (c *Client) SendResult(result Result) {

}
