package main

import (
	"math/rand"
	"time"
)

type HttpServer struct {
	Host string//域名
	Weight int//权重
}

type LoadBalance struct {
	Servers []*HttpServer
	CurIndex int//计数
}

func NewLoadBalance() *LoadBalance {
	return &LoadBalance{Servers: make([]*HttpServer, 0)}
}

func NewHttpServer(host string, weight int) *HttpServer {
	return &HttpServer{Host: host,Weight:weight}
}

func (this *LoadBalance) AddServer(server *HttpServer) *LoadBalance {
	this.Servers = append(this.Servers, server)
	return this
}

func (this *LoadBalance) SelectByRand() *HttpServer { //随机算法
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(len(this.Servers))
	return this.Servers[n]
}

func (this *LoadBalance) SelectByWeightRand() *HttpServer { //加权随机
	var serverWeight  []int
	for index,server := range this.Servers{
		for i:=0;i<server.Weight;i++ {
			serverWeight = append(serverWeight, index)
		}
	}
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(len(serverWeight))
	return this.Servers[serverWeight[n]]
}


var LB *LoadBalance
func init()  {
	LB = NewLoadBalance()
	LB.AddServer(NewHttpServer("http://localhost:9091",5)).AddServer(NewHttpServer("http://localhost:9092",2))
}