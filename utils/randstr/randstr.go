package randstr

import (
	"fmt"
	"math/rand"
	"time"
)

const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func GenerateRandomString(prefix string) string {
	src := rand.NewSource(time.Now().UnixNano())
	r := rand.New(src)

	b := make([]byte, 3)
	for i := range b {
		b[i] = charset[r.Intn(len(charset))]
	}

	return fmt.Sprintf("%s-%s-%s", prefix, time.Now().Format("020120061504"), string(b))
}
