package usecase

import (
	"crypto/sha256"
	"encoding/base64"
)

func generateShort(original, secret string) string {
	sum := sha256.Sum256([]byte(original + secret))
	// base64.URLEncoding без паддинга, берём первые 8 символов
	return base64.RawURLEncoding.EncodeToString(sum[:])[:8]
}
