package main

import (
	"net/http"

	"github.com/unrolled/render"
)

var rd = render.New()

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	rd.JSON(w, http.StatusOK, map[string]string{"hello": "tidb"})
}
