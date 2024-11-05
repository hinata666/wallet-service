/**
 * 测试转账
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

// 测试转账 go test -run TestTransfer -v
func TestTransfer(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	r.POST("/transfer", handlers.Transfer)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/transfer", bytes.NewBuffer([]byte(`{
	 "sender_id": "user3",
	 "receiver_id": "user4",
	 "amount": 150
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
