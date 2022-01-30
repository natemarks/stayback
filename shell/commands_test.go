package shell

import (
	"reflect"
	"testing"
)

func TestRunAndWait(t *testing.T) {
	type args struct {
		cmdRoot string
		cmdOpts []string
	}
	tests := []struct {
		name       string
		args       args
		wantResult RunAndWaitOutput
		wantErr    bool
	}{
		{name: "dsf", args: args{
			cmdRoot: "sh",
			cmdOpts: []string{"-c", "echo stdout; echo 1>&2 stderr"},
		}, wantErr: false, wantResult: RunAndWaitOutput{
			ExitCode:       0,
			CombinedOutput: "stdout\nstderr\n",
		}},
		{name: "sdf", args: args{
			cmdRoot: "sh",
			cmdOpts: []string{"-c", "echo stdout; echo 1>&2 stderr;exit 5"},
		}, wantErr: true, wantResult: RunAndWaitOutput{
			ExitCode:       5,
			CombinedOutput: "stdout\nstderr\n",
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult, err := RunAndWait(tt.args.cmdRoot, tt.args.cmdOpts)
			if (err != nil) != tt.wantErr {
				t.Errorf("RunAndWait() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("RunAndWait() gotResult = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}
