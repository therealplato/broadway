package server

import (
	"net/http"

	"github.com/namely/broadway/instance"
	"github.com/namely/broadway/store"

	"github.com/gin-gonic/gin"
)

type Server struct {
	store store.Store

	engine *gin.Engine
}

type ErrorResponse map[string]string

var InternalError = map[string]string{"error": "Internal Server Error"}

var UnprocessableEntity ErrorResponse = ErrorResponse{"error": "Unprocessable Entity"}

func InvalidError(message string) ErrorResponse {
	return ErrorResponse{"error": "Unprocessable Entity: " + message}
}

var NotFoundError ErrorResponse = ErrorResponse{"error": "Not Found"}

func New(s store.Store) *Server {
	srvr := &Server{store: s}
	srvr.setupHandlers()
	return srvr
}

func (s *Server) setupHandlers() {
	s.engine = gin.Default()
	s.engine.POST("/instances", s.createInstance)
	s.engine.GET("/instance/:playbookId/:instanceId", s.getInstance)
	s.engine.GET("/instances/:playbookId", s.getInstances)
}

func (s *Server) Handler() http.Handler {
	return s.engine
}

func (s *Server) Run(addr ...string) error {
	return s.engine.Run(addr...)
}

func (s *Server) createInstance(c *gin.Context) {
	var ia instance.InstanceAttributes
	var err error = c.BindJSON(&ia)
	if err != nil {
		c.JSON(422, InvalidError("Missing: "+err.Error()))
		return
	}

	i := instance.New(s.store, &ia)
	err = i.Save()
	if err != nil {
		c.JSON(http.StatusInternalServerError, InternalError)
		return
	}

	c.JSON(http.StatusCreated, i.Attributes())
}

func (s *Server) getInstance(c *gin.Context) {
	playbookId := c.Param("playbookId")
	instanceId := c.Param("instanceId")
	i, err := instance.Get(playbookId, instanceId)
	if err != nil && err.Error() == "Instance does not exist." {
		c.JSON(http.StatusNotFound, NotFoundError)
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, InternalError)
		return
	}

	c.JSON(http.StatusOK, i.Attributes())
}

func (s *Server) getInstances(c *gin.Context) {
	instances, err := instance.List(s.store, c.Param("playbookId"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, InternalError)
		return
	} else if len(instances) == 0 {
		c.JSON(http.StatusNoContent, instances)
		return
	} else {
		c.JSON(http.StatusOK, instances)
		return
	}
}
