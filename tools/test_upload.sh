#!/bin/bash - 
#===============================================================================
#
#          FILE: test_upload.sh
# 
#         USAGE: ./test_upload.sh <input-file>
# 
#        AUTHOR: Christian Heusel (christian@heusel.eu), 
#  ORGANIZATION: Fachschaft MathPhysInfo
#       CREATED: 12/16/2021 15:00
# 
#===============================================================================

set -o nounset                              # Treat unset variables as an error

# check if command line argument is empty or not present
if [ "$1" == "" ] || [ $# -ne 1 ]; then
    echo "Wrong amount of arguments!"
    echo "USAGE: ./test_upload.sh <input-file>"
    echo "Given:" $@
    exit 0
fi

TARGET_HOST='http://localhost:8080/query'
# TARGET_HOST='https://altklausuren.mathphys.info/query'
JWT_TOKEN="YOUR.JWT.TOKEN"
INPUT_FILENAME=$1

./upload_one.sh $JWT_TOKEN '{ "subject": "Info", "moduleName": "Betriebssysteme und Netzwerke", "moduleAltName": "IBN, BeNe", "year": 2021, "semester": "SoSe", "examiners": "Helene Fischer", "file": null }' $INPUT_FILENAME
