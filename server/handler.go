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
	"io/ioutil"
	"net/http"
	"pastebin/common"
	"pastebin/server/driver"
	"strconv"
	"time"
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

func ResponseJSONError(w http.ResponseWriter, httpCode int, err error) {
	//公共的响应头设置
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Content-Type", "application/json;charset=utf-8")

	if httpCode == 0 {
		httpCode = http.StatusInternalServerError
	}

	var resp = &Response{
		Code:    httpCode,
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

func ResponseHTMLError(w http.ResponseWriter, httpCode int, err error) {
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

func GetRecord(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var sk = r.URL.Query().Get("sk")
	var format = r.URL.Query().Get("format")
	if sk == "" {
		log.Errorf("sk empty")
		ResponseJSONError(w, 400, errors.New("short_key empty"))
		return
	}

	var req driver.CodeSegmentRecord
	err := driver.GetCollection().FindOne(context.Background(), bson.M{"short_key": sk}).Decode(&req)
	if err != nil {
		log.Errorf("find record err: %+v", err)
		ResponseJSONError(w, 400, err)
		return
	}

	if format == "raw" {
		ResponseHTML(w, http.StatusOK, []byte(req.Content))
		return
	}

	ResponseJSON(w, http.StatusOK, req)
}

func AddRecord(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var format = r.URL.Query().Get("format")
	var req driver.CodeSegmentRecord
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Errorf("ioutil.ReadAll err: %+v", err)
		ResponseJSONError(w, 400, err)
		return
	}

	if err := json.Unmarshal(body, &req); err != nil {
		log.Errorf("unmarshal err: %+v", err)
		ResponseJSONError(w, 400, err)
		return
	}

	if req.Title == "" {
		log.Errorf("get title empty")
		ResponseJSONError(w, 400, errors.New("req title empty"))
		return
	}

	if req.Content == "" {
		log.Errorf("get content empty")
		ResponseJSONError(w, 400, errors.New("req content empty"))
		return
	}

	req.CreatedAt = time.Now().Unix()

	if req.ShortKey == "" {
		sk, err := genSK(req.Content)
		if err != nil {
			log.Errorf("get sk err: %v", err)
			ResponseJSONError(w, 500, err)
			return
		}
		req.ShortKey = sk
	}

	_, err = driver.GetCollection().InsertOne(context.Background(), req)
	if err != nil {
		log.Errorf("insert err: %+v", err)
		ResponseJSONError(w, 500, err)
		return
	}

	if format == "" {
		// 输出某些数据
	}

	ResponseJSON(w, http.StatusOK, req)
}

func DelRecord(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var sk = r.URL.Query().Get("sk")
	if sk == "" {
		log.Errorf("sk empty")
		ResponseJSONError(w, 400, errors.New("short_key empty"))
		return
	}
	err := driver.GetCollection().FindOneAndDelete(context.Background(), bson.M{"short_key": sk}).Err()
	if err != nil {
		log.Errorf("find and delete err: %v", err)
		ResponseJSONError(w, 400, err)
		return
	}

	ResponseJSON(w, http.StatusOK, Response{})
}

func SetRecord(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var sk = r.URL.Query().Get("sk")
	if sk == "" {
		log.Errorf("get sk empty")
		ResponseJSONError(w, 400, errors.New("short_key empty"))
		return
	}

	var req driver.CodeSegmentRecord
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Errorf("ioutil.ReadAll err: %v", err)
		ResponseJSONError(w, 400, err)
		return
	}

	if err := json.Unmarshal(body, &req); err != nil {
		log.Errorf("Unmarshal err: %v", err)
		ResponseJSONError(w, 400, err)
		return
	}

	if req.Title == "" {
		log.Errorf("get title empty")
		ResponseJSONError(w, 400, errors.New("req title empty"))
		return
	}

	if req.Content == "" {
		log.Errorf("get content empty")
		ResponseJSONError(w, 400, errors.New("req content empty"))
		return
	}

	// need find first
	var old driver.CodeSegmentRecord
	err = driver.GetCollection().FindOne(context.Background(), bson.M{"short_key": sk}).Decode(&old)
	if err != nil {
		log.Errorf("find record err: %+v", err)
		ResponseJSONError(w, 400, err)
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
			ResponseJSONError(w, 500, err)
			return
		}
		req.ShortKey = sk
	}

	err = driver.GetCollection().FindOneAndReplace(context.Background(), bson.M{"short_key": sk}, req).Err()
	if err != nil {
		log.Errorf("find and replace err: %+v", err)
		ResponseJSONError(w, 400, err)
		return
	}

	ResponseJSON(w, http.StatusOK, Response{
		Code:    http.StatusOK,
		Message: "ok",
		Data:    req,
	})
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
