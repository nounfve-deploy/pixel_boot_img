#!/bin/sh

zipFile=$1
pattern=$2

if [ -z "$pattern" ]; then
  echo "missing pattern"
  exit -1
fi
echo "${zipFile} ${pattern}"

flist=$(unzip -o ${zipFile} "${pattern}"|grep -v "Archive:"|awk '{print $2}')
echo "${flist}"

echo "file list in one line"
flist=$(echo ${flist} | tr '\n' ' ')
echo "${flist}"
