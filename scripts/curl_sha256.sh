#!/bin/sh
url="https://www.minecraft.net/bedrockdedicatedserver/bin-linux/bedrock-server-1.21.132.3.zip"

fname=$(curl --http1.1 --head --remote-name --write-out '%{filename_effective}\n' ${url})
curl --http1.1 -o ${fname} ${url}
sha_sum=$(sha256sum ./${fname} | awk '{print $1}')
echo ${sha_sum}