#/bin/sh

gitRoot=$(git rev-parse --show-toplevel)
export gitRoot="${gitRoot}"

cd ${gitRoot}/working_dir/
exec go run .. $@