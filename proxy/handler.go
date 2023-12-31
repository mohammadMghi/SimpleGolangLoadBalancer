package proxy

import (
	"crypto/md5"
	"math/rand"
	"time"

	"fmt"
	"log"

	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sort"
	"sync"
	"sync/atomic"

	"github.com/redis/go-redis/v9"
	"google.golang.org/genproto/googleapis/cloud/redis/v1"
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


func RondomSeverHandler(w http.ResponseWriter , r *http.Request){
	  
	config := NewConfig()
	nodes := config.GetConfig().Nodes
 
	rand.Seed(time.Now().UnixNano())
	randomIndex := rand.Intn(len(nodes))

	selectedNode := nodes[randomIndex]

	url , e:=url.Parse(selectedNode.URL)

	if e != nil{
		log.Fatal(e.Error())
	}
 

	revProxy := httputil.NewSingleHostReverseProxy(url)
	
 

	revProxy.ServeHTTP(w,r)

}

//Hash ip is why to routing incomming request and detemine which server should handle the request base on sessions and cookies
func ipHashHandler(w http.ResponseWriter , r *http.Request){

	 // Connect to Redis cache
	 client := redis.NewClient(&redis.Options{
        Addr:     "localhost:6379",
        Password: "", // Provide password if required
        DB:       0,  // Select the appropriate database number
    })

	//determine which parameters you need to check for pass incoming request to appropriate server
	//here...
	//for example ... 
    // Retrieve data from cache
   	//get user sessions and cookies to detemine which server should handle the request
   result, err := client.Get(r.Context(),"key").Result()
   if err == redis.Nil {
	   fmt.Println("Key does not exist")
   } else if err != nil {
	   panic(err)
   } else {
	   fmt.Println("Value:", result)
   }
 

	ipAddress := r.RemoteAddr
	 


	ip, _, err := net.SplitHostPort(ipAddress)
		if err != nil {
			log.Fatal(err)
			return
		}
	ipByte := []byte(ip)
	_ = md5.Sum(ipByte)

	//
}
 
//any request call this func and this func checks which server is least connection and send request to it
var detectedLowestConnectionsURL string
var detectedLowestConnectionsMap = make(map[int]string)
var lowServer = make(map[int]string)
func (cw *ConnectionWatcher)leastHandler(w http.ResponseWriter , r *http.Request){
  
	config := NewConfig()
	nodes := config.GetConfig().Nodes
  
 
	//loop on the list config and checks every request for find least
	for _ , server := range nodes{

		mu.Lock()
		_ = &http.Server{
			ConnState: cw.OnStateChange,
			Addr:  server.URL,
		 }
 
		 cw.Count()
		 detectedLowestConnectionsMap[ cw.Count() ] =  server.URL
		 mu.Unlock()
	
	}
 
 

	keys := make([]int, 0)
    for k, _ := range detectedLowestConnectionsMap {
        keys = append(keys, k)
    }
    sort.Ints(keys)
    for _, k := range keys {


        detectedLowestConnectionsURL  =  detectedLowestConnectionsMap[k]
        break
    }

 
 
  
	url , e:=url.Parse(detectedLowestConnectionsURL)

	if e != nil{
		log.Fatal(e.Error())
	}
 

	revProxy := httputil.NewSingleHostReverseProxy(url)
	
 

	revProxy.ServeHTTP(w,r)


}


//This func checks if server faild then sends request to other(in the config -> next node)
var mu sync.Mutex
var index = 0 
func roundRobinHandler(w http.ResponseWriter , r *http.Request){
	config := NewConfig()
	nodes := config.GetConfig().Nodes
	lenNodes := len(nodes)
	if(index == lenNodes){
		index = 0
	}
	cNode := &config.GetConfig().Nodes[index]

 
	if(lenNodes == 0){
		log.Fatal("Len of nodes is zero ... see config.json and add your servers")
	}

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


	