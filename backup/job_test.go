package backup

import (
	"reflect"
	"testing"
)

// Test_checkDirsExist make sure the checkDirsExist returns a map correctly
func Test_checkDirsExist(t *testing.T) {
	realDir := t.TempDir()
	type args struct {
		dirs []string
	}
	tests := []struct {
		name       string
		args       args
		wantResult map[string]bool
	}{
		{
			name: "valid",
			args: args{dirs: []string{"/aa/bb/cc", realDir}},
			wantResult: map[string]bool{
				"/aa/bb/cc": false,
				realDir:     true,
			}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotResult := checkDirsExist(tt.args.dirs); !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("checkDirsExist() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func Test_cleanTargets(t *testing.T) {
	type args struct {
		tList       []string
		defaultRoot string
	}
	tests := []struct {
		name      string
		args      args
		wantOList []string
		wantErr   bool
	}{
		// give a list that's already absolute and sorted with no duplicates
		{
			name: "absolute and sorted",
			args: args{
				tList: []string{
					"/aa/my/dir/.hidden",
					"/zz/my/dir/.hidden",
				},
				defaultRoot: "/aaa/bbb/ccc",
			},
			wantErr: false,
			wantOList: []string{
				"/aa/my/dir/.hidden",
				"/zz/my/dir/.hidden",
			}},
		// give a list that's already absolute but unsorted with no duplicates
		{
			name: "absolute but unsorted",
			args: args{
				tList: []string{
					"/zz/my/dir/.hidden",
					"/aa/my/dir/.hidden",
				},
				defaultRoot: "/aaa/bbb/ccc",
			},
			wantErr: false,
			wantOList: []string{
				"/aa/my/dir/.hidden",
				"/zz/my/dir/.hidden",
			}},
		// give a list that's mixed relative and absolute and also unsorted with no duplicates
		{
			name: "mixed and unsorted",
			args: args{
				tList: []string{
					"/aa/my/dir/.hidden",
					"zz/my/dir/.hidden",
					"ccc/zz/my/dir/.hidden",
					"00/zz/my/dir/.hidden",
				},
				defaultRoot: "/aaa/bbb/ccc",
			},
			wantErr: false,
			wantOList: []string{
				"/aa/my/dir/.hidden",
				"/aaa/bbb/ccc/00/zz/my/dir/.hidden",
				"/aaa/bbb/ccc/ccc/zz/my/dir/.hidden",
				"/aaa/bbb/ccc/zz/my/dir/.hidden",
			}},
		// give a list that's mixed relative and absolute and also unsorted with  duplicates
		{
			name: "mixed and unsorted with duplicates",
			args: args{
				tList: []string{
					"/aa/my/dir/.hidden",
					"zz/my/dir/.hidden",
					"/aaa/bbb/ccc/zz/my/dir/.hidden",
					"ccc/zz/my/dir/.hidden",
					"00/zz/my/dir/.hidden",
				},
				defaultRoot: "/aaa/bbb/ccc",
			},
			wantErr: false,
			wantOList: []string{
				"/aa/my/dir/.hidden",
				"/aaa/bbb/ccc/00/zz/my/dir/.hidden",
				"/aaa/bbb/ccc/ccc/zz/my/dir/.hidden",
				"/aaa/bbb/ccc/zz/my/dir/.hidden",
			}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotOList, err := cleanTargets(tt.args.tList, tt.args.defaultRoot)
			if (err != nil) != tt.wantErr {
				t.Errorf("cleanTargets() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotOList, tt.wantOList) {
				t.Errorf("cleanTargets() gotOList = %v, want %v", gotOList, tt.wantOList)
			}
		})
	}
}

func Test_makeAbsolute(t *testing.T) {
	type args struct {
		dir         string
		defaultRoot string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// give absolute path. get same path
		{
			name: "absolute",
			args: args{
				dir:         "/aa/bb/cc",
				defaultRoot: "/dd/gg/hh/",
			}, want: "/aa/bb/cc"},
		// give relative path, get joined path
		{
			name: "absolute",
			args: args{
				dir:         "aa/bb/cc",
				defaultRoot: "/dd/gg/hh/",
			}, want: "/dd/gg/hh/aa/bb/cc"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := makeAbsolute(tt.args.dir, tt.args.defaultRoot); got != tt.want {
				t.Errorf("makeAbsolute() = %v, want %v", got, tt.want)
			}
		})
	}
}
