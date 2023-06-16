package proxy

import (
	"log"
	"net/http"
	"net/http/httputil"
)

func Server(){

	d := func(requst *http.Request){
		requst.URL.Scheme ="http"
		requst.URL.Host =":8081"
	}

	revProxy := &httputil.ReverseProxy{
		Director: d,
	}

	ser := http.Server{
		Addr: ":8080",
		Handler : revProxy,
	}


	err := ser.ListenAndServe()

	if err != nil{
		log.Fatal(err.Error())
	}
	

}