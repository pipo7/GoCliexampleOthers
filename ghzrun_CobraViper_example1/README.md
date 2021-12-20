## Prerequisites

1.	Download the project, navigate to perftest directory, run go mod tidy and then build it as:
  	"make"
2.	Enter the perftest arguments that is the number of operations in operations key value and the number of workers in workers key value in inputsets.json file. Multiple pair values can be entered, each in {} block followed by "," unless it’s the last block

Example of .json input file format :
{
    "inputsets": [
      {
        "operations": 100,
        "workers": 1
      }
	]
}	

3.	Specify the ks deployment replicas count as per flag  -k or --ks.
4.	Enter the output directory name without path using -d flag. If it does not exists then it will be created at the current working directory path. The default output directory name is results/ 
5.	Refer perftest --help.  For ease, install the perftest binary in system PATH. 

## Execution steps

1.	Refer perftest --help 
2.	Run with options presented
# Some Examples
For automated run only pass the flag --ks with --outputformat:
  perftest --ks 4 --outputformat 'html'
  OR use shorthand -k flag and specify keyspace replicas count as:
  perftest -k 8 -o 'csv'
For manual run pass flags --manual, --ks --operations and --workers:
  perftest --manual --ks 4 --operations 1000 --workers 10
  OR use shorthand flags as:
  perftest -m -k 4 -n 100 -c 1
6.	## Output generated files from perftest run will be in output directory
3.	Navigate to output directory for generated files.


