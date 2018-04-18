package log

import "github.com/chanxuehong/uuid"

func NewRequestId() string {
	return string(uuid.NewV1().HexEncode())
}
