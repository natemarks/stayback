package backup

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"sort"
)

// Job is the backup job definition
type Job struct {
	Source          string   `json:"source"`          // ex. file, url, telepathy
	Id              string   `json:"id"`              // ex. 20220103-080910
	HomeDirectory   string   `json:"homeDirectory"`   // assumed root for relative backup target paths
	BackupDirectory string   `json:"backupDirectory"` // local backup location with a .tmp working ssub-dir
	S3Bucket        string   `json:"s3Bucket"`        // backup s3 bucket
	EncryptedDirs   []string `json:"encryptedDirs"`   // list of abs or relative dirs to back up wit encryption
	UnEncryptedDirs []string `json:"unEncryptedDirs"` // list of absolute or relative dirs to back up without encryption
}

// Report returns a config report as a list of strings. This makes it a little easier to target parts of the report for testing.
// It can be used to describe the job at run time and give the user a chance to cancel
// report[0] is a string that describes the source of the job definition (ex. $HOME/.stayback/default.json)
func (c Job) Report() (report []string, err error) {

	return report, err
}

// checkDirsExist given a list of directories, return a string/bool map to indicate whether each exists
// and is in fact a directory
func checkDirsExist(dirs []string) (result map[string]bool) {
	result = make(map[string]bool)
	for _, v := range dirs {
		fi, err := os.Stat(v)
		if err != nil {
			result[v] = false
			continue
		}
		result[v] = fi.IsDir()
	}
	return result

}

// TargetDirsExist Returns a merge map of the absolute directories and the boolean result of their existence
// We want to check every directory and report all problems, so they can all be solved at once
func (c Job) TargetDirsExist() map[string]bool {
	encryptedDirs := checkDirsExist(c.EncryptedDirs)
	unencryptedDirs := checkDirsExist(c.UnEncryptedDirs)
	// merge the maps and return
	for k, v := range encryptedDirs {
		unencryptedDirs[k] = v
	}
	return unencryptedDirs
}

// Execute returns testable output
func (c Job) Execute() (err error) {

	return err
}

// NewJobFromFile returns a JOb object from the json file path provided
func NewJobFromFile(fPath string) (result Job, err error) {
	if fPath == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return result, err
		}
		fPath = path.Join(homeDir, ".stayback")
	}

	file, err := ioutil.ReadFile(fPath)
	if err != nil {
		return result, err
	}

	err = json.Unmarshal(file, &result)
	return result, err
}

// makeAbsolute if the give dri doesn't begin with a path separator, join it onto the default root
func makeAbsolute(dir, defaultRoot string) string {
	var first rune
	for _, c := range dir {
		first = c
		break
	}
	if first == os.PathSeparator {
		return dir
	}
	return path.Join(defaultRoot, dir)
}

// cleanTargets converts each directory entry to an absolute path using a default root for relative directories
// then it sorts the list and removes the duplicates
func cleanTargets(tList []string, defaultRoot string) (oList []string, err error) {
	var absList []string

	// go through the targets and ensure each is an absolute path
	for _, v := range tList {
		abs := makeAbsolute(v, defaultRoot)
		absList = append(absList, abs)
	}

	// sort the targets so duplicates will all be grouped together
	sort.Strings(absList)

	// loop through the absList and only copy each entry from absList to oList if it
	// is different from the previous entry
	for i, v := range absList {
		// skip the predecessor comparison on the first item
		if i == 0 {
			oList = append(oList, v)
		} else {
			if v != absList[i-1] {
				oList = append(oList, v)
			}
		}
	}

	return oList, err
}
