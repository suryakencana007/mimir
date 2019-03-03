/*  rql_test.go
*
* @Author:             Nanang Suryadi
* @Date:               March 04, 2019
* @Last Modified by:   @suryakencana007
* @Last Modified time: 2019-03-04 00:28 
 */

package sql

import (
    "net/url"
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestNewQueryFilter(t *testing.T) {
    // rql
    if rql, err := NewQueryFilter(group{}); err == nil {
        rql.Args = []interface{}{}
        rql, err := rql.QueryStringParser(params, aka, allowFields)
        if err != nil {
            t.Error(err.Error())
        }
        actual := rql.Filter.(string)
        assert.Exactly(t,
            "g.activated = $1 AND g.category = $2 AND g.confirmation_date >= $3 AND g.confirmation_date < $4 AND g.group_id = $5 AND g.group_id > $6 AND g.group_id <= $7 AND g.group_id = $8",
            actual,
        )
        args := []interface{}{
            true,
            "RG",
            66,
            "2019-03-05",
            55,
            "2019-03-04",
            77,
            77,
        }
        assert.ElementsMatch(t, args, rql.Args)
    }
}

var params = url.Values{
    "q":                              []string{""},
    "sort":                           []string{"name,category"},
    "fields":                         []string{"name"},
    "filters[activated:eq]":          []string{"true"},
    "filters[category:eq]":           []string{"RG"},
    "filters[group_id:eq]":           []string{"66"},
    "filters[confirmation_date:lt]":  []string{"2019-03-05"},
    "filters[group_id:lte]":          []string{"55"},
    "filters[confirmation_date:gte]": []string{"2019-03-04"},
    "filters[group_id:gt]":           []string{"77"},
    "filters[group_id:neq]":          []string{"77"},
}

var (
    allowFields = []string{"group_id", "category", "name", "activated", "confirmation_date"}
    aka         = "g"
)
