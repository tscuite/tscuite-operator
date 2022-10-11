#!/usr/bin/env bash

set -e

if [[ "$1" = "aws" ]]; then
  echo "shell aws"
elif [[ "$1" = "rules" ]]; then
  echo "shell rules"
else
  echo "Invalid argument"
  exit 1
fi

# //如果修改了aws, 运行
# git status -s terraform 

# for f in terraform/*; do
#     if [ -d $f ]; then
#         cd $f/pro
#         # This takes into account we always use squash and this runs on push even
#         CHANGES=$(git diff --name-only HEAD..HEAD~1 ../)
#         # if there are any changes run terraform apply
#         if [ "$CHANGES" != "" ]; then
#             echo "detected changes in $f. Running terraform apply..."
#             #terraform init -no-color
#             #terraform apply -no-color -auto-approve
#             echo "$CHANGES"
#             ../../selefra init
#         fi
#         cd -
#     fi
# done


# //如果修改了规则, 运行



# for f in 规则/*; do
#     if [ -d $f ]; then
#         cd $f/
#         CHANGES=$(git diff --name-only HEAD..HEAD~1 ../)
#         if [ "$CHANGES" != "" ]; then
#             echo "detected changes in $f. Running terraform apply..."
#             echo "$CHANGES"
#             ../../selefra init
#         fi
#         cd -
#     fi
# done