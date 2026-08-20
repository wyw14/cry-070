package masking

import "strings"

func SanitizeLog(value string) string {
	lower := strings.ToLower(value)
	for _, secret := range []string{"password", "token", "dsn", "secret"} {
		if strings.Contains(lower, secret) {
			return "[redacted]"
		}
	}
	return value
}
func SanitizeMap(input map[string]string) map[string]string {
	out := map[string]string{}
	for k, v := range input {
		out[k] = SanitizeLog(v)
	}
	return out
}
