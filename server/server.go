package server

import (
	"net/http"

	"github.com/namely/broadway/broadway"
	"github.com/namely/broadway/instance"
	"github.com/namely/broadway/playbook"
	"github.com/namely/broadway/services"
	"github.com/namely/broadway/store"

	"github.com/gin-gonic/gin"
)

// Server provides an HTTP interface to manipulate Playbooks and Instances
type Server struct {
	store     store.Store
	playbooks []playbook.Playbook
	engine    *gin.Engine
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
	gin.SetMode(gin.ReleaseMode)
	s.engine.POST("/instances", s.createInstance)
	s.engine.GET("/instance/:playbookID/:instanceID", s.getInstance)
	s.engine.GET("/instances/:playbookID", s.getInstances)
	s.engine.GET("/status", s.getStatus400)
	s.engine.GET("/status/:playbookID", s.getStatus400)
	s.engine.GET("/status/:playbookID/:instanceID", s.getStatus)
	s.engine.POST("/deploy/:playbookID/:instanceID", s.deployInstance)
}

// SetPlaybooks passes playbooks to the server (from main.go)
func (s *Server) SetPlaybooks(pbs []playbook.Playbook) {
	s.playbooks = pbs
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

	service := services.NewInstanceService(store.New())
	err := service.Create(i)

	if err != nil {
		c.JSON(http.StatusInternalServerError, InternalError)
		return
	}

	c.JSON(http.StatusCreated, i)
}

func (s *Server) getInstance(c *gin.Context) {
	service := services.NewInstanceService(s.store)
	i, err := service.Show(c.Param("playbookID"), c.Param("instanceID"))

	if err != nil {
		switch err.(type) {
		case broadway.InstanceNotFoundError:
			c.JSON(http.StatusNotFound, NotFoundError)
			return
		default:
			c.JSON(http.StatusInternalServerError, InternalError)
			return
		}
	}
	c.JSON(http.StatusOK, i)
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

func (s *Server) getStatus400(c *gin.Context) {
	c.JSON(http.StatusBadRequest, ErrorResponse{
		"error": "Use GET /status/yourPlaybookId/yourInstanceId",
	})
}

func (s *Server) getStatus(c *gin.Context) {
	service := services.NewInstanceService(s.store)
	instance, err := service.Show(c.Param("playbookID"), c.Param("instanceID"))

	if err != nil {
		switch err.(type) {
		case broadway.InstanceNotFoundError:
			c.JSON(http.StatusNotFound, NotFoundError)
			return
		default:
			c.JSON(http.StatusInternalServerError, InternalError)
			return
		}
	}
	c.JSON(http.StatusOK, map[string]string{
		"status": string(instance.Status),
	})
}

func (s *Server) deployInstance(c *gin.Context) {
	service := services.NewInstanceService(s.store)
	_, err := service.Show(c.Param("playbookID"), c.Param("instanceID"))

	if err != nil {
		switch err.(type) {
		case broadway.InstanceNotFoundError:
			c.JSON(http.StatusNotFound, NotFoundError)
			return
		default:
			c.JSON(http.StatusInternalServerError, InternalError)
			return
		}
	}
	// instance, err = service.Deploy(instance)
	// c.JSON(http.StatusOK, instance)
}
