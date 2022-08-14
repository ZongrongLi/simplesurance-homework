package main

import (
	"example.com/m/v2/config"
	"example.com/m/v2/etcd"
	"example.com/m/v2/utils"
	"github.com/gin-gonic/gin"
)

func main() {
	utils.CreateDataFileIfNotExist(config.DataPath)
	etcd.InitEtcd(config.EtcdEndpoints, config.EtcdTTl)
	r := gin.Default()

	routes(r)

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
