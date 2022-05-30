NEW_VERSION=$1

echo "curent path: $(pwd), change version to $NEW_VERSION"

          git branch
          git config --global user.name "tscuite"
          git config --global user.email "tscuite@qq.com"
          git checkout -b pr@20220530@${NEW_VERSION}
          git remote set-url origin --push --add 'https://tscuite:$cdzGITHUB_TOKEN@github.com/tscuite/tscuite.git'
          git push --set-upstream origin pr@20220530@${NEW_VERSION}
