package filesyncstate

type State string

func (s State) String() string {
	return string(s)
}

const (
	Waiting   State = "waiting"
	Syncing   State = "syncing"
	Pausing   State = "pausing"
	Paused    State = "paused"
	Resuming  State = "resuming"
	Completed State = "completed"
	Failed    State = "failed"
)
