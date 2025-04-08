package consts

import (
	"reflect"
	"testing"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
)

func TestNewStateReason(t *testing.T) {
	sr := NewStateReason(SubStateScheduling, "Scheduling")
	t.Log(sr.String())

	time.Sleep(time.Second)

	sr.Update(SubStateFileUploading, "FileUploading")
	t.Log(sr.String())

	time.Sleep(time.Second)

	sr.Update(SubStateHpcWaiting, "HpcWaiting")
	t.Log(sr.String())

	time.Sleep(time.Second)

	sr.Update(SubStateRunning, "Running")
	t.Log(sr.String())

	time.Sleep(time.Second)

	sr2, err := ParseStateReasonString(sr.String())
	if err != nil {
		t.Fatal(err)
	}

	t.Log(sr2.String())

	// s := `Failed:fsm:input:download: download all inputs failed, fulfill download map failed, call storage client file stat failed, endpoint: http://10.0.4.55:8899, path: /ys_id/input/4RBveK3ibGy/12342891-9c6b-4ad9-96db-6bb212345678/Blade.sim, err : call api stat failed, GetQuota "http://10.0.4.55:8899/api/storage/stat?AppKey=fkF9FPujdbqVvciEfnQElYP7ef0zkstncn0c&Path=%2Fys_id%2Finput%2F4RBveK3ibGy%2F12342891-9c6b-4ad9-96db-6bb212345678%2FBlade.sim&Signature=7566afacf26625f6092972a512ff8b2f&Timestamp=1686040320": dial tcp 10.0.4.55:8899: connect: connection refused | Pending[Scheduling]:Pending:Job Scheduling`
	// sr2, err := ParseStateReasonString(s)
	// if err != nil {
	// 	t.Fatal(err)
	// }

	// t.Log(sr2.State)
	// t.Log(sr2.Reason)

	// a := "Completed:Completed | Pending[HpcWaiting]:fsm:running:watch: fsm: job canceled"
	// sr2, err := ParseStateReasonString(a)
	// // sr2, err := ParseStateReasonString("Pending[HpcWaiting]:")
	// if err != nil {
	// 	t.Fatal(err)
	// }

	// assert.Equal(t, a, sr2.String())

	// spew.Dump(sr2)

	// t.Log(sr2.State)
	// t.Log(sr2.Reason)

	// t.Log(sr2.String())

	// sr2.Update(SubStateCompleted, "Completed")
	// t.Log(sr2.String())

	// sr3, err := ParseStateReasonString(sr2.String())
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// t.Log(sr3.State)
	// t.Log(sr3.Reason)

	// sr3.Update(SubStateFailed, "Failed")
	// t.Log(sr3.String())
}

func MustParse(s string) time.Time {
	t, _ := time.Parse(time.RFC3339, s)
	return t
}

func TestParseStateReasonString(t *testing.T) {

	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    *StateReason
		wantErr bool
	}{
		{
			name: "normal",
			args: args{"[2023-06-07T16:03:07+08:00]Pending[Scheduling]:Scheduling"},
			want: &StateReason{
				State:  SubStateScheduling,
				Reason: "Scheduling",
				Time:   MustParse("2023-06-07T16:03:07+08:00"),
			},
			wantErr: false,
		},
		{
			name: "normal2",
			args: args{"[2023-06-07T16:03:08+08:00]Pending[FileUploading]:FileUploading\n[2023-06-07T16:03:07+08:00]Pending[Scheduling]:Scheduling"},
			want: &StateReason{
				State:  SubStateFileUploading,
				Reason: "FileUploading",
				history: &StateReason{
					State:  SubStateScheduling,
					Reason: "Scheduling",
					Time:   MustParse("2023-06-07T16:03:07+08:00"),
				},
				Time: MustParse("2023-06-07T16:03:08+08:00"),
			},
			wantErr: false,
		},
		{
			name: "normal3",
			args: args{"[2023-06-07T16:03:09+08:00]Pending[HpcWaiting]:HpcWaiting\n[2023-06-07T16:03:08+08:00]Pending[FileUploading]:FileUploading\n[2023-06-07T16:03:07+08:00]Pending[Scheduling]:Scheduling"},
			want: &StateReason{
				State:  SubStateHpcWaiting,
				Reason: "HpcWaiting",
				history: &StateReason{
					State:  SubStateFileUploading,
					Reason: "FileUploading",
					history: &StateReason{
						State:  SubStateScheduling,
						Reason: "Scheduling",
						Time:   MustParse("2023-06-07T16:03:07+08:00"),
					},
					Time: MustParse("2023-06-07T16:03:08+08:00"),
				},
				Time: MustParse("2023-06-07T16:03:09+08:00"),
			},
			wantErr: false,
		},
		{
			name: "normal-long",
			args: args{`[2023-06-07T16:03:09+08:00]Failed:fsm:input:download: download all inputs failed, fulfill download map failed, call storage client file stat failed, endpoint: http://10.0.4.55:8899, path: /ys_id/input/4RBveK3ibGy/12342891-9c6b-4ad9-96db-6bb212345678/Blade.sim, err : call api stat failed, GetQuota "http://10.0.4.55:8899/api/storage/stat?AppKey=fkF9FPujdbqVvciEfnQElYP7ef0zkstncn0c&Path=%2Fys_id%2Finput%2F4RBveK3ibGy%2F12342891-9c6b-4ad9-96db-6bb212345678%2FBlade.sim&Signature=7566afacf26625f6092972a512ff8b2f&Timestamp=1686040320": dial tcp 10.0.4.55:8899: connect: connection refused
[2023-06-07T16:03:08+08:00]Pending[Scheduling]:Pending:Job Scheduling`},
			want: &StateReason{
				State:  SubStateFailed,
				Reason: `fsm:input:download: download all inputs failed, fulfill download map failed, call storage client file stat failed, endpoint: http://10.0.4.55:8899, path: /ys_id/input/4RBveK3ibGy/12342891-9c6b-4ad9-96db-6bb212345678/Blade.sim, err : call api stat failed, GetQuota "http://10.0.4.55:8899/api/storage/stat?AppKey=fkF9FPujdbqVvciEfnQElYP7ef0zkstncn0c&Path=%2Fys_id%2Finput%2F4RBveK3ibGy%2F12342891-9c6b-4ad9-96db-6bb212345678%2FBlade.sim&Signature=7566afacf26625f6092972a512ff8b2f&Timestamp=1686040320": dial tcp 10.0.4.55:8899: connect: connection refused`,
				history: &StateReason{
					State:   SubStateScheduling,
					Reason:  "Pending:Job Scheduling",
					history: nil,
					Time:    MustParse("2023-06-07T16:03:08+08:00"),
				},
				Time: MustParse("2023-06-07T16:03:09+08:00"),
			},
			wantErr: false,
		},
		{
			name: "other1",
			args: args{"[2023-06-07T16:03:09+08:00]Pending[HpcWaiting]:fsm:running:watch: fsm: job canceled"},
			want: &StateReason{
				State:  SubStateHpcWaiting,
				Reason: "fsm:running:watch: fsm: job canceled",
				Time:   MustParse("2023-06-07T16:03:09+08:00"),
			},
			wantErr: false,
		},
		{
			name: "other2",
			args: args{"[2023-06-07T16:03:09+08:00]Completed:Completed\n[2023-06-07T16:03:08+08:00]Pending[HpcWaiting]:fsm:running:watch: fsm: job canceled"},
			want: &StateReason{
				State:  SubStateCompleted,
				Reason: "Completed",
				history: &StateReason{
					State:  SubStateHpcWaiting,
					Reason: "fsm:running:watch: fsm: job canceled",
					Time:   MustParse("2023-06-07T16:03:08+08:00"),
				},
				Time: MustParse("2023-06-07T16:03:09+08:00"),
			},
			wantErr: false,
		},
		{
			name:    "error: empty second line",
			args:    args{"[2023-06-07T16:03:09+08:00]Completed:Completed\n"},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "error: wrong state",
			args:    args{"[2023-06-07T16:03:09+08:00]FFF:fff"},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "error: wrong second line state",
			args:    args{"[2023-06-07T16:03:09+08:00]Completed:Completed\n[2023-06-07T16:03:09+08:00]AAA"},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "error: empty time",
			args:    args{"Completed:Completed"},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "error: wrong time",
			args:    args{"[wrong time]Completed:Completed"},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "error:empty",
			args:    args{""},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseStateReasonString(tt.args.s)

			if tt.wantErr {
				if assert.Error(t, err) {
					t.Log(err.Error())
				}
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want.String(), got.String())
			t.Log(got.String())
		})
	}
}

func TestGetStateValue(t *testing.T) {
	type args struct {
		stringState string
	}
	tests := []struct {
		name  string
		args  args
		want  int
		want1 bool
	}{
		{
			name:  "normal",
			args:  args{"Pending"},
			want:  Pending,
			want1: true,
		},
		{
			name:  "error",
			args:  args{"FFF"},
			want:  0,
			want1: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := GetStateValue(tt.args.stringState)
			assert.Equal(t, tt.want1, got1)

			if got1 {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestStateReason_Diff(t *testing.T) {
	type fields struct {
		State   State
		Reason  string
		history *StateReason
	}
	type args struct {
		state  State
		reason string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name:   "normal",
			fields: fields{State: SubStateCompleted, Reason: "Completed", history: &StateReason{State: SubStateHpcWaiting, Reason: "waiting schedule", history: nil}},
			args:   args{state: SubStateRunning, reason: "waiting for the result of the job"},
			want:   true,
		},
		{
			name:   "diff",
			fields: fields{State: SubStateRunning, Reason: "waiting for the result of the job", history: &StateReason{State: SubStateCompleted, Reason: "Completed", history: &StateReason{State: SubStateHpcWaiting, Reason: "waiting schedule", history: nil}}},
			args:   args{state: SubStateRunning, reason: "waiting for the result of the job"},
			want:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sr := &StateReason{
				State:   tt.fields.State,
				Reason:  tt.fields.Reason,
				history: tt.fields.history,
			}
			if got := sr.Diff(tt.args.state, tt.args.reason); got != tt.want {
				t.Errorf("StateReason.Diff() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInit(t *testing.T) {
	spew.Dump(allStates)
	spew.Dump(allSubStates)
	spew.Dump(mapSubState)
	spew.Dump(mapStateString)
	spew.Dump(AllStateString())
}

func TestGetStateBySubState(t *testing.T) {
	type args struct {
		subState int
	}
	tests := []struct {
		name  string
		args  args
		want  State
		want1 bool
	}{
		{
			name: "normal",
			args: args{
				subState: subStateInitiallySuspended,
			},
			want:  SubStateInitiallySuspended,
			want1: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := GetStateBySubState(tt.args.subState)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetStateBySubState() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("GetStateBySubState() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
