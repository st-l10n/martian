// Package cli implements command line interface for martian.
//
//go:generate go run -tags=dev config_generate.go
package cli

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

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

type Language struct {
	Code   string `mapstructure:"code"`
	Name   string `mapstructure:"name"`
	Prefix string `mapstructure:"prefix"`
	Locale string `mapstructure:"locale"`
	Font   string `mapstructure:"font"`
}

func (l Language) GetPrefix() string {
	if l.Prefix != "" {
		return l.Prefix
	}
	return strings.ToLower(
		strings.Replace(l.Name, " ", "_", -1),
	)
}

type Languages []Language

var rootCmd = &cobra.Command{
	Use:   "martian",
	Short: "Stationeers Localization toolset",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("hello world")
	},
}

func init() {
	cobra.OnInitialize(initConfig)
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
	if _, ok := cfgErr.(viper.ConfigFileNotFoundError); ok {
		f, err := Config.Open("martian.yml")
		if err != nil {
			log.Fatalf("failed to open default config: %v", err)
		}
		defer f.Close()
		cfgErr = viper.ReadConfig(f)
	}
	if cfgErr != nil {
		log.Fatalln("failed to read config:", cfgErr)
	}
}

// Execute starts root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
