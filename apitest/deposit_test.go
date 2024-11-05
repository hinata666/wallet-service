/**
 * 测试存款
 */

package apitest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"wallet-service/handlers"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
)

// 测试存款 go test -run TestDeposit -v
func TestDeposit(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	r.POST("/deposit", handlers.Deposit)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/deposit", bytes.NewBuffer([]byte(`{"user_id": "user3", "amount": 100}`)))

	r.ServeHTTP(w, req)

	// 检查响应状态码
	assert.Equal(t, http.StatusOK, w.Code)

	// 解析响应体为 JSON
	var resp map[string]interface{}
	if err := json.Unmarshal([]byte(w.Body.String()), &resp); err != nil {
		t.Errorf("Failed to parse response body as JSON: %v", err)
	}

	// 输出 JSON 响应
	respJSON, _ := json.MarshalIndent(resp, "", "  ")
	fmt.Println(string(respJSON))
}
