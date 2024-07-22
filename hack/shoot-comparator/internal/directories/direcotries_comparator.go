package directories

type Result struct {
	Equal              bool
	LeftOnly           []string
	RightOnly          []string
	LeftDirFilesCount  int
	RightDirFilesCount int
	Diff               []Difference
}

type Difference struct {
	LeftFile  string
	RightFile string
	Message   string
}

func CompareDirectories(leftDir, rightDir string) (Result, error) {
	return Result{}, nil
}
