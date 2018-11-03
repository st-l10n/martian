package cli

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/st-10n/martian/resource"
)

var diffCmd = &cobra.Command{
	Use: "diff",
	Aliases: []string{
		"d",
	},
	Short: "Difference",
	RunE: func(cmd *cobra.Command, args []string) error {
		var (
			f = cmd.Flags()

			bDir, aDir string
			templates  []originalFile
			limit      []string
			err        error
			languages  Languages
			english    Language
		)
		if aDir, err = f.GetString("original"); err != nil {
			return err
		}
		if bDir, err = f.GetString("modified"); err != nil {
			return err
		}
		if len(bDir) == 0 {
			return errors.New("blank modified dir")
		}
		if err = viper.UnmarshalKey("languages", &languages); err != nil {
			return err
		}
		if limit, err = f.GetStringSlice("limit"); err != nil {
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
		if err = filepath.Walk(aDir, func(path string, info os.FileInfo, err error) error {
			base := filepath.Base(path)
			relative, err := filepath.Rel(aDir, filepath.Dir(path))
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
			return errors.New("no english files found in input folder")
		}
		fmt.Println("templates:", templates)
	Loop:
		for _, lang := range languages {
			if len(limit) > 0 {
				isSelected := false
				for _, limitLang := range limit {
					if strings.ToLower(limitLang) == strings.ToLower(lang.Name) {
						isSelected = true
					}
					if strings.ToLower(limitLang) == strings.ToLower(lang.Code) {
						isSelected = true
					}
				}
				if !isSelected {
					continue Loop
				}
			}
			fmt.Println("Language:", lang.Name)
			if lang.Prefix == "" {
				lang.Prefix = strings.ToLower(
					strings.Replace(lang.Name, " ", "_", -1),
				)
			}
			if lang.Locale == "" {
				lang.Locale = strings.ToLower(lang.Code)
			}
			//fmt.Printf("  prefix: %s\n", lang.Prefix)
			//fmt.Printf("  code: %s\n", lang.Code)
			//fmt.Printf("  locale: %s\n", lang.Locale)
			var aEntries, bEntries resource.Entries
			gen := func(t originalFile, dir string, e resource.Entries) (resource.Entries, error) {
				name := lang.Prefix + t.Postfix
				translatedPath := filepath.Join(dir, t.Path, name)
				origPath := filepath.Join(aDir, t.Path, "english"+t.Postfix)
				original, err := readFile(origPath)
				if err != nil {
					return nil, fmt.Errorf("failed to read english translation file: %v", err)
				}
				translated, err := readFile(translatedPath)
				if err != nil {
					if !(os.IsNotExist(err) && t.Postfix != ".xml") {
						return nil, fmt.Errorf("failed to find translated file for %s", lang.Code)
					}
				}
				o := resource.GenOptions{
					Original:   original,
					Translated: translated,
					Simplified: viper.GetStringSlice("simplified"),
				}
				// Scenario/EscapeFromMars/Language/english_mars_mission.xml -> EscapeFromMars
				// Language -> ""
				for _, s := range strings.Split(t.Path, string(filepath.Separator)) {
					if s != "Language" {
						o.FilePrefix = s
					}
				}
				gotEntries, err := resource.Gen(o)
				if err != nil {
					return nil, err
				}
				return append(e, gotEntries...), nil
			}
			for _, t := range templates {
				aEntries, err = gen(t, aDir, aEntries)
				if err != nil {
					return err
				}
				bEntries, err = gen(t, bDir, bEntries)
				if err != nil {
					return err
				}
			}
			fmt.Printf("  original (diff)  %d\n", len(aEntries.DifferentFromOriginal()))
			fmt.Printf("  automated        %d\n", len(bEntries.DifferentFromOriginal()))
			fmt.Printf("  original (total) %d\n", aEntries.TranslatedCount())
		}
		return nil
	},
}

func init() {
	{
		f := diffCmd.Flags()
		f.StringP("original", "o", ".", "original directory")
		f.StringP("modified", "m", ".", "modified directory")
		f.StringSlice("limit", nil, "limit languages")
		f.StringP("prefix", "p", "", "filename prefix")
	}
	rootCmd.AddCommand(
		diffCmd,
	)
}
