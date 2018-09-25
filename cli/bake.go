package cli

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/st-10n/martian/resource"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func stringIn(s string, list []string) bool {
	for _, l := range list {
		if l == s {
			return true
		}
	}
	return false
}

type originalFile struct {
	Postfix string
	Path    string
}

func (o originalFile) String() string {
	return path.Join(o.Path, "english"+o.Postfix)
}

var bakeCmd = &cobra.Command{
	Use: "bake",
	Aliases: []string{
		"b",
	},
	Short: "Translate the .xml files using the .po ones",
	RunE: func(cmd *cobra.Command, args []string) error {
		var (
			f = cmd.Flags()

			outDir, inDir string
			templates     []originalFile
			limit         []string
			ignore        []string
			err           error
			languages     Languages
			assetsName    string
			english       Language
		)
		if inDir, err = f.GetString("input"); err != nil {
			return err
		}
		if assetsName, err = f.GetString("list"); err != nil {
			return err
		}
		if !path.IsAbs(assetsName) {
			assetsName = path.Join(outDir, assetsName)
		}
		assetsF, err := os.Create(assetsName)
		if err != nil {
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
		if limit, err = f.GetStringSlice("limit"); err != nil {
			return err
		}
		if ignore, err = f.GetStringSlice("ignore"); err != nil {
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
		if err = filepath.Walk(outDir, func(path string, info os.FileInfo, err error) error {
			if info.IsDir() {
				for _, i := range ignore {
					if stringIn(path, []string{i, filepath.Join(outDir, i)}) {
						fmt.Println("skipping the", path)
						return filepath.SkipDir
					}
				}
			}
			base := filepath.Base(path)
			relative, err := filepath.Rel(outDir, filepath.Dir(path))
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
		simplified := viper.GetStringSlice("simplified")
		fmt.Println("templates:", templates)
		fmt.Println("limit:", limit)
		fmt.Println("simplified:", simplified)
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
			fmt.Printf("  prefix: %s\n", lang.Prefix)
			fmt.Printf("  code: %s\n", lang.Code)
			fmt.Printf("  locale: %s\n", lang.Locale)
			localeDir := filepath.Join(inDir, lang.Locale)
			var localizations [][]byte
			if err = filepath.Walk(localeDir, func(path string, info os.FileInfo, err error) error {
				if !strings.HasSuffix(path, ".po") {
					return nil
				}
				locF, locErr := os.Open(path)
				if locErr != nil {
					return locErr
				}
				defer locF.Close()
				buf, readErr := ioutil.ReadAll(locF)
				if readErr != nil {
					return readErr
				}
				localizations = append(localizations, buf)
				return nil
			}); err != nil {
				return err
			}
			if len(localizations) == 0 {
				return fmt.Errorf("failed to found .po files in %s", localeDir)
			}
			if lang.Code == "EN" || lang.Locale == "en" {
				fmt.Println("skipping english as readonly")
				continue
			}
			for _, t := range templates {
				name := lang.Prefix + t.Postfix
				outName := filepath.Join(outDir, t.Path, name)
				origName := filepath.Join(outDir, t.Path, "english"+t.Postfix)
				orig, err := readFile(origName)
				if err != nil {
					return err
				}
				opt := resource.Options{
					Simplified:  simplified,
					Code:        lang.Code,
					Name:        lang.Name,
					Original:    orig,
					Translation: localizations,
				}
				if lang.Font != "" {
					opt.Font = "font_" + lang.Font
				}
				if !stringIn(t.Postfix, []string{".xml", "_tutorial.xml", "_mars_mission.xml"}) {
					opt.Font = ""
					opt.Name = ""
				}
				out, err := resource.Bake(opt)
				if err != nil {
					return err
				}
				outF, err := os.Create(outName)
				if err != nil {
					return err
				}
				if _, err = outF.Write(out); err != nil {
					return err
				}
				if err = outF.Close(); err != nil {
					return err
				}
				fmt.Println(outName)
				fmt.Fprintln(assetsF, outName)
			}
		}
		return assetsF.Close()
	},
}

func init() {
	{
		f := bakeCmd.Flags()
		f.StringP("output", "o", ".", "output directory (StreamingAssets)")
		f.StringP("list", "l", "assets.txt", "output file list")
		f.StringP("input", "i", "locales", "input directory (locales)")
		f.StringSlice("limit", nil, "limit languages")
		f.StringSlice("ignore", []string{"game"}, "ignore directories")
	}
	rootCmd.AddCommand(
		bakeCmd,
	)
}
