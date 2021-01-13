package detailedlogger

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// A helpful wrapper for parsing a stream of bytes comprising an ssh message
type message struct {
	*bytes.Buffer
}

// When the RFC specifies a byte value
func (s *message) byte() (byte, error) {
	return s.Buffer.ReadByte()
}

// When the RFC specifies a boolean value
func (s *message) bool() (bool, error) {
	b, err := s.Buffer.ReadByte()
	if 1 == b {
		return true, err
	}
	return false, err
}

// When the RFC specifies a string value
func (s *message) str() ([]byte, error) {
	l, _ := s.uint32()
	if l > 5242880 {
		return []byte{}, fmt.Errorf("string size value too big")
	}
	v := make([]byte, l)
	err := binary.Read(s.Buffer, binary.BigEndian, &v)
	return v, err
}

// When the RFC specifies a uint32 value
func (s *message) uint32() (uint32, error) {
	var v uint32
	err := binary.Read(s.Buffer, binary.BigEndian, &v)
	return v, err
}
