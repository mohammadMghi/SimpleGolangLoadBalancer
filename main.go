package main

import "github.com/mohammadmghi/simplegGolangLoadBalancer/proxy"

func main(){
	proxy.Server("roundRobin")	
}