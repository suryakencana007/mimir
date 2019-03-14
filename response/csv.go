/*  csv.go
*
* @Author:             Nanang Suryadi
* @Date:               March 15, 2019
* @Last Modified by:   @suryakencana007
* @Last Modified time: 2019-03-15 04:38 
 */

package response

import (
    "bytes"
    "encoding/csv"
    "fmt"
    "log"
    "net/http"

    "github.com/suryakencana007/mimir/constant"
)

func WriteCSV(w http.ResponseWriter, r *http.Request, rows [][]string, filename string) {
    buf := &bytes.Buffer{}
    xCsv := csv.NewWriter(buf)

    for _, row := range rows {
        if err := xCsv.Write(row); err != nil {
            log.Println("error writing record to csv:", err)
            http.Error(w, err.Error(), http.StatusInternalServerError)
        }
    }
    xCsv.Flush()

    if err := xCsv.Error(); err != nil {
        log.Println("error writing record to csv:", err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
    w.Header().Set("Content-Description", "File Transfer")
    w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s.csv", filename))
    w.Header().Set("Content-Type", "text/csv; charset=utf-8")
    if status, ok := r.Context().Value(constant.StatusCtxKey).(int); ok {
        w.WriteHeader(status)
    }
    _, err := w.Write(buf.Bytes())
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
}
