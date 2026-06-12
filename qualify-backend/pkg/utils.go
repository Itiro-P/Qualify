package pkg

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Função auxiliar que realiza o parse do parâmetro ID.
// Retorna -1 se ocorrer erros.
func ParseIdParam(c *gin.Context) (int, error) {
	_id := c.Param("id")
	id, errID := strconv.Atoi(_id)
	if errID != nil {
		c.JSON(http.StatusBadRequest, BadRequest(c.FullPath(), errID.Error()))
		return -1, errID
	}
	return id, nil
}

// Função auxiliar que realiza o parse da query.
// Retorna -1 se ocorrer erros.
func ParsePathQuery(c *gin.Context, query string) (int, error) {
	_id := c.Query(query)
	id, errID := strconv.Atoi(_id)
	if errID != nil {
		c.JSON(http.StatusBadRequest, BadRequest(c.FullPath(), errID.Error()))
		return -1, errID
	}
	return id, nil
}

// Função auxiliar que realiza o parse do path.
// Retorna -1 se ocorrer erros.
func ParsePathParam(c *gin.Context, param string) (int, error) {
	id, err := strconv.Atoi(c.Param(param))
	if err != nil {
		c.JSON(http.StatusBadRequest, BadRequest(c.FullPath(), err.Error()))
		return -1, err
	}
	return id, nil
}
