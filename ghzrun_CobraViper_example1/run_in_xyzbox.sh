#!/bin/bash -e

## algo by default is : aesgcm

## generate new KSCKID and KSID 
echo "Generating KSCKID, and KSID"
# Create new templates



## generate first time KSCKID and KID and then re-use these
KSCKID=1234 && echo $KSCKID
KSID=5678 && echo $KSID

# Check that KSCK and KS ID are genereated ok
if [ -z "$KSCKID" ] || [ -z "$KSID" ]
then
      echo "Either KSCKID OR KSID is empty or not generated due to an error"
      exit 1
else
     echo "KSCKID is $KSCKID and KSID is $KSID"
     export KSCKID
     export KSID
fi

## Generating ghzencrypt.json file with config parameters
echo "generated ghzencryptjson file"
if ! [ $? -eq 0 ]
then
echo "An error while generating ghzencrypt.json config file"
exit $?
fi

