/*  crypto_test.go
*
* @Author:             Nanang Suryadi
* @Date:               November 28, 2019
* @Last Modified by:   @suryakencana007
* @Last Modified time: 28/11/19 18:20
 */

package mimir

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBase64(t *testing.T) {
	var flagtests = []struct {
		title string
		word  string
	}{
		{"hash password", "sekretuwhw8w8ewyewueibxw74h747hwuwywe74wuey7273y23ebebyd6773yye3456"},
		{"ascii word", "Hello, 世界"},
	}
	for _, tt := range flagtests {
		tt := tt // pin it
		t.Run(tt.title, func(t *testing.T) {
			encode := Base64Encode([]byte(tt.word))
			// assert.NoError(t, err)
			decode, err := Base64Decode(encode)
			assert.NoError(t, err)
			assert.Equal(t, tt.word, string(decode))
		})
	}
}

func TestPassLibBase64(t *testing.T) {
	var flagtests = []struct {
		title string
		word  string
	}{
		{"hash password", "sekretuwhw8w8ewyewueibxw74h747hwuwywe74wuey7273y23ebebyd6773yye3456"},
		{"ascii word", "Hello, 世界"},
	}
	for _, tt := range flagtests {
		tt := tt // pin it
		t.Run(tt.title, func(t *testing.T) {
			encode := PassLibBase64Encode([]byte(tt.word))
			// assert.NoError(t, err)
			decode, err := PassLibBase64Decode(encode)
			assert.NoError(t, err)
			assert.Equal(t, tt.word, string(decode))
		})
	}
}
