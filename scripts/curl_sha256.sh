#!/bin/sh
url="https://fill-data.papermc.io/v1/objects/4a558a00005d33dafa4c4d5f9e47b3bd47d92311fceccd9c9754ee6b913f8649/paper-1.21.11-100.jar"

fname=$(curl --head --remote-name --write-out '%{filename_effective}\n' ${url})
curl -o ${fname} ${url}
sha_sum=$(sha256sum ./${fname} | awk '{print $1}')
echo ${sha_sum}