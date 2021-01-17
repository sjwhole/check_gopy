package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

type User struct {
	id          string
	UserId      string
	Password    string
	StudentCode string
}

func (u *User) SetInfo(uid, upw string) {
	u.UserId = uid
	u.Password = upw
}

func (u *User) SetStudentCode(stdCode string) {
	u.StudentCode = stdCode
}

func InitDB(file string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", file)
	if err != nil {
		return nil, err
	}

	createTableQuery1 := `create table IF NOT EXISTS info(id integer PRIMARY KEY autoincrement, userId text not null, password text not null, studentId text not null)`
	createTableQuery2 := `create table IF NOT EXISTS lecture(id integer PRIMARY KEY autoincrement,name text not null,code text not null)`
	_, e := db.Exec(createTableQuery1)
	if e != nil {
		return nil, e
	}
	_, e = db.Exec(createTableQuery2)
	if e != nil {
		return nil, e
	}
	return db, nil
}

func AddUser(db *sql.DB, user *User) error {
	tx, _ := db.Begin()
	stmt, _ := tx.Prepare("insert into info (userId,password,studentId) values (?,?,?)")
	_, err := stmt.Exec(user.UserId, user.Password, user.StudentCode)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	tx.Commit()
	return nil
}

func AddLecture(db *sql.DB, lectureName string, lectureCode string) error {
	tx, _ := db.Begin()
	stmt, _ := tx.Prepare("insert into lecture (name, code) values (?,?)")
	_, err := stmt.Exec(lectureName, lectureCode)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	tx.Commit()
	return nil
}

func ReturnUser(db *sql.DB, userId string) (User, error) {
	var user User
	rows := db.QueryRow("select * from info where id=1")
	err := rows.Scan(&user.id, &user.UserId, &user.Password, &user.StudentCode)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

//func ReturnLecture(studentDB *sql.DB) ([]string, []string) {
//	var lectureNameList, lectureCodeList []string
//	rows := studentDB.QueryRow("select * from lecture")
//}
