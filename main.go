package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.Use(cors.Default())

	r.POST("/webhook", handleHook)

	return r
}

func handleHook(c *gin.Context) {
	// https: //cloud.docker.com/api/build/v1/source/83212efd-e24b-48b9-b7a4-464389836ffc/trigger/e4c79dc1-6e56-416f-bcd9-f92c86d69277/call/
	triggered := 0
	for _, k := range os.Environ() {
		if !strings.HasPrefix(k, "DOCKER_HOOK_") {
			continue
		}

		env := os.Getenv(k)
		if !strings.HasPrefix(env, "https://cloud.docker.com/api/build/") {
			continue
		}

		resp, err := http.DefaultClient.Post(env, "", nil)
		if err != nil {
			log.Printf("Fail to invoke build: %s, %s", env, err)
			continue
		}

		if resp.StatusCode != http.StatusAccepted {
			log.Printf("Fail to invoke build with status %d", resp.StatusCode)
			continue
		}

		log.Printf("Build Triggered: %s", env)
		triggered++
	}

	if triggered == 0 {
		log.Printf("Fail to trigger build: no hook")
	}

	c.Status(http.StatusOK)
}

func main() {
	r := setupRouter()
	if err := http.ListenAndServe("127.0.0.1:9997", r); err != nil {
		panic(err)
	}
}
