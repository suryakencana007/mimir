/*  pagination.go
*
* @Author:             Nanang Suryadi
* @Date:               December 07, 2018
* @Last Modified by:   @suryakencana007
* @Last Modified time: 07/12/18 17:45
 */

package sql

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/suryakencana007/mimir/log"
)

type Pagination struct {
	Query        string
	Params       url.Values
	Model        interface{}
	AllowFields  []string
	DefaultValue string
	InjectWhere  string
	Aka          string
	Args         []interface{}
	QueryAll     string
	Limit        int
	Page         int
	Total        int
}

func FieldAlias(val string, a string) string {
	if a != "" {
		return fmt.Sprintf(`%s.%s`, a, val)
	}
	return string(val)
}

func GetPagination(p *Pagination) (string, error) {
	// filters
	filters, sep := "", ""
	if fields, ok := p.Params["fields"]; ok {
		f := make([]string, 0)
		for _, val := range strings.Split(fields[0], ",") {
			if Contains(p.AllowFields, val) {
				f = append(f, strings.Join([]string{
					FieldAlias(val, p.Aka),
					strings.Replace(
						`iLIKE '%' || ? || '%'`, "?",
						fmt.Sprintf(`$%d`, len(p.Args)+1), 1,
					),
				}, " "))
				p.Args = append(p.Args, p.Params.Get("q"))
			}
		}
		filters = fmt.Sprintf(`( %s )`, strings.Join(f, " OR "))
	}
	if len(filters) > 0 {
		sep = " AND "
	}
	// rql
	if rql, err := NewQueryFilter(p.Model); err == nil {
		rql.Args = p.Args
		rql, err := rql.QueryStringParser(p.Params, p.Aka, p.AllowFields)
		if err != nil {
			log.Error("QueryStringParser",
				log.Field("RQL", err.Error()))
			return "", err
		}

		if rql.Filter.(string) != "" {
			filters = strings.Join([]string{filters, rql.Filter.(string)}, sep)
		}
		p.Args = rql.Args
	}
	if len(filters) > 0 {
		sep = " AND "
	}
	if len(p.InjectWhere) > 0 {
		filters = strings.Join([]string{p.InjectWhere, filters}, sep)
	}
	sort := GetParam(p.Params, "sort", p.DefaultValue)
	sorted := make([]string, 0)
	for _, val := range strings.Split(sort[0], ",") {
		if Contains(p.AllowFields,
			strings.Replace(val, "-", "", 1)) {
			sorting := "ASC"
			if strings.Contains(val, "-") {
				sorting = "DESC"
				val = strings.Replace(val, "-", "", 1)
			}
			sorted = append(sorted, strings.Join(
				[]string{FieldAlias(val, p.Aka), sorting}, " "))
		}
	}
	sortedBy := strings.Join(sorted, ", ")
	q := []string{fmt.Sprintf(`SELECT * FROM (%s) as %s`, p.Query, p.Aka)}
	if len(filters) > 0 {
		q = append(q, strings.Replace(`WHERE ?`, "?", filters, 1)) // for Query Filtering
	}
	if len(sortedBy) > 0 {
		q = append(q, strings.Replace(`ORDER BY ?`, "?", sortedBy, 1)) // for Query Sorting
	}
	// query total count
	p.QueryAll = fmt.Sprintf(`SELECT count(1) FROM (%s) as tb_count`, strings.Join(q, " "))

	numPage := GetParam(p.Params, "page[number]", "1")
	page := 0
	if page, _ = strconv.Atoi(strings.Join(numPage, "")); page < 1 {
		page = 1 // default page
	}
	l := GetParam(p.Params, "page[size]", "20")
	limit, _ := strconv.Atoi(strings.Join(l, ""))
	p.Limit = limit
	p.Page = page
	q = append(q, strings.Replace(`LIMIT ?`, "?", strconv.Itoa(limit), 1))           // for Query Pagination
	q = append(q, strings.Replace(`OFFSET ?`, "?", strconv.Itoa(limit*(page-1)), 1)) // for Query Pagination
	return strings.Join(q, " "), nil
}

func Contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}

func GetParam(params url.Values, key string, defaultVal string) []string {
	v, ok := params[key]
	if !ok {
		v = []string{defaultVal}
	}
	return v
}
