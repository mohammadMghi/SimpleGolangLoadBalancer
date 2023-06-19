package proxy

import (
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"sync/atomic"
)


type ConnectionWatcher struct{
	number int64
}

func NewCountWacher ()ConnectionWatcher{
	return ConnectionWatcher{}
}


func (cw *ConnectionWatcher) OnStateChange(conn net.Conn, state http.ConnState) {
    switch state {
    case http.StateNew:
        cw.Add(1)
    case http.StateHijacked, http.StateClosed:
        cw.Add(-1)
    }
}

func (cw *ConnectionWatcher) Count() int {
    return int(atomic.LoadInt64(&cw.number))
}
func (cw *ConnectionWatcher) Add(c int64) {
    atomic.AddInt64(&cw.number, c)
}


//any request call this func and this func checks which server is least connection and send request to it
var detectedLowestConnectionsURL string
var detectedLowestConnectionsMap = make(map[string]int)
func (cw *ConnectionWatcher)leastHandler(w http.ResponseWriter , r *http.Request){
  
	config := NewConfig()
	nodes := config.getConfig().Nodes
	cNode := &config.getConfig().Nodes[index]

 
 

	

	 
	
	//loop on the list config and checks every request for find least
	for _ , server := range nodes{
	
		mu.Lock()
		_ = &http.Server{
			ConnState: cw.OnStateChange,
		 }
		 cw.Count()
		 detectedLowestConnectionsMap[server.URL] =   cw.Count() 
		 mu.Unlock()
	
	}

	for i , detect :=range detectedLowestConnectionsMap{
		if(detect > 1){

		}
		

	}


	revProxy := httputil.NewSingleHostReverseProxy(url.Parse(detectedLowestConnectionsURL))
	



	revProxy.ErrorHandler =  func(http.ResponseWriter,  *http.Request ,error){
		cNode.Up = false
	

		 
		 cw.lowHandler(w,r)
	}

	//checks which server is least conntions



	revProxy.ServeHTTP(w,r)


}


//This func checks if server faild then sends request to other(in the config -> next node)
var mu sync.Mutex
var index = 0 
func roundRobinHandler(w http.ResponseWriter , r *http.Request){
	config := NewConfig()
	nodes := config.getConfig().Nodes
	cNode := &config.getConfig().Nodes[index]
	mu.Lock()

	url  , err:= url.Parse(nodes[index].URL) 
	if err != nil{
		log.Fatal(err.Error())
	}
	index ++

	mu.Unlock()
	revProxy := httputil.NewSingleHostReverseProxy(url)
	revProxy.ErrorHandler =  func(w http.ResponseWriter,r  *http.Request ,e error){

		log.Fatalf(e.Error())
		cNode.Up = false

		//if faild then retry and call func again
		roundRobinHandler(w ,r )
	}
	revProxy.ServeHTTP(w,r)


}


	