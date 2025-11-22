package dsl

import (
	"crypto/sha256"
	"encoding/binary"
)

// deriveSeed 将 (seedStr, templateID, version, serverSalt) 映射为稳定的 int64 种子
func deriveSeed(seedStr, templateID, version, serverSalt string) int64 {
	h := sha256.Sum256([]byte(seedStr + "|" + templateID + "|" + version + "|" + serverSalt))
	v := int64(binary.LittleEndian.Uint64(h[:8]))
	if v < 0 {
		v = -v
	}
	return v
}
