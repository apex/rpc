package format_test

import (
	"testing"

	"github.com/apex/rpc/internal/format"
	"github.com/tj/assert"
)

// Test go name formatting.
func TestGoName(t *testing.T) {
	assert.Equal(t, "UpdateUser", format.GoName("update_user"))
	assert.Equal(t, "UserID", format.GoName("user_id"))
	assert.Equal(t, "IP", format.GoName("ip"))
}

// Test js name formatting.
func TestJsName(t *testing.T) {
	assert.Equal(t, "updateUserSettings", format.JsName("update_user_settings"))
	assert.Equal(t, "userId", format.JsName("user_id"))
}
