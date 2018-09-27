package cli

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/st-10n/martian/resource"
)

var mergeCmd = &cobra.Command{
	Use: "merge",
	Aliases: []string{
		"m",
	},
	Short: "Merge old .po with new .po (legacy only)",
	RunE: func(cmd *cobra.Command, args []string) error {
		var (
			f = cmd.Flags()

			outDir, inDir string
			err           error
		)
		if inDir, err = f.GetString("input"); err != nil {
			return err
		}
		fmt.Println("inputDir:", inDir)
		if outDir, err = f.GetString("output"); err != nil {
			return err
		}
		if len(outDir) == 0 {
			return errors.New("blank output dir")
		}
		if err = filepath.Walk(outDir, func(path string, info os.FileInfo, err error) error {
			relative, err := filepath.Rel(outDir, path)
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}
			if !strings.HasSuffix(path, ".po") {
				return nil
			}
			var (
				inputName  = filepath.Join(inDir, relative)
				mergedName = filepath.Join(outDir, relative)
			)
			fmt.Println(inputName, mergedName)
			return resource.Merge(inputName, mergedName, "")
		}); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	{
		f := mergeCmd.Flags()
		f.StringP("output", "o", ".", "output directory")
		f.StringP("input", "i", "game", "input directory")
	}
	rootCmd.AddCommand(
		mergeCmd,
	)
}
