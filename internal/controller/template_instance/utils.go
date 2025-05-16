package template_instance

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math/rand"
	"time"
)

func generateDiskID() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	raw := fmt.Sprintf("%d-%d", time.Now().UnixNano(), r.Int())
	sum := md5.Sum([]byte(raw))
	return hex.EncodeToString(sum[:])[:8]
}
