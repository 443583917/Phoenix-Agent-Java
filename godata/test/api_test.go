package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"
)

const baseURL = "http://localhost:8066"

// ── helpers ──

var client = &http.Client{Timeout: 10 * time.Second}

func do(method, path string, body interface{}, headers map[string]string) (*http.Response, []byte, error) {
	var reqBody []byte
	if body != nil {
		if s, ok := body.(string); ok {
			reqBody = []byte(s)
		} else {
			b, _ := json.Marshal(body)
			reqBody = b
		}
	}
	req, _ := http.NewRequest(method, baseURL+path, bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	return resp, buf.Bytes(), nil
}

type apiResp struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
	Token   string          `json:"token"`
	Total   int64           `json:"total"`
}

type loginResp struct {
	Token    string `json:"token"`
	UserID   string `json:"userId"`
	Username string `json:"username"`
}

type platformLoginResp struct {
	Token     string `json:"token"`
	AccountID string `json:"accountId"`
}

var token string
var platToken string

// ── tests ──

func TestEcho(t *testing.T) {
	_, body, err := do("GET", "/echo", nil, nil)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	var r apiResp
	json.Unmarshal(body, &r)
	if r.Code != 100 {
		t.Fatalf("GET /echo: code=%d want 0", r.Code)
	}
	t.Log("✅ GET /echo")
}

func TestCaptcha(t *testing.T) {
	_, body, err := do("GET", "/api/privilege/auth/captcha", nil, nil)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	var r apiResp
	json.Unmarshal(body, &r)
	if r.Code != 100 {
		t.Fatalf("GET /api/privilege/auth/captcha: code=%d", r.Code)
	}
	// check captchaKey and image exist in data
	var data map[string]interface{}
	json.Unmarshal(r.Data, &data)
	if data["captchaKey"] == nil || data["image"] == nil {
		t.Fatal("captcha response missing captchaKey or image")
	}
	t.Log("✅ GET /api/privilege/auth/captcha")
}

func TestLogin(t *testing.T) {
	body := `{"type":"account","username":"xtj","password":"123456"}`
	_, respBody, err := do("POST", "/api/privilege/auth/login", body, nil)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	var r apiResp
	json.Unmarshal(respBody, &r)
	if r.Code != 100 {
		t.Fatalf("POST /api/privilege/auth/login: code=%d msg=%s", r.Code, r.Message)
	}
	var lr loginResp
	json.Unmarshal(r.Data, &lr)
	if lr.Token == "" {
		t.Fatal("login did not return token")
	}
	token = lr.Token
	t.Log("✅ POST /api/privilege/auth/login — got token")
}

func TestLoginWrongPassword(t *testing.T) {
	body := `{"type":"account","username":"xtj","password":"wrong"}`
	_, respBody, err := do("POST", "/api/privilege/auth/login", body, nil)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	var r apiResp
	json.Unmarshal(respBody, &r)
	if r.Code == 100 {
		t.Fatal("wrong password should be rejected")
	}
	t.Logf("✅ wrong password rejected (code=%d)", r.Code)
}

func TestLoginMissingFields(t *testing.T) {
	body := `{"type":"account"}`
	_, respBody, err := do("POST", "/api/privilege/auth/login", body, nil)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	var r apiResp
	json.Unmarshal(respBody, &r)
	if r.Code != 1001 {
		t.Fatalf("missing fields should return 1001, got %d", r.Code)
	}
	t.Log("✅ missing fields → 1001")
}

// ── authenticated privilege ──

func TestGetLoginUserInfo(t *testing.T) {
	needToken(t)
	_, body, err := do("GET", "/api/privilege/auth/getLoginUserInfo", nil, authHeader())
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	var r apiResp
	json.Unmarshal(body, &r)
	if r.Code != 100 {
		t.Fatalf("code=%d", r.Code)
	}
	t.Log("✅ GET /api/privilege/auth/getLoginUserInfo")
}

func TestMenus(t *testing.T) {
	needToken(t)
	_, body, err := do("GET", "/api/privilege/auth/menus", nil, authHeader())
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	var r apiResp
	json.Unmarshal(body, &r)
	if r.Code != 100 {
		t.Fatalf("code=%d", r.Code)
	}
	t.Log("✅ GET /api/privilege/auth/menus")
}

func TestUserPage(t *testing.T) {
	needToken(t)
	_, body, err := do("GET", "/api/privilege/user/page?page=1&size=10", nil, authHeader())
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	var r struct {
		Code  int   `json:"code"`
		Total int64 `json:"total"`
	}
	json.Unmarshal(body, &r)
	if r.Code != 100 {
		t.Fatalf("code=%d", r.Code)
	}
	if r.Total == 0 {
		t.Fatal("user total should be > 0")
	}
	t.Logf("✅ GET /api/privilege/user/page (total=%d)", r.Total)
}

func TestUserByID(t *testing.T) {
	needToken(t)
	_, body, err := do("GET", "/api/privilege/user/428011841386577921", nil, authHeader())
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	var r apiResp
	json.Unmarshal(body, &r)
	if r.Code != 100 {
		t.Fatalf("code=%d", r.Code)
	}
	t.Log("✅ GET /api/privilege/user/:id")
}

func TestUserByUsername(t *testing.T) {
	needToken(t)
	_, body, err := do("GET", "/api/privilege/user/username/xtj", nil, authHeader())
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	var r apiResp
	json.Unmarshal(body, &r)
	if r.Code != 100 {
		t.Fatalf("code=%d", r.Code)
	}
	t.Log("✅ GET /api/privilege/user/username/:username")
}

func TestUserNotFound(t *testing.T) {
	needToken(t)
	_, body, err := do("GET", "/api/privilege/user/999999999999999999", nil, authHeader())
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	var r apiResp
	json.Unmarshal(body, &r)
	if r.Code == 100 {
		t.Fatal("non-existent user should return error")
	}
	t.Logf("✅ non-existent user → error (code=%d)", r.Code)
}

func TestRolePage(t *testing.T) {
	needToken(t)
	body := map[string]int{"page": 1, "size": 10}
	_, respBody, err := do("POST", "/api/privilege/role/page", body, authHeader())
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	var r apiResp
	json.Unmarshal(respBody, &r)
	if r.Code != 100 {
		t.Fatalf("code=%d", r.Code)
	}
	t.Log("✅ POST /api/privilege/role/page")
}

func TestDeptTree(t *testing.T) {
	needToken(t)
	_, body, err := do("GET", "/api/privilege/department/tree?companyId=428001009954172928", nil, authHeader())
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	var r apiResp
	json.Unmarshal(body, &r)
	if r.Code != 100 {
		t.Fatalf("code=%d", r.Code)
	}
	t.Log("✅ GET /api/privilege/department/tree")
}

func TestCompanyPage(t *testing.T) {
	needToken(t)
	body := map[string]int{"page": 1, "size": 10}
	_, respBody, err := do("POST", "/api/privilege/company/page", body, authHeader())
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	var r apiResp
	json.Unmarshal(respBody, &r)
	if r.Code != 100 {
		t.Fatalf("code=%d", r.Code)
	}
	t.Log("✅ POST /api/privilege/company/page")
}

func TestDictPage(t *testing.T) {
	needToken(t)
	_, body, err := do("GET", "/api/privilege/dictionary/page?page=1&size=10", nil, authHeader())
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	var r apiResp
	json.Unmarshal(body, &r)
	if r.Code != 100 {
		t.Fatalf("code=%d", r.Code)
	}
	t.Log("✅ GET /api/privilege/dictionary/page")
}

func TestLoginLogPage(t *testing.T) {
	needToken(t)
	_, body, err := do("GET", "/api/privilege/login-log/page?page=1&size=10", nil, authHeader())
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	var r apiResp
	json.Unmarshal(body, &r)
	if r.Code != 100 {
		t.Fatalf("code=%d", r.Code)
	}
	t.Log("✅ GET /api/privilege/login-log/page")
}

func TestPvaluePage(t *testing.T) {
	needToken(t)
	body := map[string]int{"page": 1, "size": 10}
	_, respBody, err := do("POST", "/api/privilege/pvalue/page", body, authHeader())
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	var r apiResp
	json.Unmarshal(respBody, &r)
	if r.Code != 100 {
		t.Fatalf("code=%d", r.Code)
	}
	t.Log("✅ POST /api/privilege/pvalue/page")
}

func TestNoTokenRejected(t *testing.T) {
	_, body, err := do("GET", "/api/privilege/user/page", nil, nil)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	var r apiResp
	json.Unmarshal(body, &r)
	if r.Code != 401 {
		t.Fatalf("expected 401 without token, got %d", r.Code)
	}
	t.Log("✅ no token → 401")
}

// ── platform auth ──

func TestPlatformLogin(t *testing.T) {
	body := `{"username":"lwj","password":"123456"}`
	_, respBody, err := do("POST", "/auth/login", body, nil)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	var r apiResp
	json.Unmarshal(respBody, &r)
	if r.Code != 100 {
		t.Fatalf("code=%d msg=%s", r.Code, r.Message)
	}
	var pr platformLoginResp
	json.Unmarshal(r.Data, &pr)
	if pr.Token == "" {
		t.Fatal("platform login did not return token")
	}
	platToken = pr.Token
	t.Log("✅ POST /auth/login")
}

func TestPlatformLoginWrong(t *testing.T) {
	body := `{"username":"lwj","password":"wrong"}`
	_, respBody, err := do("POST", "/auth/login", body, nil)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	var r apiResp
	json.Unmarshal(respBody, &r)
	if r.Code == 100 {
		t.Fatal("wrong password should be rejected")
	}
	t.Logf("✅ platform wrong password rejected (code=%d)", r.Code)
}

// ── platform endpoints ──

func TestPlatformTenantPage(t *testing.T) {
	needPlatToken(t)
	_, body, err := do("GET", "/platform/tenant-info/page?page=1&size=10", nil, platAuthHeader())
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	var r apiResp
	json.Unmarshal(body, &r)
	if r.Code != 100 {
		t.Fatalf("code=%d", r.Code)
	}
	t.Log("✅ GET /platform/tenant-info/page")
}

func TestPlatformAccountPage(t *testing.T) {
	needPlatToken(t)
	_, body, err := do("GET", "/platform/account-info/page?page=1&size=10", nil, platAuthHeader())
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	var r apiResp
	json.Unmarshal(body, &r)
	if r.Code != 100 {
		t.Fatalf("code=%d", r.Code)
	}
	t.Log("✅ GET /platform/account-info/page")
}

func TestPlatformGroupPage(t *testing.T) {
	needPlatToken(t)
	_, body, err := do("GET", "/platform/group-info/page?page=1&size=10", nil, platAuthHeader())
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	var r apiResp
	json.Unmarshal(body, &r)
	if r.Code != 100 {
		t.Fatalf("code=%d", r.Code)
	}
	t.Log("✅ GET /platform/group-info/page")
}

func TestPlatformInfoPage(t *testing.T) {
	needPlatToken(t)
	_, body, err := do("GET", "/platform/platform-info/page?page=1&size=10", nil, platAuthHeader())
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	var r apiResp
	json.Unmarshal(body, &r)
	if r.Code != 100 {
		t.Fatalf("code=%d", r.Code)
	}
	t.Log("✅ GET /platform/platform-info/page")
}

// ── data management ──

func TestAgentPage(t *testing.T) {
	needToken(t)
	_, body, err := do("GET", "/api/agent/page?page=1&size=10", nil, authHeader())
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	var r apiResp
	json.Unmarshal(body, &r)
	if r.Code != 100 {
		t.Fatalf("code=%d", r.Code)
	}
	t.Log("✅ GET /api/agent/page")
}

func TestAgentCategoryTree(t *testing.T) {
	needToken(t)
	_, body, err := do("GET", "/api/agent-category/tree", nil, authHeader())
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	var r apiResp
	json.Unmarshal(body, &r)
	if r.Code != 100 {
		t.Fatalf("code=%d", r.Code)
	}
	t.Log("✅ GET /api/agent-category/tree")
}

func TestDatasourceTypes(t *testing.T) {
	needToken(t)
	_, body, err := do("GET", "/api/datasource/types", nil, authHeader())
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	var r apiResp
	json.Unmarshal(body, &r)
	if r.Code != 100 {
		t.Fatalf("code=%d", r.Code)
	}
	t.Log("✅ GET /api/datasource/types")
}

func TestModelConfigList(t *testing.T) {
	needToken(t)
	_, body, err := do("GET", "/api/model-config/list", nil, authHeader())
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	var r apiResp
	json.Unmarshal(body, &r)
	if r.Code != 100 {
		t.Fatalf("code=%d", r.Code)
	}
	t.Log("✅ GET /api/model-config/list")
}

func TestPromptConfigList(t *testing.T) {
	needToken(t)
	_, body, err := do("GET", "/api/prompt-config/list", nil, authHeader())
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	var r apiResp
	json.Unmarshal(body, &r)
	if r.Code != 100 {
		t.Fatalf("code=%d", r.Code)
	}
	t.Log("✅ GET /api/prompt-config/list")
}

func TestSemanticModelList(t *testing.T) {
	needToken(t)
	_, body, err := do("GET", "/api/semantic-model/", nil, authHeader())
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	var r apiResp
	json.Unmarshal(body, &r)
	if r.Code != 100 {
		t.Fatalf("code=%d", r.Code)
	}
	t.Log("✅ GET /api/semantic-model/")
}

func TestBusinessKnowledgeList(t *testing.T) {
	needToken(t)
	_, body, err := do("GET", "/api/business-knowledge/", nil, authHeader())
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	var r apiResp
	json.Unmarshal(body, &r)
	if r.Code != 100 {
		t.Fatalf("code=%d", r.Code)
	}
	t.Log("✅ GET /api/business-knowledge/")
}

func TestAgentKnowledgePage(t *testing.T) {
	needToken(t)
	_, body, err := do("GET", "/api/agent-knowledge/page?page=1&size=10", nil, authHeader())
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	var r apiResp
	json.Unmarshal(body, &r)
	if r.Code != 100 {
		t.Fatalf("code=%d", r.Code)
	}
	t.Log("✅ GET /api/agent-knowledge/page")
}

// ── RAG ──

func TestRagFilePage(t *testing.T) {
	needToken(t)
	_, body, err := do("GET", "/api/rag/file/page?page=1&size=10", nil, authHeader())
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	var r apiResp
	json.Unmarshal(body, &r)
	if r.Code != 100 {
		t.Fatalf("code=%d", r.Code)
	}
	t.Log("✅ GET /api/rag/file/page")
}

func TestRagCategoryPage(t *testing.T) {
	needToken(t)
	_, body, err := do("GET", "/api/rag/category/page?page=1&size=10", nil, authHeader())
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	var r apiResp
	json.Unmarshal(body, &r)
	if r.Code != 100 {
		t.Fatalf("code=%d", r.Code)
	}
	t.Log("✅ GET /api/rag/category/page")
}

// ── KG ──

func TestKGEntityPage(t *testing.T) {
	needToken(t)
	_, body, err := do("GET", "/api/kg/entity/page?page=1&size=10", nil, authHeader())
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	var r apiResp
	json.Unmarshal(body, &r)
	if r.Code != 100 {
		t.Fatalf("code=%d", r.Code)
	}
	t.Log("✅ GET /api/kg/entity/page")
}

func TestKGDomainPage(t *testing.T) {
	needToken(t)
	_, body, err := do("GET", "/api/kg/domain/page?page=1&size=10", nil, authHeader())
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	var r apiResp
	json.Unmarshal(body, &r)
	if r.Code != 100 {
		t.Fatalf("code=%d", r.Code)
	}
	t.Log("✅ GET /api/kg/domain/page")
}

// ── CRUD ──

func TestCreateRole(t *testing.T) {
	needToken(t)
	ts := time.Now().UnixNano()
	body := map[string]interface{}{
		"name":      fmt.Sprintf("test%d", ts),
		"sn":        fmt.Sprintf("TEST%d", ts),
		"companyId": 0,
		"status":    0,
	}
	_, respBody, err := do("POST", "/api/privilege/role", body, authHeader())
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	var r apiResp
	json.Unmarshal(respBody, &r)
	if r.Code != 100 {
		t.Fatalf("code=%d msg=%s", r.Code, r.Message)
	}
	t.Log("✅ POST /api/privilege/role (create)")
}

func TestCreateCompany(t *testing.T) {
	needToken(t)
	ts := time.Now().UnixNano()
	body := fmt.Sprintf(`{"cname":"test%d","code":"C%d","status":0}`, ts, ts)
	_, respBody, err := do("POST", "/api/privilege/company", body, authHeader())
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	var r apiResp
	json.Unmarshal(respBody, &r)
	if r.Code != 100 {
		t.Fatalf("code=%d msg=%s", r.Code, r.Message)
	}
	t.Log("✅ POST /api/privilege/company (create)")
}

func TestCreateDepartment(t *testing.T) {
	needToken(t)
	ts := time.Now().UnixNano()
	body := fmt.Sprintf(`{"name":"test%d","code":"D%d","companyId":"428001009954172928","pid":"0","status":0}`, ts, ts)
	_, respBody, err := do("POST", "/api/privilege/department", body, authHeader())
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	var r apiResp
	json.Unmarshal(respBody, &r)
	if r.Code != 100 {
		t.Fatalf("code=%d msg=%s", r.Code, r.Message)
	}
	t.Log("✅ POST /api/privilege/department (create)")
}

func TestCreateDictionary(t *testing.T) {
	needToken(t)
	ts := time.Now().UnixNano()
	body := fmt.Sprintf(`{"name":"test%d","code":"DICT%d","systemSn":"test","sn":"SN%d","sort":1}`, ts, ts, ts)
	_, respBody, err := do("POST", "/api/privilege/dictionary", body, authHeader())
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	var r apiResp
	json.Unmarshal(respBody, &r)
	if r.Code != 100 {
		t.Fatalf("code=%d msg=%s", r.Code, r.Message)
	}
	t.Log("✅ POST /api/privilege/dictionary (create)")
}

func TestCreateUser(t *testing.T) {
	needToken(t)
	ts := time.Now().UnixNano()
	body := fmt.Sprintf(`{"username":"utest%d","code":"UC%d","realName":"Test","companyId":"428001009954172928","deptId":"428003302434869248","status":0}`, ts, ts)
	_, respBody, err := do("POST", "/api/privilege/user", body, authHeader())
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	var r apiResp
	json.Unmarshal(respBody, &r)
	if r.Code != 100 {
		t.Fatalf("code=%d msg=%s", r.Code, r.Message)
	}
	t.Log("✅ POST /api/privilege/user (create)")
}

// ── helpers ──

func authHeader() map[string]string {
	return map[string]string{"Authorization": "Bearer " + token}
}

func platAuthHeader() map[string]string {
	return map[string]string{"Authorization": "Bearer " + platToken}
}

func needToken(t *testing.T) {
	if token == "" {
		t.Skip("no token available — login must run first")
	}
}

func needPlatToken(t *testing.T) {
	if platToken == "" {
		t.Skip("no platform token available — platform login must run first")
	}
}

// ── run order: tests execute alphabetically; TestA* ensures login runs first ──

// TestALoginGroup runs first (alphabetically) and obtains tokens for subsequent tests.
func TestALoginGroup(t *testing.T) {
	t.Run("captcha", TestCaptcha)
	t.Run("login", TestLogin)
	t.Run("platformLogin", TestPlatformLogin)
}

// TestZPrintSummary runs last and prints a summary line.
func TestZPrintSummary(t *testing.T) {
	// Just a marker — actual counts come from `go test -v`
	t.Logf("── All API tests completed ──")
	println("")
	println(strings.Repeat("=", 50))
	println("  Run:  go test ./test/ -v -count=1")
	println("  To see per-test pass/fail output above")
	println(strings.Repeat("=", 50))
}
