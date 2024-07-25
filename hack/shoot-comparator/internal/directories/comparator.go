package directories

import (
	"github.com/kyma-project/infrastructure-manager/tools/shoot-comparator/internal/files"
	"os"
	"path"
	"slices"
	"time"
)

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

	leftFileNames, err := getFileNames(leftDir, olderThan)
	if err != nil {
		return Result{}, err
	}

	rightFileNames, err := getFileNames(rightDir, olderThan)
	if err != nil {
		return Result{}, err
	}

	fileNamesToCompare := getIntersection(leftFileNames, rightFileNames)
	differences, err := compare(fileNamesToCompare, leftDir, rightDir)

	onlyLeftFiles := filterOut(leftFileNames, fileNamesToCompare)
	onlyRightFiles := filterOut(rightFileNames, fileNamesToCompare)

	return Result{
		LeftDir:            leftDir,
		RightDir:           rightDir,
		Diff:               differences,
		RightOnly:          onlyRightFiles,
		LeftOnly:           onlyLeftFiles,
		LeftDirFilesCount:  len(onlyRightFiles) + len(fileNamesToCompare),
		RightDirFilesCount: len(onlyLeftFiles) + len(fileNamesToCompare),
	}, nil
}

func getFileNames(dir string, olderThan time.Time) ([]string, error) {
	dirEntries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var fileNames []string

	for _, dirEntry := range dirEntries {
		if dirEntry.IsDir() {
			continue
		}

		fileInfo, err := dirEntry.Info()
		if err != nil {
			return nil, err
		}

		if fileInfo.ModTime().After(olderThan) {
			fileNames = append(fileNames, dirEntry.Name())
		}
	}

	slices.Sort(fileNames)

	return fileNames, nil
}

func getIntersection(leftFiles []string, rightFiles []string) []string {
	var intersection []string
	for _, leftFile := range leftFiles {
		_, found := slices.BinarySearch(rightFiles, leftFile)
		if found {
			intersection = append(intersection, leftFile)
		}
	}

	return intersection
}

func filterOut(fullFileList []string, filesToFilterOut []string) []string {
	var result []string
	for _, file := range fullFileList {
		_, found := slices.BinarySearch(filesToFilterOut, file)

		if !found {
			result = append(result, file)
		}
	}

	return result
}

func compare(filesNames []string, leftDir, rightDir string) ([]Difference, error) {
	var differences []Difference

	for _, fileName := range filesNames {
		leftFilePath := path.Join(leftDir, fileName)
		rightFilePath := path.Join(rightDir, fileName)

		equal, diffMessage, err := files.CompareFiles(leftFilePath, rightFilePath)
		if err != nil {
			return nil, err
		}

		if !equal {
			diff := Difference{
				Filename:  fileName,
				LeftFile:  leftFilePath,
				RightFile: rightFilePath,
				Message:   diffMessage,
			}
			differences = append(differences, diff)
		}
	}

	return differences, nil
}
