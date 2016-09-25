package main

import (
	"gopkg.in/mgo.v2/bson"
	"time"
	"net/http"
	"github.com/mholt/binding"
	"github.com/julienschmidt/httprouter"
)

type Notice struct {
	ID bson.ObjectId `bson: "_id" json: "id"`
	TITLE string `bson: "title" json: "title"`
	DESCRIPTION string `bson: "description" json: "description"`
	CREATED time.Time `bson: "created" json: "created"`
}

const DATABASE = "blog"
const TABLE_NOTICE = "notices"

func timeStamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func (n *Notice) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{&n.TITLE: "title", &n.DESCRIPTION: "description", &n.CREATED: timeStamp()}
}

func getNotices(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	// 몽고DB 세션 생성
	session := mongoSession.Copy()

	// 몽고 DB 닫는 세션 defer로 등록
	defer  session.Close()

	var notices []Notice
	// 모든 게시물 조회
	err := session.DB(DATABASE).C(TABLE_NOTICE).Find(nil).All(&notices)
	if err != nil {
		// 오류 발생 시
		renderer.JSON(w, http.StatusInternalServerError, err)
		return
	}

	renderer.JSON(w, http.StatusOK, notices)
}

func createNotice(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	// binding 패키지로 notice 생성 요청 정보를 Notice 타입으로 변환
	n := new(Notice)
	errs := binding.Bind(req, n)
	if errs.Handle(w) {
		return
	}

	// 몽고 DB 세션 생성
	session := mongoSession.Copy()

	defer session.Close()

	// 몽고DB 아이디 생성
	n.ID = bson.NewObjectId()
	n.TITLE = "title"
	n.DESCRIPTION = "description"
	n.CREATED = time.Now()

	// notice 정보 저장을 위한 몽고DB 컬렉션 객체 생성
	c := session.DB(DATABASE).C(TABLE_NOTICE)

	// notices 컬렉션에 notice 정보 저장
	if err := c.Insert(n); err != nil {
		// 오류 발생 시 500 에러 반환
		renderer.JSON(w, http.StatusInternalServerError, err)
		return
	}

	// 처리 결과 반환
	renderer.JSON(w, http.StatusCreated, n)
}