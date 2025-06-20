
function updateMesh() {
cluster=$1
namespace=$2
service=$3
status=$4
if [[ -z $1 || -z $2 || -z $3 || -z $4 ]]; then
  echo "Usage: updateMes.sh <cluster> <namespace> <service> <status>"
  echo "Example: "
  echo "开启: updateMesh.sh oci-wus-k8s-base-test-01 test-xdh demo-goweb true"
  echo "关闭: updateMesh.sh oci-wus-k8s-base-test-01 test-xdh demo-goweb false"
  return 1
fi

if [[ $status == "true" ]];then
status=1
else
status=0
fi
echo "status: $status"
bk_token='bkcrypt%24gAAAAABoUVJoyZt11G383igAMzbF-JldBauCzQ8aJ3o8sqyv-Su9mDWWnVuZ-TiWa8_qd8E_ocyWye658ic58eBtl_FE5uILs1qQ0WrV5cyWrLErdtL3D-E%3D'

curl "https://beta-pcs-api.bkce7.int.pixocial.com/api/pcs/clusters/${cluster}/namespaces/${namespace}/services/${service}/mesh/status" \
  -X 'PUT' \
  -H 'accept: application/json, text/plain, */*' \
  -H 'accept-language: zh-CN,zh;q=0.9,en;q=0.8' \
  -H 'content-type: application/json' \
  -H "bk_token: ${bk_token}" \
  -H 'dnt: 1' \
  -H 'origin: https://test-ocean.int.pixocial.com' \
  -H 'priority: u=1, i' \
  -H 'referer: https://test-ocean.int.pixocial.com/' \
  -H 'sec-ch-ua: "Google Chrome";v="137", "Chromium";v="137", "Not/A)Brand";v="24"' \
  -H 'sec-ch-ua-mobile: ?0' \
  -H 'sec-ch-ua-platform: "macOS"' \
  -H 'sec-fetch-dest: empty' \
  -H 'sec-fetch-mode: cors' \
  -H 'sec-fetch-site: same-site' \
  -H 'user-agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36' \
  --data-raw "{\"status\":${status}}"
}

updateMesh $@