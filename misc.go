/*  misc.go
*
* @Author:             Nanang Suryadi
* @Date:               November 22, 2019
* @Last Modified by:   @suryakencana007
* @Last Modified time: 22/11/19 14:01
 */

package mimir

import (
	"encoding/json"
	"errors"
	"math/rand"
	"reflect"
	"regexp"
	"strings"
	"time"
	"unsafe"

	"github.com/oklog/ulid/v2"
	"github.com/satori/go.uuid"
)

const (
	letterBytes   = "ABCDEFGHIJKLMNPQRSTUVWXYZ123456789" // 34 possibilities
	letterIdxBits = 6                                    // 6 bits to represent 64 possibilities / indexes
	letterIdxMask = 1<<letterIdxBits - 1                 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits                   // # of letter indices fitting in 63 bits
)

var src = rand.NewSource(time.Now().UnixNano())

// GenerateVoucher Generate Voucher using Alphanumeric except O & 0 return as String
func GenerateChar(length int) string {
	b := make([]byte, length)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := length-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return *(*string)(unsafe.Pointer(&b))
}

var numberSequence = regexp.MustCompile(`([a-zA-Z])(\d+)([a-zA-Z]?)`)
var numberReplacement = []byte(`$1 $2 $3`)

func addWordBoundariesToNumbers(s string) string {
	b := []byte(s)
	b = numberSequence.ReplaceAll(b, numberReplacement)
	return string(b)
}

// Converts a string to CamelCase
func toCamelInitCase(s string, initCase bool) string {
	s = addWordBoundariesToNumbers(s)
	s = strings.Trim(s, " ")
	n := ""
	capNext := initCase
	for _, v := range s {
		if v >= 'A' && v <= 'Z' {
			n += string(v)
		}
		if v >= '0' && v <= '9' {
			n += string(v)
		}
		if v >= 'a' && v <= 'z' {
			if capNext {
				n += strings.ToUpper(string(v))
			} else {
				n += string(v)
			}
		}
		if v == '_' || v == ' ' || v == '-' || v == '.' {
			capNext = true
		} else {
			capNext = false
		}
	}
	return n
}

// ToCamel converts a string to CamelCase
func ToCamel(s string) string {
	return toCamelInitCase(s, false)
}

// UUID returns a newly initialized string object that implements the UUID
// interface.
func UUID() string {
	return uuid.NewV4().String()
}

// Universally Unique Lexicographically Sortable Identifier
func ULID() ulid.ULID {
	t := time.Unix(1000000, 0)
	entropy := ulid.Monotonic(rand.New(rand.NewSource(t.UnixNano())), 0)
	return ulid.MustNew(ulid.Timestamp(t), entropy)
}

var (
	NotPointer = errors.New("must pass a pointer, not a value")
	NilPointer = errors.New("nil pointer passed")
)

func CloneStruct(src, dest interface{}) error {
	d := reflect.ValueOf(dest)
	if d.Kind() != reflect.Ptr {
		return NotPointer
	}
	if d.IsNil() {
		return NilPointer
	}
	s := reflect.ValueOf(src)
	if s.Kind() != reflect.Ptr {
		return NotPointer
	}
	if s.IsNil() {
		return NilPointer
	}

	b, err := json.Marshal(src)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, dest)
}
