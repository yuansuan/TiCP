package mock

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// httpRequestExample 示例
type httpRequestExample struct {
	ID   int64  `json:"id" query:"id" form:"id"`
	Name string `json:"name" query:"name" form:"name"`
	Age  int    `json:"age" query:"age" form:"age"`
}

var requestExample = &httpRequestExample{
	ID:   1,
	Name: "xx",
	Age:  18,
}

// getRequestExample 示例
func getRequestExample(c *gin.Context) {
	// logger := logging.GetLogger(c)

	idstr := c.Param("id")
	namestr := c.Query("name")

	// 转换成int64
	id, err := strconv.ParseInt(idstr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "id is invalid",
		})
		return
	}

	if id == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "id is empty",
		})
		return
	}

	if namestr != "" && namestr != requestExample.Name {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "name is invalid",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": requestExample,
	})
}

// TestGetRequestExample 示例
func TestGetRequestExample(t *testing.T) {
	/* ------------------------------- normal test ------------------------------ */
	// mock gin context
	w := httptest.NewRecorder()
	c := GinContext(w)

	// set your logger or some other things
	// c.Set(logging.LoggerName,logger)

	// mock path params like /api/v1/example/:id
	p := gin.Params{
		{Key: "id", Value: "1"},
	}

	// mock query params like /api/v1/example?name=xx
	u := url.Values{}
	u.Add("name", "xx")

	HTTPRequest(c, http.MethodGet, nil, p, u)

	// call api
	getRequestExample(c)

	// check response
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "{\"data\":{\"id\":1,\"name\":\"xx\",\"age\":18}}", w.Body.String())

	/* ------------------------------- error test ------------------------------ */
	w = httptest.NewRecorder() // need to reset recorder
	c = GinContext(w)          // need to reset recorder

	p2 := gin.Params{
		{Key: "id", Value: "0"},
	}

	HTTPRequest(c, http.MethodGet, nil, p2, u)

	// call api
	getRequestExample(c)

	// check response
	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, "{\"error\":\"id is empty\"}", w.Body.String())

	w = httptest.NewRecorder()
	c = GinContext(w)

	u2 := url.Values{
		"name": []string{"yy"},
	}

	HTTPRequest(c, http.MethodGet, nil, p, u2)

	// call api
	getRequestExample(c)

	// check response
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "{\"error\":\"name is invalid\"}", w.Body.String())
}

// postRequestExample 示例
func postRequestExample(c *gin.Context) {
	idstr := c.Param("id")
	// 转换成int64
	id, err := strconv.ParseInt(idstr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "id is invalid",
		})
		return
	}

	var example httpRequestExample
	err = c.ShouldBindJSON(&example)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "request body is invalid",
		})
		return
	}
	example.ID = id

	c.JSON(http.StatusOK, gin.H{
		"data": example,
	})
}

func TestPostRequestExample(t *testing.T) {
	// mock gin context
	w := httptest.NewRecorder()
	c := GinContext(w)

	// set your logger or some other things
	// c.Set(logging.LoggerName,logger)

	// mock path params like /api/v1/example/:id
	p := gin.Params{
		{Key: "id", Value: "2"},
	}

	HTTPRequest(c, http.MethodPost, requestExample, p, nil)

	// call api
	postRequestExample(c)

	// check response
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "{\"data\":{\"id\":2,\"name\":\"xx\",\"age\":18}}", w.Body.String())
}
