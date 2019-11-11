package db

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"gochat.udp/config"
	"gochat.udp/logger"
)

var (
	mysqlquerycli *sql.DB
)

func init() {
	var err error
	connStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		config.GetKey("mysql_user"),
		config.GetKey("mysql_passwd"),
		config.GetKey("mysql_host"),
		config.GetKey("mysql_port"),
		config.GetKey("mysql_db"),
	)
	mysqlquerycli, err = sql.Open("mysql", connStr)
	if err != nil {
		logger.Error("open mysql err, ", err)
	}
	err = mysqlquerycli.Ping()
	if err != nil {
		logger.Error("ping mysql err, ", err)
	}
}

func Insert(SQL string) bool {
	stmt, err := mysqlquerycli.Prepare(SQL)

	if err != nil {
		logger.Error("prepare sql error, ", err)
	}
	res, err := stmt.Exec()
	if err != nil {
		logger.Error("stmt.Exec err:", err)
		return false
	}
	id, err := res.LastInsertId()
	if err != nil {
		logger.Errorf("Insert err:%v\n\t%s\n", err, SQL)
		return false
	} else {
		logger.Info("Insert successfully!", id)
	}

	return true
}

func Update(SQL string) bool {
	stmt, err := mysqlquerycli.Prepare(SQL)

	if err != nil {
		logger.Error("prepare sql error, ", err)
	}
	res, err := stmt.Exec()
	if err != nil {
		logger.Error("stmt.Exec err:", err)
		return false
	}
	id, err := res.LastInsertId()
	if err != nil {
		logger.Errorf("Insert err:%v\n\t%s\n", err, SQL)
		return false
	} else {
		logger.Info("Insert successfully!", id)
	}

	return true
}

func Delete(SQL string) bool {
	stmt, err := mysqlquerycli.Prepare(SQL)

	if err != nil {
		logger.Error("prepare sql error, ", err)
	}
	res, err := stmt.Exec()
	if err != nil {
		logger.Error("stmt.Exec err:", err)
		return false
	}
	id, err := res.LastInsertId()
	if err != nil {
		logger.Errorf("Insert err:%v\n\t%s\n", err, SQL)
		return false
	} else {
		logger.Info("Insert successfully!", id)
	}

	return true
}

func Query(SQL string) sql.Result {
	stmt, err := mysqlquerycli.Prepare(SQL)

	if err != nil {
		logger.Error("prepare sql error, ", err)
	}
	res, err := stmt.Exec()
	if err != nil {
		logger.Error("stmt.Exec err:", err)
		return res
	}

	return res
}

func QueryAllFansIDByID(id int) *sql.Rows {
	sql := fmt.Sprintf("SELECT fansid FROM fans WHERE userid=%d", id)
	res, err := mysqlquerycli.Query(sql)
	if err != nil {
		logger.Error("stmt.Exec err:", err)
		return res
	}

	return res
}

func QueryNameByID(id int) (name string) {
	_ = mysqlquerycli.QueryRow("SELECT username FROM user WHERE id=?", id).Scan(&name)
	logger.Info("mysql ", name, "id ", id)
	return name
}

func VarifyUserByPasswd(username, password string) (id int, name string) {
	_ = mysqlquerycli.QueryRow("SELECT id, username FROM user WHERE username=? AND password=?",
		username, password).Scan(&id, &name)
	return
}

func ModifyName(name string, id int) bool {
	res, err := mysqlquerycli.Exec("UPDATE user set username=? WHERE id=?", name, id)
	if err != nil {
		logger.Error("update error ", err)
		return false
	}
	if ra, err := res.RowsAffected(); ra == 0 || err != nil {
		logger.Error("update error ", err)
		return false
	}
	return true
}
