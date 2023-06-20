package proxy

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
 
	"net/url"
	"time"
)
 

 

func Server(blancerType string){
 
	var config  = NewConfig()
	var serve http.Server
    data, err := ioutil.ReadFile("./config.json")

    if err != nil {

        log.Fatal(err.Error())
    }

	cw := NewCountWacher()
    json.Unmarshal(data, &config)

    go helthChecker()

	if blancerType == "roundRobin"{
		serve  = http.Server{
			Addr:    ":" + config.Proxy.Port,
			Handler: http.HandlerFunc(roundRobinHandler),
		}
	}

	if blancerType == "LeastConnection"{
		serve  = http.Server{
			Addr:    ":" + config.Proxy.Port,
			Handler: http.HandlerFunc(cw.leastHandler),
		}
	}

 

	log.Print("Load balancer runing on port " + config.Proxy.Port +"...")

    err = serve.ListenAndServe() 


	if err != nil {
 
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
	config :=  NewConfig() 
	nodes := config.GetConfig()
 
	for{
	
		select{
	
			case <-time.C:
				for _ , node := range nodes.Nodes {
					nnode := &node
				
					url, err := url.Parse(nnode.URL); if err != nil{
						log.Fatal(err.Error())
					} 
					err = isUpService(url) ; if err != nil{
						log.Fatal("node is down :  " ,  err.Error() , nnode.URL )
					}
					fmt.Println("node is up ...  "  , nnode.URL )


				}
		}
	}
}