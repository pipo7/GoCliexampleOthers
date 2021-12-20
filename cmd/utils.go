package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
)

// InputSets used in base command to parse input data
type InputSets struct {
	InputSets []InputSet `json:"inputsets"`
}

// InputSet is one set of data
type InputSet struct {
	Operations int `json:"operations"`
	Workers    int `json:"workers"`
}

// Parses the json file and return custom type i.e []InputSet struct.
func parseJSONFile(filename string) (InputSets, error) {

	// We initialize our inputsets array
	var inputsets InputSets

	// Open our jsonFile
	fmt.Println("Opening " + filename + "...")
	jsonFile, err := os.Open(filename)
	defer func() {
		_ = jsonFile.Close()
	}()

	if err != nil {
		return inputsets, fmt.Errorf("Error while opening %s json file: %v", filename, err)
	}

	// Reading input pair values from json file
	fmt.Println("Reading input pairs...")
	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return inputsets, fmt.Errorf("Error while reading %s json file: %v", filename, err)
	}

	// Unmarshal our byteArray which contains our file content into InputSets format
	fmt.Println("Unmarshaling input pairs...")
	err = json.Unmarshal(byteValue, &inputsets)
	if err != nil {
		return inputsets, fmt.Errorf("Error while unmarshaling %s json file: %v", filename, err)
	}
	return inputsets, nil
}

func executeCommand(command string, arg ...string) (string, error) {
	c := exec.Command(command, arg...)
	stderr := &bytes.Buffer{}
	stdout := &bytes.Buffer{}
	c.Stderr = stderr
	c.Stdout = stdout
	if err := c.Run(); err != nil {
		err := fmt.Errorf("\n\t%v%v", stderr.String(), err)
		return "", err
	}
	return stdout.String(), nil

}

func cleanup() {
	fmt.Println("\nCleaning deployment created by run...")
	_, _ = executeCommand("bash", "-c", "kubectl get deployment --no-headers -o custom-columns=:metadata.name | xargs kubectl delete deployment --wait=true; exit 0 ")
	_, _ = executeCommand("bash", "-c", "kubectl get nodes -l node-role.kubernetes.io/master=true --no-headers -o custom-columns=\":metadata.name\" | xargs kubectl uncordon")
}

// Create the output directory if it does not exists
func createOutputDir(dir string) error {

	if _, err := os.Stat(dir); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(dir, os.ModePerm)
		if err != nil {
			return fmt.Errorf("Error created '%s' output directory: %v", dir, err)
		}
		fmt.Printf("Output directory '%s' created\n", dir)
	} else {
		fmt.Printf("Output directory '%s' already exists\n", dir)
	}
	return nil
}

// RemoveGlob used in main program to clear the output directory's contents
func removeGlob(dir string) (err error) {
	contents, err := filepath.Glob(filepath.Join(dir, "*"))
	if err != nil {
		return err
	}
	for _, item := range contents {
		err = os.RemoveAll(item)
		if err != nil {
			return err
		}
	}
	return nil
}
