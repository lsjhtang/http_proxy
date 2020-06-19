package main

import (
	"hash/crc32"
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

func (this *LoadBalance) SelectByWeightRand2() *HttpServer { //加权随机(改良版)
	rand.Seed(time.Now().UnixNano())
	var sums int
	for _,server := range this.Servers{
		sums += server.Weight
	}
	n := rand.Intn(sums)
	sums=0
	for _,server := range this.Servers{
		sums += server.Weight
		if n < sums {
			return server
		}
	}
	return this.Servers[0]
}


func (this *LoadBalance) RandRobin() *HttpServer { //轮询算法
	server := this.Servers[this.CurIndex]
	this.CurIndex = (this.CurIndex+1) % len(this.Servers)
	return server
}

func (this *LoadBalance) RandRobin2() *HttpServer { //加权轮询算法
	server := this.Servers[0]
	sums:=0
	for i:=0;i<len(this.Servers);i++{
		sums += this.Servers[i].Weight
		if this.CurIndex < sums {
			server = this.Servers[i]
			if i != len(this.Servers)-1 && this.CurIndex+1 == sums {//到达最后一轮且循环到最后一次
				this.CurIndex ++
			}else {
				this.CurIndex = (this.CurIndex+1) % sums
			}
			break
		}
	}
	return server
}

func (this *LoadBalance) SelectByIpHash(ip string) *HttpServer { //ip取余
	index := int(crc32.ChecksumIEEE([]byte(ip))) % len(this.Servers)
	return this.Servers[index]
}

var LB *LoadBalance
func init()  {
	LB = NewLoadBalance()
	LB.AddServer(NewHttpServer("http://localhost:9091",5)).AddServer(NewHttpServer("http://localhost:9092",2))
}