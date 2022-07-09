package web

import (
    "context"
    "log"
    "net/http"
    "{{ .GoModulePath }}/web/handlers"

    "github.com/IfanTsai/go-lib/gin/middlewares"
    "github.com/IfanTsai/go-lib/logger"
    "github.com/gin-gonic/gin"
    "github.com/pkg/errors"
)

// GinServer serves HTTP requests for web service.
type GinServer struct {
    server  *http.Server
    address string
}

// NewGinServer creates a new HTTP gin server and setup routing.
func NewGinServer(address string) (*GinServer, error) {
    server := &GinServer{
        address: address,
    }

    server.setupRouter()

    return server, nil
}

func (s *GinServer) setupRouter() {
    gin.SetMode(gin.ReleaseMode)

    version := "1.0.0"

    jsonLogger, err := logger.NewJSONLogger(
        logger.WithDisableConsole(),
        logger.WithFileRotationP("./logs/{{ .ProjectName }}.log"),
    )
    if err != nil {
        log.Fatalln("cannot new json logger, err:", err)
    }

    router := gin.New()

    s.server = &http.Server{
        Addr:    s.address,
        Handler: router,
    }

    router.Use(
        middlewares.Recovery(version, jsonLogger, true),
        middlewares.Logger(jsonLogger),
        middlewares.Jsonifier(version),
    )

    middlewares.NewPrometheus("{{ .ProjectName }}", "{{ .ModuleName }}").Use(router)

    v1API := router.Group("/apis/{{ .ProjectName }}/v1")
    v1API.GET("/user", handlers.GetUser)
}

// Start runs the HTTP server on a specific address.
func (s *GinServer) Start() error {
    log.Println("http server is listening on", s.address)

    if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
        return errors.Wrap(err, "failed to start Gin server")
    }

    return nil
}

// Stop stops the HTTP server.
func (s *GinServer) Stop(ctx context.Context) error {
    return errors.Wrap(s.server.Shutdown(ctx), "failed to shutdown http server")
}
