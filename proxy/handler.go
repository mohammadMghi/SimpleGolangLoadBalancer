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


var detectedLowestConnections string
func (cw *ConnectionWatcher)lowHandler(w http.ResponseWriter , r *http.Request){
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
 


	revProxy.ErrorHandler =  func(http.ResponseWriter,  *http.Request ,error){
		cNode.Up = false
		 _ = &http.Server{
			ConnState: cw.OnStateChange,
		 }
		 cw.Count()
		 cw.lowHandler(w,r)
	}
	revProxy.ServeHTTP(w,r)


}


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
	revProxy.ErrorHandler =  func(http.ResponseWriter,  *http.Request ,error){
		cNode.Up = false
		roundRobinHandler(w ,r )
	}
	revProxy.ServeHTTP(w,r)


}


	