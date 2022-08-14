package config

import "time"

// configuration
var (
	DataPath          = "/data/data.csv"
	DataLimit         = 5  // Scroll to delete the oldest data, keep maximum $dataLimit row of data.
	MovingWindowLimit = 60 // second, keep maximum $movingWindow of data.
	LockKey           = "/distributed-lock/"
	EtcdEndpoints     = []string{"http://43.132.128.25:2379", // for debug
		"http://etcd.etcd.svc.cluster.local:2379", // in k8s cluster
	}
	EtcdTTl = time.Second * 3
)
