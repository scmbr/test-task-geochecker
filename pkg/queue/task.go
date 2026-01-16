package queue

type Task struct {
	TaskID    string
	Payload   string
	TargetURL string
	Attempts  uint
}
