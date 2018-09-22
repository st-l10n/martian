package cli

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/st-10n/martian/resource"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

func readFile(name string) ([]byte, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	b, readErr := ioutil.ReadAll(f)
	if readErr != nil {
		return nil, err
	}
	return b, f.Close()
}

var genCmd = &cobra.Command{
	Use: "generate",
	Aliases: []string{
		"gen", "g",
	},
	Short: "Generate po file",
	RunE: func(cmd *cobra.Command, args []string) error {
		var (
			f = cmd.Flags()

			origName       string
			outName        string
			translatedName string
			lang           string
			err            error
		)
		if origName, err = f.GetString("original"); err != nil {
			return err
		}
		if outName, err = f.GetString("output"); err != nil {
			return err
		}
		if lang, err = f.GetString("language"); err != nil {
			return err
		}
		if translatedName, err = f.GetString("translated"); err != nil {
			return err
		}
		if outName == defaultOutputPO {
			outName = fmt.Sprintf("%s.po", strings.ToLower(lang))
		}
		if translatedName == defaultTranslatedXML {
			translatedName = fmt.Sprintf("%s.xml", strings.ToLower(lang))
		}
		var (
			orig, translated, out []byte
		)
		if orig, err = readFile(origName); err != nil {
			return err
		}
		if translated, err = readFile(translatedName); err != nil {
			return err
		}
		if out, err = resource.Gen(resource.GenOptions{
			Original:   orig,
			Translated: translated,
			Language:   lang,
		}); err != nil {
			return err
		}
		outFile, createErr := os.Create(outName)
		if createErr != nil {
			return createErr
		}
		if _, err = outFile.Write(out); err != nil {
			return err
		}
		return outFile.Close()
	},
}

var rootCmd = &cobra.Command{
	Use:   "martian",
	Short: "Stationeers Localization toolset",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("hello world")
	},
}

const defaultOutputPO = "{language_name}.po"
const defaultTranslatedXML = "{language_name}.xml"

func init() {
	cobra.OnInitialize(initConfig)
	{
		f := genCmd.Flags()
		f.StringP("original", "e", "english.xml", "path to original english file")
		f.StringP("translated", "t", defaultTranslatedXML, "path to previous translated xml file")
		f.StringP("output", "o", defaultOutputPO, "path to output po file")
		f.StringP("language", "l", "Language", "full language name")
	}
	rootCmd.AddCommand(
		genCmd,
	)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/martian.yml)")
}

func initConfigCommon() {
	home, err := homedir.Dir()
	if err != nil {
		log.Fatalln("failed to find home directory:", err)
	}
	viper.AddConfigPath(".")
	viper.AddConfigPath("/etc/st-l10n-martian/")
	viper.AddConfigPath(home)
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		initConfigCommon()
		viper.SetConfigName("martian")
		viper.SetConfigType("yaml")
	}
	cfgErr := viper.ReadInConfig()
	if cfgErr != nil {
		// log.Fatalln("failed to read config:", cfgErr)
	}
}

// Execute starts root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
