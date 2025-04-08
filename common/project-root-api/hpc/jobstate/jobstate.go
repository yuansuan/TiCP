package jobstate

type State string

func (s State) String() string {
	return string(s)
}

const (
	Unknown    State = "unknown"
	Preparing  State = "preparing"
	Pending    State = "pending"
	Running    State = "running"
	Completing State = "completing"
	Completed  State = "completed"
	Failed     State = "failed"
	Cancel     State = "cancel"
	Cancelling State = "cancelling"
	Canceled   State = "canceled"
)
