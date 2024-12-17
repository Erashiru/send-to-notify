package v1

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

func parsePaging(c *gin.Context) (page int64, limit int64, err error) {

	pageStr, ok := c.GetQuery("page")
	if ok {
		page, err = strconv.ParseInt(pageStr, 10, 64)
		if err != nil {
			return
		}
	}

	limitStr, ok := c.GetQuery("limit")
	if ok {
		limit, err = strconv.ParseInt(limitStr, 10, 64)
		if err != nil {
			return
		}
	}

	return
}
