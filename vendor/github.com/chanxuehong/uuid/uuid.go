package uuid

import (
	"github.com/chanxuehong/uuid/internal/v1"
	"github.com/chanxuehong/uuid/internal/v5"
)

//   +------ 0 ------+------ 1 ------+------ 2 ------+------ 3 ------+
//   +0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1+
//   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//   |                          time_low                             |
//   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//   |       time_mid                |         time_hi_and_version   |
//   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//   |clk_seq_hi_res |  clk_seq_low  |         node (0-1)            |
//   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//   |                         node (2-5)                            |
//   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+

type UUID [16]byte

// NewV1 returns a STANDARD version 1 UUID.
func NewV1() UUID {
	return v1.New()
}

// NewV5 returns a STANDARD version 5 UUID.
func NewV5(ns UUID, name []byte) UUID {
	return v5.New(ns, name)
}

// NewV1x returns a NONSTANDARD UUID(lower probability of conflict).
func NewV1x() UUID {
	return v1.Newx()
}

func (uuid UUID) Version() byte {
	return uuid[6] >> 4
}

func (uuid UUID) HexEncode() []byte {
	return HexEncode(uuid)
}

// Encode encodes UUID to "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx" format.
func (uuid UUID) Encode() []byte {
	return Encode(uuid)
}

// String returns "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx" format.
func (uuid UUID) String() string {
	return string(Encode(uuid))
}

// Decode decodes data with "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx" format into UUID.
func (uuid *UUID) Decode(data []byte) (err error) {
	u, err := Decode(data)
	if err != nil {
		return
	}
	*uuid = u
	return
}
