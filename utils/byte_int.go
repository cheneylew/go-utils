package utils

import (
	"bytes"
	"encoding/binary"
)

func IntToBytes(n int) []byte {
	x := int32(n)
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, x)

	return bytesBuffer.Bytes()
}

func BytesToInt(b []byte) int {
	bytesBuffer := bytes.NewBuffer(b)

	var x int32
	binary.Read(bytesBuffer, binary.BigEndian, &x)

	return int(x)
}

func UInt8ToBytes(n uint8) []byte {
	x := uint8(n)

	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, x)

	return bytesBuffer.Bytes()
}

func BytesToUInt8(b []byte) uint8 {
	bytesBuffer := bytes.NewBuffer(b)

	var x uint8
	binary.Read(bytesBuffer, binary.BigEndian, &x)

	return uint8(x)
}

func UInt16ToBytes(n uint16) []byte {
	x := uint16(n)

	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, x)

	return bytesBuffer.Bytes()
}

func BytesToUInt16(b []byte) uint16 {
	bytesBuffer := bytes.NewBuffer(b)

	var x uint16
	binary.Read(bytesBuffer, binary.BigEndian, &x)

	return uint16(x)
}

func UInt64ToBytes(n uint64) []byte {
	x := uint64(n)

	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, x)

	return bytesBuffer.Bytes()
}

func BytesToUInt64(b []byte) uint64 {
	bytesBuffer := bytes.NewBuffer(b)

	var x uint64
	binary.Read(bytesBuffer, binary.BigEndian, &x)

	return uint64(x)
}