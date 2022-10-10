#!/usr/bin/env bash

set -e

echo git status -s terraform | wc -l

for f in terraform/*; do
    if [ -d $f ]; then
        cd $f/pro
        # This takes into account we always use squash and this runs on push even
        CHANGES=$(git diff --name-only HEAD..HEAD~1 ../)
        # if there are any changes run terraform apply
        if [ "$CHANGES" != "" ]; then
            echo "detected changes in $f. Running terraform apply..."
            #terraform init -no-color
            #terraform apply -no-color -auto-approve
            echo "$CHANGES"
        fi
        cd -
    fi
done