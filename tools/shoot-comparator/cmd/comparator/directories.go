package comparator

import (
	"fmt"
	"github.com/kyma-project/infrastructure-manager/tools/shoot-comparator/internal/directories"
	"github.com/spf13/cobra"
)

var directoriesCmd = &cobra.Command{
	Use:     "dirs",
	Aliases: []string{"d"},
	Short:   "Compare directories",
	Args:    cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		leftDir := args[0]
		rightDir := args[1]

		fmt.Printf("Comparing directories: %s and %s \n", leftDir, rightDir)
		result, err := directories.CompareDirectories(leftDir, rightDir)
		if err != nil {
			fmt.Printf("Failed to compare directories: %s \n", err.Error())
			return
		}
		fmt.Printf("Numer of files in %s directory = %d \n", leftDir, result.LeftDirFilesCount)
		fmt.Printf("Numer of files in %s directory = %d \n", rightDir, result.RightDirFilesCount)

		if result.Equal {
			fmt.Print("Directories are equal \n")
		} else {
			fmt.Print("Directories are NOT equal \n")
			if len(result.LeftOnly) != 0 {
				fmt.Printf("Files existing in %s folder only: %s", leftDir, result.LeftOnly)
			}

			if len(result.RightOnly) != 0 {
				fmt.Printf("Files existing in %s folder only: %s", rightDir, result.RightOnly)
			}

			if len(result.Diff) != 0 {
				fmt.Printf("Differences found: %s", result.Diff)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(directoriesCmd)
}
