/*  csv.go
*
* @Author:             Nanang Suryadi
* @Date:               November 21, 2019
* @Last Modified by:   @suryakencana007
* @Last Modified time: 21/11/19 22:36
 */

package mimir

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
)

func (r *responseWriter) WriteCSV(rows [][]string, filename string) {
	buf := &bytes.Buffer{}
	xCsv := csv.NewWriter(buf)

	for _, row := range rows {
		if err := xCsv.Write(row); err != nil {
			log.Println("error writing record to csv:", err)
			http.Error(r.Writer, err.Error(), http.StatusInternalServerError)
		}
	}
	xCsv.Flush()

	if err := xCsv.Error(); err != nil {
		log.Println("error writing record to csv:", err)
		http.Error(r.Writer, err.Error(), http.StatusInternalServerError)
	}
	r.Writer.Header().Set("Content-Description", "File Transfer")
	r.Writer.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s.csv", filename))
	r.Writer.Header().Set("Content-Type", "text/csv; charset=utf-8")
	if status, ok := r.Request.Context().Value(CtxResponse).(int); ok {
		r.Writer.WriteHeader(status)
	}
	_, err := r.Writer.Write(buf.Bytes())
	if err != nil {
		http.Error(r.Writer, err.Error(), http.StatusInternalServerError)
		return
	}
}
