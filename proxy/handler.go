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
	cNode := &config.getConfig().Nodes[index]
	mu.Lock()

	url  , err:= url.Parse(nodes[index].URL) 
	if err != nil{
		log.Fatal(err.Error())
	}
	index ++

	mu.Unlock()
	revProxy := httputil.NewSingleHostReverseProxy(url)
	revProxy.ErrorHandler =  ErrProxyHanlder(cNode)
	revProxy.ServeHTTP(w,r)


}
 
func ErrProxyHanlder(cNode *Nodes) func(http.ResponseWriter, *http.Request, error){
	return func(http.ResponseWriter, *http.Request, error){
		cNode.Up = false
	}
}