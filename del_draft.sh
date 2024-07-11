#!/bin/bash

# GitHub 用户名和仓库名
USERNAME="bizvip"
REPO="go-utils"
TOKEN="*******"

# 获取所有草稿版本的ID
draft_releases=$(curl -s -H "Authorization: token $TOKEN" "https://api.github.com/repos/$USERNAME/$REPO/releases" | jq '.[] | select(.draft == true) | .id')

# 删除每个草稿版本
for release_id in $draft_releases; do
  echo "Deleting draft release ID: $release_id"
  curl -s -X DELETE -H "Authorization: token $TOKEN" "https://api.github.com/repos/$USERNAME/$REPO/releases/$release_id"
done