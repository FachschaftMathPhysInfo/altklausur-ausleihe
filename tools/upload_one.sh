#!/bin/bash -
#===============================================================================
#
#          FILE: upload.sh
#
#         USAGE: ./upload_one.sh <JWT-Token> <meta-data-string> <input-file>
#  DEPENDENCIES: curl, jq
#
#        AUTHOR: Christian Heusel (christian@heusel.eu),
#  ORGANIZATION: Fachschaft MathPhysInfo
#       CREATED: 02/05/21 20:59
#
#===============================================================================


# check if command line argument is empty or not present
if [ "$1" == "" ] || [ $# -ne 3 ]; then
    echo "Wrong amount of arguments!"
    echo "USAGE: ./upload_one.sh <JWT-Token> <meta-data-string> <input-file>"
    echo "Given:" $@
    exit 0
fi

TARGET_HOST="${TARGET_HOST_ENV:-https://altklausuren.mathphys.info/query}"
JWT_TOKEN=$1
METADATA_STRING=$2
INPUT_FILENAME=$3

OPERATION_STRING='{ "query": "mutation createNewExam($input: NewExam!) {createExam(input: $input) {UUID} }", "variables": { "input": '$METADATA_STRING' } }'
echo $OPERATION_STRING | jq

set -o nounset                               # Treat unset variables as an error

curl --silent $TARGET_HOST \
    -H 'Cookie: jwt='"$JWT_TOKEN" \
    -F operations="$OPERATION_STRING" \
    -F map='{ "0": ["variables.input.file"] }' \
    -F 0=@$INPUT_FILENAME
