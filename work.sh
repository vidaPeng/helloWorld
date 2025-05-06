#!/bin/bash

# 配置你的 APISIX 地址和管理密钥
# OCI TEST
#APISIX_HOST="http://test-wus-api.gw.int.pixocial.com"
#API_KEY="edd1c9f034kmFId7Gd5RD7RIMiJv0IgL"

# OCI prod
#APISIX_HOST="http://base-wus-api.gw.int.pixocial.com"
#API_KEY="eddd5f136fkj5f136f5f1387ad84b628f1"

# 1. 获取所有 upstream 列表
echo "Fetching all upstreams..."
response=$(curl -s -H "X-API-KEY: ${API_KEY}" "${APISIX_HOST}/apisix/admin/upstreams")

# 2. 解析出所有 upstream 的 ID
# 返回的 JSON key 是 /apisix/upstreams/{id}，我们用 jq 处理
# 如果你没有安装 jq，可以用简单 grep+sed 方式
upstream_ids=$(echo "$response" | jq -r '.node.nodes[].key' | awk -F'/' '{print $4}')

# 检查有没有找到 upstream
if [[ -z "$upstream_ids" ]]; then
    echo "No upstreams found."
    exit 0
fi

echo "Found upstream IDs:"
echo "$upstream_ids"

# 3. 循环删除每个 upstream
for id in $upstream_ids; do
    echo "Deleting upstream ID: $id"
    curl -s -X DELETE -H "X-API-KEY: ${API_KEY}" "${APISIX_HOST}/apisix/admin/upstreams/${id}"
    echo "Deleted upstream ID: $id"
done

echo "All upstreams deleted."

