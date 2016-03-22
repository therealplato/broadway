package server

import (
	"log"
	"net/http"
	"os"

	"github.com/namely/broadway/instance"
	"github.com/namely/broadway/store"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// Server provides an HTTP interface to manipulate Playbooks and Instances
type Server struct {
	store      store.Store
	slackToken string
	engine     *gin.Engine
}

const slackTokenENV string = "SLACK_VERIFICATION_TOKEN"

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
	actualToken := os.Getenv(slackTokenENV)
	srvr := &Server{
		store:      s,
		slackToken: actualToken,
	}
	log.Printf("Started new server with token %s", actualToken)
	srvr.setupHandlers()
	return srvr
}

func (s *Server) setupHandlers() {
	s.engine = gin.Default()
	s.engine.POST("/instances", s.createInstance)
	s.engine.GET("/instance/:playbookID/:instanceID", s.getInstance)
	s.engine.GET("/instances/:playbookID", s.getInstances)
	s.engine.GET("/status", s.getStatus400)
	s.engine.GET("/status/:playbookID", s.getStatus400)
	s.engine.GET("/status/:playbookID/:instanceID", s.getStatus)
	s.engine.GET("/command", s.getCommand)
	s.engine.POST("/command", s.postCommand)
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
	var ia instance.Attributes
	var err = c.BindJSON(&ia)
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

func (s *Server) getStatus400(c *gin.Context) {
	c.JSON(http.StatusBadRequest, ErrorResponse{
		"error": "Use GET /status/yourPlaybookId/yourInstanceId",
	})
}

func (s *Server) getStatus(c *gin.Context) {
	status, err := instance.GetStatus(s.store, c.Param("playbookID"), c.Param("instanceID"))
	if err != nil {
		if err.Error() == "Instance does not exist." {
			c.JSON(http.StatusNotFound, ErrorResponse{
				"error": err.Error(),
			})
		} else {
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				"error": err.Error(),
			})
		}
		return
	}
	c.JSON(http.StatusOK, map[string]string{
		"status": string(status),
	})
}

func (s *Server) getCommand(c *gin.Context) {
	ssl := c.Query("ssl_check")
	log.Println(ssl)
	if ssl == "1" {
		c.String(http.StatusOK, "")
	} else {
		c.String(http.StatusBadRequest, "Use POST /command")
	}
}

// SlackCommand ...
type SlackCommand struct {
	Token       string `form:"token" binding:"required"`
	TeamID      string `form:"team_id"`
	TeamDomain  string `form:"team_domain"`
	ChannelID   string `form:"channel_id"`
	ChannelName string `form:"channel_name"`
	UserID      string `form:"user_id"`
	UserName    string `form:"user_name"`
	Command     string `form:"command"`
	Text        string `form:"text"`
	ResponseURL string `form:"response_url"`
}

func (s *Server) postCommand(c *gin.Context) {
	var form SlackCommand
	if err := c.BindWith(&form, binding.Form); err != nil {
		c.JSON(http.StatusInternalServerError, InternalError)
		return
	}
	log.Println(form.Token)
	if form.Token != s.slackToken {
		log.Printf("Expected token %s. Actual token %s.", s.slackToken, form.Token)
		c.JSON(http.StatusUnauthorized, InternalError)
		return
	}
	if form.Command != "/broadway" {
		c.JSON(http.StatusBadRequest, InternalError)
		return
	}
	if form.Text == "help" {
		c.String(http.StatusOK, "/broadway status playbook1 instance1: Check the status of instance1\n /broadway deploy playbook1 instance1: Deploy instance1")
		return
	}
	output, err := helperRunCommand(form.Text)
	if err != nil {
		c.JSON(http.StatusInternalServerError, InternalError)
		return
	}
	c.String(http.StatusOK, output)
	return
}

func helperRunCommand(text string) (string, error) {
	return "unimplemented :sadpanda:", nil
}
