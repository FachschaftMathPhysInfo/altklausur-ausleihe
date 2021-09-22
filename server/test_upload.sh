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

JWT_TOKEN="ADD.YOUR.TOKEN"

set -o nounset                              # Treat unset variables as an error

# mark Exam
curl 'https://altklausuren.mathphys.info/query' -H 'Accept-Encoding: gzip, deflate, br' -H 'Content-Type: application/json' -H 'Accept: application/json' -H "Cookie: jwt=$JWT_TOKEN" --data-binary '{"query":"mutation requestMarkedExam {\n  requestMarkedExam(StringUUID: \"85975cf6-b3da-4ef0-9e66-7ae44d47635e\",)\n}\n"}' --compressed
# get URL
curl 'https://altklausuren.mathphys.info/query' -H 'Accept-Encoding: gzip, deflate, br' -H 'Content-Type: application/json' -H 'Accept: application/json' -H "Cookie: jwt=$JWT_TOKEN" --data-binary '{"query":"query($UUID: String!) {\n  getExam(StringUUID: $UUID) {\n    viewUrl\n    downloadUrl\n  }\n}\n","variables":{"UUID":"81d6750d-a535-46b6-8491-accd6b1eb0c7"}}' --compressed

# Original:
# curl localhost:8081/query \
#   -F operations='{ "query": "mutation ($file: Upload!) { createExam(file: $file) { id, name, content } }", "variables": { "file": null } }' \
#   -F map='{ "0": ["variables.file"] }' \
#   -F 0=@./gqlgen.yml
# curl 'http://localhost:8081/query' -H 'Accept-Encoding: gzip, deflate, br' -H 'Content-Type: application/json' -H 'Accept: application/json' -H 'Connection: keep-alive' -H 'DNT: 1' -H 'Origin: http://localhost:8081' -H 'input: 044fb0f6-48de-4c27-b778-8cff8db7c67c' --data-binary 
