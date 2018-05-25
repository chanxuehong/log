package trace

import "github.com/chanxuehong/uuid"

func NewTraceId() string {
	return string(uuid.NewV1().HexEncode())
}
