package cli

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var updateCmd = &cobra.Command{
	Use: "update",
	Aliases: []string{
		"u", "update",
	},
	Short: "Copy the english xml files from game",
	RunE: func(cmd *cobra.Command, args []string) error {
		var (
			f = cmd.Flags()

			outDir, inDir string
			templates     []originalFile
			err           error
			languages     Languages
			english       Language
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
		if err = viper.UnmarshalKey("languages", &languages); err != nil {
			return err
		}
		for _, lang := range languages {
			if lang.Code == "EN" {
				english = lang
				break
			}
		}
		if english.Code == "" {
			return errors.New("no english language configured (code=EN)")
		}
		if err = filepath.Walk(inDir, func(path string, info os.FileInfo, err error) error {
			base := filepath.Base(path)
			relative, err := filepath.Rel(inDir, filepath.Dir(path))
			if err != nil {
				return err
			}
			if strings.HasPrefix(base, "english") && strings.HasSuffix(base, ".xml") {
				templates = append(templates, originalFile{
					Postfix: strings.TrimPrefix(base, "english"),
					Path:    relative,
				})
			}
			return nil
		}); err != nil {
			return err
		}
		if len(templates) == 0 {
			return errors.New("no english files found in output folder")
		}
		fmt.Println("templates:", templates)
		for _, t := range templates {
			name := "english" + t.Postfix
			origName := filepath.Join(inDir, t.Path, name)
			orig, err := readFile(origName)
			if err != nil {
				return err
			}
			outName := filepath.Join(outDir, t.Path, name)
			outF, err := os.Create(outName)
			if err != nil {
				return err
			}
			if _, err = outF.Write(orig); err != nil {
				return err
			}
			if err = outF.Close(); err != nil {
				return err
			}
			fmt.Println(origName, "->", outName)
		}
		return nil
	},
}

func init() {
	{
		f := updateCmd.Flags()
		f.StringP("output", "o", ".", "output directory (StreamingAssets repo)")
		f.StringP("input", "i", "game", "input directory (StreamingAssets from game)")
	}
	rootCmd.AddCommand(
		updateCmd,
	)
}
