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
	cw := NewConfig()
    json.Unmarshal(data, &config)

    go helthChecker()


	switch blancerType{
		case "roundRobin":
			serve  = http.Server{
				Addr:    ":" + config.Proxy.Port,
				Handler: http.HandlerFunc(roundRobinHandler),
			}
		case "LeastConnection:":
			serve  = http.Server{
				Addr:    ":" + config.Proxy.Port,
			 
				ConnState: cw.OnStateChange,
			}
	}
 
    if err = serve.ListenAndServe(); err != nil {
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