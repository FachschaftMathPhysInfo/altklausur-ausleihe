#!/bin/bash -
#===============================================================================
#
#          FILE: query_prod.sh
#
#         USAGE: ./query_prod.sh
#
#        AUTHOR: Christian Heusel (christian@heusel.eu),
#  ORGANIZATION: Fachschaft MathPhysInfo
#       CREATED: 22/09/21 20:19
#
#===============================================================================

JWT_TOKEN="PASTE.YOUR.TOKEN"
UUID="81d6750d-a535-46b6-8491-accd6b1eb0c7"

set -o nounset                              # Treat unset variables as an error

# get exam list
curl --silent 'https://altklausuren.mathphys.info/query' \
    -H 'Accept-Encoding: gzip, deflate, br' -H 'Content-Type: application/json' -H 'Accept: application/json' \
    -H "Cookie: jwt=$JWT_TOKEN" \
    --data-binary '{"query":"query {\n  exams { UUID }\n}\n","variables":{}}' --compressed | \
    jq

# mark Exam
curl --silent 'https://altklausuren.mathphys.info/query' \
    -H 'Accept-Encoding: gzip, deflate, br' -H 'Content-Type: application/json' -H 'Accept: application/json' \
    -H "Cookie: jwt=$JWT_TOKEN" \
    --data-binary '{"query":"mutation requestMarkedExam {\n  requestMarkedExam(StringUUID: \"'$UUID'\",)\n}\n"}' --compressed | \
    jq

# get URL
curl --silent 'https://altklausuren.mathphys.info/query' \
    -H 'Accept-Encoding: gzip, deflate, br' -H 'Content-Type: application/json' -H 'Accept: application/json' \
    -H "Cookie: jwt=$JWT_TOKEN" \
    --data-binary '{"query":"query($UUID: String!) {\n  getExam(StringUUID: $UUID) {\n    viewUrl\n    downloadUrl\n  }\n}\n","variables":{"UUID":"'$UUID'"}}' --compressed | \
    jq
