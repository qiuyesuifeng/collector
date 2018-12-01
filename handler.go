package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/unrolled/render"
)

var rd = render.New()

type ErrorResponse struct {
	ErrorCode int64  `json:"error_code"`
	ErrorMsg  string `json:"error_msg"`
}

func (r *ErrorResponse) String() string {
	if r == nil {
		return "<nil>"
	}
	return fmt.Sprintf("[ErrorResponse](%+v)", *r)
}

type KafkaStreamData struct {
	Data   string `json:"data"`
	Offset int64  `json:"offset"`
}

type KafkaStreamResponse struct {
	ErrorResponse

	Msgs []KafkaStreamData `json:"msgs"`
}

func (r *KafkaStreamResponse) String() string {
	if r == nil {
		return "<nil>"
	}
	return fmt.Sprintf("[KafkaStreamResponse](%+v)", *r)
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	rd.JSON(w, http.StatusOK, map[string]string{"hello": "tidb"})
}

func DataGetHandler(w http.ResponseWriter, r *http.Request) {
	response := &KafkaStreamResponse{}

	values := r.URL.Query()
	topic := values.Get("topic")
	if topic == "" {
		Log.Warn("[DataGetHandler]Parse topic argument fail[url]%s", r.RequestURI)
		response.ErrorCode = 1
		response.ErrorMsg = "Parse topic argument fail"
		rd.JSON(w, http.StatusOK, response)
		return
	}

	offset := values.Get("offset")
	offsetNum, err := strconv.ParseInt(offset, 10, 64)
	if err != nil {
		Log.Warn("[DataGetHandler]Parse offset argument fail[url]%s[offset]%s[error]%s", r.RequestURI, offset, err.Error())
		response.ErrorCode = 1
		response.ErrorMsg = "Parse offset argument fail"
		rd.JSON(w, http.StatusOK, response)
		return
	}

	batch := values.Get("batch")
	batchNum, err := strconv.ParseInt(batch, 10, 64)
	if err != nil {
		Log.Warn("[DataGetHandler]Parse batch argument fail[url]%s[batch]%s[error]%s", r.RequestURI, batch, err.Error())
		response.ErrorCode = 1
		response.ErrorMsg = "Parse batch argument fail"
		rd.JSON(w, http.StatusOK, response)
		return
	}

	// set response data info
	data, err := service.FetchKafkaMsgs(topic, offsetNum, batchNum)
	if err != nil {
		Log.Warn("[DataGetHandler]FetchKafkaConsume fail[url]%s[error]%s", r.RequestURI, err.Error())
		response.ErrorCode = 1
		response.ErrorMsg = "FetchKafkaConsume fail"
		rd.JSON(w, http.StatusOK, response)
		return
	}

	response.Msgs = data

	rd.JSON(w, http.StatusOK, response)
}
