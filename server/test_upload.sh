#!/bin/bash -
#===============================================================================
#
#          FILE: test_upload.sh
#
#         USAGE: ./test_upload.sh <input-file>
#
#        AUTHOR: Christian Heusel (christian@heusel.eu),
#  ORGANIZATION: Fachschaft MathPhysInfo
#       CREATED: 02/05/21 20:59
#
#===============================================================================


# check if command line argument is empty or not present
if [ "$1" == "" ] || [ $# -ne 1 ]; then
    echo "No input file given!"
    echo "USAGE: ./test_upload.sh <input-file>"
    exit 0
fi

set -o nounset                              # Treat unset variables as an error

curl --silent 'http://localhost:8081/query' \
    -F operations='{ "query": "mutation createNewExam($input: NewExam!) {createExam(input: $input) {UUID, subject, moduleName} }", "variables": { "input": { "subject": "Tom Rix", "moduleName": "test", "file": null } } }' \
    -F map='{ "0": ["variables.input.file"] }' \
    -F 0=@$1 | \
    python3 -m json.tool

# Original:
# curl localhost:8081/query \
#   -F operations='{ "query": "mutation ($file: Upload!) { createExam(file: $file) { id, name, content } }", "variables": { "file": null } }' \
#   -F map='{ "0": ["variables.file"] }' \
#   -F 0=@./gqlgen.yml
