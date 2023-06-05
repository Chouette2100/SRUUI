/*!
Copyright © 2022 chouette.21.00@gmail.com
Released under the MIT license
https://opensource.org/licenses/mit-license.php
Ver. 0.1.0 ApiLiveCurrentUser.goをsrapi.goから分離する。ApiLiveCurrentUser()のRoomIDをstring型に変更する。
*/
package UpdateUserInfLib

import (
	"fmt"
	"log"

	//	"math"
	"sort"
	"strconv"
	"strings"
	"time"

	//	"bufio"
	"os"

	//	"runtime"

	"encoding/json"

	//	"html/template"
	"net/http"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"

	"github.com/PuerkitoBio/goquery"

	//	svg "github.com/ajstarks/svgo/float"

	"github.com/dustin/go-humanize"

)

/*

	0100L1	安定版（～2021.12.26）
	0100M0	vscodeでの指摘箇所の修正
	0101A0	LinuxとMySQL8.0に対応する。
	0101B0	OSとWebサーバに応じた処理を行うようにする。アクセスログを作成する。
	0101B1	実行時パラメータをファイルから与えるように変更する。
	0101C0	GetRoomInfoByAPI()に配信開始時刻の取得を追加する。
	0101D0	詳細なランク情報の導入（Nrank）
	0101D1	"Next Live"の表示を追加する。
	0101D2	GetScoreEvery5Minutes RU20E4 に適合するバージョン
	0101D3	ランクをshow_rank_subdividedからleague_labe + lshow_rank_subdivided にする。
	0101E1	環境設定ファイルをyaml形式に変更する。
	0101G0	配信枠別貢献ポイントを導入する。
	0101G1	list-last.gtplでは維新枠別貢献ポイントの記録があるルームのみリンクを作成する。
	0101G2	list-last.gtplにジャンルを追加した。
	0101G3	リスナー貢献ポイントの履歴の表示(list-cntrbH)を作成する。
	0101G4	一つの貢献ポイントランキングの表示(list-cntrbS)を作成する(リスナー名の突き合わせのチェックが主目的)
	0101G5	list-lasth.gtplのリロード予告の表示でデータ取得間隔が5分と固定されていたものを設定値に合わせるように変更する。
	0101G6	ModminがIntervalminに対して不適切な値のときは修正して保存する。
	0101G7	ランクに関しnext_scoreに加えprev_scoreの表示を追加する。ファンの数の表示を追加する。
	0101H0	ファンレベル(/fanlevel)に関する画面を追加する。
	0101J0	ファンダム王イベント参加者のファン数ランキングを作成する。
	0101J2	終了したイベントについては無条件にルーム詳細情報（ランキング、フォロワ、レベル、ファン数）を出力しない。
	0101J2a	"ルーム詳細情報”の説明を追加した。
	0101J3	イベントリストのRoom_IDに変えてルーム名を表示する。表示数を6から10にする。
	0101J4	NewDocument()をNewDocumentFromReader()に変更する。list-last_h.gtplにルーム情報詳細表示/非表示のボタンを追加する。
	0101J5	イベント選択（最近のイベント）にnobasisが0のイベントも表示する（テーブルuserにusernoが0のデータを追加することが必要）
	10AA00	枠別貢献ポイントの「目標値(推定)」を追加する。
	10AB00	枠別貢献ポイントのポイント、増分の表示でhumanaizeを使用する。リスナー別貢献ポイント履歴に達成状況欄を追加する。
	10AB01	WebserverをDbconfig.Webserverに置き換える。枠別貢献ポイントのGTPLを変更する。GTPLのCRLFをLFに変更する。
	-----------------------------------------------------------------------------------
	ここからUpdateUserInfLibとして
	010AC00	SelectAcqTimeList()を追加する。
	010AD00	UpdateUserInfLibから不要な関数を削除する。
	010AE00	Room_url_keyから除外する文字列を"/"から"/r/"に変更する。

*/

const Version = "010AE00"

type Event_Inf struct {
	Event_ID    string
	I_Event_ID  int
	Event_name  string
	Event_no    int
	MaxPoint    int
	Start_time  time.Time
	Sstart_time string
	Start_date  float64
	End_time    time.Time
	Send_time   string
	Period      string
	Dperiod     float64
	Intervalmin int
	Modmin      int
	Modsec      int
	Fromorder   int
	Toorder     int
	Resethh     int
	Resetmm     int
	Nobasis     int
	Maxdsp      int
	NoEntry     int
	NoRoom      int    //	ルーム数
	EventStatus string //	"Over", "BeingHeld", "NotHeldYet"
	Pntbasis    int
	Ordbasis    int
	League_ids  string
	Cmap        int
	Target      int
	Maxpoint    int
	//	Status		string		//	"Confirmed":	イベント終了日翌日に確定した獲得ポイントが反映されている。
}



type RoomLevel struct {
	User_name string
	Genre     string
	Rank      string
	Nrank     string
	Prank     string
	Level     int
	Followers int
	Fans      int
	Fans_lst  int
	ts        time.Time
	Sts       string
}

type RoomLevelInf struct {
	Userno        int
	User_name     string
	RoomLevelList []RoomLevel
}


type ColorInf struct {
	Color      string
	Colorvalue string
	Selected   string
}

type ColorInfList []ColorInf

type RoomInfo struct {
	Name      string //	ルーム名のリスト
	Longname  string
	Shortname string
	Account   string //	アカウントのリスト、アカウントは配信のURLの最後の部分の英数字です。
	ID        string //	IDのリスト、IDはプロフィールのURLの最後の部分で5～6桁の数字です。
	Userno    int
	//	APIで取得できるデータ(1)
	Genre      string
	Rank       string
	Irank      int
	Nrank      string
	Prank      string
	Followers  int
	Sfollowers string
	Fans       int
	Fans_lst   int
	Level      int
	Slevel     string
	//	APIで取得できるデータ(2)
	Order        int
	Point        int //	イベント終了後12時間〜36時間はイベントページから取得できることもある
	Spoint       string
	Istarget     string
	Graph        string
	Iscntrbpoint string
	Color        string
	Colorvalue   string
	Colorinflist ColorInfList
	Formid       string
	Eventid      string
	Status       string
	Statuscolor  string
}

type RoomInfoList []RoomInfo

//	sort.Sort()のための関数三つ
func (r RoomInfoList) Len() int {
	return len(r)
}

func (r RoomInfoList) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

func (r RoomInfoList) Choose(from, to int) (s RoomInfoList) {
	s = r[from:to]
	return
}

var SortByFollowers bool

//	降順に並べる
func (r RoomInfoList) Less(i, j int) bool {
	//	return e[i].point < e[j].point
	if SortByFollowers {
		return r[i].Followers > r[j].Followers
	} else {
		return r[i].Point > r[j].Point
	}
}

type DBConfig struct {
	WebServer string `yaml:"WebServer"`
	HTTPport  string `yaml:"HTTPport"`
	SSLcrt    string `yaml:"SSLcrt"`
	SSLkey    string `yaml:"SSLkey"`
	Dbhost    string `yaml:"Dbhost"`
	Dbname    string `yaml:"Dbname"`
	Dbuser    string `yaml:"Dbuser"`
	Dbpw      string `yaml:"Dbpw"`
}

var Dbconfig *DBConfig

var Event_inf Event_Inf

var Db *sql.DB
var Err error

var OS string

//	var WebServer string

type Color struct {
	Name  string
	Value string
}

//	https://www.fukushihoken.metro.tokyo.lg.jp/kiban/machizukuri/kanren/color.files/colorudguideline.pdf
var Colorlist2 []Color = []Color{
	{"red", "#FF2800"},
	{"yellow", "#FAF500"},
	{"green", "#35A16B"},
	{"blue", "#0041FF"},
	{"skyblue", "#66CCFF"},
	{"lightpink", "#FFD1D1"},
	{"orange", "#FF9900"},
	{"purple", "#9A0079"},
	{"brown", "#663300"},
	{"lightgreen", "#87D7B0"},
	{"white", "#FFFFFF"},
	{"gray", "#77878F"},
}

var Colorlist1 []Color = []Color{
	{"cyan", "cyan"},
	{"magenta", "magenta"},
	{"yellow", "yellow"},
	{"royalblue", "royalblue"},
	{"coral", "coral"},
	{"khaki", "khaki"},
	{"deepskyblue", "deepskyblue"},
	{"crimson", "crimson"},
	{"orange", "orange"},
	{"lightsteelblue", "lightsteelblue"},
	{"pink", "pink"},
	{"sienna", "sienna"},
	{"springgreen", "springgreen"},
	{"blueviolet", "blueviolet"},
	{"salmon", "salmon"},
	{"lime", "lime"},
	{"red", "red"},
	{"darkorange", "darkorange"},
	{"skyblue", "skyblue"},
	{"lightpink", "lightpink"},
}

type Event struct {
	EventID   string
	EventName string
	Period    string
	Starttime time.Time
	S_start   string
	Endtime   time.Time
	S_end     string
	Status    string
	Pntbasis  int
	Pbname    string
	Selected  string
}

type User struct {
	Userno       int
	Userlongname string
	Selected     string
}



func GetUserInfForHistory() (status int) {

	status = 0

	//	select distinct(nobasis) from event
	stmt, err := Db.Prepare("select distinct(nobasis) from event")
	if err != nil {
		//	log.Fatal(err)
		log.Printf("err=[%s]\n", err.Error())
		status = -1
		return
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		//	log.Fatal(err)
		log.Printf("err=[%s]\n", err.Error())
		status = -1
		return
	}
	defer rows.Close()

	var roominf RoomInfo
	var roominflist RoomInfoList

	for rows.Next() {
		err := rows.Scan(&roominf.Userno)
		if err != nil {
			//	log.Fatal(err)
			log.Printf("err=[%s]\n", err.Error())
			status = -1
			return
		}
		if roominf.Userno != 0 {
			roominf.ID = fmt.Sprintf("%d", roominf.Userno)
			roominflist = append(roominflist, roominf)
		}
	}
	if err = rows.Err(); err != nil {
		//	log.Fatal(err)
		log.Printf("err=[%s]\n", err.Error())
		status = -1
		return
	}

	eventid := ""

	//	Update user , Insert into userhistory
	for _, roominf := range roominflist {

		sql := "select currentevent from user where userno = ?"
		err := Db.QueryRow(sql, roominf.Userno).Scan(&eventid)
		if err != nil {
			log.Printf("err=[%s]\n", err.Error())
			status = -1
		}

		roominf.Genre, roominf.Rank, roominf.Nrank, roominf.Prank, roominf.Level,
			roominf.Followers, roominf.Fans, roominf.Fans_lst, roominf.Name, roominf.Account, _, status = GetRoomInfoByAPI(roominf.ID)
		InsertIntoOrUpdateUser(time.Now().Truncate(time.Second), eventid, roominf)
	}

	return
}

func GetEventListByAPI(eventinflist *[]Event_Inf) (status int) {

	status = 0

	last_page := 1
	total_count := 1

	for page := 1; page <= last_page; page++ {

		URL := "https://www.showroom-live.com/api/event/search?page=" + fmt.Sprintf("%d", page)
		log.Printf("GetEventListByAPI() URL=%s\n", URL)

		resp, err := http.Get(URL)
		if err != nil {
			//	一時的にデータが取得できない。
			log.Printf("GetEventListByAPI() err=%s\n", err.Error())
			//		panic(err)
			status = -1
			return
		}
		defer resp.Body.Close()

		//	JSONをデコードする。
		//	次の記事を参考にさせていただいております。
		//		Go言語でJSONに泣かないためのコーディングパターン
		//		https://qiita.com/msh5/items/dc524e38073ed8e3831b

		var result interface{}
		decoder := json.NewDecoder(resp.Body)
		if err := decoder.Decode(&result); err != nil {
			log.Printf("GetEventListByAPI() err=%s\n", err.Error())
			//	panic(err)
			status = -2
			return
		}

		if page == 1 {
			value, _ := result.(map[string]interface{})["last_page"].(float64)
			last_page = int(value)
			value, _ = result.(map[string]interface{})["total_count"].(float64)
			total_count = int(value)
			log.Printf("GetEventListByAPI() total_count=%d, last_page=%d\n", total_count, last_page)
		}

		noroom := 30
		if page == last_page {
			noroom = total_count % 30
			if noroom == 0 {
				noroom = 30
			}
		}

		for i := 0; i < noroom; i++ {
			var eventinf Event_Inf

			tres := result.(map[string]interface{})["event_list"].([]interface{})[i]

			ttres := tres.(map[string]interface{})["league_ids"]
			norec := len(ttres.([]interface{}))
			if norec == 0 {
				continue
			}
			log.Printf("norec =%d\n", norec)
			eventinf.League_ids = ""
			/*
				for j := 0; j < norec; j++ {
					eventinf.League_ids += ttres.([]interface{})[j].(string) + ","
				}
			*/
			eventinf.League_ids = ttres.([]interface{})[norec-1].(string)
			if eventinf.League_ids != "60" {
				continue
			}

			eventinf.Event_ID, _ = tres.(map[string]interface{})["event_url_key"].(string)
			eventinf.Event_name, _ = tres.(map[string]interface{})["event_name"].(string)
			//	log.Printf("id=%s, name=%s\n", eventinf.Event_ID, eventinf.Event_name)

			started_at, _ := tres.(map[string]interface{})["started_at"].(float64)
			eventinf.Start_time = time.Unix(int64(started_at), 0)
			eventinf.Sstart_time = eventinf.Start_time.Format("06/01/02 15:04")
			ended_at, _ := tres.(map[string]interface{})["ended_at"].(float64)
			eventinf.End_time = time.Unix(int64(ended_at), 0)
			eventinf.Send_time = eventinf.End_time.Format("06/01/02 15:04")

			(*eventinflist) = append((*eventinflist), eventinf)

		}

		//	resp.Body.Close()
	}

	return
}

//	idで指定した配信者さんの獲得ポイントを取得する。
//	戻り値は 獲得ポイント、順位、上位とのポイント差（1位の場合は2位とのポイント差）、イベント名
//	レベルイベントのときは順位、上位とのポイント差は0がセットされる。
func GetPointsByAPI(id string) (Point, Rank, Gap int, EventID string) {

	//	獲得ポイントなどの配信者情報を得るURL（このURLについては記事参照）
	URL := "https://www.showroom-live.com/api/room/event_and_support?room_id=" + id

	resp, err := http.Get(URL)
	if err != nil {
		//	一時的にデータが取得できない。
		//		panic(err)
		return 0, 0, 0, "**Error** http.Get(URL)"
	}
	defer resp.Body.Close()

	//	JSONをデコードする。
	//	次の記事を参考にさせていただいております。
	//		Go言語でJSONに泣かないためのコーディングパターン
	//		https://qiita.com/msh5/items/dc524e38073ed8e3831b

	var result interface{}
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&result); err != nil {
		//	panic(err)
		return 0, 0, 0, "**Error** http.Get(URL)"
	}

	//	イベントが終わっている、イベント参加をとりやめた、SHOWROOMをやめた、などの対応
	if result.(map[string]interface{})["event"] == nil {
		return 0, 0, 0, "not held yet./over./not entry."
	}

	if result.(map[string]interface{})["event"].(map[string]interface{})["ranking"] != nil {
		//	ランキングのあるイベントの場合
		//	（順位に応じて特典が与えられるイベント、ただし獲得ポイントに対して特典が与えられる場合でも順位付けがある場合はこちら）

		//	獲得ポイント
		l, _ := result.(map[string]interface{})["event"].(map[string]interface{})["ranking"].(map[string]interface{})["point"].(float64)
		//	順位
		m, _ := result.(map[string]interface{})["event"].(map[string]interface{})["ranking"].(map[string]interface{})["rank"].(float64)
		//	ポイント差
		n, _ := result.(map[string]interface{})["event"].(map[string]interface{})["ranking"].(map[string]interface{})["gap"].(float64)

		Point = int(l)
		Rank = int(m)
		Gap = int(n)

		//	イベント名
		EventID, _ = result.(map[string]interface{})["event"].(map[string]interface{})["event_url"].(string)
		EventID = strings.Replace(EventID, "https://www.showroom-live.com/event/", "", -1)

	} else if result.(map[string]interface{})["event"].(map[string]interface{})["quest"] != nil {
		//	レベルイベント（ランキングのないイベント）の場合
		//	（アバ権やステッカーなど獲得ポイントに応じて特典が与えられるイベント、ただし順位付けがある場合は除く）

		//	獲得ポイント
		l, _ := result.(map[string]interface{})["event"].(map[string]interface{})["quest"].(map[string]interface{})["support"].(map[string]interface{})["current_point"].(float64)
		//	順位
		m := 0.0
		//	ポイント差
		n := 0.0

		Point = int(l)
		Rank = int(m)
		Gap = int(n)

		//	イベント名
		EventID, _ = result.(map[string]interface{})["event"].(map[string]interface{})["event_url"].(string)
		EventID = strings.Replace(EventID, "https://www.showroom-live.com/event/", "", -1)

	} else {
		//	上記ランキングイベントでもレベルイベントでもない場合
		//	もしこのようなケースが存在するならJSONを確認して新たにコーディングする
		log.Println(" N/A")
		return 0, 0, 0, "N/A"
	}

	return
}

/*

 */
func GetIsOnliveByAPI(room_id string) (
	isonlive bool, //	true:	配信中
	startedat time.Time, //	配信開始時刻（isonliveがtrueのときだけ意味があります）
	status int,
) {

	status = 0

	//	https://qiita.com/takeru7584/items/f4ba4c31551204279ed2
	url := "https://www.showroom-live.com/api/room/profile?room_id=" + room_id

	resp, err := http.Get(url)
	if err != nil {
		//	一時的にデータが取得できない。
		//	resp.Body.Close()
		//		panic(err)
		status = -1
		return
	}
	defer resp.Body.Close()

	//	JSONをデコードする。
	//	次の記事を参考にさせていただいております。
	//		Go言語でJSONに泣かないためのコーディングパターン
	//		https://qiita.com/msh5/items/dc524e38073ed8e3831b

	var result interface{}
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&result); err != nil {
		//	panic(err)
		status = -2
		return
	}

	//	配信中か？
	isonlive, _ = result.(map[string]interface{})["is_onlive"].(bool)

	if isonlive {
		//	配信開始時刻の取得
		value, _ := result.(map[string]interface{})["current_live_started_at"].(float64)
		startedat = time.Unix(int64(value), 0).Truncate(time.Second)
		//	log.Printf("current_live_stared_at %f %v\n", value, startedat)
	}

	return

}

func GetAciveFanByAPI(room_id string, yyyymm string) (nofan int) {

	nofan = -1

	url := "https://www.showroom-live.com/api/active_fan/room?room_id=" + room_id + "&ym=" + yyyymm

	resp, err := http.Get(url)
	if err != nil {
		//	一時的にデータが取得できない。
		//	resp.Body.Close()
		//		panic(err)
		nofan = -1
		return
	}
	defer resp.Body.Close()

	//	JSONをデコードする。
	//	次の記事を参考にさせていただいております。
	//		Go言語でJSONに泣かないためのコーディングパターン
	//		https://qiita.com/msh5/items/dc524e38073ed8e3831b

	var result interface{}
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&result); err != nil {
		//	panic(err)
		nofan = -2
		return
	}

	value, _ := result.(map[string]interface{})["total_user_count"].(float64)
	nofan = int(value)

	return
}
func GetRoomInfoByAPI(room_id string) (
	genre string,
	rank string,
	nrank string,
	prank string,
	level int,
	followers int,
	fans int,
	fans_lst int,
	roomname string,
	roomurlkey string,
	startedat time.Time,
	status int,
) {

	status = 0

	//	https://qiita.com/takeru7584/items/f4ba4c31551204279ed2
	url := "https://www.showroom-live.com/api/room/profile?room_id=" + room_id

	resp, err := http.Get(url)
	if err != nil {
		//	一時的にデータが取得できない。
		//	resp.Body.Close()
		//		panic(err)
		status = -1
		return
	}
	defer resp.Body.Close()

	//	JSONをデコードする。
	//	次の記事を参考にさせていただいております。
	//		Go言語でJSONに泣かないためのコーディングパターン
	//		https://qiita.com/msh5/items/dc524e38073ed8e3831b

	var result interface{}
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&result); err != nil {
		//	panic(err)
		status = -2
		return
	}

	value, _ := result.(map[string]interface{})["follower_num"].(float64)
	followers = int(value)

	tnow := time.Now()
	fans = GetAciveFanByAPI(room_id, tnow.Format("200601"))
	yy := tnow.Year()
	mm := tnow.Month() - 1
	if mm < 0 {
		yy -= 1
		mm = 12
	}
	fans_lst = GetAciveFanByAPI(room_id, fmt.Sprintf("%04d%02d", yy, mm))

	genre, _ = result.(map[string]interface{})["genre_name"].(string)

	rank, _ = result.(map[string]interface{})["league_label"].(string)
	ranks, _ := result.(map[string]interface{})["show_rank_subdivided"].(string)
	rank = rank + " | " + ranks

	value, _ = result.(map[string]interface{})["next_score"].(float64)
	nrank = humanize.Comma(int64(value))
	value, _ = result.(map[string]interface{})["prev_score"].(float64)
	prank = humanize.Comma(int64(value))

	value, _ = result.(map[string]interface{})["room_level"].(float64)
	level = int(value)

	roomname, _ = result.(map[string]interface{})["room_name"].(string)

	roomurlkey, _ = result.(map[string]interface{})["room_url_key"].(string)

	//	配信開始時刻の取得
	value, _ = result.(map[string]interface{})["current_live_started_at"].(float64)
	startedat = time.Unix(int64(value), 0).Truncate(time.Second)
	//	log.Printf("current_live_stared_at %f %v\n", value, startedat)

	return

}

func GetNextliveByAPI(room_id string) (
	nextlive string,
	status int,
) {

	status = 0

	//	https://qiita.com/takeru7584/items/f4ba4c31551204279ed2
	url := "https://www.showroom-live.com/api/room/next_live?room_id=" + room_id

	resp, err := http.Get(url)
	if err != nil {
		//	一時的にデータが取得できない。
		//	resp.Body.Close()
		//		panic(err)
		status = -1
		return
	}
	defer resp.Body.Close()

	//	JSONをデコードする。
	//	次の記事を参考にさせていただいております。
	//		Go言語でJSONに泣かないためのコーディングパターン
	//		https://qiita.com/msh5/items/dc524e38073ed8e3831b

	var result interface{}
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&result); err != nil {
		//	panic(err)
		status = -2
		return
	}

	nextlive, _ = result.(map[string]interface{})["text"].(string)

	return

}

func SelectRoomInf(
	userno int,
) (
	roominf RoomInfo,
	status int,
) {

	status = 0

	sql := "select distinct u.userno, userid, user_name, longname, shortname, genre, nrank, prank, level, followers, fans, fans_lst, e.istarget,e.graph, e.color, e.iscntrbpoints, e.point "
	sql += " from user u join eventuser e "
	//	sql += " where u.userno = e.userno and u.userno = " + fmt.Sprintf("%d", userno)
	sql += " where u.userno = e.userno and u.userno = ? "

	stmt, err := Db.Prepare(sql)
	if err != nil {
		log.Printf("SelectRoomInf() Prepare() err=%s\n", err.Error())
		status = -5
		return
	}
	defer stmt.Close()

	err = stmt.QueryRow(userno).Scan(&roominf.Userno,
		&roominf.Account,
		&roominf.Name,
		&roominf.Longname,
		&roominf.Shortname,
		&roominf.Genre,
		&roominf.Nrank,
		&roominf.Prank,
		&roominf.Level,
		&roominf.Followers,
		&roominf.Fans,
		&roominf.Fans_lst,
		&roominf.Istarget,
		&roominf.Graph,
		&roominf.Color,
		&roominf.Iscntrbpoint,
		&roominf.Point,
	)
	if err != nil {
		log.Printf("SelectRoomInf() Query() (6) err=%s\n", err.Error())
		status = -6
		return
	}
	if roominf.Istarget == "Y" {
		roominf.Istarget = "Checked"
	} else {
		roominf.Istarget = ""
	}
	if roominf.Graph == "Y" {
		roominf.Graph = "Checked"
	} else {
		roominf.Graph = ""
	}
	if roominf.Iscntrbpoint == "Y" {
		roominf.Iscntrbpoint = "Checked"
	} else {
		roominf.Iscntrbpoint = ""
	}
	roominf.Slevel = humanize.Comma(int64(roominf.Level))
	roominf.Sfollowers = humanize.Comma(int64(roominf.Followers))
	roominf.Spoint = humanize.Comma(int64(roominf.Point))
	roominf.Name = strings.ReplaceAll(roominf.Name, "'", "’")

	return
}

func SelectEventRoomInfList(
	eventid string,
	roominfolist *RoomInfoList,
) (
	eventname string,
	status int,
) {

	status = 0

	//	eventno := 0
	//	eventno, eventname, _ = SelectEventNoAndName(eventid)
	Event_inf, _ = SelectEventInf(eventid)

	//	eventno := Event_inf.Event_no
	eventname = Event_inf.Event_name

	sql := "select distinct u.userno, userid, user_name, longname, shortname, genre, `rank`, nrank, prank, level, followers, fans, fans_lst, e.istarget, e.graph, e.color, e.iscntrbpoints, e.point "
	sql += " from user u join eventuser e "
	sql += " where u.userno = e.userno and e.eventid= ?"
	if Event_inf.Start_time.After(time.Now()) {
		sql += " order by followers desc"
	} else {
		sql += " order by e.point desc"
	}

	stmt, err := Db.Prepare(sql)
	if err != nil {
		log.Printf("SelectEventRoomInfList() Prepare() err=%s\n", err.Error())
		status = -5
		return
	}
	defer stmt.Close()

	rows, err := stmt.Query(eventid)
	if err != nil {
		log.Printf("SelectRoomIn() Query() (6) err=%s\n", err.Error())
		status = -6
		return
	}
	defer rows.Close()

	ColorlistA := Colorlist2
	ColorlistB := Colorlist1
	if Event_inf.Cmap == 1 {
		ColorlistA = Colorlist1
		ColorlistB = Colorlist2
	}

	colormap := make(map[string]int)

	for i := 0; i < len(ColorlistA); i++ {
		colormap[ColorlistA[i].Name] = i
	}

	var roominf RoomInfo

	i := 0
	for rows.Next() {
		err := rows.Scan(&roominf.Userno,
			&roominf.Account,
			&roominf.Name,
			&roominf.Longname,
			&roominf.Shortname,
			&roominf.Genre,
			&roominf.Rank,
			&roominf.Nrank,
			&roominf.Prank,
			&roominf.Level,
			&roominf.Followers,
			&roominf.Fans,
			&roominf.Fans_lst,
			&roominf.Istarget,
			&roominf.Graph,
			&roominf.Color,
			&roominf.Iscntrbpoint,
			&roominf.Point,
		)

		ci := 0
		for ; ci < len(ColorlistA); ci++ {
			if ColorlistA[ci].Name == roominf.Color {
				roominf.Colorvalue = ColorlistA[ci].Value
				break
			}
		}
		if ci == len(ColorlistA) {
			ci := 0
			for ; ci < len(ColorlistB); ci++ {
				if ColorlistB[ci].Name == roominf.Color {
					roominf.Colorvalue = ColorlistB[ci].Value
					break
				}
			}
			if ci == len(ColorlistB) {
				roominf.Colorvalue = roominf.Color
			}
		}

		if roominf.Istarget == "Y" {
			roominf.Istarget = "Checked"
		} else {
			roominf.Istarget = ""
		}
		if roominf.Graph == "Y" {
			roominf.Graph = "Checked"
		} else {
			roominf.Graph = ""
		}
		if roominf.Iscntrbpoint == "Y" {
			roominf.Iscntrbpoint = "Checked"
		} else {
			roominf.Iscntrbpoint = ""
		}
		roominf.Slevel = humanize.Comma(int64(roominf.Level))
		roominf.Sfollowers = humanize.Comma(int64(roominf.Followers))
		if roominf.Point < 0 {
			roominf.Spoint = ""
		} else {
			roominf.Spoint = humanize.Comma(int64(roominf.Point))
		}
		roominf.Formid = "Form" + fmt.Sprintf("%d", i)
		roominf.Eventid = eventid
		roominf.Name = strings.ReplaceAll(roominf.Name, "'", "’")
		if err != nil {
			log.Printf("SelectEventRoomInfList() Scan() err=%s\n", err.Error())
			status = -7
			return
		}
		//	var colorinf ColorInf
		colorinflist := make([]ColorInf, len(ColorlistA))

		for i := 0; i < len(ColorlistA); i++ {
			colorinflist[i].Color = ColorlistA[i].Name
			colorinflist[i].Colorvalue = ColorlistA[i].Value
		}

		roominf.Colorinflist = colorinflist
		if cidx, ok := colormap[roominf.Color]; ok {
			roominf.Colorinflist[cidx].Selected = "Selected"
		}
		*roominfolist = append(*roominfolist, roominf)

		i++
	}

	if err = rows.Err(); err != nil {
		log.Printf("SelectEventRoomInfList() rows err=%s\n", err.Error())
		status = -8
		return
	}

	if Event_inf.Start_time.After(time.Now()) {
		SortByFollowers = true
	} else {
		SortByFollowers = false
	}
	sort.Sort(*roominfolist)

	/*
		for i := 0; i < len(*roominfolist); i++ {

			sql = "select max(point) from points where "
			sql += " user_id = " + fmt.Sprintf("%d", (*roominfolist)[i].Userno)
			//	sql += " and event_id = " + fmt.Sprintf("%d", eventno)
			sql += " and event_id = " + eventid

			err = Db.QueryRow(sql).Scan(&(*roominfolist)[i].Point)
			(*roominfolist)[i].Spoint = humanize.Comma(int64((*roominfolist)[i].Point))

			if err == nil {
				continue
			} else {
				log.Printf("err=[%s]\n", err.Error())
				if err.Error() != "sql: no rows in result set" {
					eventno = -2
					continue
				} else {
					(*roominfolist)[i].Point = -1
					(*roominfolist)[i].Spoint = ""
				}
			}
		}
	*/

	return
}

func UpdateRoomInf(eventid, suserno, longname, shortname, istarget, graph, color, iscntrbpoint string) (status int) {

	status = 0

	userno, _ := strconv.Atoi(suserno)

	sql := "update user set longname=?, shortname=? where userno = ?"

	row, err := Db.Prepare(sql)
	if err != nil {
		log.Printf("UpdateRoomInf() error(Update/Prepare) err=%s\n", err.Error())
		status = -1
		return
	}

	_, err = row.Exec(longname, shortname, userno)

	if err != nil {
		log.Printf("UpdateRoomInf() error(InsertIntoOrUpdateUser() Update/Exec) err=%s\n", err.Error())
		status = -2
		return
	}

	//	eventno, _, _ := SelectEventNoAndName(eventid)

	if istarget == "1" {
		istarget = "Y"
	} else {
		istarget = "N"
	}

	if graph == "1" {
		graph = "Y"
	} else {
		graph = "N"
	}

	if iscntrbpoint == "1" {
		iscntrbpoint = "Y"
	} else {
		iscntrbpoint = "N"
	}

	//	sql = "update eventuser set istarget=?, graph=?, color=? where eventno=? and userno=?"
	sql = "update eventuser set istarget=?, graph=?, color=?, iscntrbpoints=? where eventid=? and userno=?"

	row, err = Db.Prepare(sql)
	if err != nil {
		log.Printf("UpdateRoomInf() error(Update/Prepare) err=%s\n", err.Error())
		status = -1
		return
	}

	_, err = row.Exec(istarget, graph, color, iscntrbpoint, eventid, userno)

	if err != nil {
		log.Printf("error(InsertIntoOrUpdateUser() Update/Exec) err=%s\n", err.Error())
		status = -2
	}

	return

}

func UpdateEventuserSetPoint(eventid, userid string, point int) (status int) {
	status = 0

	//	eventno, _, _ := SelectEventNoAndName(eventid)
	userno, _ := strconv.Atoi(userid)

	sql := "update eventuser set point=? where eventid = ? and userno = ?"
	row, err := Db.Prepare(sql)
	if err != nil {
		log.Printf("UpdateEventuserSetPoint() error (Update/Prepare) err=%s\n", err.Error())
		status = -1
		return
	}

	_, err = row.Exec(point, eventid, userno)

	if err != nil {
		log.Printf("error(UpdateEventuserSetPoint() Update/Exec) err=%s\n", err.Error())
		status = -2
	}

	return
}
func GetRoomInfoAndPoint(
	eventid string,
	roomid string,
	idbasis string,
) (
	roominf RoomInfo,
	status int,
) {

	status = 0

	roominf.ID = roomid
	roominf.Userno, _ = strconv.Atoi(roomid)

	Event_inf, _ = SelectEventInf(eventid)

	roominf.Genre, roominf.Rank, roominf.Nrank, roominf.Prank, roominf.Level, roominf.Followers,
		roominf.Fans,
		roominf.Fans_lst,
		roominf.Name, roominf.Account, _, status =
		GetRoomInfoByAPI(roomid)

	point, _, _, peventid := GetPointsByAPI(roominf.ID)
	if peventid == Event_inf.Event_ID {
		roominf.Point = point
		UpdateEventuserSetPoint(peventid, roominf.ID, point)
	} else {
		log.Printf(" %s %s %d\n", Event_inf.Event_ID, peventid, point)
	}

	/*
		if (*roominfolist)[i].ID == idbasis {
			(*eventinfo).Pntbasis = point
			(*eventinfo).Ordbasis = i
		}
	*/

	//	log.Printf(" followers=<%d> level=<%d> nrank=<%s> genre=<%s> point=%d\n",
	//	(*roominfolist)[i].Followers,
	//	(*roominfolist)[i].Level,
	//	(*roominfolist)[i].Nrank,
	//	(*roominfolist)[i].Genre,
	//	(*roominfolist)[i].Point)

	return
}

func GetAndInsertEventRoomInfo(
	eventid string,
	breg int,
	ereg int,
	eventinfo *Event_Inf,
	roominfolist *RoomInfoList,
) (
	starttimeafternow bool,
	status int,
) {

	log.Println("GetAndInsertEventRoomInfo() Called.")
	log.Println(*eventinfo)

	status = 0
	starttimeafternow = false

	//	イベントに参加しているルームの一覧を取得します。
	//	ルーム名、ID、URLを取得しますが、イベント終了直後の場合の最終獲得ポイントが表示されている場合はそれも取得します。
	status = GetEventInfAndRoomList(eventid, breg, ereg, eventinfo, roominfolist)

	if status != 0 {
		log.Printf("GetEventInfAndRoomList() returned %d\n", status)
		return
	}

	//	各ルームのジャンル、ランク、レベル、フォロワー数を取得します。
	for i := 0; i < (*eventinfo).NoRoom; i++ {
		(*roominfolist)[i].Genre, (*roominfolist)[i].Rank,
			(*roominfolist)[i].Nrank,
			(*roominfolist)[i].Prank,
			(*roominfolist)[i].Level,
			(*roominfolist)[i].Followers,
			(*roominfolist)[i].Fans,
			(*roominfolist)[i].Fans_lst,
			_, _, _, _ = GetRoomInfoByAPI((*roominfolist)[i].ID)

	}

	//	各ルームの獲得ポイントを取得します。
	for i := 0; i < (*eventinfo).NoRoom; i++ {
		point, _, _, eventid := GetPointsByAPI((*roominfolist)[i].ID)
		if eventid == (*eventinfo).Event_ID {
			(*roominfolist)[i].Point = point
			UpdateEventuserSetPoint(eventid, (*roominfolist)[i].ID, point)
			if point < 0 {
				(*roominfolist)[i].Spoint = ""
			} else {
				(*roominfolist)[i].Spoint = humanize.Comma(int64(point))
			}
		} else {
			log.Printf(" %s %s %d\n", (*eventinfo).Event_ID, eventid, point)
		}

		if (*roominfolist)[i].ID == fmt.Sprintf("%d", (*eventinfo).Nobasis) {
			(*eventinfo).Pntbasis = point
			(*eventinfo).Ordbasis = i
		}

		//	log.Printf(" followers=<%d> level=<%d> nrank=<%s> genre=<%s> point=%d\n",
		//	(*roominfolist)[i].Followers,
		//	(*roominfolist)[i].Level,
		//	(*roominfolist)[i].Nrank,
		//	(*roominfolist)[i].Genre,
		//	(*roominfolist)[i].Point)
	}

	if (*eventinfo).Start_time.After(time.Now()) {
		SortByFollowers = true
		sort.Sort(*roominfolist)
		if ereg > len(*roominfolist) {
			ereg = len(*roominfolist)
		}
		r := (*roominfolist).Choose(breg-1, ereg)
		roominfolist = &r
		starttimeafternow = true
	}

	log.Printf(" GetEventRoomInfo() len(*roominfolist)=%d\n", len(*roominfolist))

	log.Println("GetAndInsertEventRoomInfo() before InsertEventIinf()")
	log.Println(*eventinfo)
	status = InsertEventInf(eventinfo)

	if status == 1 {
		log.Println("InsertEventInf() returned 1.")
		UpdateEventInf(eventinfo)
		status = 0
	}
	log.Println("GetAndInsertEventRoomInfo() after InsertEventIinf() or UpdateEventInf")
	log.Println(*eventinfo)

	_, _, status = SelectEventNoAndName(eventid)

	if status == 0 {
		//	InsertRoomInf(eventno, eventid, roominfolist)
		InsertRoomInf(eventid, roominfolist)
	}

	return
}

func InsertEventInf(eventinf *Event_Inf) (
	status int,
) {

	if _, _, status = SelectEventNoAndName((*eventinf).Event_ID); status != 0 {
		sql := "INSERT INTO event(eventid, event_name, period, starttime, endtime, noentry,"
		sql += " intervalmin, modmin, modsec, "
		sql += " Fromorder, Toorder, Resethh, Resetmm, Nobasis, Maxdsp, Cmap, target, maxpoint "
		sql += ") VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"
		log.Printf("db.Prepare(sql)\n")
		row, err := Db.Prepare(sql)
		if err != nil {
			log.Printf("error InsertEventInf() (INSERT/Prepare) err=%s\n", err.Error())
			status = -1
			return
		}

		log.Printf("row.Exec()\n")
		_, err = row.Exec(
			(*eventinf).Event_ID,
			(*eventinf).Event_name,
			(*eventinf).Period,
			(*eventinf).Start_time,
			(*eventinf).End_time,
			(*eventinf).NoEntry,
			(*eventinf).Intervalmin,
			(*eventinf).Modmin,
			(*eventinf).Modsec,
			(*eventinf).Fromorder,
			(*eventinf).Toorder,
			(*eventinf).Resethh,
			(*eventinf).Resetmm,
			(*eventinf).Nobasis,
			(*eventinf).Maxdsp,
			(*eventinf).Cmap,
			(*eventinf).Target,
			(*eventinf).Maxpoint,
		)

		if err != nil {
			log.Printf("error InsertEventInf() (INSERT/Exec) err=%s\n", err.Error())
			status = -2
		}
	} else {
		status = 1
	}

	return
}

func UpdateEventInf(eventinf *Event_Inf) (
	status int,
) {

	if _, _, status = SelectEventNoAndName((*eventinf).Event_ID); status == 0 {
		sql := "Update event set "
		sql += " event_name=?,"
		sql += " period=?,"
		sql += " starttime=?,"
		sql += " endtime=?,"
		sql += " noentry=?,"
		sql += " intervalmin=?,"
		sql += " modmin=?,"
		sql += " modsec=?,"
		sql += " Fromorder=?,"
		sql += " Toorder=?,"
		sql += " Resethh=?,"
		sql += " Resetmm=?,"
		sql += " Nobasis=?,"
		sql += " Target=?,"
		sql += " Maxdsp=?, "
		sql += " cmap=?, "
		sql += " maxpoint=? "
		//	sql += " where eventno = ?"
		sql += " where eventid = ?"
		log.Printf("db.Prepare(sql)\n")
		row, err := Db.Prepare(sql)
		if err != nil {
			log.Printf("UpdateEventInf() error (Update/Prepare) err=%s\n", err.Error())
			status = -1
			return
		}

		log.Printf("row.Exec()\n")
		_, err = row.Exec(
			(*eventinf).Event_name,
			(*eventinf).Period,
			(*eventinf).Start_time,
			(*eventinf).End_time,
			(*eventinf).NoEntry,
			(*eventinf).Intervalmin,
			(*eventinf).Modmin,
			(*eventinf).Modsec,
			(*eventinf).Fromorder,
			(*eventinf).Toorder,
			(*eventinf).Resethh,
			(*eventinf).Resetmm,
			(*eventinf).Nobasis,
			(*eventinf).Target,
			(*eventinf).Maxdsp,
			(*eventinf).Cmap,
			(*eventinf).Maxpoint,
			(*eventinf).Event_ID,
		)

		if err != nil {
			log.Printf("error UpdateEventInf() (update/Exec) err=%s\n", err.Error())
			status = -2
		}
	} else {
		status = 1
	}

	return
}

func InsertRoomInf(eventid string, roominfolist *RoomInfoList) {

	log.Printf("***** InsertRoomInf() ***********  NoRoom=%d\n", len(*roominfolist))
	tnow := time.Now().Truncate(time.Second)
	for i := 0; i < len(*roominfolist); i++ {
		log.Printf("   ** InsertRoomInf() ***********  i=%d\n", i)
		InsertIntoOrUpdateUser(tnow, eventid, (*roominfolist)[i])
		status := InsertIntoEventUser(i, eventid, (*roominfolist)[i])
		if status == 0 {
			(*roominfolist)[i].Status = "更新"
			(*roominfolist)[i].Statuscolor = "black"
		} else if status == 1 {
			(*roominfolist)[i].Status = "新規"
			(*roominfolist)[i].Statuscolor = "green"
		} else {
			(*roominfolist)[i].Status = "エラー"
			(*roominfolist)[i].Statuscolor = "red"
		}
	}
	log.Printf("***** end of InsertRoomInf() ***********\n")
}

func InsertIntoOrUpdateUser(tnow time.Time, eventid string, roominf RoomInfo) (status int) {

	status = 0

	isnew := false

	userno, _ := strconv.Atoi(roominf.ID)
	log.Printf("  *** InsertIntoOrUpdateUser() *** userno=%d\n", userno)

	nrow := 0
	err := Db.QueryRow("select count(*) from user where userno =" + roominf.ID).Scan(&nrow)

	if err != nil {
		log.Printf("select count(*) from user ... err=[%s]\n", err.Error())
		status = -1
		return
	}

	name := ""
	genre := ""
	rank := ""
	nrank := ""
	prank := ""
	level := 0
	followers := 0
	fans := -1
	fans_lst := -1

	if nrow == 0 {

		isnew = true

		log.Printf("insert into userhistory(*new*) userno=%d rank=<%s> nrank=<%s> prank=<%s> level=%d, followers=%d, fans=%d, fans_lst=%d\n",
			userno, roominf.Rank, roominf.Nrank, roominf.Prank, roominf.Level, roominf.Followers, fans, fans_lst)

		sql := "INSERT INTO user(userno, userid, user_name, longname, shortname, genre, `rank`, nrank, prank, level, followers, fans, fans_lst, ts, currentevent)"
		sql += " VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"

		//	log.Printf("sql=%s\n", sql)
		row, err := Db.Prepare(sql)
		if err != nil {
			log.Printf("InsertIntoOrUpdateUser() error() (INSERT/Prepare) err=%s\n", err.Error())
			status = -1
			return
		}

		lenid := len(roominf.ID)
		_, err = row.Exec(
			userno,
			roominf.Account,
			roominf.Name,
			roominf.ID,
			roominf.ID[lenid-2:lenid],
			roominf.Genre,
			roominf.Rank,
			roominf.Nrank,
			roominf.Prank,
			roominf.Level,
			roominf.Followers,
			roominf.Fans,
			roominf.Fans_lst,
			tnow,
			eventid,
		)

		if err != nil {
			log.Printf("error(InsertIntoOrUpdateUser() INSERT/Exec) err=%s\n", err.Error())
			//	status = -2
			_, err = row.Exec(
				userno,
				roominf.Account,
				roominf.Account,
				roominf.ID,
				roominf.ID[lenid-2:lenid],
				roominf.Genre,
				roominf.Rank,
				roominf.Nrank,
				roominf.Prank,
				roominf.Level,
				roominf.Followers,
				roominf.Fans,
				roominf.Fans_lst,
				tnow,
				eventid,
			)
			if err != nil {
				log.Printf("error(InsertIntoOrUpdateUser() INSERT/Exec) err=%s\n", err.Error())
				status = -2
			}
		}
	} else {

		sql := "select user_name, genre, `rank`, nrank, prank, level, followers, fans, fans_lst from user where userno = ?"
		err = Db.QueryRow(sql, userno).Scan(&name, &genre, &rank, &nrank, &prank, &level, &followers, &fans, &fans_lst)
		if err != nil {
			log.Printf("err=[%s]\n", err.Error())
			status = -1
		}
		//	log.Printf("current userno=%d name=%s, nrank=%s, prank=%s level=%d, followers=%d\n", userno, name, nrank, prank, level, followers)

		if roominf.Genre != genre ||
			roominf.Rank != rank ||
			//	roominf.Nrank != nrank ||
			//	roominf.Prank != prank ||
			roominf.Level != level ||
			roominf.Followers != followers ||
			roominf.Fans != fans {

			isnew = true

			log.Printf("insert into userhistory(*changed*) userno=%d level=%d, followers=%d, fans=%d\n",
				userno, roominf.Level, roominf.Followers, roominf.Fans)
			sql := "update user set userid=?,"
			sql += "user_name=?,"
			sql += "genre=?,"
			sql += "`rank`=?,"
			sql += "nrank=?,"
			sql += "prank=?,"
			sql += "level=?,"
			sql += "followers=?,"
			sql += "fans=?,"
			sql += "fans_lst=?,"
			sql += "ts=?,"
			sql += "currentevent=? "
			sql += "where userno=?"
			row, err := Db.Prepare(sql)

			if err != nil {
				log.Printf("InsertIntoOrUpdateUser() error(Update/Prepare) err=%s\n", err.Error())
				status = -1
				return
			}
			_, err = row.Exec(
				roominf.Account,
				roominf.Name,
				roominf.Genre,
				roominf.Rank,
				roominf.Nrank,
				roominf.Prank,
				roominf.Level,
				roominf.Followers,
				roominf.Fans,
				roominf.Fans_lst,
				tnow,
				eventid,
				roominf.ID,
			)

			if err != nil {
				log.Printf("error(InsertIntoOrUpdateUser() Update/Exec) err=%s\n", err.Error())
				status = -2
			}
		}
		/* else {
			//	log.Printf("not insert into userhistory(*same*) userno=%d level=%d, followers=%d\n", userno, roominf.Level, roominf.Followers)
		}
		*/

	}

	if isnew {
		sql := "INSERT INTO userhistory(userno, user_name, genre, `rank`, nrank, prank, level, followers, fans, fans_lst, ts)"
		sql += " VALUES(?,?,?,?,?,?,?,?,?,?,?)"
		//	log.Printf("sql=%s\n", sql)
		row, err := Db.Prepare(sql)
		if err != nil {
			log.Printf("error(INSERT into userhistory/Prepare) err=%s\n", err.Error())
			status = -1
			return
		}

		_, err = row.Exec(
			userno,
			roominf.Name,
			roominf.Genre,
			roominf.Rank,
			roominf.Nrank,
			roominf.Prank,
			roominf.Level,
			roominf.Followers,
			roominf.Fans,
			roominf.Fans_lst,
			tnow,
		)

		if err != nil {
			log.Printf("error(Insert Into into userhistory INSERT/Exec) err=%s\n", err.Error())
			//	status = -2
			_, err = row.Exec(
				userno,
				roominf.Account,
				roominf.Genre,
				roominf.Rank,
				roominf.Nrank,
				roominf.Prank,
				roominf.Level,
				roominf.Followers,
				roominf.Fans,
				roominf.Fans_lst,
				tnow,
			)
			if err != nil {
				log.Printf("error(Insert Into into userhistory INSERT/Exec) err=%s\n", err.Error())
				status = -2
			}
		}

	}

	return

}
func InsertIntoEventUser(i int, eventid string, roominf RoomInfo) (status int) {

	status = 0

	userno, _ := strconv.Atoi(roominf.ID)

	nrow := 0
	/*
		sql := "select count(*) from eventuser where "
		sql += "userno =" + roominf.ID + " and "
		//	sql += "eventno = " + fmt.Sprintf("%d", eventno)
		sql += "eventid = " + eventid
		//	log.Printf("sql=%s\n", sql)
		err := Db.QueryRow(sql).Scan(&nrow)
	*/
	sql := "select count(*) from eventuser where userno =? and eventid = ?"
	err := Db.QueryRow(sql, roominf.ID, eventid).Scan(&nrow)

	if err != nil {
		log.Printf("select count(*) from user ... err=[%s]\n", err.Error())
		status = -1
		return
	}

	Colorlist := Colorlist2
	if Event_inf.Cmap == 1 {
		Colorlist = Colorlist1
	}

	if nrow == 0 {
		sql := "INSERT INTO eventuser(eventid, userno, istarget, graph, color, iscntrbpoints, point) VALUES(?,?,?,?,?,?,?)"
		row, err := Db.Prepare(sql)
		if err != nil {
			log.Printf("error(INSERT/Prepare) err=%s\n", err.Error())
			status = -1
			return
		}
		if i < 10 {
			_, err = row.Exec(
				eventid,
				userno,
				"Y",
				"Y",
				Colorlist[i%len(Colorlist)].Name,
				"N",
				roominf.Point,
			)
		} else {
			_, err = row.Exec(
				eventid,
				userno,
				"N",
				"N",
				Colorlist[i%len(Colorlist)].Name,
				"N",
				roominf.Point,
			)
		}

		if err != nil {
			log.Printf("error(InsertIntoOrUpdateUser() INSERT/Exec) err=%s\n", err.Error())
			status = -2
		}
		status = 1
	}
	return

}

func GetEventInfAndRoomList(
	eventid string,
	breg int,
	ereg int,
	eventinfo *Event_Inf,
	roominfolist *RoomInfoList,
) (
	status int,
) {

	//	画面からのデータ取得部分は次を参考にしました。
	//		はじめてのGo言語：Golangでスクレイピングをしてみた
	//		https://qiita.com/ryo_naka/items/a08d70f003fac7fb0808

	//	_url := "https://www.showroom-live.com/event/" + EventID
	//	_url = "file:///C:/Users/kohei47/Go/src/EventRoomList03/20210128-1143.html"
	//	_url = "file:20210128-1143.html"

	var doc *goquery.Document
	var err error

	inputmode := "url"
	eventidorfilename := eventid
	maxroom := ereg

	status = 0

	if inputmode == "file" {

		//	ファイルからドキュメントを作成します
		f, e := os.Open(eventidorfilename)
		if e != nil {
			//	log.Fatal(e)
			log.Printf("err=[%s]\n", e.Error())
			status = -1
			return
		}
		defer f.Close()
		doc, err = goquery.NewDocumentFromReader(f)
		if err != nil {
			//	log.Fatal(err)
			log.Printf("err=[%s]\n", err.Error())
			status = -1
			return
		}

		content, _ := doc.Find("head > meta:nth-child(6)").Attr("content")
		content_div := strings.Split(content, "/")
		(*eventinfo).Event_ID = content_div[len(content_div)-1]

	} else {
		//	URLからドキュメントを作成します
		_url := "https://www.showroom-live.com/event/" + eventidorfilename
		/*
			doc, err = goquery.NewDocument(_url)
		*/
		resp, error := http.Get(_url)
		if error != nil {
			log.Printf("GetEventInfAndRoomList() http.Get() err=%s\n", error.Error())
			status = 1
			return
		}
		defer resp.Body.Close()

		doc, error = goquery.NewDocumentFromReader(resp.Body)
		if error != nil {
			log.Printf("GetEventInfAndRoomList() goquery.NewDocumentFromReader() err=<%s>.\n", error.Error())
			status = 1
			return
		}

		(*eventinfo).Event_ID = eventidorfilename
	}
	//	log.Printf(" eventid=%s\n", (*eventinfo).Event_ID)

	selector := doc.Find(".detail")
	(*eventinfo).Event_name = selector.Find(".tx-title").Text()
	if (*eventinfo).Event_name == "" {
		log.Printf("Event not found. Event_ID=%s\n", (*eventinfo).Event_ID)
		status = -1
		return
	}
	(*eventinfo).Period = selector.Find(".info").Text()
	period := strings.Split((*eventinfo).Period, " - ")
	if inputmode == "url" {
		(*eventinfo).Start_time, _ = time.Parse("Jan 2, 2006 3:04 PM MST", period[0]+" JST")
		(*eventinfo).End_time, _ = time.Parse("Jan 2, 2006 3:04 PM MST", period[1]+" JST")
	} else {
		(*eventinfo).Start_time, _ = time.Parse("2006/01/02 15:04 MST", period[0]+" JST")
		(*eventinfo).End_time, _ = time.Parse("2006/01/02 15:04 MST", period[1]+" JST")
	}

	(*eventinfo).EventStatus = "BeingHeld"
	if (*eventinfo).Start_time.After(time.Now()) {
		(*eventinfo).EventStatus = "NotHeldYet"
	} else if (*eventinfo).End_time.Before(time.Now()) {
		(*eventinfo).EventStatus = "Over"
	}

	//	イベントに参加しているルームの数を求めます。
	//	参加ルーム数と表示されているルームの数は違うので、ここで取得したルームの数を以下の処理で使うわけではありません。
	SNoEntry := doc.Find("p.ta-r").Text()
	fmt.Sscanf(SNoEntry, "%d", &((*eventinfo).NoEntry))
	log.Printf("[%s]\n[%s] [%s] (*event).EventStatus=%s NoEntry=%d\n",
		(*eventinfo).Event_name,
		(*eventinfo).Start_time.Format("2006/01/02 15:04 MST"),
		(*eventinfo).End_time.Format("2006/01/02 15:04 MST"),
		(*eventinfo).EventStatus, (*eventinfo).NoEntry)
	log.Printf("breg=%d ereg=%d\n", breg, ereg)

	//	eventno, _, _ := SelectEventNoAndName(eventidorfilename)
	//	log.Printf(" eventno=%d\n", eventno)
	//	(*eventinfo).Event_no = eventno

	//	抽出したルームすべてに対して処理を繰り返す(が、イベント開始後の場合の処理はルーム数がbreg、eregの範囲に限定）
	//	イベント開始前のときはすべて取得し、ソートしたあてで範囲を限定する）
	doc.Find(".listcardinfo").EachWithBreak(func(i int, s *goquery.Selection) bool {
		//	log.Printf("i=%d\n", i)
		if (*eventinfo).Start_time.Before(time.Now()) {
			if i < breg-1 {
				return true
			}
			if i == maxroom {
				return false
			}
		}

		var roominfo RoomInfo

		roominfo.Name = s.Find(".listcardinfo-main-text").Text()

		spoint1 := strings.Split(s.Find(".listcardinfo-sub-single-right-text").Text(), ": ")

		var point int64
		if spoint1[0] != "" {
			spoint2 := strings.Split(spoint1[1], "pt")
			fmt.Sscanf(spoint2[0], "%d", &point)

		} else {
			point = -1
		}
		roominfo.Point = int(point)

		ReplaceString := ""

		selection_c := s.Find(".listcardinfo-menu")

		account, _ := selection_c.Find(".room-url").Attr("href")
		if inputmode == "file" {
			ReplaceString = "https://www.showroom-live.com/"
		} else {
			ReplaceString = "/r/"
		}
		roominfo.Account = strings.Replace(account, ReplaceString, "", -1)
		roominfo.Account = strings.Replace(roominfo.Account, "/", "", -1)

		roominfo.ID, _ = selection_c.Find(".js-follow-btn").Attr("data-room-id")
		roominfo.Userno, _ = strconv.Atoi(roominfo.ID)

		*roominfolist = append(*roominfolist, roominfo)

		//	log.Printf("%11s %-20s %-10s %s\n",
		//		humanize.Comma(int64(roominfo.Point)), roominfo.Account, roominfo.ID, roominfo.Name)
		return true

	})

	(*eventinfo).NoRoom = len(*roominfolist)

	log.Printf(" GetEventInfAndRoomList() len(*roominfolist)=%d\n", len(*roominfolist))

	return
}

func GetEventInf(
	eventid string,
	eventinfo *Event_Inf,
) (
	status int,
) {

	//	画面からのデータ取得部分は次を参考にしました。
	//		はじめてのGo言語：Golangでスクレイピングをしてみた
	//		https://qiita.com/ryo_naka/items/a08d70f003fac7fb0808

	//	_url := "https://www.showroom-live.com/event/" + EventID
	//	_url = "file:///C:/Users/kohei47/Go/src/EventRoomList03/20210128-1143.html"
	//	_url = "file:20210128-1143.html"

	var doc *goquery.Document
	var err error

	inputmode := "url"
	eventidorfilename := eventid

	status = 0

	/*
		_, _, status := SelectEventNoAndName(eventidorfilename)
		log.Printf(" status=%d\n", status)
		if status != 0 {
			return
		}
		(*eventinfo).Event_no = eventno
	*/

	if inputmode == "file" {

		//	ファイルからドキュメントを作成します
		f, e := os.Open(eventidorfilename)
		if e != nil {
			//	log.Fatal(e)
			log.Printf("err=[%s]\n", e.Error())
			status = -1
			return
		}
		defer f.Close()
		doc, err = goquery.NewDocumentFromReader(f)
		if err != nil {
			status = -4
			return
		}

		content, _ := doc.Find("head > meta:nth-child(6)").Attr("content")
		content_div := strings.Split(content, "/")
		(*eventinfo).Event_ID = content_div[len(content_div)-1]

	} else {
		//	URLからドキュメントを作成します
		_url := "https://www.showroom-live.com/event/" + eventidorfilename
		/*
		doc, err = goquery.NewDocument(_url)
		*/
		resp, error := http.Get(_url)
		if error != nil {
			log.Printf("GetEventInfAndRoomList() http.Get() err=%s\n", error.Error())
			status = 1
			return
		}
		defer resp.Body.Close()

		doc, error = goquery.NewDocumentFromReader(resp.Body)
		if error != nil {
			log.Printf("GetEventInfAndRoomList() goquery.NewDocumentFromReader() err=<%s>.\n", error.Error())
			status = 1
			return
		}

		(*eventinfo).Event_ID = eventidorfilename
	}
	value, _ := doc.Find("#eventDetail").Attr("data-event-id")
	(*eventinfo).I_Event_ID, _ = strconv.Atoi(value)

	log.Printf(" eventid=%s (%d)\n", (*eventinfo).Event_ID, (*eventinfo).I_Event_ID)

	selector := doc.Find(".detail")
	(*eventinfo).Event_name = selector.Find(".tx-title").Text()
	if (*eventinfo).Event_name == "" {
		log.Printf("Event not found. Event_ID=%s\n", (*eventinfo).Event_ID)
		status = -2
		return
	}
	(*eventinfo).Period = selector.Find(".info").Text()
	period := strings.Split((*eventinfo).Period, " - ")
	if inputmode == "url" {
		(*eventinfo).Start_time, _ = time.Parse("Jan 2, 2006 3:04 PM MST", period[0]+" JST")
		(*eventinfo).End_time, _ = time.Parse("Jan 2, 2006 3:04 PM MST", period[1]+" JST")
	} else {
		(*eventinfo).Start_time, _ = time.Parse("2006/01/02 15:04 MST", period[0]+" JST")
		(*eventinfo).End_time, _ = time.Parse("2006/01/02 15:04 MST", period[1]+" JST")
	}

	(*eventinfo).EventStatus = "BeingHeld"
	if (*eventinfo).Start_time.After(time.Now()) {
		(*eventinfo).EventStatus = "NotHeldYet"
	} else if (*eventinfo).End_time.Before(time.Now()) {
		(*eventinfo).EventStatus = "Over"
	}

	//	イベントに参加しているルームの数を求めます。
	//	参加ルーム数と表示されているルームの数は違うので注意。ここで取得しているのは参加ルーム数。
	SNoEntry := doc.Find("p.ta-r").Text()
	fmt.Sscanf(SNoEntry, "%d", &((*eventinfo).NoEntry))
	log.Printf("[%s]\n[%s] [%s] (*event).EventStatus=%s NoEntry=%d\n",
		(*eventinfo).Event_name,
		(*eventinfo).Start_time.Format("2006/01/02 15:04 MST"),
		(*eventinfo).End_time.Format("2006/01/02 15:04 MST"),
		(*eventinfo).EventStatus, (*eventinfo).NoEntry)

	return
}

func SelectEventNoAndName(eventid string) (
	eventname string,
	period string,
	status int,
) {

	status = 0

	err := Db.QueryRow("select event_name, period from event where eventid ='"+eventid+"'").Scan(&eventname, &period)

	if err == nil {
		return
	} else {
		log.Printf("err=[%s]\n", err.Error())
		if err.Error() != "sql: no rows in result set" {
			status = -2
			return
		}
	}

	status = -1
	return
}



func SelectRoomLevel(userno int, levelonly int) (roomlevelinf RoomLevelInf, status int) {

	var stmt *sql.Stmt
	var rows *sql.Rows

	status = 0

	sqlstmt := "select user_name, genre, `rank`, nrank, prank, level, followers, fans, fans_lst, ts from userhistory where userno = ? order by ts desc"
	stmt, Err = Db.Prepare(sqlstmt)
	if Err != nil {
		log.Printf("SelectRoomLevel() (3) err=%s\n", Err.Error())
		status = -3
		return
	}
	defer stmt.Close()

	rows, Err = stmt.Query(userno)
	if Err != nil {
		log.Printf("SelectRoomLevel() (6) err=%s\n", Err.Error())
		status = -6
		return
	}
	defer rows.Close()

	/*
	   type RoomLevel struct {
	   	User_name  string
	   	Genre      string
	   	Rank       string
	   	Nrank       string
	   	Level      int
	   	Followeres int
	   	Sts        string
	   }

	   type RoomLevelInf struct {
	   	Userno        int
	   	User_name      string
	   	RoomLevelList []RoomLevel
	   }
	*/

	var roomlevel RoomLevel

	roomlevelinf.Userno = userno

	lastlevel := 0

	for rows.Next() {
		Err = rows.Scan(&roomlevel.User_name, &roomlevel.Genre, &roomlevel.Rank,
			&roomlevel.Nrank,
			&roomlevel.Prank,
			&roomlevel.Level,
			&roomlevel.Followers,
			&roomlevel.Fans,
			&roomlevel.Fans_lst,
			&roomlevel.ts)
		if Err != nil {
			log.Printf("GetCurrentScore() (7) err=%s\n", Err.Error())
			status = -7
			return
		}

		if lastlevel == 0 {
			roomlevelinf.User_name = roomlevel.User_name
		}

		if levelonly == 1 && roomlevel.Level == lastlevel {
			continue
		}
		lastlevel = roomlevel.Level

		//	roomlevel.Sfollowers = humanize.Comma(int64(roomlevel.Followers))
		roomlevel.Sts = roomlevel.ts.Format("2006/01/02 15:04")

		roomlevelinf.RoomLevelList = append(roomlevelinf.RoomLevelList, roomlevel)

	}

	return
}




func SelectLastEventList() (eventlist []Event, status int) {

	var stmt *sql.Stmt
	var rows *sql.Rows

	sql := "select eventid, event_name, period, starttime, endtime, nobasis, longname from event join user "
	sql += " where nobasis = userno and endtime IS not null order by endtime desc "
	stmt, Err = Db.Prepare(sql)
	if Err != nil {
		log.Printf("err=[%s]\n", Err.Error())
		status = -1
		return
	}
	defer stmt.Close()

	rows, Err = stmt.Query()
	if Err != nil {
		log.Printf("err=[%s]\n", Err.Error())
		status = -1
		return
	}
	defer rows.Close()

	var event Event
	i := 0
	for rows.Next() {
		Err = rows.Scan(&event.EventID, &event.EventName, &event.Period, &event.Starttime, &event.Endtime, &event.Pntbasis, &event.Pbname)
		if Err != nil {
			log.Printf("err=[%s]\n", Err.Error())
			status = -1
			return
		}
		eventlist = append(eventlist, event)
		i++
		if i == 10 {
			break
		}
	}
	if Err = rows.Err(); Err != nil {
		log.Printf("err=[%s]\n", Err.Error())
		status = -1
		return
	}

	tnow := time.Now()
	for i = 0; i < len(eventlist); i++ {
		eventlist[i].S_start = eventlist[i].Starttime.Format("2006-01-02 15:04")
		eventlist[i].S_end = eventlist[i].Endtime.Format("2006-01-02 15:04")

		if eventlist[i].Starttime.After(tnow) {
			eventlist[i].Status = "これから開催"
		} else if eventlist[i].Endtime.Before(tnow) {
			eventlist[i].Status = "終了"
		} else {
			eventlist[i].Status = "開催中"
		}

	}

	return

}

func OpenDb() (status int) {

	status = 0

	//	https://leben.mobi/go/mysql-connect/practice/
	//	OS := runtime.GOOS

	//	https://ssabcire.hatenablog.com/entry/2019/02/13/000722
	//	https://konboi.hatenablog.com/entry/2016/04/12/100903
	/*
		switch OS {
		case "windows":
			Db, Err = sql.Open("mysql", wuser+":"+wpw+"@/"+wdb+"?parseTime=true&loc=Asia%2FTokyo")
		case "linux":
			Db, Err = sql.Open("mysql", luser+":"+lpw+"@/"+ldb+"?parseTime=true&loc=Asia%2FTokyo")
		case "freebsd":
			//	https://leben.mobi/go/mysql-connect/practice/
			Db, Err = sql.Open("mysql", buser+":"+bpw+"@tcp("+bhost+":3306)/"+bdb+"?parseTime=true&loc=Asia%2FTokyo")
		default:
			log.Printf("%s is not supported.\n", OS)
			status = -2
		}
	*/

	if (*Dbconfig).Dbhost == "" {
		Db, Err = sql.Open("mysql", (*Dbconfig).Dbuser+":"+(*Dbconfig).Dbpw+"@/"+(*Dbconfig).Dbname+"?parseTime=true&loc=Asia%2FTokyo")
	} else {
		Db, Err = sql.Open("mysql", (*Dbconfig).Dbuser+":"+(*Dbconfig).Dbpw+"@tcp("+(*Dbconfig).Dbhost+":3306)/"+(*Dbconfig).Dbname+"?parseTime=true&loc=Asia%2FTokyo")
	}

	if Err != nil {
		status = -1
	}
	return
}

func SelectEventInfAndRoomList() (IDlist []int, status int) {

	status = 0

	/*
		//	sql := "select eventno, event_name, period, starttime, endtime from event where eventid ='"+Event_inf.Event_ID+"'"
		sql := "select eventno, event_name, period, starttime, endtime from event where eventid = ?"
		err := Db.QueryRow(sql, Event_inf.Event_ID).Scan(&Event_inf.Event_no, &Event_inf.Event_name, &Event_inf.Period, &Event_inf.Start_time, &Event_inf.End_time)

		if err != nil {
			log.Printf("select eventno, starttime, endtime from event where eventid ='%s'\n", Event_inf.Event_ID)
			log.Printf("err=[%s]\n", err.Error())
			//	if err.Error() != "sql: no rows in result set" {
			status = -1
			return
			//	}
		}
	*/

	Event_inf, _ = SelectEventInf(Event_inf.Event_ID)

	//	log.Printf("eventno=%d\n", Event_inf.Event_no)

	start_date := Event_inf.Start_time.Truncate(time.Hour).Add(-time.Duration(Event_inf.Start_time.Hour()) * time.Hour)
	end_date := Event_inf.End_time.Truncate(time.Hour).Add(-time.Duration(Event_inf.End_time.Hour())*time.Hour).AddDate(0, 0, 1)

	//	log.Printf("start_t=%v\nstart_d=%v\nend_t=%v\nend_t=%v\n", Event_inf.Start_time, start_date, Event_inf.End_time, end_date)

	Event_inf.Start_date = float64(start_date.Unix()) / 60.0 / 60.0 / 24.0
	Event_inf.Dperiod = float64(end_date.Unix())/60.0/60.0/24.0 - Event_inf.Start_date

	//	log.Printf("Start_data=%f Dperiod=%f\n", Event_inf.Start_date, Event_inf.Dperiod)

	//	err = Db.QueryRow("select max(point) from points where event_id = '" + fmt.Sprintf("%d", Event_inf.Event_no) + "'").Scan(&Event_inf.MaxPoint)
	//	sql := "select max(point) from eventuser where eventno = ? and graph = 'Y'"
	sql := "select max(point) from eventuser where eventid = ? and graph = 'Y'"
	err := Db.QueryRow(sql, Event_inf.Event_ID).Scan(&Event_inf.MaxPoint)

	if err != nil {
		log.Printf("select max(point) from eventuser where eventid = '%s'\n", Event_inf.Event_ID)
		log.Printf("err=[%s]\n", err.Error())
		status = -2
		return
	}

	//	log.Printf("MaxPoint=%d\n", Event_inf.MaxPoint)

	//	-------------------------------------------------------------------
	//	sql := "select user_id from points where event_id = ? and idx = ( select max(idx) from points where event_id = ? ) order by point desc"
	sql = " select userno from eventuser "
	sql += " where graph = 'Y' "
	//	sql += " and eventno = ? "
	sql += " and eventid = ? "
	sql += " order by point desc"
	stmt, err := Db.Prepare(sql)
	if err != nil {
		//	log.Fatal(err)
		log.Printf("err=[%s]\n", err.Error())
		status = -1
		return
	}
	defer stmt.Close()

	rows, err := stmt.Query(Event_inf.Event_ID)
	if err != nil {
		//	log.Fatal(err)
		log.Printf("err=[%s]\n", err.Error())
		status = -1
		return
	}
	defer rows.Close()

	id := 0
	i := 0
	for rows.Next() {
		err := rows.Scan(&id)
		if err != nil {
			//	log.Fatal(err)
			log.Printf("err=[%s]\n", err.Error())
			status = -1
			return
		}
		IDlist = append(IDlist, id)
		i++
		if i == Event_inf.Maxdsp {
			break
		}
	}
	if err = rows.Err(); err != nil {
		//	log.Fatal(err)
		log.Printf("err=[%s]\n", err.Error())
		status = -1
		return
	}

	return
}

func SelectEventInf(eventid string) (eventinf Event_Inf, status int) {

	status = 0

	sql := "select eventid,event_name, period, starttime, endtime, noentry, intervalmin, modmin, modsec, "
	sql += " Fromorder, Toorder, Resethh, Resetmm, Nobasis, Maxdsp, cmap, target, maxpoint "
	sql += " from event where eventid = ?"
	err := Db.QueryRow(sql, eventid).Scan(
		&eventinf.Event_ID,
		&eventinf.Event_name,
		&eventinf.Period,
		&eventinf.Start_time,
		&eventinf.End_time,
		&eventinf.NoEntry,
		&eventinf.Intervalmin,
		&eventinf.Modmin,
		&eventinf.Modsec,
		&eventinf.Fromorder,
		&eventinf.Toorder,
		&eventinf.Resethh,
		&eventinf.Resetmm,
		&eventinf.Nobasis,
		&eventinf.Maxdsp,
		&eventinf.Cmap,
		&eventinf.Target,
		&eventinf.Maxpoint,
	)

	if err != nil {
		log.Printf("%s\n", sql)
		log.Printf("err=[%s]\n", err.Error())
		//	if err.Error() != "sql: no rows in result set" {
		status = -1
		return
		//	}
	}

	//	log.Printf("eventno=%d\n", Event_inf.Event_no)

	start_date := eventinf.Start_time.Truncate(time.Hour).Add(-time.Duration(eventinf.Start_time.Hour()) * time.Hour)
	end_date := eventinf.End_time.Truncate(time.Hour).Add(-time.Duration(eventinf.End_time.Hour())*time.Hour).AddDate(0, 0, 1)

	//	log.Printf("start_t=%v\nstart_d=%v\nend_t=%v\nend_t=%v\n", Event_inf.Start_time, start_date, Event_inf.End_time, end_date)

	eventinf.Start_date = float64(start_date.Unix()) / 60.0 / 60.0 / 24.0
	eventinf.Dperiod = float64(end_date.Unix())/60.0/60.0/24.0 - Event_inf.Start_date

	//	log.Printf("Start_data=%f Dperiod=%f\n", eventinf.Start_date, eventinf.Dperiod)

	return
}















/*
        SelectAcqTimeList()
		指定したイベント、ユーザーの貢献ランキングを取得した時刻の一覧を取得する。

        引数
		eventid			string			イベントID
		userno			int				ユーザーID

        戻り値
        acqtimelist		[] time.Time	取得時刻一覧

*/
func SelectAcqTimeList(eventid string, userno int) (acqtimelist []time.Time, status int) {

	var stmt *sql.Stmt
	var rows *sql.Rows

	status = 0

	//	貢献ポイントランキングを取得した時刻の一覧を取得する。
	sql := "select sampletm2 from timetable where eventid = ? and userid = ? and status = 1 order by sampletm2"
	stmt, Err = Db.Prepare(sql)

	if Err != nil {
		log.Printf("SelectAcqTimeList() (5) err=%s\n", Err.Error())
		status = -5
		return
	}
	defer stmt.Close()

	rows, Err = stmt.Query(eventid, userno)
	if Err != nil {
		log.Printf("SelectAcqTimeList() (6) err=%s\n", Err.Error())
		status = -6
		return
	}
	defer rows.Close()

	var ts time.Time

	for rows.Next() {
		Err = rows.Scan(&ts)
		if Err != nil {
			log.Printf("SelectAcqTimeList() (7) err=%s\n", Err.Error())
			status = -7
			return
		}
		//	log.Printf("%+v\n", cntrbinf)
		acqtimelist = append(acqtimelist, ts)
	}
	if Err = rows.Err(); Err != nil {
		log.Printf("SelectAcqTimeList() (8) err=%s\n", Err.Error())
		status = -8
		return
	}

	return

}

