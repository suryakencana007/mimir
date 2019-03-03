/*  rql.go
*
* @Author:             Nanang Suryadi
* @Date:               December 07, 2018
* @Last Modified by:   @suryakencana007
* @Last Modified time: 07/12/18 16:07 
 */

package sql

import (
    "bytes"
    "fmt"
    "net/url"
    "regexp"
    "sort"
    "strconv"
    "strings"
)

type QueryFilter struct {
    *bytes.Buffer
    Model    interface{}
    TagName  string
    FieldSep string
    Args     []interface{}
    Filter   interface{} `json:"filter"`
    Log      func(string, ...interface{})
}

// NewParser creates a new Parser. it fails if the configuration is invalid.
func NewQueryFilter(model interface{}) (*QueryFilter, error) {
    p := &QueryFilter{
        Model:    model,
        TagName:  "rql",
        FieldSep: ".",
    }
    return p, nil
}

func (p *QueryFilter) QueryStringParser(params url.Values, alias string, allowFields []string) (*QueryFilter, error) {
    var keys []string
    for k := range params {
        keys = append(keys, k)
    }
    sort.Strings(keys)
    f := make([]string, 0)
    for _, k := range keys {
        rx := regexp.MustCompile(`filters\[(.*?)]`)
        if ok := rx.MatchString(k); ok {
            matchGroup := rx.FindStringSubmatch(k)
            m := strings.Split(matchGroup[1], ":")
            if len(m) < 2 {
                continue
            }
            filter, operator := m[0], m[1]
            if Contains(allowFields, filter) {
                switch Op(operator) {
                case EQ:
                    operator = EQ.SQL()
                case LT:
                    operator = LT.SQL()
                case LTE:
                    operator = LTE.SQL()
                case GTE:
                    operator = GTE.SQL()
                case GT:
                    operator = GT.SQL()
                default:
                    operator = EQ.SQL()
                }
                filter = FieldAlias(filter, alias)
                f = append(f, fmt.Sprintf(`%s %s $%d`, filter, operator, len(p.Args)+1))
                b, err := strconv.ParseBool(params[k][0])
                if err == nil {
                    p.Args = append(p.Args, b)
                    continue
                }

                integer, err := strconv.Atoi(params[k][0])
                if err == nil {
                    p.Args = append(p.Args, integer)
                    continue
                }
                p.Args = append(p.Args, params[k][0])
            }
        }
    }
    p.Filter = strings.Join(f, " AND ")
    return p, nil
}

// Op is a filter operator used by rql.
type Op string

// SQL returns the SQL representation of the operator.
func (o Op) SQL() string {
    return opFormat[o]
}

// Operators that support by rql.
const (
    EQ   = Op("eq")   // =
    NEQ  = Op("neq")  // <>
    LT   = Op("lt")   // <
    GT   = Op("gt")   // >
    LTE  = Op("lte")  // <=
    GTE  = Op("gte")  // >=
    LIKE = Op("like") // LIKE "PATTERN"
    OR   = Op("or")   // disjunction
    AND  = Op("and")  // conjunction
)

var (
    opFormat = map[Op]string{
        EQ:   "=",
        NEQ:  "<>",
        LT:   "<",
        GT:   ">",
        LTE:  "<=",
        GTE:  ">=",
        LIKE: "LIKE",
        OR:   "OR",
        AND:  "AND",
    }
)
