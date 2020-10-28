package v1

import (
	"codisgraph/src/hsc"
	"codisgraph/src/middleware/mysql"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetAll(c *gin.Context) {
	code := hsc.SUCCESS
	Graph := mysql.DB.GraphGetAll()
	//log.Println( gettagtotal, count)
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  hsc.GetMsg(code),
		"data": Graph,
	})

}
