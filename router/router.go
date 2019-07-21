package router

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/wicoady1/gtdr-score-parser/imageextractor"
)

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if r.Method == "GET" {
		crutime := time.Now().Unix()
		h := md5.New()
		io.WriteString(h, strconv.FormatInt(crutime, 10))
		token := fmt.Sprintf("%x", h.Sum(nil))

		err := RenderPage(w, "imagemaker", map[string]string{
			"Token": token,
		})
		if err != nil {
			log.Println(err)
		}

	}
}

func ResultImage(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	err := RenderPage(w, "imageresult", map[string]string{
		"ImageResult": "/assets/images/output.png",
	})
	if err != nil {
		log.Println(err)
	}
}

func UploadFile(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	// Allow access from another host
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method == "POST" {

		r.ParseMultipartForm(32 << 20)
		file, handler, err := r.FormFile("upload_file")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()

		newFilePath := "asset/" + handler.Filename
		destFile, err := os.Create(newFilePath)
		defer destFile.Close()

		io.Copy(destFile, file)
		destFile.Sync()

		//-----

		language := r.FormValue("title_language")

		resp, err := imageextractor.ParseImage(newFilePath, language)
		if err != nil {
			fmt.Println(err)

			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		}

		responseEncoded, err := json.Marshal(resp)
		if err != nil {
			fmt.Println(err)

			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		}

		log.Printf("%+v %s", resp, string(responseEncoded))

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(responseEncoded)
	}
}
