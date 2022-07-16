package web

import (
    "context"
    "github.com/IfanTsai/go-lib/user/token"
    "log"
    "net/http"
    "{{ .GoModulePath }}/internal/db"
    "{{ .GoModulePath }}/web/handlers"

    "github.com/IfanTsai/go-lib/config"
    "github.com/IfanTsai/go-lib/gin/middlewares"
    "github.com/IfanTsai/go-lib/logger"
    "github.com/gin-gonic/gin"
    "github.com/pkg/errors"
)

// GinServer serves HTTP requests for web service.
type GinServer struct {
    server        *http.Server
    address       string
    tokenMaker    token.Maker
    accessLogPath string
    panicLogPath  string
    dbLogPath     string
}

// NewGinServer creates a new HTTP gin server and setup routing.
func NewGinServer(address string) (*GinServer, error) {
    tokenMaker, err := token.NewJWTMaker(config.GetTokenSymmetricKey())
    if err != nil {
        return nil, errors.Wrap(err, "cannot create token maker")
    }

    server := &GinServer{
        address:       address,
        tokenMaker:    tokenMaker,
        accessLogPath: config.GetString("log.access_path"),
        panicLogPath:  config.GetString("log.panic_path"),
        dbLogPath:     config.GetString("log.db_path"),
    }

    server.setupRouter()

    return server, nil
}

func (s *GinServer) setupRouter() {
    version := config.GetVersion()

    if err := db.Init(s.dbLogPath); err != nil {
        log.Fatal(err)
    }

    if !config.IsDebugMode() {
        gin.SetMode(gin.ReleaseMode)
    }

    router := gin.New()

    s.server = &http.Server{
        Addr:    s.address,
        Handler: router,
    }

    accessLogger := logger.NewJSONLogger(logger.WithFileRotationP(s.accessLogPath))
    panicLogger := logger.NewJSONLogger(logger.WithFileRotationP(s.panicLogPath))

    router.Use(
        middlewares.Recovery(version, panicLogger, true),
        middlewares.SetTokenMaker(s.tokenMaker),
        middlewares.Logger(accessLogger),
        middlewares.Jsonifier(version),
    )

    middlewares.NewPrometheus("", "{{ .ProjectName }}").Use(router)

    v1API := router.Group("/apis/{{ .ProjectName }}/v1")
    v1API.POST("/user/create", handlers.CreateUser)
    v1API.POST("/user/login", handlers.LoginUser)

    authAPI := v1API.Use(middlewares.Authorization(version, s.tokenMaker))
    authAPI.GET("/user", handlers.GetUser)
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
