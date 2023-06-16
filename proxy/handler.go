package proxy

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
)
var mu sync.Mutex
var index = 0 
func loadBalancerHandler(w http.ResponseWriter , r *http.Request){
	config := NewConfig()
	nodes := config.getConfig().Nodes
	mu.Lock()

	url  , err:= url.Parse(nodes[index].URL) 
	if err != nil{
		log.Fatal(err.Error())
	}
	index ++

	mu.Unlock()
	revProxy := httputil.NewSingleHostReverseProxy(url)
	revProxy.ServeHTTP(w,r)


}