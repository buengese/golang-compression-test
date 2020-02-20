#!/bin/sh
set -e

go build compress.go

level=1

./compress -in=$1 -out=* -stats -header=true -w="gzstd" -l=$level >>results.txt
./compress -in=$1 -out=* -stats -header=false -w="gzkp" -l=$level >>results.txt
./compress -in=$1 -out=* -stats -header=false -w="pgzip" -l=$level >>results.txt
./compress -in=$1 -out=* -stats -header=false -w="bgzf" -l=$level >>results.txt
./compress -in=$1 -out=* -stats -header=false -w="cgzip" -l=$level >>results.txt
./compress -in=$1 -out=* -stats -header=false -w="clzma" -l=$level >>results.txt

level=2
echo "." >> results.txt

./compress -in=$1 -out=* -stats -header=false -w="gzstd" -l=$level >>results.txt
./compress -in=$1 -out=* -stats -header=false -w="gzkp" -l=$level >>results.txt
./compress -in=$1 -out=* -stats -header=false -w="pgzip" -l=$level >>results.txt
./compress -in=$1 -out=* -stats -header=false -w="bgzf" -l=$level >>results.txt
./compress -in=$1 -out=* -stats -header=false -w="cgzip" -l=$level >>results.txt
./compress -in=$1 -out=* -stats -header=false -w="clzma" -l=$level >>results.txt

level=3
echo "." >> results.txt

./compress -in=$1 -out=* -stats -header=false -w="gzstd" -l=$level >>results.txt
./compress -in=$1 -out=* -stats -header=false -w="gzkp" -l=$level >>results.txt
./compress -in=$1 -out=* -stats -header=false -w="pgzip" -l=$level >>results.txt
./compress -in=$1 -out=* -stats -header=false -w="bgzf" -l=$level >>results.txt
./compress -in=$1 -out=* -stats -header=false -w="cgzip" -l=$level >>results.txt
./compress -in=$1 -out=* -stats -header=false -w="clzma" -l=$level >>results.txt

level=4
echo "." >> results.txt

./compress -in=$1 -out=* -stats -header=false -w="gzstd" -l=$level >>results.txt
./compress -in=$1 -out=* -stats -header=false -w="gzkp" -l=$level >>results.txt
./compress -in=$1 -out=* -stats -header=false -w="pgzip" -l=$level >>results.txt
./compress -in=$1 -out=* -stats -header=false -w="bgzf" -l=$level >>results.txt
./compress -in=$1 -out=* -stats -header=false -w="cgzip" -l=$level >>results.txt
./compress -in=$1 -out=* -stats -header=false -w="clzma" -l=$level >>results.txt

level=5
echo "." >> results.txt

./compress -in=$1 -out=* -stats -header=false -w="gzstd" -l=$level >>results.txt
./compress -in=$1 -out=* -stats -header=false -w="gzkp" -l=$level >>results.txt
./compress -in=$1 -out=* -stats -header=false -w="pgzip" -l=$level >>results.txt
./compress -in=$1 -out=* -stats -header=false -w="bgzf" -l=$level >>results.txt
./compress -in=$1 -out=* -stats -header=false -w="cgzip" -l=$level >>results.txt
./compress -in=$1 -out=* -stats -header=false -w="clzma" -l=$level >>results.txt

level=6
echo "." >> results.txt

./compress -in=$1 -out=* -stats -header=false -w="gzstd" -l=$level >>results.txt
./compress -in=$1 -out=* -stats -header=false -w="gzkp" -l=$level >>results.txt
./compress -in=$1 -out=* -stats -header=false -w="pgzip" -l=$level >>results.txt
./compress -in=$1 -out=* -stats -header=false -w="bgzf" -l=$level >>results.txt
./compress -in=$1 -out=* -stats -header=false -w="cgzip" -l=$level >>results.txt
./compress -in=$1 -out=* -stats -header=false -w="clzma" -l=$level >>results.txt

level=7
echo "." >> results.txt

./compress -in=$1 -out=* -stats -header=false -w="gzstd" -l=$level >>results.txt
./compress -in=$1 -out=* -stats -header=false -w="gzkp" -l=$level >>results.txt
./compress -in=$1 -out=* -stats -header=false -w="pgzip" -l=$level >>results.txt
./compress -in=$1 -out=* -stats -header=false -w="bgzf" -l=$level >>results.txt
./compress -in=$1 -out=* -stats -header=false -w="cgzip" -l=$level >>results.txt
./compress -in=$1 -out=* -stats -header=false -w="clzma" -l=$level >>results.txt

level=8
echo "." >> results.txt

./compress -in=$1 -out=* -stats -header=false -w="gzstd" -l=$level >>results.txt
./compress -in=$1 -out=* -stats -header=false -w="gzkp" -l=$level >>results.txt
./compress -in=$1 -out=* -stats -header=false -w="pgzip" -l=$level >>results.txt
./compress -in=$1 -out=* -stats -header=false -w="bgzf" -l=$level >>results.txt
./compress -in=$1 -out=* -stats -header=false -w="cgzip" -l=$level >>results.txt
./compress -in=$1 -out=* -stats -header=false -w="clzma" -l=$level >>results.txt

level=9
echo "." >> results.txt

./compress -in=$1 -out=* -stats -header=false -w="gzstd" -l=$level >>results.txt
./compress -in=$1 -out=* -stats -header=false -w="gzkp" -l=$level >>results.txt
./compress -in=$1 -out=* -stats -header=false -w="pgzip" -l=$level >>results.txt
./compress -in=$1 -out=* -stats -header=false -w="bgzf" -l=$level >>results.txt
./compress -in=$1 -out=* -stats -header=false -w="cgzip" -l=$level >>results.txt
./compress -in=$1 -out=* -stats -header=false -w="clzma" -l=$level >>results.txt


level=-2
echo "." >> results.txt

./compress -in=$1 -out=* -stats -header=false -w="gzkp" -l=$level >>results.txt
./compress -in=$1 -out=* -stats -header=false -w="pgzip" -l=$level >>results.txt

level=0 # no level supported beyond here
echo "." >> results.txt
./compress -in=$1 -out=* -stats -header=false -w="snappy" -l=$level >>results.txt
./compress -in=$1 -out=* -stats -header=false -w="lz4" -l=$level >>results.txt
./compress -in=$1 -out=* -stats -header=false -w="lz4p" -l=$level >>results.txt
./compress -in=$1 -out=* -stats -header=false -w="clz4" -l=$level >>results.txt
./compress -in=$1 -out=* -stats -header=false -w="zstdfast" -l=$level >>results.txt
./compress -in=$1 -out=* -stats -header=false -w="zstddefault" -l=$level >>results.txt
./compress -in=$1 -out=* -stats -header=false -w="gxz" -l=$level >>results.txt
./compress -in=$1 -out=* -stats -header=false -w="gxz2" -l=$level >>results.txt
echo "." >> results.txt