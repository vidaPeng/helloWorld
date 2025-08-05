#!/bin/bash

# ==============================================================================
# 【高兼容性版本】批量为GCP项目启用指定的API
#
# 使用方法:
# 1. 确保 gcloud CLI 已安装并已通过服务账号认证。
# 2. 赋予此脚本执行权限: chmod +x gcp_work.sh
# 3. 直接运行: ./gcp_work.sh 或 bash ./gcp_work.sh
# ==============================================================================

# 设置脚本在遇到错误时立即退出
set -e

# --- 需要启用的API服务名称列表 ---
APIS_TO_ENABLE=(
  "redis.googleapis.com"
  "compute.googleapis.com"
  "storage.googleapis.com"
)

echo "脚本将为所有可访问的活跃项目，检查并启用以下API:"
for api in "${APIS_TO_ENABLE[@]}"; do
  echo " - ${api}"
done
echo "--------------------------------------------------------"
read -p "确认继续吗? (y/n): " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "操作已取消。"
    exit 1
fi
echo "--------------------------------------------------------"


# --- 获取所有活跃的项目ID列表 ---
echo "正在获取所有活跃的项目列表..."
PROJECT_IDS=$(gcloud projects list --filter="lifecycleState=ACTIVE" --format="value(projectId)")

if [ -z "$PROJECT_IDS" ]; then
    echo "错误：未能获取到任何项目，请检查您的服务账号权限。"
    exit 1
fi

# --- 主循环：遍历所有项目和API ---
for project in $PROJECT_IDS; do
  echo ""
  echo "========================================================"
  echo "正在处理项目: ${project}"
  echo "========================================================"

  # ⭐⭐⭐【代码修改部分】⭐⭐⭐
  # 不再使用关联数组，而是将已启用的服务列表获取到一个长字符串中
  echo "  -> 正在查询当前已启用的API..."
  ENABLED_SERVICES_LIST=$(gcloud services list --project="${project}" --enabled --format="value(config.name)")

  # 遍历我们需要启用的API列表
  for api in "${APIS_TO_ENABLE[@]}"; do
    # 使用grep命令来检查API是否已经存在于列表中
    # grep -q 表示静默模式，只通过退出状态码告诉我们结果
    # ^${api}$ 表示精确匹配一整行，避免部分匹配 (例如 redis.com 匹配 my-redis.com)
    if echo "${ENABLED_SERVICES_LIST}" | grep -q "^${api}$"; then
      echo "  ✅ API '${api}' 在项目 '${project}' 中已经启用，跳过。"
    else
      echo "  🔥 API '${api}' 在项目 '${project}' 中未启用，现在开始启用..."
      # 执行启用命令
      gcloud services enable ${api} --project="${project}"
      echo "  ✅ API '${api}' 已成功为项目 '${project}' 启用。"
    fi
  done
done

echo ""
echo "========================================================"
echo "所有项目的API检查和启用任务已完成！"
echo "========================================================"