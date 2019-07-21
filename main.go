package main

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	linkRouter "github.com/wicoady1/gtdr-score-parser/router"
)

func main() {
	router := httprouter.New()
	router.GET("/", linkRouter.Index)
	router.POST("/uploadfile", linkRouter.UploadFile)
	router.GET("/resultimage", linkRouter.ResultImage)
	router.POST("/resultimage", linkRouter.ResultImage)
	router.ServeFiles("/sources/*filepath", http.Dir("sources"))
	router.ServeFiles("/asset/*filepath", http.Dir("asset"))

	log.Println("Serving on 8080")

	log.Fatal(http.ListenAndServe(":8080", router))
}
