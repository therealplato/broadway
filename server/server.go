package server

import (
	"net/http"

	"github.com/namely/broadway/broadway"
	"github.com/namely/broadway/instance"
	"github.com/namely/broadway/services"
	"github.com/namely/broadway/store"

	"github.com/gin-gonic/gin"
)

// Server provides an HTTP interface to manipulate Playbooks and Instances
type Server struct {
	store store.Store

	engine *gin.Engine
}

// ErrorResponse represents a JSON response to be returned in failure cases
type ErrorResponse map[string]string

// InternalError represents a JSON response for status 500
var InternalError = map[string]string{"error": "Internal Server Error"}

// UnprocessableEntity represents a generic JSON response for bad requests
var UnprocessableEntity = ErrorResponse{"error": "Unprocessable Entity"}

// InvalidError creates an ErrorResponse with a custom message
func InvalidError(message string) ErrorResponse {
	return ErrorResponse{"error": "Unprocessable Entity: " + message}
}

// NotFoundError represents a JSON response for status 404
var NotFoundError = ErrorResponse{"error": "Not Found"}

// New instantiates a new Server and binds its handlers. The Server will look
// for playbooks and instances in store `s`
func New(s store.Store) *Server {
	srvr := &Server{store: s}
	srvr.setupHandlers()
	return srvr
}

func (s *Server) setupHandlers() {
	s.engine = gin.Default()
	s.engine.POST("/instances", s.createInstance)
	s.engine.GET("/instance/:playbookID/:instanceID", s.getInstance)
	s.engine.GET("/instances/:playbookID", s.getInstances)
}

// Handler returns a reference to the Gin engine that powers Server
func (s *Server) Handler() http.Handler {
	return s.engine
}

// Run starts the server on the specified address
func (s *Server) Run(addr ...string) error {
	return s.engine.Run(addr...)
}

func (s *Server) createInstance(c *gin.Context) {
	var i broadway.Instance
	if err := c.BindJSON(&i); err != nil {
		c.JSON(422, InvalidError("Missing: "+err.Error()))
		return
	}

	repo := broadway.NewInstanceRepo(store.New())
	service := services.NewInstanceService(repo)
	err := service.Create(i)

	if err != nil {
		c.JSON(http.StatusInternalServerError, InternalError)
		return
	}

	c.JSON(http.StatusCreated, i)
}

func (s *Server) getInstance(c *gin.Context) {
	playbookID := c.Param("playbookID")
	instanceID := c.Param("instanceID")
	i, err := instance.Get(playbookID, instanceID)
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
	instances, err := instance.List(s.store, c.Param("playbookID"))
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
