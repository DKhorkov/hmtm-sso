package security

import "encoding/base64"

func Encode(data []byte) string {
	return base64.URLEncoding.EncodeToString(data)
}

func Decode(encoded string) ([]byte, error) {
	decoded, err := base64.URLEncoding.DecodeString(encoded)
	if err != nil {
		return nil, err
	}

	return decoded, nil
}
