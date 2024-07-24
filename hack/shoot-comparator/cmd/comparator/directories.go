package comparator

import (
	"fmt"
	"github.com/kyma-project/infrastructure-manager/tools/shoot-comparator/internal/directories"
	"github.com/spf13/cobra"
	"time"
)

func init() {
	rootCmd.AddCommand(directoriesCmd)
	directoriesCmd.Flags().String("fromDate", "", "Files older than specified date will not be compared")
	directoriesCmd.Flags().String("outputDir", "", "Directory storing comparison results")
}

var directoriesCmd = &cobra.Command{
	Use:     "dirs",
	Aliases: []string{"d"},
	Short:   "Compare directories",
	Args:    cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		leftDir := args[0]
		rightDir := args[1]

		fromDateString, err := cmd.Flags().GetString("fromDate")
		if err != nil {
			fmt.Printf("Error occurred when parsing command line arguments:%q", err)
			return
		}

		fromDate, err := parseStartFromDate(fromDateString)
		if err != nil {
			fmt.Printf("Error occurred when parsing command line arguments: %q", err)
			return
		}

		outputDir, err := cmd.Flags().GetString("outputDir")
		if err != nil {
			fmt.Printf("Error occurred when parsing command line arguments: %q", err)
			return
		}

		if fromDate.IsZero() {
			fmt.Printf("Only files created after the following date: %v will be compared.\n", fromDate)
		}

		fmt.Printf("Comparing directories: %s and %s \n", leftDir, rightDir)
		result, err := directories.CompareDirectories(leftDir, rightDir, time.Time{})
		if err != nil {
			fmt.Printf("Failed to compare directories: %s \n", err.Error())
			return
		}
		logComparisonResults(result, leftDir, rightDir)

		if outputDir != "" {
			fmt.Printf("Saving comparison details in the directory: %q", outputDir)
			err := directories.SaveComparisonResults(result, outputDir, fromDate)
			if err != nil {
				fmt.Printf("Failed to compare directories: %s \n", err.Error())
			}
			return
		}
	},
}

func parseStartFromDate(fromDateString string) (time.Time, error) {
	if fromDateString != "" {
		fromDate, err := time.Parse(time.RFC3339, fromDateString)
		if err != nil {
			return time.Time{}, fmt.Errorf("failed to parse `fromDate': %q", err)
		}

		return fromDate, nil
	}

	return time.Time{}, nil
}

func logComparisonResults(comparisonResult directories.Result, leftDir, rightDir string) {
	fmt.Printf("Numer of files in %s directory = %d \n", leftDir, comparisonResult.LeftDirFilesCount)
	fmt.Printf("Numer of files in %s directory = %d \n", rightDir, comparisonResult.RightDirFilesCount)

	if comparisonResult.Equal {
		fmt.Println("Directories are equal")
	} else {
		fmt.Println("Directories are NOT equal")
	}
}
