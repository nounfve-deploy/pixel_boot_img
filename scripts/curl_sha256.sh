#!/bin/sh
url="https://dl.google.com/dl/android/aosp/akita-bp4a.260105.004.e1-factory-d9bc5fb8.zip"

fname=$(curl --head --remote-name --write-out '%{filename_effective}\n' ${url})
curl -o ${fname} ${url}
sha_sum=$(sha256sum ./${fname} | awk '{print $1}')
echo ${fname}
echo ${sha_sum}