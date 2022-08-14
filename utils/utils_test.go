package utils

import (
	"fmt"
	"os"
	"testing"
	"time"

	"example.com/m/v2/config"
	"example.com/m/v2/etcd"
	"example.com/m/v2/store"
	"github.com/jszwec/csvutil"
	"gopkg.in/go-playground/assert.v1"
)

func TestCSVReadNoExist(t *testing.T) {
	filename := "./data/" + fmt.Sprintf("%d", time.Now().UnixMicro()) + ".csv"
	_, err := GetCountLists(filename)
	if err == nil {
		t.Fatal("file not eixst report error")
	}
}

func TestCSVReadEmpty(t *testing.T) {
	filename := "../data/" + fmt.Sprintf("%d", time.Now().UnixMicro()) + ".csv"

	err := WriteCsv(filename, []byte{})
	if err != nil {
		t.Errorf("%+v\n", err)
	}

	defer func() {
		err = os.Remove(filename)
		if err != nil {
			t.Errorf("%+v\n", err)
		}
	}()

	got, err := GetCountLists(filename)
	if err == nil {
		t.Fatal("file not eixst report error")
	}
	assert.Equal(t, len(got), 0)
}

func TestCSVReadWrite(t *testing.T) {
	filename := "../data/" + fmt.Sprintf("%d", time.Now().UnixMicro()) + ".csv"
	countList := []*store.CountStatistic{
		{
			CreatedAt: time.Now(),
			Count:     1,
		},
	}

	b, err := csvutil.Marshal(countList)
	if err != nil {
		t.Errorf("%+v\n", err)
	}

	err = WriteCsv(filename, b)
	if err != nil {
		t.Errorf("%+v\n", err)
	}

	defer func() {
		err = os.Remove(filename)
		if err != nil {
			t.Errorf("%+v\n", err)
		}
	}()

	got, err := GetCountLists(filename)
	if err != nil {
		t.Errorf("%+v\n", err)
	}
	assert.Equal(t, len(got), 1)

	countList = append(countList, &store.CountStatistic{
		CreatedAt: time.Now(),
		Count:     1,
	})

	b, err = csvutil.Marshal(countList)
	if err != nil {
		t.Errorf("%+v\n", err)
	}

	err = WriteCsv(filename, b)
	if err != nil {
		t.Errorf("%+v\n", err)
	}

	got, err = GetCountLists(filename)
	if err != nil {
		t.Errorf("%+v\n", err)
	}
	assert.Equal(t, len(got), 2)
}

func TestWriteCSV(t *testing.T) {
	filename := "../data/notexist" + fmt.Sprintf("%d", time.Now().UnixMicro())
	err := WriteCsv(filename, []byte{})
	if err != nil {
		t.Errorf("%+v\n", err)
	}
	defer func() {
		err = os.Remove(filename)
		if err != nil {
			t.Errorf("%+v\n", err)
		}
	}()
	err = WriteCsv(filename, []byte("write the same file again"))
	if err != nil {
		t.Errorf("%+v\n", err)
	}
}

func TestCreateDataFileIfNotExist(t *testing.T) {
	filename := "../data/CreateDataFileIfNotExist.txt" + fmt.Sprintf("%d", time.Now().UnixMicro())
	CreateDataFileIfNotExist(filename)
	fs, err := os.Open(filename)
	if err != nil {
		t.Errorf("%+v\n", err)
	}
	defer fs.Close()
	defer func() {
		err = os.Remove(filename)
		if err != nil {
			t.Errorf("%+v\n", err)
		}
	}()
}

func TestRunInTX(t *testing.T) {
	etcd.InitEtcd(config.EtcdEndpoints, config.EtcdTTl)

	got, err := RunInTX(func() (interface{}, error) {
		return "hello world", nil
	})
	if err != nil {
		t.Errorf("%+v\n", err)
	}
	assert.Equal(t, got.(string), "hello world")
}
