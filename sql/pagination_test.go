/*  pagination_test.go
*
* @Author:             Nanang Suryadi
* @Date:               December 07, 2018
* @Last Modified by:   @suryakencana007
* @Last Modified time: 07/12/18 19:41 
 */

package sql

import (
    "database/sql"
    "fmt"
    "net/url"
    "testing"
    "time"

    "github.com/stretchr/testify/assert"
    "github.com/suryakencana007/mimir/log"
)

func init() {
    log.ZapInit()
}

func TestFieldAlias(t *testing.T) {
    assert := assert.New(t)
    a, v := "g", "value"
    expected := fmt.Sprintf("%s.%s", a, v)
    actual := FieldAlias(v, a)
    assert.Equal(expected, actual)
}

func TestFieldNoAlias(t *testing.T) {
    assert := assert.New(t)
    v := "value"
    expected := fmt.Sprintf("%s", v)
    actual := FieldAlias(v, "")
    assert.Equal(expected, actual)
}

func TestGetPagination(t *testing.T) {
    assert := assert.New(t)
    pagination := &Pagination{
        Query: `SELECT group_id, name, parent_id, category, description, activated FROM groups g`,
        Params: url.Values{
            "q":                     []string{""},
            "sort":                  []string{"name,category"},
            "fields":                []string{"name"},
            "filters[activated:eq]": []string{"true"},
            "filters[category:eq]":  []string{"RG"},
            "filters[group_id:eq]":  []string{"66"},
        },
        Model:        group{},
        AllowFields:  []string{"group_id", "category", "name", "activated"},
        DefaultValue: "name",
        Aka:          "g",
    }
    actual, err := GetPagination(pagination)
    expected := `SELECT group_id, name, parent_id, category, description, activated FROM groups g WHERE g.name iLIKE '%' || $1 || '%' AND g.activated = $2 AND g.category = $3 AND g.group_id = $4 ORDER BY g.name ASC, g.category ASC LIMIT 20 OFFSET 0`
    assert.NoError(err)
    assert.Equal(expected, actual)
}

// Group representation of group
type group struct {
    ID          uint16         `json:"group_id"`
    Parent      sql.NullInt64  `json:"parent_id"`
    Category    string         `json:"category" rql:"filter"`
    Name        string         `json:"name" rql:"filter,sort"`
    Description sql.NullString `json:"description"`
    Activated   sql.NullBool   `json:"activated" rql:"filter"`
    CreatedAt   *time.Time     `json:"created_at" rql:"filter"`
    CreatedBy   string         `json:"created_by" rql:"filter"`
    UpdatedAt   *time.Time     `json:"updated_at" rql:"filter"`
    UpdatedBy   string         `json:"updated_by" rql:"filter"`
    DeletedAt   *time.Time     `json:"deleted_at"`
}
