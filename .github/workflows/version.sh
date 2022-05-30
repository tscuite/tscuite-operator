git branch
          git config --global user.name "tscuite"
          git config --global user.email "tscuite@qq.com"
          git checkout -b pr@20220530@${{github.run_number}}
          git remote set-url origin --push --add 'https://tscuite:$cdzGITHUB_TOKEN@github.com/tscuite/tscuite.git'
          git push --set-upstream origin pr@20220530@${{github.run_number}}
