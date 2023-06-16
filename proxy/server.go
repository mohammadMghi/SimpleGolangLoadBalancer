package proxy

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"golang.org/x/tools/go/pointer"
)

func Server(){

	d := func(requst *http.Request){
		requst.URL.Scheme ="http"
		requst.URL.Host =":8081"
	}

	revProxy := &httputil.ReverseProxy{
		Director: d,
	}
	
	go helthChecker()

	ser := http.Server{
		Addr: ":8080",
		Handler : revProxy,
	}



	err := ser.ListenAndServe()

	if err != nil{
		log.Fatal(err.Error())
	}
	

}

func isUpService(url *url.URL) error{
	connection , err := net.DialTimeout("tcp" , url.Host , time.Minute * 1) ; if err != nil {
		return err
	}
	defer connection.Close()
	return nil
}

func helthChecker (){
	time := time.NewTicker(time.Minute * 1)
	nodes :=  NewConfig().Nodes
	for{
		select{
			case <-time.C:
				for _ , node := range nodes {
					nnode := &node
					url, err := url.Parse(nnode.URL); if err != nil{
						log.Fatal(err.Error())
					} 
					err = isUpService(url) ; if err != nil{
						log.Fatal("node is down : %v :: %v" ,  err.Error() , nnode.URL )
					}
					fmt.Println("node is up ... %v"  , nnode.URL )


				}
		}
	}
}