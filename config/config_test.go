package config

import (
	"reflect"
	"testing"
)

func Test_validateDir(t *testing.T) {
	type args struct {
		iDir string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "valid_absolute",
			args:    args{iDir: "/etc"},
			want:    "/etc",
			wantErr: false,
		},
		{
			name:    "valid_relative",
			args:    args{iDir: ".ssh"},
			want:    "/Users/nmarks/.ssh",
			wantErr: false,
		},
		{
			name:    "invalid_relative",
			args:    args{iDir: "aaa/bbb"},
			want:    "",
			wantErr: true,
		},
		{
			name:    "invalid_absolute",
			args:    args{iDir: "/opt/bbb"},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := validateDir(tt.args.iDir)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateDir() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("validateDir() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_sortUnique(t *testing.T) {
	type args struct {
		iList []string
	}
	tests := []struct {
		name      string
		args      args
		wantOList []string
		wantErr   bool
	}{
		{
			name: "valid",
			args: args{
				iList: []string{
					".ssh", "/opt/cisco", ".aws",
				},
			},
			wantErr: false,
			wantOList: []string{
				"/Users/nmarks/.aws", "/Users/nmarks/.ssh", "/opt/cisco",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotOList, err := sortUnique(tt.args.iList)
			if (err != nil) != tt.wantErr {
				t.Errorf("sortUnique() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotOList, tt.wantOList) {
				t.Errorf("sortUnique() gotOList = %v, want %v", gotOList, tt.wantOList)
			}
		})
	}
}
