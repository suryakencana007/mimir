/*  generator.go
*
* @Date:               March 05, 2019
* @Last Modified by:   @suryakencana007
* @Last Modified time: 2019-03-05 22:55 
 */

package generator

import (
    "math/rand"
    "time"
)

const (
    letterBytes   = "ABCDEFGHIJKLMNPQRSTUVWXYZ123456789" // 34 possibilities
    letterIdxBits = 6                                    // 6 bits to represent 64 possibilities / indexes
    letterIdxMask = 1<<letterIdxBits - 1                 // All 1-bits, as many as letterIdxBits
    letterIdxMax  = 63 / letterIdxBits                   // # of letter indices fitting in 63 bits
)

func init() {
    rand.Seed(time.Now().UTC().UnixNano())
}

// GenerateVoucher Generate Voucher using Alphanumeric except O & 0 return as String
func RandomChar(length int) string {

    b := make([]byte, length)
    // A rand.Int63() generates 63 random bits, enough for letterIdxMax letters!
    for i, cache, remain := length-1, rand.Int63(), letterIdxMax; i >= 0; {
        if remain == 0 {
            cache, remain = rand.Int63(), letterIdxMax
        }
        if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
            b[i] = letterBytes[idx]
            i--
        }
        cache >>= letterIdxBits
        remain--
    }

    return string(b)
}
