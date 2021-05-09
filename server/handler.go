package server

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/julienschmidt/httprouter"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io/ioutil"
	"net/http"
	"pastebin/common"
	"pastebin/server/driver"
	"strconv"
	"time"
)

const (
	ErrHttpCodeOk      = 200
	ErrRecordNotFound  = 4001
	ErrRecordNeedPass  = 4002
	ErrArgsMissing     = 4003
	ErrArgsFormatError = 4004
	ErrRecordExpired   = 4005
	ErrServerInner     = 5000
)

// ResponseJSON 输出JSON结果
func ResponseJSON(w http.ResponseWriter, httpCode int, response interface{}) {
	var resp = Response{
		Code:    httpCode,
		Message: "ok",
		Data:    response,
	}
	body, err := json.Marshal(resp)
	if err != nil {
		log.Errorf("response marshal err: %v", err)
		return
	}
	//公共的响应头设置
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.Header().Set("Content-Length", strconv.Itoa(len(string(body))))

	if httpCode == 0 {
		httpCode = http.StatusOK
	}
	w.WriteHeader(httpCode)
	_, _ = w.Write(body)
}

func ResponseJSONError(w http.ResponseWriter, httpCode, bizCode int, err error) {
	//公共的响应头设置
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Content-Type", "application/json;charset=utf-8")

	if httpCode == 0 {
		httpCode = http.StatusOK
	}

	var resp = &Response{
		Code:    bizCode,
		Message: err.Error(),
		Data:    err,
	}

	respBody, err := json.Marshal(resp)
	if err != nil {
		log.Errorf("marshal err: %+v", err)
	}

	w.Header().Set("Content-Length", strconv.Itoa(len(string(respBody))))

	w.WriteHeader(httpCode)
	_, err = w.Write(respBody)

	if err != nil {
		log.Errorf("Write err: %+v", err)
	}
}

// ResponseHTML 输出原始的结果
func ResponseHTML(w http.ResponseWriter, httpCode int, response []byte) {
	//公共的响应头设置
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Content-Type", "text/plain;charset=utf-8;")
	w.Header().Set("Content-Length", strconv.Itoa(len(string(response))))

	if httpCode == 0 {
		httpCode = http.StatusOK
	}

	w.WriteHeader(httpCode)
	_, _ = w.Write(response)
}

func ResponseHTMLError(w http.ResponseWriter, httpCode, bizCode int, err error) {
	//公共的响应头设置
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Content-Type", "text/plain;charset=utf-8;")
	w.Header().Set("Content-Length", strconv.Itoa(len(err.Error())))

	if httpCode == 0 {
		httpCode = http.StatusInternalServerError
	}

	w.WriteHeader(httpCode)
	_, _ = w.Write([]byte(err.Error()))
}

func GetRecord(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var skraw = p.ByName("sk")
	var sk = r.URL.Query().Get("sk")
	var password = r.URL.Query().Get("password")
	var format = r.URL.Query().Get("format")

	if skraw != "" {
		sk = skraw
		format = "raw"
	}

	if sk == "" {
		log.Errorf("sk empty")
		ResponseJSONError(w, ErrHttpCodeOk, ErrArgsMissing, errors.New("short_key empty"))
		return
	}

	var record driver.CodeSegmentRecord
	err := driver.GetCollection().FindOne(context.Background(), bson.M{"short_key": sk}).Decode(&record)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			err = errors.New("record not found")
			ResponseJSONError(w, ErrHttpCodeOk, ErrRecordNotFound, err)
			return
		}
		log.Errorf("find record err: %+v", err)
		ResponseJSONError(w, ErrHttpCodeOk, ErrServerInner, err)
		return
	}

	// 检查是否需要密码
	if record.Password != "" && record.Password != password {
		ResponseJSONError(w, ErrHttpCodeOk, ErrRecordNeedPass, errors.New("you need a password to view this record"))
		return
	}
	// 不允许 password 字段暴露
	record.Password = ""

	// 检查过期时间
	var uptime = record.CreatedAt
	var nowTime = time.Now().Unix()
	if record.UpdatedAt > record.CreatedAt {
		uptime = record.UpdatedAt
	}
	var sub = nowTime - uptime

	var isExp bool
	var expSecond int64
	switch record.Lifecycle {
	case driver.LifeCycleOneDay:
		isExp = sub > 86400
		expSecond = 86400
	case driver.LifeCycleOneWeek:
		isExp = sub > 86400*7
		expSecond = 86400 * 7
	case driver.LifeCycleOneMonth:
		isExp = sub > 86400*30
		expSecond = 86400 * 30
	case driver.LifeCycleOneYear:
		isExp = sub > 86400*365
		expSecond = 86400 * 365
	default:
		isExp = false
	}
	if expSecond != 0 {
		record.ExpiredAt = uptime + expSecond
	}

	if isExp {
		ResponseJSONError(w, ErrHttpCodeOk, ErrRecordExpired, errors.New("record expired"))
		return
	}

	if format == "raw" {
		ResponseHTML(w, http.StatusOK, []byte(record.Content))
		return
	}

	ResponseJSON(w, http.StatusOK, record)
}

func GetPublicRecordList(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var ctx = context.Background()
	var records []*driver.CodeSegmentRecord

	var limit int64 = 100
	var sort = bson.D{
		{"updated_at", -1},
	}
	var opts = options.FindOptions{
		Limit: &limit,
		Sort:  sort,
	}
	cursor, err := driver.GetCollection().Find(ctx, bson.M{"password": ""}, &opts)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			err = errors.New("record not found")
			ResponseJSONError(w, ErrHttpCodeOk, ErrRecordNotFound, err)
			return
		}
		log.Errorf("find record err: %+v", err)
		ResponseJSONError(w, ErrHttpCodeOk, ErrRecordNotFound, err)
		return
	}

	if cursor.Err() != nil {
		log.Errorf("find record err: %+v", err)
		ResponseJSONError(w, ErrHttpCodeOk, ErrServerInner, err)
		return
	}

	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &records); err != nil {
		log.Errorf("find record err: %+v", err)
		ResponseJSONError(w, ErrHttpCodeOk, ErrServerInner, err)
		return
	}

	ResponseJSON(w, http.StatusOK, records)
}

func AddRecord(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var req driver.CodeSegmentRecord
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Errorf("ioutil.ReadAll err: %+v", err)
		ResponseJSONError(w, ErrRecordExpired, ErrArgsFormatError, err)
		return
	}

	if err := json.Unmarshal(body, &req); err != nil {
		log.Errorf("unmarshal err: %+v", err)
		ResponseJSONError(w, ErrRecordExpired, ErrArgsFormatError, err)
		return
	}

	if req.Title == "" {
		log.Errorf("get title empty")
		ResponseJSONError(w, ErrRecordExpired, ErrArgsMissing, errors.New("req title empty"))
		return
	}

	if req.Content == "" {
		log.Errorf("get content empty")
		ResponseJSONError(w, ErrRecordExpired, ErrArgsMissing, errors.New("req content empty"))
		return
	}

	if req.Author == "" {
		req.Author = "Anonymous"
	}

	var tm = time.Now().Unix()
	req.CreatedAt = tm
	req.UpdatedAt = tm

	if req.ShortKey == "" {
		sk, err := genSK(req.Content)
		if err != nil {
			log.Errorf("get sk err: %v", err)
			ResponseJSONError(w, ErrHttpCodeOk, ErrServerInner, err)
			return
		}
		req.ShortKey = sk
	}

	if req.Lifecycle == 0 {
		req.Lifecycle = 1
	}

	_, err = driver.GetCollection().InsertOne(context.Background(), req)
	if err != nil {
		log.Errorf("insert err: %+v", err)
		ResponseJSONError(w, ErrHttpCodeOk, ErrServerInner, err)
		return
	}

	ResponseJSON(w, http.StatusOK, req)
}

func DelRecord(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var sk = r.URL.Query().Get("sk")
	if sk == "" {
		log.Errorf("sk empty")
		ResponseJSONError(w, ErrRecordExpired, ErrArgsMissing, errors.New("short_key empty"))
		return
	}
	err := driver.GetCollection().FindOneAndDelete(context.Background(), bson.M{"short_key": sk}).Err()
	if err != nil {
		log.Errorf("find and delete err: %v", err)
		ResponseJSONError(w, ErrRecordExpired, ErrRecordNotFound, err)
		return
	}

	ResponseJSON(w, http.StatusOK, nil)
}

func SetRecord(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var sk = r.URL.Query().Get("sk")

	if sk == "" {
		log.Errorf("get sk empty")
		ResponseJSONError(w, ErrRecordExpired, ErrArgsMissing, errors.New("short_key empty"))
		return
	}

	var req driver.CodeSegmentRecord
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Errorf("ioutil.ReadAll err: %v", err)
		ResponseJSONError(w, ErrRecordExpired, ErrArgsFormatError, err)
		return
	}

	if err := json.Unmarshal(body, &req); err != nil {
		log.Errorf("Unmarshal err: %v", err)
		ResponseJSONError(w, ErrRecordExpired, ErrArgsFormatError, err)
		return
	}

	if req.Title == "" {
		log.Errorf("get title empty")
		ResponseJSONError(w, ErrRecordExpired, ErrArgsMissing, errors.New("req title empty"))
		return
	}

	if req.Content == "" {
		log.Errorf("get content empty")
		ResponseJSONError(w, ErrRecordExpired, ErrArgsMissing, errors.New("req content empty"))
		return
	}

	if req.Lifecycle == 0 {
		req.Lifecycle = 1
	}

	// need find first
	var old driver.CodeSegmentRecord
	err = driver.GetCollection().FindOne(context.Background(), bson.M{"short_key": sk}).Decode(&old)
	if err != nil {
		log.Errorf("find record err: %+v", err)
		ResponseJSONError(w, ErrRecordExpired, ErrRecordNotFound, err)
		return
	}

	req.UpdatedAt = time.Now().Unix()
	req.CreatedAt = old.CreatedAt
	if req.ShortKey == "" && old.ShortKey != "" {
		req.ShortKey = old.ShortKey
	} else if req.ShortKey == "" {
		sk, err := genSK(req.Content)
		if err != nil {
			log.Errorf("get sk err: %v", err)
			ResponseJSONError(w, ErrHttpCodeOk, ErrServerInner, err)
			return
		}
		req.ShortKey = sk
	}

	err = driver.GetCollection().FindOneAndReplace(context.Background(), bson.M{"short_key": sk}, req).Err()
	if err != nil {
		log.Errorf("find and replace err: %+v", err)
		ResponseJSONError(w, ErrHttpCodeOk, ErrRecordNotFound, err)
		return
	}

	ResponseJSON(w, http.StatusOK, req)
}

func genSK(content string) (string, error) {
	var _t = fmt.Sprintf(".%d.", time.Now().Unix())
	var h = md5.New()
	_, _ = h.Write([]byte(content + _t))
	var digest = hex.EncodeToString(h.Sum(nil))[8:24]
	front, err := strconv.ParseInt(digest[:len(digest)/2], 16, 64)
	if err != nil {
		log.Errorf("gen sk err: %v", err)
		return "", err
	}

	end, err := strconv.ParseInt(digest[len(digest)/2:], 16, 64)
	if err != nil {
		log.Errorf("gen sk err: %v", err)
		return "", err
	}

	frontDigest := common.Encode62(front)
	endDigest := common.Encode62(end)
	return fmt.Sprintf("%s%s", frontDigest, endDigest), nil
}
