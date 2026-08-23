package security

import (
	"strings"
	"testing"
)

func TestRedactSecrets(t *testing.T) {
	input := `
  api_key = "sk-12345678901234567890"
  Authorization: Bearer my-secret-bearer-token-123
  password: "this-is-a-secret-password"
  -----BEGIN RSA PRIVATE KEY-----
  secret private key content
	  -----END RSA PRIVATE KEY-----
	  `
	if strings.Contains(input, "[REDACTED_") {
		t.Fatal("测试输入不应预先包含脱敏占位符")
	}

	output := RedactSecrets(input)
	if output == input {
		t.Fatal("RedactSecrets 未改变包含敏感信息的输入")
	}
	for _, secret := range []string{
		"sk-12345678901234567890",
		"my-secret-bearer-token-123",
		"this-is-a-secret-password",
		"secret private key content",
	} {
		if strings.Contains(output, secret) {
			t.Fatalf("敏感信息未被脱敏: %q\noutput: %s", secret, output)
		}
	}
	if !strings.Contains(output, "[REDACTED_SECRET]") || !strings.Contains(output, "[REDACTED_BEARER_TOKEN]") {
		t.Fatalf("脱敏占位符不完整: %s", output)
	}
}
