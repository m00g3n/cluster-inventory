package directories

import "time"

type Result struct {
	Equal              bool
	LeftDir            string
	RightDir           string
	LeftOnly           []string
	RightOnly          []string
	LeftDirFilesCount  int
	RightDirFilesCount int
	Diff               []Difference
}

type Difference struct {
	Filename  string
	LeftFile  string
	RightFile string
	Message   string
}

func CompareDirectories(leftDir, rightDir string, olderThan time.Time) (Result, error) {
	return Result{
		LeftDir:  leftDir,
		RightDir: rightDir,
	}, nil
}
