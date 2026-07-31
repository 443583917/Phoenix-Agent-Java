#!/bin/bash
BASE="http://localhost:8066"
PASS=0
FAIL=0

ok() { echo "  ✅ $1"; PASS=$((PASS+1)); }
fail() { echo "  ❌ $1 (expected: $2, got: $3)"; FAIL=$((FAIL+1)); }
skip() { echo "  ⏭️  $1 (跳过: $2)"; PASS=$((PASS+1)); }

assert_code() { [ "$1" = "0" ] && ok "$2" || fail "$2" "code=0" "code=$1"; }

echo ""
echo "============================================="
echo "  Phoenix API 接口测试"
echo "============================================="

# ==================== 1. Health Check ====================
echo ""
echo "[1] 健康检查"
resp=$(curl -s "$BASE/echo")
code=$(echo "$resp" | grep -o '"code":[0-9]*' | head -1 | cut -d: -f2)
assert_code "$code" "GET /echo"

# ==================== 2. Auth ====================
echo ""
echo "[2] 验证码"
resp=$(curl -s "$BASE/api/privilege/auth/captcha")
code=$(echo "$resp" | grep -o '"code":[0-9]*' | head -1 | cut -d: -f2)
hasImage=$(echo "$resp" | grep -o '"image":"[^"]*"' | cut -d'"' -f4)
assert_code "$code" "GET /api/privilege/auth/captcha"
[ -n "$hasImage" ] && ok "验证码图片生成" || fail "验证码" "base64" "空"

echo ""
echo "[3] 登录"
resp=$(curl -s -X POST "$BASE/api/privilege/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"type":"account","username":"xtj","password":"123456"}')
code=$(echo "$resp" | grep -o '"code":[0-9]*' | head -1 | cut -d: -f2)
token=$(echo "$resp" | grep -o '"token":"[^"]*"' | cut -d'"' -f4)
assert_code "$code" "POST /api/privilege/auth/login"
[ -n "$token" ] && ok "JWT Token" || fail "Token" "非空" "空"

echo ""
echo "[4] 错误密码"
resp=$(curl -s -X POST "$BASE/api/privilege/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"type":"account","username":"xtj","password":"wrong"}')
code=$(echo "$resp" | grep -o '"code":[0-9]*' | head -1 | cut -d: -f2)
[ "$code" != "0" ] && ok "拒绝错误密码 (code=$code)" || fail "应拒绝" "非0" "0"

echo ""
echo "[5] 参数校验"
resp=$(curl -s -X POST "$BASE/api/privilege/auth/login" \
  -H "Content-Type: application/json" -d '{"type":"account"}')
code=$(echo "$resp" | grep -o '"code":[0-9]*' | head -1 | cut -d: -f2)
[ "$code" = "1001" ] && ok "缺少必填字段->1001" || fail "参数校验" "1001" "$code"

if [ -z "$token" ]; then
  echo "  ERROR: 未获取Token，无法继续测试"
  echo "  测试结果: PASS=$PASS FAIL=$FAIL"
  exit 1
fi

AUTH="Authorization: Bearer $token"

# ==================== 3. Authenticated Privilege ====================
echo ""
echo "[6] 当前用户信息"
resp=$(curl -s "$BASE/api/privilege/auth/getLoginUserInfo" -H "$AUTH")
code=$(echo "$resp" | grep -o '"code":[0-9]*' | head -1 | cut -d: -f2)
assert_code "$code" "GET /api/privilege/auth/getLoginUserInfo"

echo ""
echo "[7] 菜单"
resp=$(curl -s "$BASE/api/privilege/auth/menus" -H "$AUTH")
code=$(echo "$resp" | grep -o '"code":[0-9]*' | head -1 | cut -d: -f2)
assert_code "$code" "GET /api/privilege/auth/menus"

echo ""
echo "[8] 用户分页"
resp=$(curl -s "$BASE/api/privilege/user/page?page=1&size=10" -H "$AUTH")
code=$(echo "$resp" | grep -o '"code":[0-9]*' | head -1 | cut -d: -f2)
total=$(echo "$resp" | grep -o '"total":[0-9]*' | head -1 | cut -d: -f2)
assert_code "$code" "GET /api/privilege/user/page"
[ "$total" -gt 0 ] 2>/dev/null && ok "用户总数: $total" || fail "用户分页" "total>0" "$total"

echo ""
echo "[9] 用户详情"
resp=$(curl -s "$BASE/api/privilege/user/428011841386577921" -H "$AUTH")
code=$(echo "$resp" | grep -o '"code":[0-9]*' | head -1 | cut -d: -f2)
assert_code "$code" "GET /api/privilege/user/:id"

echo ""
echo "[10] 按用户名查询"
resp=$(curl -s "$BASE/api/privilege/user/username/xtj" -H "$AUTH")
code=$(echo "$resp" | grep -o '"code":[0-9]*' | head -1 | cut -d: -f2)
assert_code "$code" "GET /api/privilege/user/username/:username"

echo ""
echo "[11] 不存在的用户"
resp=$(curl -s "$BASE/api/privilege/user/999999999999999999" -H "$AUTH")
code=$(echo "$resp" | grep -o '"code":[0-9]*' | head -1 | cut -d: -f2)
[ "$code" != "0" ] && ok "不存在用户->错误 (code=$code)" || fail "不存在用户" "非0" "0"

echo ""
echo "[12] 角色分页(JSON body)"
resp=$(curl -s -X POST "$BASE/api/privilege/role/page" \
  -H "$AUTH" -H "Content-Type: application/json" \
  -d '{"page":1,"size":10}')
code=$(echo "$resp" | grep -o '"code":[0-9]*' | head -1 | cut -d: -f2)
assert_code "$code" "POST /api/privilege/role/page"

echo ""
echo "[13] 部门树"
resp=$(curl -s "$BASE/api/privilege/department/tree?companyId=428001009954172928" -H "$AUTH")
code=$(echo "$resp" | grep -o '"code":[0-9]*' | head -1 | cut -d: -f2)
assert_code "$code" "GET /api/privilege/department/tree"

echo ""
echo "[14] 公司分页(JSON body)"
resp=$(curl -s -X POST "$BASE/api/privilege/company/page" \
  -H "$AUTH" -H "Content-Type: application/json" \
  -d '{"page":1,"size":10}')
code=$(echo "$resp" | grep -o '"code":[0-9]*' | head -1 | cut -d: -f2)
assert_code "$code" "POST /api/privilege/company/page"

echo ""
echo "[15] 字典分页"
resp=$(curl -s "$BASE/api/privilege/dictionary/page?page=1&size=10" -H "$AUTH")
code=$(echo "$resp" | grep -o '"code":[0-9]*' | head -1 | cut -d: -f2)
assert_code "$code" "GET /api/privilege/dictionary/page"

echo ""
echo "[16] 登录日志分页"
resp=$(curl -s "$BASE/api/privilege/login-log/page?page=1&size=10" -H "$AUTH")
code=$(echo "$resp" | grep -o '"code":[0-9]*' | head -1 | cut -d: -f2)
assert_code "$code" "GET /api/privilege/login-log/page"

echo ""
echo "[17] 权限值分页(JSON body)"
resp=$(curl -s -X POST "$BASE/api/privilege/pvalue/page" \
  -H "$AUTH" -H "Content-Type: application/json" \
  -d '{"page":1,"size":10}')
code=$(echo "$resp" | grep -o '"code":[0-9]*' | head -1 | cut -d: -f2)
assert_code "$code" "POST /api/privilege/pvalue/page"

echo ""
echo "[18] 无Token->401"
resp=$(curl -s "$BASE/api/privilege/user/page")
code=$(echo "$resp" | grep -o '"code":[0-9]*' | head -1 | cut -d: -f2)
[ "$code" = "401" ] && ok "无Token->401" || fail "无Token" "401" "$code"

# ==================== 4. Platform ====================
echo ""
echo "[19] 平台登录"
resp=$(curl -s -X POST "$BASE/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"username":"lwj","password":"123456"}')
code=$(echo "$resp" | grep -o '"code":[0-9]*' | head -1 | cut -d: -f2)
platToken=$(echo "$resp" | grep -o '"token":"[^"]*"' | cut -d'"' -f4)
assert_code "$code" "POST /auth/login"
[ -n "$platToken" ] && ok "平台Token" || fail "平台Token" "非空" "空"

echo ""
echo "[20] 平台错误密码"
resp=$(curl -s -X POST "$BASE/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"username":"lwj","password":"wrong"}')
code=$(echo "$resp" | grep -o '"code":[0-9]*' | head -1 | cut -d: -f2)
[ "$code" != "0" ] && ok "拒绝错误密码 (code=$code)" || fail "应拒绝" "非0" "0"

if [ -n "$platToken" ]; then
  P_AUTH="Authorization: Bearer $platToken"
  echo ""
  echo "[21] 平台-租户分页"
  resp=$(curl -s "$BASE/platform/tenant-info/page?page=1&size=10" -H "$P_AUTH")
  code=$(echo "$resp" | grep -o '"code":[0-9]*' | head -1 | cut -d: -f2)
  assert_code "$code" "GET /platform/tenant-info/page"

  echo ""
  echo "[22] 平台-账号分页"
  resp=$(curl -s "$BASE/platform/account-info/page?page=1&size=10" -H "$P_AUTH")
  code=$(echo "$resp" | grep -o '"code":[0-9]*' | head -1 | cut -d: -f2)
  assert_code "$code" "GET /platform/account-info/page"

  echo ""
  echo "[23] 平台-分组分页"
  resp=$(curl -s "$BASE/platform/group-info/page?page=1&size=10" -H "$P_AUTH")
  code=$(echo "$resp" | grep -o '"code":[0-9]*' | head -1 | cut -d: -f2)
  assert_code "$code" "GET /platform/group-info/page"

  echo ""
  echo "[24] 平台-平台信息分页"
  resp=$(curl -s "$BASE/platform/platform-info/page?page=1&size=10" -H "$P_AUTH")
  code=$(echo "$resp" | grep -o '"code":[0-9]*' | head -1 | cut -d: -f2)
  assert_code "$code" "GET /platform/platform-info/page"
fi

# ==================== 5. Data Management ====================
echo ""
echo "[25] Agent分页"
resp=$(curl -s "$BASE/api/agent/page?page=1&size=10" -H "$AUTH")
code=$(echo "$resp" | grep -o '"code":[0-9]*' | head -1 | cut -d: -f2)
assert_code "$code" "GET /api/agent/page"

echo ""
echo "[26] Agent分类树"
resp=$(curl -s "$BASE/api/agent-category/tree" -H "$AUTH")
code=$(echo "$resp" | grep -o '"code":[0-9]*' | head -1 | cut -d: -f2)
assert_code "$code" "GET /api/agent-category/tree"

echo ""
echo "[27] 数据源类型"
resp=$(curl -s "$BASE/api/datasource/types" -H "$AUTH")
code=$(echo "$resp" | grep -o '"code":[0-9]*' | head -1 | cut -d: -f2)
assert_code "$code" "GET /api/datasource/types"

echo ""
echo "[28] 模型配置列表"
resp=$(curl -s "$BASE/api/model-config/list" -H "$AUTH")
code=$(echo "$resp" | grep -o '"code":[0-9]*' | head -1 | cut -d: -f2)
assert_code "$code" "GET /api/model-config/list"

echo ""
echo "[29] Prompt配置列表"
resp=$(curl -s "$BASE/api/prompt-config/list" -H "$AUTH")
code=$(echo "$resp" | grep -o '"code":[0-9]*' | head -1 | cut -d: -f2)
assert_code "$code" "GET /api/prompt-config/list"

echo ""
echo "[30] 语义模型列表"
resp=$(curl -s "$BASE/api/semantic-model/" -H "$AUTH")
code=$(echo "$resp" | grep -o '"code":[0-9]*' | head -1 | cut -d: -f2)
assert_code "$code" "GET /api/semantic-model/"

echo ""
echo "[31] 业务知识列表"
resp=$(curl -s "$BASE/api/business-knowledge/" -H "$AUTH")
code=$(echo "$resp" | grep -o '"code":[0-9]*' | head -1 | cut -d: -f2)
assert_code "$code" "GET /api/business-knowledge/"

echo ""
echo "[32] Agent知识分页"
resp=$(curl -s "$BASE/api/agent-knowledge/page?page=1&size=10" -H "$AUTH")
code=$(echo "$resp" | grep -o '"code":[0-9]*' | head -1 | cut -d: -f2)
assert_code "$code" "GET /api/agent-knowledge/page"

# ==================== 6. RAG & KG (tables may not exist) ====================
echo ""
echo "[33] RAG文件分页"
resp=$(curl -s "$BASE/api/rag/file/page?page=1&size=10" -H "$AUTH")
code=$(echo "$resp" | grep -o '"code":[0-9]*' | head -1 | cut -d: -f2)
assert_code "$code" "GET /api/rag/file/page"

echo ""
echo "[34] RAG分类分页"
resp=$(curl -s "$BASE/api/rag/category/page?page=1&size=10" -H "$AUTH")
code=$(echo "$resp" | grep -o '"code":[0-9]*' | head -1 | cut -d: -f2)
assert_code "$code" "GET /api/rag/category/page"

echo ""
echo "[35] KG实体分页"
resp=$(curl -s "$BASE/api/kg/entity/page?page=1&size=10" -H "$AUTH")
code=$(echo "$resp" | grep -o '"code":[0-9]*' | head -1 | cut -d: -f2)
assert_code "$code" "GET /api/kg/entity/page"

echo ""
echo "[36] KG领域分页"
resp=$(curl -s "$BASE/api/kg/domain/page?page=1&size=10" -H "$AUTH")
code=$(echo "$resp" | grep -o '"code":[0-9]*' | head -1 | cut -d: -f2)
assert_code "$code" "GET /api/kg/domain/page"

# ==================== 7. CRUD Tests ====================
TS=$(date +%s)

echo ""
echo "[37] 创建角色"
resp=$(curl -s -X POST "$BASE/api/privilege/role" \
  -H "$AUTH" -H "Content-Type: application/json" \
  -d '{"name":"test'$TS'","sn":"TEST'$TS'","companyId":0,"status":0}')
code=$(echo "$resp" | grep -o '"code":[0-9]*' | head -1 | cut -d: -f2)
assert_code "$code" "POST /api/privilege/role"

echo ""
echo "[38] 创建公司"
resp=$(curl -s -X POST "$BASE/api/privilege/company" \
  -H "$AUTH" -H "Content-Type: application/json" \
  -d '{"cname":"test'$TS'","code":"C'$TS'","status":0}')
code=$(echo "$resp" | grep -o '"code":[0-9]*' | head -1 | cut -d: -f2)
assert_code "$code" "POST /api/privilege/company"

echo ""
echo "[39] 创建部门"
resp=$(curl -s -X POST "$BASE/api/privilege/department" \
  -H "$AUTH" -H "Content-Type: application/json" \
  -d '{"name":"test'$TS'","code":"D'$TS'","companyId":"428001009954172928","pid":"0","status":0}')
code=$(echo "$resp" | grep -o '"code":[0-9]*' | head -1 | cut -d: -f2)
assert_code "$code" "POST /api/privilege/department"

echo ""
echo "[40] 创建字典"
resp=$(curl -s -X POST "$BASE/api/privilege/dictionary" \
  -H "$AUTH" -H "Content-Type: application/json" \
  -d '{"name":"test'$TS'","code":"DICT'$TS'","systemSn":"test","pcode":"0","sn":"SN'$TS'","sort":1}')
code=$(echo "$resp" | grep -o '"code":[0-9]*' | head -1 | cut -d: -f2)
assert_code "$code" "POST /api/privilege/dictionary"

echo ""
echo "[41] 创建用户"
resp=$(curl -s -X POST "$BASE/api/privilege/user" \
  -H "$AUTH" -H "Content-Type: application/json" \
  -d '{"username":"utest'$TS'","code":"UC'$TS'","realName":"Test","companyId":"428001009954172928","deptId":"428003302434869248","status":0}')
code=$(echo "$resp" | grep -o '"code":[0-9]*' | head -1 | cut -d: -f2)
assert_code "$code" "POST /api/privilege/user"

# ==================== Summary ====================
echo ""
echo "============================================="
echo "  测试结果: PASS=$PASS  FAIL=$FAIL"
echo "============================================="

[ "$FAIL" -gt 0 ] && exit 1 || exit 0
