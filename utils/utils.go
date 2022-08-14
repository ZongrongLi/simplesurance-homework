package utils

import (
	"context"
	"encoding/csv"
	"io"
	"io/ioutil"
	"log"
	"os"
	"time"

	"example.com/m/v2/config"
	"example.com/m/v2/etcd"
	"example.com/m/v2/store"
	"github.com/jszwec/csvutil"
)

func GetCountLists(filename string) ([]*store.CountStatistic, error) {
	fs, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	defer fs.Close()

	r := csv.NewReader(fs)

	dec, err := csvutil.NewDecoder(r)
	if err != nil {
		return nil, err
	}

	var counts []*store.CountStatistic
	for {
		c := store.CountStatistic{}

		if err := dec.Decode(&c); err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		counts = append(counts, &c)
	}

	return counts, nil
}

func WriteCsv(fileName string, content []byte) error {
	return ioutil.WriteFile(fileName, content, 0666)
}

func CreateDataFileIfNotExist(filename string) {
	fs, err := os.Open(filename)
	if err != nil {
		if !os.IsNotExist(err) {
			panic(err)
		}

		countList := []*store.CountStatistic{
			{
				CreatedAt: time.Now(),
				Count:     0,
			},
		}
		b, err := csvutil.Marshal(countList)
		if err != nil {
			log.Fatalf("csvutil.Marshal, err is %+v", err)
		}
		err = WriteCsv(filename, b)
		if err != nil {
			log.Fatalf("csvutil.Marshal, err is %+v", err)
		}

		return
	}

	defer func() {
		err := fs.Close()
		if err != nil {
			log.Fatalf("can not open the file, err is %+v", err)
		}
	}()
}

func RunInTX(f func() (interface{}, error)) (interface{}, error) {
	l, err := etcd.MustGetEtcd().GetLock(config.LockKey)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := l.Lock(ctx); err != nil {
		return nil, err
	}

	defer func() {
		if err := l.Unlock(ctx); err != nil {
			log.Fatal(err)
		}
	}()
	return f()
}
