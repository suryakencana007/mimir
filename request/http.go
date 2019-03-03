/*  http.go
*
* @Author:             Nanang Suryadi
* @Date:               February 21, 2019
* @Last Modified by:   @suryakencana007
* @Last Modified time: 2019-02-21 01:27 
 */

package request

import (
    "bytes"
    "encoding/json"
    "io"
    "strings"
)

/**
* Convert Body to Json
* :body: interface
* :return: io.Reader
*/
func BodyToJson(body interface{}) (io.Reader, error) {
    if body == nil {
        return nil, nil
    }

    switch v := body.(type) {
    case string:
        // return as is
        return strings.NewReader(v), nil
    default:
        b, err := json.Marshal(v)
        if err != nil {
            return nil, err
        }

        return bytes.NewReader(b), nil
    }
}
