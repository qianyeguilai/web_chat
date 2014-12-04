package main

import (
	"bufio"
	"code.google.com/p/go.net/websocket"
	"container/list"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"mysql/mysqlfunc"
	"net/http"
	"os"
	"path"
)

const (
	TEMPLATE_DIR = "../views"
)

var connid int
var conns *list.List

var templates = make(map[string]*template.Template)
var mysqlinfo = mysqlfunc.Db_info{"127.0.0.1", "root", "", "web_chat_go"}
var mysql_db *sql.DB

func init() {
	fileInfoArr, err := ioutil.ReadDir(TEMPLATE_DIR)
	check(err)
	var templateName, templatePath string
	for _, fileInfo := range fileInfoArr {
		templateName = fileInfo.Name()
		if ext := path.Ext(templateName); ext != ".html" {
			continue
		}
		templatePath = TEMPLATE_DIR + "/" + templateName
		log.Println("Loading template:", templatePath)
		t := template.Must(template.ParseFiles(templatePath))
		templates[templateName] = t
	}

	x, err1 := mysqlfunc.ConnectDb(mysqlinfo)
	if err1 != nil {
		fmt.Println("ConnectDb  error:", err1)
		return
	}
	check(err1)
	mysql_db = x
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func ChatroomServer(ws *websocket.Conn) {
	defer ws.Close()

	connid++
	id := connid

	fmt.Printf("connection id: %d\n", id)

	item := conns.PushBack(ws)
	defer conns.Remove(item)

	name := fmt.Sprintf("user%d", id)

	SendMessage(nil, fmt.Sprintf("welcome %s join\n", name))

	r := bufio.NewReader(ws)

	for {
		data, err := r.ReadBytes('\n')
		if err != nil {
			fmt.Printf("disconnected id: %d\n", id)
			SendMessage(item, fmt.Sprintf("%s offline\n", name))
			break
		}

		fmt.Printf("%s: %s", name, data)

		SendMessage(item, fmt.Sprintf("%s\t> %s", name, data))
	}
}

func SendMessage(self *list.Element, data string) {
	// for _, item := range conns {
	for item := conns.Front(); item != nil; item = item.Next() {
		ws, ok := item.Value.(*websocket.Conn)
		if !ok {
			panic("item not *websocket.Conn")
		}

		if item == self {
			continue
		}

		io.WriteString(ws, data)
	}
}

/**********************************************************/
/*              write  the   html   page                  */
/*                                                        */
/**********************************************************/
func renderHtml(w http.ResponseWriter, tmpl string, locals map[string]interface{}) {
	tmpl += ".html"
	err := templates[tmpl].Execute(w, locals)
	check(err)
}

/*********************************************************/
/*               judge  the file  path  exist            */
/*                                                       */
/*********************************************************/
func isExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	return os.IsExist(err)
}

/*********************************************************/
/*              user   load  func                        */
/*                                                       */
/*********************************************************/
func userLoad(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		err := templates["userload.html"].Execute(w, nil)
		check(err)
	}
	if r.Method == "POST" {
		r.ParseForm()
		ok := mysqlfunc.Checkuserok(mysql_db, r.Form["userName"], r.Form["password"])
		if !ok {
			err := templates["userload.html"].Execute(w, "用户名或密码错误")
			check(err)
		} else {
			renderHtml(w, "web_js", nil)
		}
	}
}

/*********************************************************/
/*            download the  html picture                 */
/*                                                       */
/*********************************************************/
func picture_down(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	imageId := r.Form["pic_name"]
	imgpath := "../img/" + imageId[0]

	if exist := isExists(imgpath); !exist {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "image")
	http.ServeFile(w, r, imgpath)
}
func js_down(w http.ResponseWriter, r *http.Request) {
	imgpath := "../img/js/jquery.min.js"
	if exist := isExists(imgpath); !exist {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "image")
	http.ServeFile(w, r, imgpath)
}

/**********************************************************/
/*             user  register  func                       */
/*                                                        */
/**********************************************************/
func register_user(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		err := templates["register.html"].Execute(w, nil)
		check(err)
	}

	if r.Method == "POST" {
		r.ParseForm()
		username := r.Form["username"]
		password1 := r.Form["password1"]
		password2 := r.Form["password2"]
		sex := r.Form["sex"]

		if username[0] == "" || password1[0] == "" || sex[0] == "" || password2[0] == "" {
			fmt.Println("---------------------------")
			err := templates["register.html"].Execute(w, "请填写完整信息")
			check(err)
			return
		}

		if password1[0] != password2[0] {
			err := templates["register.html"].Execute(w, "两次密码不一致，请重试")
			check(err)
		} else {
			exist, err := mysqlfunc.Judge_user_exist(mysql_db, username[0])
			if err != nil {
				check(err)
				return
			}

			if exist != true {
				err = templates["register.html"].Execute(w, "该用户名被占用")
				check(err)
			}
			err = mysqlfunc.Db_register_user(mysql_db, username[0], password1[0], sex[0])
			if err != nil {
				panic(err)
			}
			temp_str := fmt.Sprintf("注册成功:帐号%s,密码%s", username[0], password1[0])
			err = templates["register.html"].Execute(w, temp_str)
			check(err)
		}
	}
}

/*
func safeHandler(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if e, ok := recover().(error); ok {
				http.Error(w, e.Error(), http.StatusInternalServerError)

				// 或者输出自定义的 50x 错误页面
				// w.WriteHeader(http.StatusInternalServerError)
				// renderHtml(w, "error", e.Error())

				// logging
				log.Println("WARN: panic fired in %v.panic - %v", fn, e)
				log.Println(string(debug.Stack()))
			}
		}()
		fn(w, r)
	}
}
*/
func main() {

	connid = 0
	conns = list.New()

	/*	mux := http.NewServeMux()
		mux.HandleFunc("/", safeHandler(userLoad))
		mux.HandleFunc("/img", safeHandler(picture_down))
		mux.HandleFunc("js",safeHandler(js_down))
		mux.HandleFunc("/register", safeHandler(register_user))
		mux.Handle("/chatroom",websocket.Handler(ChatroomServer))

		err := http.ListenAndServe(":7777", mux)
		if err != nil {
			log.Fatal("ListenAndServe: ", err.Error())
		}

	*/
	http.Handle("/chatroom", websocket.Handler(ChatroomServer))
	http.HandleFunc("/", userLoad)
	http.HandleFunc("/img", picture_down)
	http.HandleFunc("/js", js_down)
	http.HandleFunc("/register", register_user)

	err := http.ListenAndServe(":7777", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}

}
