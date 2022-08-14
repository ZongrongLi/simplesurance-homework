# simplesurance-homework
- It took me about a couple of days, maybe 6 hours in total to finish this, it took me some time to configure k8s from scratch, apply SSL certificate and configure HTTPS, write Dockerfile and YAML
- I will read all the records from the CSV file and check if the create time of the latest one since from is exceeded 60s, if not, increment by 1 of the latest record, and create a new record otherwise, finaly write the whole list to CSV file, which is to be truncated to 5 rows.
- host on k3s to do HA, traffic from ingress and then service and then pod, 3 pods to prevent a single point of failure, pod data dir is mounted as volume map to host which is '/data'
- use etcd lock to prevent ABA problem of concurrency, requests will be recorded one by one with reading and write to be atomic.


## demand:
```
Using only the standard library, create a Go HTTP server that on each request responds with a counter of the total number of requests that it has received during the previous 60 seconds (moving window). The server should continue to the return the correct numbers after restarting it, by persisting data to a file.
```


## file structure
```
├── config
│   └── config.go
├── data
│   └── data.csv
├── Dockerfile
├── etcd
│   ├── etcd.go
│   └── etcd_test.go
├── k8s
│   ├── ca # https ca
│   │   ├── tls.crt
│   │   └── tls.key
│   ├── deployment.yaml
│   └── secret.yaml
├── main.go
├── README.md
├── routes.go
├── store
│   └── schema.go
└── utils
    ├── utils.go
    └── utils_test.go
```







## API
v1/count: counter of the total number of requests that it has received during the previous 60 seconds
https://lizongrong.xyz/v1/count
```json
{
    "count":2
}
```

v1/history: history list of count of recent 5 record  
https://lizongrong.xyz/v1/history

```json
{
    "history":[
        {
            "CreatedAt":"2022-08-14T21:18:07.256638186Z",
            "Count":2
        },
        {
            "CreatedAt":"2022-08-14T20:48:18.004547032Z",
            "Count":1
        },
        {
            "CreatedAt":"2022-08-14T20:46:10.846339144Z",
            "Count":5
        },
        {
            "CreatedAt":"2022-08-14T20:42:45.887815604Z",
            "Count":8
        },
        {
            "CreatedAt":"2022-08-14T20:18:31.066873122Z",
            "Count":3
        }
    ]
}
```
