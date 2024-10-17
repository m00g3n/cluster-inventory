package comparator

import (
	"fmt"
	"log"
	"os"
	"path"
)

type Report struct {
	reportFile *os.File
	contents   string
}

func SaveComparisonReport(comparisonResult Result, outputDir string, shootName string) (string, error) {
	resultsDir, err := createOutputDir(outputDir, shootName)
	if err != nil {
		return "", err
	}

	report, err := newReport(resultsDir)
	if err != nil {
		return "", err
	}
	defer func() {
		if err := report.Close(); err != nil {
			log.Printf("Failed to close report file: %q", err)
		}
	}()

	writeSummary(&report, comparisonResult, shootName)

	if !comparisonResult.Equal {
		err = writeDifferencesToReport(&report, comparisonResult.Diff)
		if err != nil {
			return "", fmt.Errorf("failed to write results to file: %v", err)
		}

		err = writeResultsToDiffFiles(comparisonResult.Diff, resultsDir)
		if err != nil {
			return "", fmt.Errorf("failed to write files with detected differences: %v", err)
		}
	}

	return resultsDir, report.Save()
}

func newReport(resultsDir string) (Report, error) {
	resultsFile := path.Join(resultsDir, "results.txt")

	file, err := os.Create(resultsFile)
	if err != nil {
		return Report{}, fmt.Errorf("failed to create results file: %v", err)
	}

	return Report{
		reportFile: file,
	}, nil
}

func createOutputDir(outputDir, shootName string) (string, error) {
	resultsDir := path.Join(outputDir, shootName)

	err := os.MkdirAll(resultsDir, os.ModePerm)
	if err != nil {
		return "", fmt.Errorf("failed to create results directory: %v", err)
	}

	return resultsDir, nil
}

func (rw *Report) Save() error {
	_, err := rw.reportFile.Write([]byte(rw.contents))

	return err
}

func (rw *Report) Close() error {
	if rw.reportFile != nil {
		return rw.reportFile.Close()
	}

	return nil
}

func (rw *Report) AddLine(line string) {
	rw.contents += line + "\n"
}

func writeSummary(report *Report, comparisonResult Result, shootName string) {
	report.AddLine(fmt.Sprintf("Comparing generated Shoot with Shoot from Gardener, name: %s", shootName))
	report.AddLine("")

	if comparisonResult.Equal {
		report.AddLine("No differences found.")
	} else {
		report.AddLine("Differences found.")
	}
}

func writeDifferencesToReport(report *Report, differences []Difference) error {
	if len(differences) == 0 {
		return nil
	}

	report.AddLine("")
	report.AddLine("------------------------------------------------------------------------------------------")
	report.AddLine("Shoot that differ: ")

	for _, difference := range differences {
		msg := fmt.Sprintf("-%s", difference.ShootName)

		report.AddLine(msg)
	}

	report.AddLine("------------------------------------------------------------------------------------------")

	return nil
}

func writeResultsToDiffFiles(differences []Difference, resultsDir string) error {
	writeAndCloseFunc := func(filePath string, text string) error {
		file, err := os.Create(filePath)
		if err != nil {
			return err
		}
		defer func() {
			if file != nil {
				err := file.Close()
				if err != nil {
					fmt.Printf("failed to close file: %v", err)
				}
			}
		}()

		_, err = file.Write([]byte(text))

		return err
	}

	for _, difference := range differences {
		diffFile := path.Join(resultsDir, fmt.Sprintf("%s.diff", difference.ShootName))

		err := writeAndCloseFunc(diffFile, difference.Message)
		if err != nil {
			return err
		}
	}

	return nil
}
