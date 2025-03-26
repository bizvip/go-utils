#!/bin/bash

# 获取所有以 v1.5 开头的标签
tags=$(git tag | grep '^v1.3')

# 删除本地标签
for tag in $tags; do
    git tag -d $tag
done

# 推送删除操作到远程仓库
for tag in $tags; do
    git push origin --delete $tag
done
