#!/usr/bin/env bash

set -e

if [[ "$1" = "aws" ]]; then
  echo "shell aws"
  git status -s terraform 
  for f in terraform/*; do
      if [ -d $f ]; then
          cd $f/pro
          CHANGES=$(git diff --name-only HEAD..HEAD~1 ../)
          if [ "$CHANGES" != "" ]; then
              echo "detected changes in $f. Running terraform apply..."
              terraform init -no-color
              terraform apply -no-color -auto-approve
              echo "$CHANGES"
          fi
          cd -
      fi
  done
elif [[ "$1" = "rules" ]]; then
  echo "shell rules"
else
  echo "Invalid argument"
  exit 1
fi