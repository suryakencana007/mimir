/*  http.go
*
* @Author:             Nanang Suryadi
* @Date:               February 13, 2019
* @Last Modified by:   @suryakencana007
* @Last Modified time: 2019-02-13 13:51 
 */

package constant

import (
    "net/http"
)

const (
    StatusCtxKey                = 0
    StatusSuccess               = http.StatusOK
    StatusErrorForm             = http.StatusBadRequest
    StatusErrorUnknown          = http.StatusBadGateway
    StatusInternalError         = http.StatusInternalServerError
    StatusUnauthorized          = http.StatusUnauthorized
    StatusCreated               = http.StatusCreated
    StatusForbidden             = http.StatusForbidden
    StatusInvalidAuthentication = http.StatusProxyAuthRequired
)

var statusMap = map[int][]string{
    StatusSuccess:               {"STATUS_OK", "Success"},
    StatusErrorForm:             {"STATUS_BAD_REQUEST", "Invalid data request"},
    StatusErrorUnknown:          {"STATUS_BAG_GATEWAY", "Oops something went wrong"},
    StatusInternalError:         {"INTERNAL_SERVER_ERROR", "Oops something went wrong"},
    StatusUnauthorized:          {"STATUS_UNAUTHORIZED", "Not authorized to access the service"},
    StatusCreated:               {"STATUS_CREATED", "Resource has been created"},
    StatusForbidden:             {"STATUS_FORBIDDEN", "Forbidden access the resource "},
    StatusInvalidAuthentication: {"STATUS_INVALID_AUTHENTICATION", "The resource owner or authorization server denied the request"},
}

func StatusCode(code int) string {
    return statusMap[code][0]
}

func StatusText(code int) string {
    return statusMap[code][1]
}
