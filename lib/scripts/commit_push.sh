#!/bin/sh
git diff --quiet HEAD --

exit_status=$?
if [ $exit_status -ne 0 ]; then
    exit 0
fi

git add -A
git commit -m "Automated changes from workflow"
git push