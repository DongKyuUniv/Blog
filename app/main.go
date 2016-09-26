package main

import (
	"net/http"

	"github.com/codegangsta/negroni"
	"github.com/julienschmidt/httprouter"
	"github.com/unrolled/render"

	"gopkg.in/mgo.v2"
)


var (
	renderer *render.Render
	mongoSession *mgo.Session
)

func init() {
	renderer = render.New()

	s, err := mgo.Dial("mongodb://localhost")
	if err != nil {
		panic(err)
	}

	mongoSession = s
}

func main() {
	// 렌더러 생성
	router := httprouter.New()

	router.GET("/", func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		// 템플릿 렌더링
		renderer.HTML(w, http.StatusOK, "index", map[string]string{"title": "BangUl Blog!"})
	})

	router.GET("/notices", getNotices)
	router.GET("/notice", getNotice)

	router.GET("/createNotice", func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		renderer.HTML(w, http.StatusOK, "createNotice", map[string]string{})
	})

	router.POST("/createNotice", createNotice)

	// negroni 미들웨어 생성
	n:= negroni.Classic()

	// negroni에 router를 핸들러로 등록
	n.UseHandler(router)

	n.Run(":3000")
}
