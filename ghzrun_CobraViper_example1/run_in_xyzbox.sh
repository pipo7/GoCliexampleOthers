#!/bin/bash -e

## algo by default is : aesgcm

## generate new ID1, and ID2
echo "Generating ID1, and ID2"
# Create new templates



## generate first time ID1 and ID2 and then re-use these
ID1=1234 && echo $ID1
ID2=5678 && echo $ID2

# Check that ID1 and ID2 ID are genereated ok
if [ -z "$ID1" ] || [ -z "$ID2" ]
then
      echo "Either ID1 OR ID2 is empty or not generated due to an error"
      exit 1
else
     echo "ID1 is $ID1 and ID2 is $ID2"
     export ID1
     export ID2
fi

## Generating ghzencrypt.json file with config parameters
echo "generated ghzencryptjson file"
if ! [ $? -eq 0 ]
then
echo "An error while generating ghzencrypt.json config file"
exit $?
fi

