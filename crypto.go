/*  crypto.go
*
* @Author:             Nanang Suryadi
* @Date:               November 28, 2019
* @Last Modified by:   @suryakencana007
* @Last Modified time: 28/11/19 18:20
 */

package mimir

import (
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"hash"
	"strconv"
	"strings"

	"golang.org/x/crypto/pbkdf2"
)

const (
	RecommendedRoundsSHA1   = 131000
	RecommendedRoundsSHA256 = 29000
	RecommendedRoundsSHA512 = 25000
)

var b64 = base64.RawStdEncoding

// PassLibBase64Encode encodes using a variant of base64, like Passlib.
// Check https://pythonhosted.org/passlib/lib/passlib.utils.html#passlib.utils.ab64_encode
func PassLibBase64Encode(src []byte) (dst string) {
	dst = b64.EncodeToString(src)
	dst = strings.Replace(dst, "+", ".", -1)
	return
}

// PassLibBase64Decode decodes using a variant of base64, like Passlib.
// Check https://pythonhosted.org/passlib/lib/passlib.utils.html#passlib.utils.ab64_decode
func PassLibBase64Decode(src string) (dst []byte, err error) {
	src = strings.Replace(src, ".", "+", -1)
	dst, err = b64.DecodeString(src)
	return
}

// Base64Encode encodes using a Standard of base64.
// return string base64 encode
func Base64Encode(src []byte) (dst string) {
	return base64.StdEncoding.EncodeToString(src)
}

// Base64Encode decodes using a Standard of base64.
// return string base64 encode
func Base64Decode(src string) (dst []byte, err error) {
	return base64.StdEncoding.DecodeString(src)
}

func HashPassword(password, salt string) string {
	return fmt.Sprintf(
		"$pbkdf2-sha512$%d$%s$%v",
		RecommendedRoundsSHA512,
		PassLibBase64Encode([]byte(salt)),
		PassLibBase64Encode(
			pbkdf2.Key(
				[]byte(password),
				[]byte(salt),
				RecommendedRoundsSHA512,
				sha512.Size, sha512.New,
			),
		),
	)
}

func VerifyPassword(hashpassword, password string) (bool, error) {
	// only pbkdf2 supported
	if !strings.HasPrefix(hashpassword, "$pbkdf2-") {
		return false, fmt.Errorf("invalid hashPass")
	}
	// five fields expected: $pbkdf2-digest$rounds$salt$checksum
	fields := strings.Split(hashpassword, "$")
	if len(fields) != 5 {
		return false, fmt.Errorf("invalid hashPass format")
	}
	// extract digest
	hdr := strings.Split(fields[1], "-")
	if len(hdr) != 2 {
		return false, fmt.Errorf("invalid digest")
	}
	var (
		keyLen   int
		hashFunc func() hash.Hash
	)
	switch hdr[1] {
	case "sha256":
		keyLen = sha256.Size
		hashFunc = sha256.New
	case "sha512":
		keyLen = sha512.Size
		hashFunc = sha512.New
	default:
		return false, fmt.Errorf("invalid hashPass func")
	}
	// get remaining fields
	rounds, err := strconv.Atoi(fields[2])
	if err != nil {
		return false, fmt.Errorf("invalid hashPass roound")
	}
	salt, err := PassLibBase64Decode(fields[3])
	if err != nil {
		return false, fmt.Errorf("invalid hashPass salt")
	}
	key := pbkdf2.Key([]byte(password), salt, rounds, keyLen, hashFunc)
	return fields[4] == PassLibBase64Encode(key), nil
}
