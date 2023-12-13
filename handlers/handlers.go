package handlers

import (
	"math/rand"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/twistingmercury/monitoring/health"
)

func PingHandler(c *gin.Context) {
	time.Sleep(time.Duration(sleepTime(5, 250)) * time.Millisecond)
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func PongHandler(c *gin.Context) {
	time.Sleep(time.Duration(sleepTime(10, 500)) * time.Millisecond)
	c.JSON(200, gin.H{
		"message": "ping",
	})
}

func sleepTime(min, max int) int {
	return rand.Intn(max-min) + min
}

func CheckMSSQL() (hsr health.StatusResult) {
	// connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d", "myServer", "user", "pwd", 1234)
	// db, err := sql.Open("mssql", connString)
	// if err != nil {
	// 	hsr.Status = healthcheck.HealthStatusCritical
	// 	hsr.Message = err.Error()
	// 	return
	// }
	// defer db.Close()

	// query, err := db.Prepare("SELECT 1;") // --> 'SELECT 1;' is the fastest query that can be returned from a working MSSQL database.
	// if err != nil {
	// 	hsr.Status = healthcheck.HealthStatusCritical
	// 	hsr.Message = err.Error()
	// 	return
	// }
	// defer query.Close()
	// r := query.QueryRow()
	// var ans int
	// err = r.Scan(&ans)
	// if err != nil {
	// 	hsr.Status = healthcheck.HealthStatusCritical
	// 	hsr.Message = err.Error()
	// 	return
	// }
	// if ans == 1 {
	// 	hsr.Status = healthcheck.HealthStatusOK
	// 	hsr.Message = "ok"
	// }

	hsr.Status = health.HealthStatusOK
	hsr.Message = "ok"
	return
}