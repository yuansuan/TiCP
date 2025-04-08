package consts

import (
	"fmt"
	"strings"
	"time"

	"github.com/yuansuan/ticp/common/go-kit/logging"
	"github.com/yuansuan/ticp/common/project-root-api/hpc/jobstate"
)

// 作业状态
// Initiated (0值)
// InitiallySuspended (作业提交时指定了暂停,作业状态为InitiallySuspended)
// Pending (中央调度器等待、数据上传、超算调度器等待都在这个状态,具体原因检查StateReason)
// Running (作业运行中)
// Suspending (暂停操作中)
// Suspended (已暂停,可通过resume操作恢复)
// // Transmitting (数据回传中) !废弃
// // TransmitSuspended (回传暂停) !废弃
// Terminating (中止操作中)
// Terminated (已中止,无法通过resume操作恢复)
// Completed (成功结束)
// Failed (包括系统失败和用户程序失败,具体原因可以检查StateReason)
const (
	Initiated          int = 0
	InitiallySuspended     = 1
	Pending                = 2
	Running                = 3
	Suspending             = 4
	Suspended              = 5
	Terminating            = 8
	Terminated             = 9
	Completed              = 10
	Failed                 = 11
	UNKNOWN                = 99
)

const (
	subStateInitiated int = Initiated * 100 // 0值

	subStateInitiallySuspended int = InitiallySuspended * 100

	subStateScheduling    int = Pending * 100
	subStateFileUploading int = subStateScheduling + 1
	subStateHpcWaiting    int = subStateFileUploading + 1

	subStateRunning     int = Running * 100
	subStateSuspending  int = Suspending * 100
	subStateSuspended   int = Suspended * 100
	subStateTerminating int = Terminating * 100
	subStateTerminated  int = Terminated * 100
	subStateCompleted   int = Completed * 100
	subStateFailed      int = Failed * 100
	subStateUnknown     int = 9999
)

var (
	// SubStateInitiated 初始状态(预留,暂时不用)
	SubStateInitiated = State{State: Initiated, SubState: subStateInitiated}

	// SubStateInitiallySuspended  初始暂停(提交作业时指定了暂停，进入该状态)
	SubStateInitiallySuspended = State{State: InitiallySuspended, SubState: subStateInitiallySuspended}

	// SubStateScheduling 作业在中央调度器等待
	SubStateScheduling = State{State: Pending, SubState: subStateScheduling}
	// SubStateFileUploading 作业在数据上传
	SubStateFileUploading = State{State: Pending, SubState: subStateFileUploading}
	// SubStateHpcWaiting 作业在超算调度器等待
	SubStateHpcWaiting = State{State: Pending, SubState: subStateHpcWaiting}

	// SubStateRunning 作业运行中
	SubStateRunning = State{State: Running, SubState: subStateRunning}

	// SubStateSuspending 暂停操作中
	SubStateSuspending = State{State: Suspending, SubState: subStateSuspending}

	// SubStateSuspended 已暂停,可通过resume操作恢复
	SubStateSuspended = State{State: Suspended, SubState: subStateSuspended}

	// SubStateTerminating 中止操作中
	SubStateTerminating = State{State: Terminating, SubState: subStateTerminating}

	// SubStateTerminated 已中止,无法通过resume操作恢复 [终态]
	SubStateTerminated = State{State: Terminated, SubState: subStateTerminated}

	// SubStateCompleted 成功结束 [终态]
	SubStateCompleted = State{State: Completed, SubState: subStateCompleted}

	// SubStateFailed 包括系统失败和用户程序失败,具体原因可以检查StateReason [终态]
	SubStateFailed = State{State: Failed, SubState: subStateFailed}

	// SubStateUnknown 未知状态
	SubStateUnknown = State{State: UNKNOWN, SubState: subStateUnknown}
)

// State 作业状态
type State struct {
	State    int `json:"state"`
	SubState int `json:"sub_state"`
}

// NewState 新建作业状态
func NewState(state, subState int) State {
	return State{State: state, SubState: subState}
}

var allStates = []int{}
var allSubStates = []State{}

var mapStateString = map[int]string{} // map[state(int)]stateString
var mapSubState = map[int]State{}     // map[subState]

func addState(state State) {
	allSubStates = append(allSubStates, state)
	mapSubState[state.SubState] = state

	if _, ok := mapStateString[state.State]; !ok {
		allStates = append(allStates, state.State)
		mapStateString[state.State] = state.StateString()
	}

}

func initState() {
	addState(SubStateInitiated)
	addState(SubStateInitiallySuspended)
	addState(SubStateScheduling)
	addState(SubStateFileUploading)
	addState(SubStateHpcWaiting)
	addState(SubStateRunning)
	addState(SubStateSuspending)
	addState(SubStateSuspended)
	addState(SubStateTerminating)
	addState(SubStateTerminated)
	addState(SubStateCompleted)
	addState(SubStateFailed)
	addState(SubStateUnknown)
}

func init() {
	initState()
}

// GetStateBySubState ...
func GetStateBySubState(subState int) (State, bool) {
	s, ok := mapSubState[subState]
	return s, ok
}

// ValidStringState ...
func ValidStringState(stringState string) bool {
	for _, state := range allSubStates {
		if stringState == state.StateString() {
			return true
		}
	}
	return false
}

// GetStateValue ...
func GetStateValue(stringState string) (int, bool) {
	for _, state := range allSubStates {
		if stringState == state.StateString() {
			return state.State, true
		}
	}
	return UNKNOWN, false
}

// AllStateString return all state string
func AllStateString() string {
	var allStateString []string
	for _, state := range allStates {
		allStateString = append(allStateString, mapStateString[state])
	}

	return fmt.Sprintf("[%s]", strings.Join(allStateString, ", "))
}

// StateString 主状态字符串
func (s State) StateString() string {
	switch s.State {
	case Initiated:
		return "Initiated"
	case InitiallySuspended:
		return "InitiallySuspended"
	case Pending:
		return "Pending"
	case Running:
		return "Running"
	case Suspending:
		return "Suspending"
	case Suspended:
		return "Suspended"
	case Terminating:
		return "Terminating"
	case Terminated:
		return "Terminated"
	case Completed:
		return "Completed"
	case Failed:
		return "Failed"
	case UNKNOWN:
		return "UNKNOWN"
	default:
		return "UNKNOWN"
	}
}

// String 转换为字符串
func (s State) String() string {
	switch s.State {
	case Initiated:
		return "Initiated"
	case InitiallySuspended:
		return "InitiallySuspended"
	case Pending:
		switch s.SubState {
		case subStateScheduling:
			return "Pending[Scheduling]"
		case subStateFileUploading:
			return "Pending[FileUploading]"
		case subStateHpcWaiting:
			return "Pending[HpcWaiting]"
		default:
			return "Pending"
		}
	case Running:
		return "Running"
	case Suspending:
		return "Suspending"
	case Suspended:
		return "Suspended"
	case Terminating:
		return "Terminating"
	case Terminated:
		return "Terminated"
	case Completed:
		return "Completed"
	case Failed:
		return "Failed"
	case UNKNOWN:
		return "Unknown"
	default:
		return "Unknown"
	}
}

// ParseStateString 解析字符串
func ParseStateString(s string) (State, error) {
	state := State{}

	switch s {
	case "Initiated":
		state = SubStateInitiated
	case "InitiallySuspended":
		state = SubStateInitiallySuspended
	case "Pending":
		state = SubStateScheduling
	case "Pending[Scheduling]":
		state = SubStateScheduling
	case "Pending[FileUploading]":
		state = SubStateFileUploading
	case "Pending[HpcWaiting]":
		state = SubStateHpcWaiting
	case "Running":
		state = SubStateRunning
	case "Suspending":
		state = SubStateSuspending
	case "Suspended":
		state = SubStateSuspended
	case "Terminating":
		state = SubStateTerminating
	case "Terminated":
		state = SubStateTerminated
	case "Completed":
		state = SubStateCompleted
	case "Failed":
		state = SubStateFailed
	case "Unknown":
		state = SubStateUnknown
	default:
		return state, fmt.Errorf("invalid State string: %s", s)
	}

	return state, nil
}

// IsFinal 是否是终态
func (s State) IsFinal() bool {
	switch s.State {
	case Terminated, Completed, Failed:
		return true
	default:
		return false
	}
}

// ToFinal 判断是否从非终态转到终态
func (s State) ToFinal(newState State) bool {
	if s.IsFinal() {
		return false
	}

	if newState.IsFinal() {
		return true
	}

	return false
}

// FinalString func return all final state string
func FinalString() []string {
	return []string{
		SubStateTerminated.String(),
		SubStateCompleted.String(),
		SubStateFailed.String(),
	}
}

// CanTerminate 能否终止
func (s State) CanTerminate() bool {
	switch s.State {
	case Initiated, InitiallySuspended, Pending, Running, Suspended, Terminating:
		return true
	default:
		return false
	}
}

// CanResume 能否恢复(根据主状态)
func (s State) CanResume() bool {
	switch s.State {
	case Suspended, InitiallySuspended:
		return true
	default:
		return false
	}
}

// IsInitiated 是否是Initiated状态
func (s State) IsInitiated() bool {
	return s.State == Initiated
}

// IsPending 是否是Pending状态
func (s State) IsPending() bool {
	return s.State == Pending
}

// IsRunning 是否是Running状态
func (s State) IsRunning() bool {
	return s.State == Running
}

// IsTerminating 是否是Terminating状态
func (s State) IsTerminating() bool {
	return s.State == Terminating
}

// IsCompleted 是否是Completed状态
func (s State) IsCompleted() bool {
	return s.State == Completed
}

// IsInitiallySuspended 是否是InitiallySuspended状态
func (s State) IsInitiallySuspended() bool {
	return s.SubState == subStateInitiallySuspended
}

// IsScheduling 是否是Scheduling状态
func (s State) IsScheduling() bool {
	return s.SubState == subStateScheduling
}

// Diff 比较两个状态是否不同
func (s State) Diff(s2 *State) bool {
	return s.State != s2.State || s.SubState != s2.SubState
}

// StateReason 状态原因
type StateReason struct {
	State   State        // 当前状态
	Reason  string       // 当前状态原因
	Time    time.Time    // 时间
	history *StateReason // 状态历史
}

// NewStateReason 新建状态原因
func NewStateReason(state State, reason string, ts ...time.Time) *StateReason {
	if len(ts) > 1 {
		panic("too many time")
	}

	t := time.Now()
	if len(ts) == 1 {
		t = ts[0]
	}

	return &StateReason{
		State:   state,
		Reason:  reason,
		Time:    t,
		history: nil,
	}
}

// Update 更新状态原因
func (sr *StateReason) Update(state State, reason string, ts ...time.Time) {
	if len(ts) > 1 {
		panic("too many time")
	}

	t := time.Now()
	if len(ts) == 1 {
		t = ts[0]
	}

	sr.history = &StateReason{
		State:   sr.State,
		Reason:  sr.Reason,
		Time:    sr.Time,
		history: sr.history,
	}
	sr.State = state
	sr.Reason = reason
	sr.Time = t
}

// String 递归打印所有状态原因
func (sr *StateReason) String() string {
	if sr == nil {
		return ""
	}

	historyString := ""
	if sr.history != nil {
		historyString = "\n" + sr.history.String()
	}

	return fmt.Sprintf("[%s]%s:%s%s", sr.Time.Format(time.RFC3339), sr.State.String(), sr.Reason, historyString)
}

// Diff 比较两个状态原因
func (sr *StateReason) Diff(state State, reason string) bool {
	if sr == nil {
		return false
	}

	if sr.State.Diff(&state) {
		return true
	}

	if sr.Reason != reason {
		return true
	}

	return false
}

// ParseStateReasonString 解析状态原因字符串
func ParseStateReasonString(s string) (*StateReason, error) {
	if s == "" {
		return nil, nil
	}

	lines := strings.Split(s, "\n")
	var sr *StateReason

	for i := len(lines) - 1; i >= 0; i-- {
		line := strings.TrimSpace(lines[i])

		parts := strings.SplitN(line, "]", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid StateReason string: %s", s)
		}

		timeString := strings.TrimPrefix(parts[0], "[")
		stateReasonString := strings.TrimSpace(parts[1])

		t, err := time.Parse(time.RFC3339, timeString)
		if err != nil {
			return nil, fmt.Errorf("invalid time format in StateReason string: %s", s)
		}

		stateReasonParts := strings.SplitN(stateReasonString, ":", 2)
		if len(stateReasonParts) != 2 {
			return nil, fmt.Errorf("invalid StateReason string: %s", s)
		}

		stateString := strings.TrimSpace(stateReasonParts[0])
		reason := strings.TrimSpace(stateReasonParts[1])

		state, err := ParseStateString(stateString)
		if err != nil {
			return nil, err
		}

		if sr == nil {
			sr = NewStateReason(state, reason, t)
		} else {
			newSR := NewStateReason(state, reason, t)
			newSR.history = sr
			sr = newSR
		}
	}

	return sr, nil
}

// ParseAndUpdateStateReasonString 解析并且更新
func ParseAndUpdateStateReasonString(oldReason string, state State, reason string) *StateReason {
	sr, err := ParseStateReasonString(oldReason)
	if err != nil {
		logging.Default().Infof("ParseStateReasonString failed,s:%s,err:%v", oldReason, err)
		return NewStateReason(state, reason)
	}

	if sr == nil {
		return NewStateReason(state, reason)
	}

	if sr.Diff(state, reason) {
		sr.Update(state, reason)
	}

	return sr
}

// 状态原因const
const (
	InitiallySuspendedReason   = "Job is initiallySuspended"
	UserSubmitReason           = "User Submit"
	SchedulingReason           = "Job Scheduling"
	FileUploadingReason        = ""
	HpcWaitingReason           = ""
	StateReasonUserCancel      = "User Cancel"
	StateReasonAdminCancel     = "Admin Cancel"
	StateReasonTransmitSuspend = "Transmit Suspend"
	StateReasonTransmitResume  = "Transmit Resume"
	StateReasonNonZeroExitCode = "NonZeroExitCode"
)

// FileSyncState 文件同步状态
type FileSyncState string

// 文件同步状态
const (
	FileSyncStateWaiting   FileSyncState = "Waiting"
	FileSyncStateSyncing   FileSyncState = "Syncing"
	FileSyncStatePausing   FileSyncState = "Pausing"
	FileSyncStatePaused    FileSyncState = "Paused"
	FileSyncStateResuming  FileSyncState = "Resuming"
	FileSyncStateCompleted FileSyncState = "Completed"
	FileSyncStateFailed    FileSyncState = "Failed"
	FileSyncStateUnknown   FileSyncState = "Unknown"
	FileSyncStateNone      FileSyncState = ""
)

// IsValid 方法用于验证输入字符串是否为有效的文件同步状态
func (f FileSyncState) IsValid() bool {
	switch f {
	case FileSyncStateWaiting, FileSyncStateSyncing, FileSyncStatePausing,
		FileSyncStatePaused, FileSyncStateResuming, FileSyncStateCompleted,
		FileSyncStateFailed, FileSyncStateUnknown, FileSyncStateNone:
		return true
	default:
		return false
	}
}

// String 转换为字符串
func (f FileSyncState) String() string {
	return string(f)
}

// IsFinal 是否是终态
func (f FileSyncState) IsFinal() bool {
	switch f {
	case FileSyncStateCompleted, FileSyncStateFailed:
		return true
	default:
		return false
	}
}

func (f FileSyncState) IsNeedFileSync() bool {
	switch f {
	case FileSyncStateSyncing, FileSyncStatePausing, FileSyncStateResuming, FileSyncStateWaiting, FileSyncStateNone:
		return true
	default:
		return false
	}
}

// CanDelete 能否删除
func (f FileSyncState) CanDelete() bool {
	return f.IsFinal() || f == FileSyncStatePaused
}

// ToFinal 判断是否从非终态转到终态
func (f FileSyncState) ToFinal(newState FileSyncState) bool {
	if f.IsFinal() {
		return false
	}

	if newState.IsFinal() {
		return true
	}

	return false
}

// CanTransmitSuspended 能否传输暂停
func (f FileSyncState) CanTransmitSuspended() bool {
	if f == FileSyncStateSyncing || f == FileSyncStatePausing {
		return true
	}
	return false
}

// CanTransmitResume 能否恢复
func (f FileSyncState) CanTransmitResume() bool {
	if f == FileSyncStatePaused || f == FileSyncStateResuming {
		return true
	}
	return false
}

// ConvertHpcStateToYsState ...
func ConvertHpcStateToYsState(hpcState jobstate.State) State {
	switch hpcState {
	case jobstate.Preparing:
		return SubStateFileUploading

	case jobstate.Pending:
		return SubStateHpcWaiting

	case jobstate.Running:
		return SubStateRunning

	case jobstate.Completing, jobstate.Completed:
		return SubStateCompleted

	case jobstate.Failed:
		return SubStateFailed

	case jobstate.Cancel, jobstate.Cancelling:
		return SubStateTerminating

	case jobstate.Canceled:
		return SubStateTerminated

	default:
		return SubStateUnknown
	}
}
