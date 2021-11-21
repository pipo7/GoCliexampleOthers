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
	"strconv"

	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "add Cmd to add integers . Give multiple integer values to be added separated by space",
	Long: `add Cmd examples:
	Valid examples
	add 2 3 
	add "2" '3'
	add 1 2 3 -4 5 
	add 3 "4" '5'
	
	Not valid examples:
	add 1.0 2
	add 1.0 2.0 `,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("add is called now")
		fstatus, _ := cmd.Flags().GetBool("float")
		if fstatus {
			addFloat(args)
		} else {
			addInt(args)
		}
	},
}

func init() {
	rootCmd.AddCommand(addCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	addCmd.Flags().BoolP("float", "f", false, "Adding floating numbers")
}

func addInt(args []string) {
	var sum int
	for _, ival := range args {
		itemp, err := strconv.Atoi(ival)
		if err != nil {
			fmt.Println(err)
		}
		sum = sum + itemp
	}
	fmt.Printf("Addition of numbers %s is %d\n", args, sum)
}

func addFloat(args []string) {
	var sum float64
	for _, fval := range args {
		// convert string to float64
		ftemp, err := strconv.ParseFloat(fval, 64)
		if err != nil {
			fmt.Println(err)
		}
		sum = sum + ftemp
	}
	fmt.Printf("Sum of floating numbers %s is %f\n", args, sum)
}
