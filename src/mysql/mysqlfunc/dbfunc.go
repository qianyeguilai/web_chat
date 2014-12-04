package mysqlfunc

import "fmt"
import "database/sql"
import _ "github.com/go-sql-driver/mysql"

type Db_info struct {
	Db_host     string
	Db_user     string
	Db_passwd   string
	Db_database string
}

func ConnectDb(info Db_info) (*sql.DB, error) {
	var connect_str string
	connect_str = fmt.Sprintf("%s@tcp(%s:3306)/%s?autocommit=true&charset=utf8", info.Db_user, info.Db_passwd, info.Db_database)
	db, err := sql.Open("mysql", connect_str)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func Checkuserok(db *sql.DB, username, passwd []string) bool {
	connect_str := fmt.Sprintf("select password  from user_base_info where cname=\"%s\";", username[0])
	fmt.Println(connect_str)
	rows, err := db.Query(connect_str)
	if err != nil {
		return false
	}
	defer rows.Close()

	var temp string
	for rows.Next() {
		err := rows.Scan(&temp)
		if err != nil {
			return false
		}
		if temp == passwd[0] {
			return true
		} else {
			return false
		}
	}
	return false
}

func Db_register_user(db *sql.DB, username, passwd string, sex string) error {
	stmt, err := db.Prepare("insert  into  user_base_info (cname,password,sex) value(?,?,?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	stmt.Exec(username, passwd, sex)
	return nil
}

func Judge_user_exist(db *sql.DB, username string) (bool, error) {
	connect_str := fmt.Sprintf("select count(*) from user_base_info where  cname=\"%s\";", username)
	rows, err := db.Query(connect_str)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	var temp string
	for rows.Next() {
		err := rows.Scan(&temp)
		if err != nil {
			return false, err
		}

		if temp == "0" {
			return true, nil
		}
	}
	return false, nil

}
