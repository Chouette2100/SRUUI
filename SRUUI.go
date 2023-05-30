/*!
Copyright © 2022 chouette.21.00@gmail.com
Released under the MIT license
https://opensource.org/licenses/mit-license.php
Ver. 0.1.0 ApiLiveCurrentUser.goをsrapi.goから分離する。ApiLiveCurrentUser()のRoomIDをstring型に変更する。
*/
package main

import (
	"fmt"
	"log"

	//	"math"
	//	"sort"
	//	"strconv"
	//	"strings"
	"time"

	//	"bufio"
	"io"
	"os"

	//	"runtime"

	//	"encoding/json"

	//	"html/template"
	//	"net/http"

	"database/sql"
	_ "github.com/go-sql-driver/mysql"

	//	"github.com/PuerkitoBio/goquery"

	//	svg "github.com/ajstarks/svgo/float"

	//	"github.com/dustin/go-humanize"

	//	scl "UpdateUserInf/ShowroomCGIlib"
	"SRUUI/UpdateUserInfLib"
	"github.com/Chouette2100/exsrapi"
)

/*
	Ver.0000A00 指定したイベントの配信者情報を取得し更新する。
	Ver.000AB00 開催中、開催前のイベントの配信者情報を取得し更新する。
*/

const version = "000AB00"

func main() {

	//	ログ出力を設定する
	logfilename := version + "_" + UpdateUserInfLib.Version + "_" + time.Now().Format("20060102") + ".txt"
	logfile, err := os.OpenFile(logfilename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		panic("cannnot open logfile: " + logfilename + err.Error())
	}
	defer logfile.Close()
	//	log.SetOutput(logfile)
	log.SetOutput(io.MultiWriter(logfile, os.Stdout))

	//	サーバー設定を読み込む
	exerr := exsrapi.LoadConfig("ServerConfig.yml", &UpdateUserInfLib.Dbconfig)
	if exerr != nil {
		log.Printf("LoadConfig: %s\n", exerr.Error())
		return
	}

	dbconfig := UpdateUserInfLib.Dbconfig

	log.Printf("%+v\n", *dbconfig)

	log.Printf("\n")
	log.Printf("\n")
	log.Printf("********** Dbhost=<%s> Dbname = <%s> Dbuser = <%s> Dbpw = <%s>\n", (*dbconfig).Dbhost, (*dbconfig).Dbname, (*dbconfig).Dbuser, (*dbconfig).Dbpw)

	//	データベースとの接続をオープンする。
	status := UpdateUserInfLib.OpenDb()
	if status != 0 {
		log.Printf("Database error.\n")
		return
	}
	defer UpdateUserInfLib.Db.Close()

	/*
	if len(os.Args) != 2 {
		log.Printf("usage: %s eventid\n", os.Args[0])
		os.Exit(1)
	}
	*/

	eventlist, status :=UpdateUserInfLib.SelectLastEventList()
	if status != 0 {
		log.Printf("status=%d.\n", status)
		os.Exit(1)
	}

	for _, event := range eventlist {
		if event.Endtime.Before(time.Now()) {
			break
		}
	UpdateUserRankInf(event.EventID)
	}


}

//	イベント参加者のリストを作る
func SelectUsernoFromEventuser(eventid string) (roominflist []UpdateUserInfLib.RoomInfo, status int) {

	var stmt *sql.Stmt
	var rows *sql.Rows

	status = 0

	sql := "select userno from eventuser where eventid = ? "
	stmt, UpdateUserInfLib.Err = UpdateUserInfLib.Db.Prepare(sql)
	if UpdateUserInfLib.Err != nil {
		log.Printf("SelectUsernoFromEventuser() (1) err=%s\n", UpdateUserInfLib.Err.Error())
		status = 1
		return
	}
	defer stmt.Close()

	rows, UpdateUserInfLib.Err = stmt.Query(eventid)
	if UpdateUserInfLib.Err != nil {
		log.Printf("SelectUsernoFromEventuser() (2) err=%s\n", UpdateUserInfLib.Err.Error())
		status = 2
		return
	}
	defer rows.Close()

	var roominf UpdateUserInfLib.RoomInfo

	for rows.Next() {
		UpdateUserInfLib.Err = rows.Scan(&roominf.Userno)
		if UpdateUserInfLib.Err = rows.Err(); UpdateUserInfLib.Err != nil {
			log.Printf("SelectUsernoFromEventuser() (3) err=%s\n", UpdateUserInfLib.Err.Error())
			status = 3
			return
		}
		roominflist = append(roominflist, roominf)
	}
	if UpdateUserInfLib.Err = rows.Err(); UpdateUserInfLib.Err != nil {
		log.Printf("SelectUsernoFromEventuser() (4) err=%s\n", UpdateUserInfLib.Err.Error())
		status = 4
		return
	}

	return
}

//	イベント参加者のリストを作る
func SelectUserRankInfFromUser(roominflist *[]UpdateUserInfLib.RoomInfo) (status int) {

	var stmt *sql.Stmt

	status = 0

	sql := "select user_name, genre, `rank`, nrank, prank, level, followers, fans, fans_lst from user where userno = ? "
	stmt, UpdateUserInfLib.Err = UpdateUserInfLib.Db.Prepare(sql)
	if UpdateUserInfLib.Err != nil {
		log.Printf("SelectUserRankInfFromUser() (1) err=%s\n", UpdateUserInfLib.Err.Error())
		status = 1
		return
	}
	defer stmt.Close()

	for i := range *roominflist {
		UpdateUserInfLib.Err = stmt.QueryRow((*roominflist)[i].Userno).Scan(
			&(*roominflist)[i].Name,
			&(*roominflist)[i].Genre,
			&(*roominflist)[i].Rank,
			&(*roominflist)[i].Nrank,
			&(*roominflist)[i].Prank,
			&(*roominflist)[i].Level,
			&(*roominflist)[i].Followers,
			&(*roominflist)[i].Fans,
			&(*roominflist)[i].Fans_lst,
		)
		if UpdateUserInfLib.Err != nil {
			log.Printf("SelectUserRankInfFromUser() (3) err=%s\n", UpdateUserInfLib.Err.Error())
			status = 3
			return
		}
		//	log.Printf("%+v\n", (*roominflist)[i])
	}

	return
}

func CompareUserAndUpdate(roominflist []UpdateUserInfLib.RoomInfo) (status int) {

	status = 0

	ts := time.Now().Truncate(time.Second)
	log.Printf(" CompareUserAndUpdate()  ts= %s\n", ts.Format("2006-01-02 15:04"))

	for _, roominf := range roominflist {

		genre, rank, nrank, prank, level, followers, fans, fans_lst, _, _, _, tstatus := UpdateUserInfLib.GetRoomInfoByAPI(fmt.Sprintf("%d", roominf.Userno))
		if tstatus != 0 {
			continue
		}

		if roominf.Genre != genre ||
			roominf.Rank != rank ||
			roominf.Nrank != nrank ||
			roominf.Prank != prank ||
			roominf.Followers != followers ||
			roominf.Level != level ||
			roominf.Fans != fans ||
			roominf.Fans_lst != fans_lst {

			log.Printf("userno=%d %s|%s, %s|%s, %s|%s, %s|%s, %d|%d, %d|%d, %d|%d, %d|%d\n",
				roominf.Userno,
				roominf.Genre, genre,
				roominf.Rank, rank,
				roominf.Nrank, nrank,
				roominf.Prank, prank,
				roominf.Followers, followers,
				roominf.Level, level,
				roominf.Fans, fans,
				roominf.Fans_lst, fans_lst,
			)

			InsertUserHistory(roominf, ts)

			roominf.Genre = genre
			roominf.Rank = rank
			roominf.Nrank = nrank
			roominf.Prank = prank
			roominf.Followers = followers
			roominf.Level = level
			roominf.Fans = fans
			roominf.Fans_lst = fans_lst

			status = UpdateUser(roominf, ts)
			if status != 0 {
				log.Printf(" UpdateUser() returned status=%d", status)
				return
			}
		} else {
			log.Printf("userno=%d ==\n", roominf.Userno)
		}

	}

	return

}

func InsertUserHistory(roominf UpdateUserInfLib.RoomInfo, ts time.Time) (status int) {

	var stmt *sql.Stmt

	sql := "INSERT INTO userhistory(userno, user_name, genre, `rank`, nrank, prank, level, followers, fans, fans_lst, ts)"
	sql += " VALUES(?,?,?,?,?,?,?,?,?,?,?)"
	//	log.Printf("sql=%s\n", sql)
	stmt, UpdateUserInfLib.Err = UpdateUserInfLib.Db.Prepare(sql)
	if UpdateUserInfLib.Err != nil {
		log.Printf("error(InsertUserHistory) err=%s\n", UpdateUserInfLib.Err.Error())
		status = -1
		return
	}

	_, UpdateUserInfLib.Err = stmt.Exec(
		roominf.Userno,
		roominf.Name,
		roominf.Genre,
		roominf.Rank,
		roominf.Nrank,
		roominf.Prank,
		roominf.Level,
		roominf.Followers,
		roominf.Fans,
		roominf.Fans_lst,
		ts,
	)

	if UpdateUserInfLib.Err != nil {
		log.Printf("error(InsertUserHistory) err=%s\n", UpdateUserInfLib.Err.Error())
		status = -2
	}

	return

}

func UpdateUser(roominf UpdateUserInfLib.RoomInfo, ts time.Time) (status int) {

	var stmt *sql.Stmt


	sql := "update user set "
	sql += " genre=?,"
	sql += " `rank`=?,"
	sql += " nrank=?,"
	sql += " prank=?,"
	sql += " level=?,"
	sql += " followers=?,"
	sql += " fans=?,"
	sql += " fans_lst=?,"
	sql += " ts=?"
	sql += " where userno=?"
	stmt, UpdateUserInfLib.Err = UpdateUserInfLib.Db.Prepare(sql)

	if UpdateUserInfLib.Err != nil {
		log.Printf("UpdateUser() (1) err=%s\n", UpdateUserInfLib.Err.Error())
		log.Printf("  sql=%s\n", sql)
		status = -1
		return
	}
	_, UpdateUserInfLib.Err = stmt.Exec(
		roominf.Genre,
		roominf.Rank,
		roominf.Nrank,
		roominf.Prank,
		roominf.Level,
		roominf.Followers,
		roominf.Fans,
		roominf.Fans_lst,
		ts,
		roominf.Userno,
	)

	if UpdateUserInfLib.Err != nil {
		log.Printf("UpdateUser() (2) err=%s\n", UpdateUserInfLib.Err.Error())
		status = -2
	}

	return

}

func UpdateUserRankInf(eventid string) (status int) {

	status = 0

	//	イベントの詳細情報を得る
	UpdateUserInfLib.Event_inf, status = UpdateUserInfLib.SelectEventInf(eventid)
	if status != 0 {
		return
	}

	eventname := UpdateUserInfLib.Event_inf.Event_name
	period := UpdateUserInfLib.Event_inf.Period
	log.Printf(" period=%s evnentid=%s eventname=%s\n", period, eventid, eventname)

	//	イベント参加者のリストを作る
	var roominflist []UpdateUserInfLib.RoomInfo
	roominflist, status = SelectUsernoFromEventuser(eventid)
	log.Printf("  SelectUsernoFromEventuser() returned status=%d\n", status)

	//	ユーザー（ランク）情報をデータベースから取得する
	status = SelectUserRankInfFromUser(&roominflist)
	log.Printf("  SelectUserRankInfFromUser() returned status=%d\n", status)

	status = CompareUserAndUpdate(roominflist)
	log.Printf(" 	CompareUserAndUpdate()  returned status=%d\n", status)
	return

}
