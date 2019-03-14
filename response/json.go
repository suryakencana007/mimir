/*  response.go
*
* @Author:             Nanang Suryadi
* @Date:               October 01, 2018
* @Last Modified by:   @suryakencana007
* @Last Modified time: 01/10/18 00:21 
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

// Write writes the data to http response writer
func WriteCSV(w http.ResponseWriter, r *http.Request, data []string, filename string) {
    buf := &bytes.Buffer{}
    xCsv := csv.NewWriter(buf)
    if err := xCsv.Write(data); err != nil {
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
