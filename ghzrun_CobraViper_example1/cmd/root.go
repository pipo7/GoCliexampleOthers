package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

// CmdFlags used in base command
type CmdFlags struct {
	manual           bool
	outputFormat     string
	waitTime         time.Duration
	Replicas int
	operations       int
	workers          int
	outputDir        string
	inputFile        string
}

var cmdFlags CmdFlags

// rootCmd represents the base command
var rootCmd = &cobra.Command{
	Use:   "perftest",
	Short: "Runs perftest in manual or automated mode",
	Long:  printLongDescText(),
	Run: func(cmd *cobra.Command, args []string) {
		// Cleans up - deployments and uncordon's the master node when exiting
		defer cleanup()
		err := run(cmd, args)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() error {
	return rootCmd.Execute()
}

func run(cmd *cobra.Command, args []string) error {

	if cmdFlags.manual {
		fmt.Println("PerfTest manual mode run in progress...")
	} else {
		fmt.Println("PerfTest automated mode run in progress...")
	}

	// Get and check the output file format to be used.
	format := strings.ToLower(cmdFlags.outputFormat)
	fmt.Println("Output format is ", format)
	switch {
	case format == "csv":
	case format == "html":
	case format == "json":
		break
	default:
		return fmt.Errorf("Error '%s' output format is not supported", cmdFlags.outputFormat)
	}

	// Create the output directory if it does not exists.
	err := createOutputDir(cmdFlags.outputDir)
	if err != nil {
		return fmt.Errorf("Error creating the '%s' output directory: %v", cmdFlags.outputDir, err)
	}

	// Clear the contents of ghzoutput directory
	fmt.Printf("\nCleaning '%s' local output directory...\n", cmdFlags.outputDir)
	err = removeGlob(cmdFlags.outputDir)
	if err != nil {
		return fmt.Errorf("Error removing files in '%s' output directory: %v", cmdFlags.outputDir, err)
	}

	// Delete any existing ID1 or ID2
	cleanup()

	// Get the name of XYZ pod in cluster
	fmt.Println("Getting the XYZ pod name...")
	xyzPodName, err := executeCommand("bash", "-c", "podname=$(kubectl get pods | grep XYZ | cut -d ' ' -f1) && echo $podname")
	if err != nil {
		return fmt.Errorf("Error while getting pod name: %v", err)

	}
	xyzPodName = strings.TrimSuffix(xyzPodName, "\n")

	//Copy script run_in_xyzbox inside the xyz-box container
	fmt.Println("\nCopying run_in_xyzbox.sh in xyz-box...")
	_, err = executeCommand("bash", "-c", "kubectl cp run_in_xyzbox.sh "+xyzPodName+":/home/xyz -c xyz-box")
	if err != nil {
		return fmt.Errorf("Error while copying bash file: %v", err)

	}

	// Execute the copied script in xyz-box
	fmt.Println("\nExecuting run_in_xyzbox.sh in xyz-box to create ID2...")
	_, err = executeCommand("bash", "-c", "kubectl exec -it "+xyzPodName+" -c xyz-box -- /bin/bash run_in_xyzbox.sh")
	if err != nil {
		return fmt.Errorf("Error while executing bash file: %v", err)

	}

	// clean output folder in xyz-box
	fmt.Printf("\nCleaning %s in xyz-box...\n", cmdFlags.outputDir)
	_, err = executeCommand("bash", "-c", "kubectl exec -t "+xyzPodName+" -c xyz-box -- rm -rf "+cmdFlags.outputDir)
	if err != nil {
		return fmt.Errorf("Error while cleaning output dir: %v", err)
	}

	// Prepare output folder in xyz-box
	fmt.Printf("\nCreating %s directory in xyz-box...\n", cmdFlags.outputDir)
	_, err = executeCommand("bash", "-c", "kubectl exec -t "+xyzPodName+" -c xyz-box -- mkdir -p "+cmdFlags.outputDir)
	if err != nil {
		return fmt.Errorf("Error while creating output dir: %v", err)
	}

	// Cordon master node if needed
	fmt.Println("\nGetting number of nodes...")
	nodesCount, err := executeCommand("bash", "-c", "kubectl get nodes --no-headers=true | wc -l")
	if err != nil {
		return fmt.Errorf("Error while creating output dir: %v", err)
	}

	if nc, _ := strconv.Atoi(nodesCount); nc > 1 {
		fmt.Println("\nCordoning Master Node...")
		_, err = executeCommand("bash", "-c", "kubectl get nodes -l node-role.kubernetes.io/master=true --no-headers -o custom-columns=\":metadata.name\" | xargs kubectl cordon")
		if err != nil {
			return fmt.Errorf("Error while cordoning Master node: %v", err)
		}
	}

	// Get the  deployment name.
	fmt.Println("Getting  deployment name...")
	deplName, err := executeCommand("bash", "-c", "deplName=$(kubectl get deployments | grep -E '^deployment-[a-z0-9].+-[a-z0-9].+-[a-z0-9].+-[a-z0-9].+-[a-z0-9].+' | cut -d ' ' -f1) && echo $deplName")
	if err != nil {
		return fmt.Errorf("Error getting deployment name: %v", err)

	}
	deplName = strings.TrimSuffix(deplName, "\n")

	// Patch with pod spread topology
	fmt.Println("Patching pod spread topology...")
	_, err = executeCommand("bash", "-c", "kubectl  patch deployment "+deplName+" --patch-file patch.yaml --type=merge")
	if err != nil {
		return fmt.Errorf("Error while patching pod spread topology: %v", err)
	}

	// Change replicas with patching
	fmt.Printf("Patching keyspace replicas count to: %d\n", cmdFlags.Replicas)
	replicas := strconv.Itoa(cmdFlags.Replicas)
	_, err = executeCommand("bash", "-c", "kubectl patch keyspace "+deplName+" --patch $'spec:\\n  replicas: "+replicas+"' --type=merge")
	if err != nil {
		return fmt.Errorf("Error while patching keyspace replicas count: %v", err)
	}

	fmt.Printf("Waiting for few seconds - changing replicas count to %d & distributing ksapp pods accross the nodes now \n", cmdFlags.Replicas)
	// Not using dynamic wait as sometimes old pods remains in Terminated state for long, and these pods are not requried.
	time.Sleep(cmdFlags.waitTime)

	// Check whether its manual or automated run and proceed accordingly
	var ips InputSets
	if cmdFlags.manual {
		// In Manual run get operations and workers from flag values
		ips = InputSets{[]InputSet{InputSet{Operations: cmdFlags.operations, Workers: cmdFlags.workers}}}
	} else {
		// In Automated run .Execute the ghz with each pair of data input defined in inputsets json file
		fmt.Printf("Parsing %s\n", cmdFlags.inputFile)
		ips, err = parseJSONFile(cmdFlags.inputFile)
		if err != nil {
			return fmt.Errorf("Error while parsing json file: %v", err)
		}
	}
	for _, v := range ips.InputSets {
		dt := time.Now()
		operations := strconv.Itoa(v.Operations)
		workers := strconv.Itoa(v.Workers)
		fmt.Printf("\nExecuting performance measurements with '%s' Operations and '%s' Workers...\n", operations, workers)
		_, err = executeCommand("bash", "-c", "kubectl exec -t "+xyzPodName+" -c xyz-box -- ./ghz --config ghzencrypt.json -n "+operations+" -c "+workers+" -O "+format+" -o "+cmdFlags.outputDir+"/result-"+dt.Format("01-02-2006-15-04-05")+"-numops"+operations+"-workers"+workers+"."+format)
		if err != nil {
			return fmt.Errorf("Error while executing: %v", err)
		}
	}

	// Copy the output files from xyz-box to host
	fmt.Printf("\nCopying output file(s) at %s directory...\n", cmdFlags.outputDir)
	_, err = executeCommand("bash", "-c", "kubectl cp "+xyzPodName+":/home/xyz/"+cmdFlags.outputDir+" $(pwd)/"+cmdFlags.outputDir+" -c xyz-box > /dev/null")
	if err != nil {
		return fmt.Errorf("Error while copying output files: %v", err)

	}
	return nil
}

func init() {

	// Initialize command line flags with values
	cobra.OnInitialize(initCmdFlags)

	// Not using viper to read JSON Array file "inputsets.json" .
	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.perftest.yaml)")

	// Cobra also supports local flags, which will only run when this action is called directly.
	rootCmd.Flags().BoolVarP(&cmdFlags.manual, "manual", "m", false, "For manual run use --manual boolean flag and ignore it for automated run")
	rootCmd.Flags().IntVarP(&cmdFlags.Replicas, "replicas", "k", 4, "You may specify the Keyspace deployment pods' replicas count ,which must be a positive integer")
	rootCmd.Flags().IntVarP(&cmdFlags.operations, "operations", "n", 200, "Only when using --manual flag , then you may also Specify the number of operations , which must be a positive integer")
	rootCmd.Flags().IntVarP(&cmdFlags.workers, "workers", "w", 50, "Only when using --manual flag , then you may also specify the number of concurrent workers , which must be a positive integer")
	rootCmd.Flags().DurationVarP(&cmdFlags.waitTime, "waittime", "t", 60*time.Second, "Time in seconds to wait until all pods are in Runnning state , which must be a positive integer. Example: 30s or 60s")
	rootCmd.Flags().StringVarP(&cmdFlags.outputFormat, "outputformat", "o", "html", "You may specify the output file format one among 'csv' or 'html' or 'json'")
	rootCmd.Flags().StringVarP(&cmdFlags.inputFile, "inputfile", "i", "inputsets.json", "You may specify the input json file")
	rootCmd.Flags().StringVarP(&cmdFlags.outputDir, "outputdir", "d", "results", "You may specify the results output directory")

	//Bind Flags using viper
	_ = viper.BindPFlag("manual", rootCmd.Flags().Lookup("manual"))
	_ = viper.BindPFlag("replicas", rootCmd.Flags().Lookup("replicas"))
	_ = viper.BindPFlag("operations", rootCmd.Flags().Lookup("operations"))
	_ = viper.BindPFlag("workers", rootCmd.Flags().Lookup("workers"))
	_ = viper.BindPFlag("waittime", rootCmd.Flags().Lookup("waittime"))
	_ = viper.BindPFlag("outputformat", rootCmd.Flags().Lookup("outputformat"))
	_ = viper.BindPFlag("inputfile", rootCmd.Flags().Lookup("inputfile"))
	_ = viper.BindPFlag("outputdir", rootCmd.Flags().Lookup("outputdir"))
}

func printLongDescText() string {
	return "\nRuns ghz in manual or automated mode, pass the boolean flag --manual for manual run, See below examples: \n\n" +
		"\tFor automated run only pass the flag --replicas with --outputformat: \n" +
		"\t\t\tperftest --replicas 4 --outputformat 'html'\n" +
		"\t\tOR use shorthand -k flag and specify keyspace replicas count as:\n" +
		"\t\t\tperftest -k 8 -o 'csv'\n" +
		"\tFor manual run pass flags --manual, --replicas --operations and --workers:\n" +
		"\t\t\tperftest --manual --replicas 4 --operations 1000 --workers 10\n" +
		"\t\tOR use shorthand flags as:\n" +
		"\t\t\tperftest -m -k 4 -n 100 -c 1\n"
}

func initCmdFlags() {

	cmdFlags := CmdFlags{}
	cmdFlags.manual = viper.GetBool("manual")
	cmdFlags.outputFormat = viper.GetString("outputformat")
	cmdFlags.waitTime = viper.GetDuration("waittime")
	cmdFlags.Replicas = viper.GetInt("replicas")
	cmdFlags.operations = viper.GetInt("operations")
	cmdFlags.workers = viper.GetInt("workers")
	cmdFlags.outputDir = viper.GetString("outputdir")
	cmdFlags.inputFile = viper.GetString("inputfile")
	if err := viper.Unmarshal(&cmdFlags); err != nil {
		log.Fatalf("Config file unmarshal error: %v", err)
	}

}
