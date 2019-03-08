/*  db.go
*
* @Author:             Nanang Suryadi
* @Date:               February 12, 2019
* @Last Modified by:   @suryakencana007
* @Last Modified time: 2019-02-12 13:45 
 */

package constant

type dbConstKey string

const (
    TxKey         = dbConstKey("TxKey")
    MYSQL         = "mysql"
    POSTGRES      = "postgres"
    POSTGRES_CONN = `host=%s port=%d user=%s password=%s dbname=%s sslmode=disable`
)
