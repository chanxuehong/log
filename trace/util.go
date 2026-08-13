package trace

import "github.com/chanxuehong/log/internal/uuid"

func NewTraceId() string {
	return string(uuid.NewV1().HexEncode())
}
