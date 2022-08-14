package main

import (
	"net/http"
	"time"

	"example.com/m/v2/config"
	"example.com/m/v2/store"
	"example.com/m/v2/utils"
	"github.com/gin-gonic/gin"
	"github.com/jszwec/csvutil"
)

func routes(r *gin.Engine) {
	v1Group := r.Group("/v1")
	{
		v1Group.GET("/count", countHandler)
		v1Group.GET("/history", historyHandler)
	}

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, "")
	})
	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, "")
	})
}

func countHandler(c *gin.Context) {
	rep, err := utils.RunInTX(calculateCount)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, ok := rep.([]*store.CountStatistic)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"data disorderd": ""})
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"count": result[0].Count,
	})
}

func historyHandler(c *gin.Context) {
	countList, err := utils.GetCountLists(config.DataPath)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"history": countList,
	})
}

func calculateCount() (interface{}, error) {
	countList, err := utils.GetCountLists(config.DataPath)
	if err != nil {
		return nil, err
	}

	if len(countList) == 0 {
		countList = []*store.CountStatistic{
			{
				CreatedAt: time.Now(),
				Count:     1,
			},
		}
	} else {
		if time.Since(countList[0].CreatedAt) < time.Duration(config.MovingWindowLimit)*time.Second {
			countList[0].Count += 1
		} else {
			countList = append([]*store.CountStatistic{
				{
					CreatedAt: time.Now(),
					Count:     1,
				},
			}, countList...)
		}
	}

	//  Scroll to delete
	if len(countList) > config.DataLimit {
		countList = countList[:config.DataLimit]
	}

	b, err := csvutil.Marshal(countList)
	if err != nil {
		return nil, err
	}

	err = utils.WriteCsv(config.DataPath, b)
	if err != nil {
		return nil, err
	}
	return countList, nil
}
