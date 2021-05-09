package server

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"pastebin/config"
	"pastebin/server/sfs"

	log "github.com/sirupsen/logrus"
)

func NewRouter() {
	confer := config.GetConfig()
	var router = httprouter.New()
	router.PanicHandler = func(w http.ResponseWriter, r *http.Request, v interface{}) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(fmt.Sprintf("err: %+v", v)))
	}

	router.GET("/ping", func(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
		ResponseJSON(w, http.StatusOK, struct {
			Ping string `json:"ping"`
		}{
			Ping: "Pong",
		})
	})

	router.GET("/v1/get", GetRecord)
	router.GET("/v1/list", GetRecordList)
	router.POST("/v1/add", AddRecord)
	router.PUT("/v1/set", SetRecord)
	router.DELETE("/v1/del", DelRecord)
	router.GET("/raw/:sk", GetRecord)

	router.NotFound = sfs.New(http.Dir("webui/dist"), func(writer http.ResponseWriter, request *http.Request) {
		http.ServeFile(writer, request, "webui/dist/index.html")
	})

	log.Infof("%s is running at http://localhost%s . Press Ctrl+C to stop.", confer.ProjectName, confer.Port)
	if confer.Https.Enable {
		log.Errorf("err: %v", http.ListenAndServeTLS(
			confer.Port, confer.Https.CrtFile, confer.Https.KeyFile, router))
	} else {
		log.Errorf("err: %v", http.ListenAndServe(
			confer.Port, router))
	}
}
