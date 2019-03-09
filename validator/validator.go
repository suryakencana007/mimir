/*  validator.go
*
* @Author:             Nanang Suryadi
* @Date:               March 09, 2019
* @Last Modified by:   @suryakencana007
* @Last Modified time: 2019-03-09 20:06 
 */

package validator

import (
    "time"

    "gopkg.in/go-playground/validator.v9"
)

func DateValidation(fl validator.FieldLevel) bool {
    _, err := time.Parse("2006-01-02", fl.Field().String())
    if err != nil {
        return false
    }
    return true
}

func DatetimeValidation(fl validator.FieldLevel) bool {
    _, err := time.Parse(time.RFC3339, fl.Field().String())
    if err != nil {
        return false
    }
    return true
}

func DateRangeValidation(fl validator.FieldLevel) bool {

    var date = fl.Field().String()
    var minDate = time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)
    var maxDate = time.Date(2100, 1, 1, 0, 0, 0, 0, time.UTC)

    datetime, err := time.Parse("2006-01-02", date)
    if err != nil {
        return false
    }
    if datetime.Before(minDate) || datetime.After(maxDate) {
        return false
    }
    return true
}

func ParseDate(dtStr string) time.Time {
    date, err := time.Parse("2006-01-02", dtStr)
    if err != nil {
        return time.Time{}
    }
    return date
}

func ParseDatetime(dtStr string) time.Time {
    dateTime, err := time.Parse(time.RFC3339, dtStr)
    if err != nil {
        return time.Time{}
    }
    return dateTime
}
