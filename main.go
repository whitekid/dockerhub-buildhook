package main

import (
	"flag"
	"net/http"
	"os"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	"github.com/whitekid/go-utils/request"
)

type apiServer struct {
}

func (s *apiServer) setupRouter() *gin.Engine {
	r := gin.Default()
	r.Use(cors.Default())

	r.POST("/webhook", s.handleHook)

	return r
}

func (s *apiServer) handleHook(c *gin.Context) {
	key := c.Query("auth")
	if key != os.Getenv("AUTH_KEY") {
		c.Status(http.StatusUnauthorized)
		glog.V(2).Info("Unauthorized request")
		return
	}

	triggered := 0
	for _, env := range os.Environ() {
		if !strings.HasPrefix(env, "DOCKER_HOOK_") {
			glog.V(2).Infof("skip %s", env)
			continue
		}

		kv := strings.SplitN(env, "=", 2)
		value := kv[1]

		if !strings.HasPrefix(value, "https://cloud.docker.com/api/build/") {
			glog.V(2).Infof("skip invalid hook url: %s", value)
			continue
		}

		resp, err := request.Post(value).Do()
		if err != nil {
			glog.Errorf("Fail to invoke build: %s, %s", value, err)
			continue
		}

		if resp.StatusCode != http.StatusAccepted {
			glog.Errorf("Fail to invoke build with status %d", resp.StatusCode)
			continue
		}

		glog.Infof("Build Triggered: %s", value)
		triggered++
	}

	if triggered == 0 {
		glog.Errorf("Fail to trigger build: no hook")
	}

	c.Status(http.StatusOK)
}

func (s *apiServer) serve() error {
	authKey := os.Getenv("AUTH_KEY")
	if authKey == "" {
		panic("AUTH_KEY required")
	}

	r := s.setupRouter()
	return http.ListenAndServe("127.0.0.1:9997", r)

}

func main() {
	flag.Parse()

	s := apiServer{}
	if err := s.serve(); err != nil {
		panic(err)
	}
}
