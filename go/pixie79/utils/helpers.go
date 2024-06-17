package utils

import (
	"encoding/base64"
	"encoding/binary"
)

// EncodeBuffer encodes the given schemaID into a byte array and returns the resulting header.
//
// The schemaID parameter is an integer representing the ID of the schema to be encoded.
// The function returns a byte array that represents the header, including the schemaID.
func EncodeBuffer(schemaID int) []byte {
	byteArray := make([]byte, 4)
	binary.BigEndian.PutUint32(byteArray, uint32(schemaID))
	header := append([]byte{0}, byteArray...)
	return header
}

// b64DecodeMsg decodes a base64 encoded key and returns a subset of the key starting from the specified offset.
//
// Parameters:
//   - b64Key: The base64 encoded key to be decoded.
//   - offsetF: An optional integer representing the offset from which to start the subset of the key. If not provided, it defaults to 7.
//
// Returns:
//   - []byte: The subset of the key starting from the specified offset.
//   - error: An error if the decoding or subset operation fails.
func b64DecodeMsg(b64Key string, offsetF ...int) ([]byte, error) {
	offset := 7
	if len(offsetF) > 0 {
		offset = offsetF[0]
	}

	key, err := base64.StdEncoding.DecodeString(b64Key)
	if err != nil {
		return nil, err
	}

	result := key[offset:]
	return result, nil
}
