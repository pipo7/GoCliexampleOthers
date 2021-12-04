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
	"gopkg.in/go-playground/validator.v9"
)

var cfgFile string

var config Config

type Config struct {
	Operations int `json:"operations"`
	Workers    int `json:"workers"`
}

var myconfig MyConfig

type MyConfig struct {
	word     string
	myint    int
	waittime time.Duration
}

var validate *validator.Validate

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cliexample2",
	Short: "cliexample2 brief description of your application",
	Long:  `cliexample2 longer description `,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {

		fmt.Println("cliexample2 with subCommands foo or bar")
		fmt.Println("String passed to flag word is :", viper.GetString("word"))
		fmt.Println("Int passed to flag myint is :", viper.GetInt("myint"))
		d := viper.GetDuration("waittime")
		fmt.Printf("d is duration %s passed via duration flag of type %T\n", d, d)
		time.Sleep(d)
		fmt.Println("Config file is : ", viper.ConfigFileUsed())
		//fmt.Println("Config file values : ", viper.GetString("DRIVER"), viper.GetString("HOST"))

		myconfig.myint = viper.GetInt("myint")
		myconfig.word = viper.GetString("word")
		myconfig.waittime = viper.GetDuration("waittime")
		if err := viper.Unmarshal(&myconfig); err != nil {
			fmt.Fprintln(os.Stderr, "Config file unmarshal error:")
			fmt.Println(err)
		} else {
			fmt.Printf("MyConfig file values : %+v\n", myconfig)
		}
		validate = validator.New()
		my := MyConfig{
			word:     viper.GetString("word"),
			myint:    viper.GetInt("myint"),
			waittime: viper.GetDuration("waittime"),
		}
		fmt.Println("my :", my)
		if errs := validate.Struct(my); errs != nil {
			fmt.Println("my errs:", errs)
		}
		err := validate.Struct(my)
		if err != nil {

			// this check is only needed when your code could produce
			// an invalid value for validation such as interface with nil
			// value most including myself do not usually have code like this.
			if _, ok := err.(*validator.InvalidValidationError); ok {
				fmt.Println(err)
				return
			}

			for _, err := range err.(validator.ValidationErrors) {

				fmt.Println(err.Namespace())
				fmt.Println(err.Field())
				fmt.Println(err.StructNamespace())
				fmt.Println(err.StructField())
				fmt.Println(err.Tag())
				fmt.Println(err.ActualTag())
				fmt.Println(err.Kind())
				fmt.Println(err.Type())
				fmt.Println(err.Value())
				fmt.Println(err.Param())
				fmt.Println()
			}

			// from here you can create your own error messages in whatever language you wish
			return
		}

		/*
			if errs := validator.Validate(myconfig.myint); errs != nil {
				// values not valid, deal with errors here
				fmt.Println("Validation errors", errs)
			} else {
				fmt.Printf("MyConfig file values : %+v\n", myconfig)
			}
		*/
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
	rootCmd.Flags().IntP("myint", "i", 10, "enter any integer")
	rootCmd.Flags().DurationP("waittime", "s", 60, "Specify time in sec like 10s or 30s")

	//Bind Flags
	viper.BindPFlag("word", rootCmd.Flags().Lookup("word"))
	viper.BindPFlag("myint", rootCmd.Flags().Lookup("myint"))
	viper.BindPFlag("waittime", rootCmd.Flags().Lookup("waittime"))
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
	/*
		if errs := validator.Validate(config); errs != nil {
			// values not valid, deal with errors here
			fmt.Println("Validation errors", errs)
		} */
	//fmt.Printf("Config file values : %+v\n", config)
	fmt.Printf("Config file values are : Operations %v Workers %v\n", config.Operations, config.Workers)
	/*if err := viper.Unmarshal(&config); err != nil {
		fmt.Fprintln(os.Stderr, "Config file unmarshal error:")
		fmt.Println(err)
	}*/

}
