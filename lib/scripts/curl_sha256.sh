#!/bin/sh

url=$1
if [ -z "$url" ]; then
  echo "empty-url"
  exit -1
fi

echo "[Url] ${url}"
if [ -z "${Mocking}"];then
  fname=$(curl --head --remote-name --write-out '%{filename_effective}\n' ${url})
  curl -o ${fname} ${url}

  sha_sum=$(sha256sum ./${fname} | awk '{print $1}')
else
  fname="akita-bp3a.250905.014-factory-47caa9a7.zip"
  sha_sum="47caa9a7e39670689f1746b7cf6c4770906eb69eff6977c5f03573fa62068697"
fi

echo ===============================================
echo ${fname} ${sha_sum}
