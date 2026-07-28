package common

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

func GenerateRoomID() string {
	b := make([]byte, 16)
	rand.Read(b)
	data := fmt.Sprintf("%d-%s", time.Now().UnixNano(), hex.EncodeToString(b))
	sum := sha256.Sum256([]byte(data))
	return hex.EncodeToString(sum[:])[:32]
}
