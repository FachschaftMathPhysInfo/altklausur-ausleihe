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

set -o nounset                              # Treat unset variables as an error

# mark Exam
curl 'https://altklausuren.mathphys.info/query' -H 'Accept-Encoding: gzip, deflate, br' -H 'Content-Type: application/json' -H 'Accept: application/json' -H "Cookie: jwt=$JWT_TOKEN" --data-binary '{"query":"mutation requestMarkedExam {\n  requestMarkedExam(StringUUID: \"85975cf6-b3da-4ef0-9e66-7ae44d47635e\",)\n}\n"}' --compressed
# get URL
curl 'https://altklausuren.mathphys.info/query' -H 'Accept-Encoding: gzip, deflate, br' -H 'Content-Type: application/json' -H 'Accept: application/json' -H "Cookie: jwt=$JWT_TOKEN" --data-binary '{"query":"query($UUID: String!) {\n  getExam(StringUUID: $UUID) {\n    viewUrl\n    downloadUrl\n  }\n}\n","variables":{"UUID":"81d6750d-a535-46b6-8491-accd6b1eb0c7"}}' --compressed
