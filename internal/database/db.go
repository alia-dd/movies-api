package database
import(
	"database/sql"
	_"github.com/mattn/go-sqlite3"
)

func OpenDB()(*sql.DB,error){
	db, err := sql.Open("sqlite3", "./movies.db")
	if err != nil{
		return db, err
	}
	return db, nil
}