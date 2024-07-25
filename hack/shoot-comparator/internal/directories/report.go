package directories

import (
	"fmt"
	"os"
	"path"
	"time"
)

func SaveComparisonReport(comparisonResult Result, outputDir string, fromDate time.Time) (string, error) {
	resultsDir, err := createOutputDir(outputDir)
	if err != nil {
		return "", err
	}

	reportFile, err := createReportFile(resultsDir)
	if err != nil {
		return "", err
	}
	defer reportFile.Close()

	writeSummary(reportFile, comparisonResult, fromDate)

	if !comparisonResult.Equal {
		err := writeResultsToTheReportFile(reportFile, comparisonResult)
		if err != nil {
			return "", err
		}

		err = writeResultsToDiffFiles(comparisonResult.Diff, resultsDir)
		if err != nil {
			return "", fmt.Errorf("failed to write files with detected differences: %v", err)
		}
	}

	return resultsDir, nil
}

func writeSummary(reportFile *os.File, comparisonResult Result, fromDate time.Time) {
	reportFile.Write([]byte(fmt.Sprintf("Comparing files older than:%v \n", fromDate)))

	reportFile.Write([]byte("\n"))
	numberOfFilesLeftMsg := fmt.Sprintf("Number of files in %s directory = %d \n", comparisonResult.LeftDir, comparisonResult.LeftDirFilesCount)
	reportFile.Write([]byte(numberOfFilesLeftMsg))

	numberOfFilesRightMsg := fmt.Sprintf("Number of files in %s directory = %d \n", comparisonResult.RightDir, comparisonResult.RightDirFilesCount)

	reportFile.Write([]byte(numberOfFilesRightMsg))

	reportFile.Write([]byte("\n"))

	if comparisonResult.Equal {
		reportFile.Write([]byte("Directories are equal \n"))
	} else {
		reportFile.Write([]byte("Directories are NOT equal \n"))
	}
}

func createOutputDir(outputDir string) (string, error) {
	resultsDir := path.Join(outputDir, time.Now().Format(time.RFC3339))

	err := os.MkdirAll(resultsDir, os.ModePerm)
	if err != nil {
		return "", fmt.Errorf("failed to create results directory: %v", err)
	}

	return resultsDir, nil
}

func createReportFile(resultsDir string) (*os.File, error) {
	resultsFile := path.Join(resultsDir, "result.txt")

	file, err := os.Create(resultsFile)
	if err != nil {
		return nil, fmt.Errorf("failed to create results file: %v", err)
	}

	return file, nil
}

func writeResultsToTheReportFile(file *os.File, comparisonResult Result) error {
	err := writeMissingFilesToReport(file, comparisonResult.LeftDir, comparisonResult.LeftOnly)
	if err != nil {
		return fmt.Errorf("failed to write results to file: %v", err)
	}

	err = writeMissingFilesToReport(file, comparisonResult.RightDir, comparisonResult.RightOnly)
	if err != nil {
		return fmt.Errorf("failed to write results to file: %v", err)
	}

	err = writeDifferencesToReport(file, comparisonResult.Diff)
	if err != nil {
		return fmt.Errorf("failed to write results to file: %v", err)
	}

	return nil
}

func writeMissingFilesToReport(file *os.File, dir string, missingFiles []string) error {
	if len(missingFiles) == 0 {
		return nil
	}
	file.Write([]byte("\n"))

	file.Write([]byte(fmt.Sprintf("---------------------------------------------\n")))
	file.Write([]byte(fmt.Sprintf("Files existing in %s folder only: \n", dir)))

	for _, missingFile := range missingFiles {
		if _, err := file.Write([]byte(missingFile + "\n")); err != nil {
			return err
		}
	}

	file.Write([]byte(fmt.Sprintf("---------------------------------------------\n")))

	return nil
}

func writeDifferencesToReport(file *os.File, differences []Difference) error {
	if len(differences) == 0 {
		return nil
	}

	file.Write([]byte("\n"))

	file.Write([]byte(fmt.Sprintf("---------------------------------------------\n")))
	file.Write([]byte(fmt.Sprintf("Files that differ: \n")))

	for _, difference := range differences {
		msg := fmt.Sprintf("Files: %q and %q differ.", difference.LeftFile, difference.RightFile)

		if _, err := file.Write([]byte(msg + "\n")); err != nil {
			return err
		}
	}

	file.Write([]byte(fmt.Sprintf("---------------------------------------------\n")))

	return nil
}

func writeResultsToDiffFiles(differences []Difference, resultsDir string) error {
	writeAndCloseFunc := func(filePath string, text string) error {
		file, err := os.Create(filePath)
		if err != nil {
			return err
		}
		defer file.Close()

		_, err = file.Write([]byte(text))

		return err
	}

	for _, difference := range differences {
		diffFile := path.Join(resultsDir, fmt.Sprintf("%s.diff", difference.Filename))

		writeAndCloseFunc(diffFile, difference.Message)
	}

	return nil
}
