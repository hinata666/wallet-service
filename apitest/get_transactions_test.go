/**
 * 测试获取指定用户的交易历史
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

// 测试获取指定用户的交易历史 go test -run TestGetTransactions -v
func TestGetTransactions(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	r.POST("/transactions", handlers.GetTransactions)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/transactions", bytes.NewBuffer([]byte(`{
	 "user_id": "user1",
	 "page": 1,
	 "page_size": 2
 }`)))

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
