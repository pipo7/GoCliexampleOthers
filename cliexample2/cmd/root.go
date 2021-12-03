/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/validator.v2"
)

var cfgFile string

var config Config

type Config struct {
	Operations int `json:"operations" validate:"min=10"`
	Workers    int `json:"workers" validate:"min=10"`
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cliexample2",
	Short: "cliexample2 brief description of your application",
	Long:  `cliexample2 longer description `,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("cliexample2 with subCommands foo or bar")
		flagstring, _ := cmd.Flags().GetString("word")
		fmt.Println("String passed to flag word is :", flagstring)
		flagint, _ := cmd.Flags().GetInt("myint")
		fmt.Println("Int passed to flag myint is :", flagint)
		d, _ := cmd.Flags().GetDuration("waittime")
		fmt.Printf("d is duration %s passed via duration flag and %T\n", d, d)
		time.Sleep(d)
		fmt.Println("Config file is : ", viper.ConfigFileUsed())

		//fmt.Println("Config file values : ", viper.GetString("DRIVER"), viper.GetString("HOST"))
		fmt.Println("End of run")

	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "/home/ps/cliexample2/input.json", "config file (default is $HOME/.cliexample2.yaml)")
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cliexample2.yaml)")
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.Flags().StringP("word", "w", "meaning", "enter a word to find its meaning")
	rootCmd.Flags().IntP("myint", "i", 0, "enter any integer")
	rootCmd.Flags().DurationP("waittime", "s", 10, "Time in sec")

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".cliexample2" (without extension).
		viper.AddConfigPath(strings.Join([]string{home, "/cliexample2/"}, ""))
		viper.SetConfigType("json")
		viper.SetConfigName("input")

	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
		fmt.Println(err)
	}

	if err := viper.Unmarshal(&config); err != nil {
		fmt.Fprintln(os.Stderr, "Config file unmarshal error:")
		fmt.Println(err)
	}
	if errs := validator.Validate(config); errs != nil {
		// values not valid, deal with errors here
		fmt.Println("Validation errors", errs)

	}
	fmt.Printf("Config file values : %+v\n", config)

	/*if err := viper.Unmarshal(&config); err != nil {
		fmt.Fprintln(os.Stderr, "Config file unmarshal error:")
		fmt.Println(err)
	}*/

}
