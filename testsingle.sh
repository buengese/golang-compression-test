#!/bin/sh 
set -e

go build compress.go 

level=-2
./compress -in=$1 -out=* -stats -header=true -w=$2 -l=$level >>results.txt

level=1
compress -in=$1 -out=* -stats -header=false -w=$2 -l=$level >>results.txt

level=2
compress -in=$1 -out=* -stats -header=false -w=$2 -l=$level >>results.txt

level=3
compress -in=$1 -out=* -stats -header=false -w=$2 -l=$level >>results.txt

level=4
compress -in=$1 -out=* -stats -header=false -w=$2 -l=$level >>results.txt

level=5
compress -in=$1 -out=* -stats -header=false -w=$2 -l=$level >>results.txt

level=6
compress -in=$1 -out=* -stats -header=false -w=$2 -l=$level >>results.txt

level=7
compress -in=$1 -out=* -stats -header=false -w=$2 -l=$level >>results.txt

level=8
compress -in=%1 -out=* -stats -header=false -w=$2 -l=$level >>results.txt

level=9
compress -in=$1 -out=* -stats -header=false -w=$2 -l=$level >>results.txt
