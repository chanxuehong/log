package log

import "github.com/chanxuehong/uuid"

func newRequestId() string {
	return string(uuid.NewV1().HexEncode())
}
