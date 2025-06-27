package httpapp

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"thingify/internal/http/middleware"
	"thingify/internal/http/webhook"
	"thingify/internal/service/issue"

	"github.com/gofiber/fiber/v2"
)

const (
	srvPortDefault           = 8080
	srvReadTimeoutDefault    = 10 * time.Second
	srvWriteTimeoutDefault   = 10 * time.Second
	srvGatewayTimeoutDefault = 10 * time.Second
)

// App представляет HTTP приложение, которое управляет сервером и его состоянием.
type App struct {
	log *slog.Logger

	issueSvc *issue.Service

	webhookSecret string

	port           int
	readTimeout    time.Duration
	writeTimeout   time.Duration
	requestTimeout time.Duration

	mu   sync.Mutex
	fapp *fiber.App
}

type Option func(*App)

func WithPort(port int) Option {
	return func(a *App) {
		if port <= 0 || port > 65535 {
			a.log.Error("invalid port number, using default port", slog.Int("default_port", srvPortDefault))
			a.port = srvPortDefault
		} else {
			a.port = port
		}
	}
}

func WithReadTimeout(timeout time.Duration) Option {
	return func(a *App) {
		a.readTimeout = timeout
	}
}

func WithWriteTimeout(timeout time.Duration) Option {
	return func(a *App) {
		a.writeTimeout = timeout
	}
}

func WithRequestTimeout(timeout time.Duration) Option {
	return func(a *App) {
		a.requestTimeout = timeout
	}
}

// New создает новое HTTP приложение.
func New(
	log *slog.Logger,
	issueSvc *issue.Service,
	webhookSecret string,
	opts ...Option,
) *App {
	app := &App{
		log:            log,
		issueSvc:       issueSvc,
		webhookSecret:  webhookSecret,
		port:           srvPortDefault,
		readTimeout:    srvReadTimeoutDefault,
		writeTimeout:   srvWriteTimeoutDefault,
		requestTimeout: srvGatewayTimeoutDefault,
	}

	for _, opt := range opts {
		opt(app)
	}

	return app
}

// MustRun запускает HTTP сервер и вызывает панику в случае ошибки.
func (a *App) MustRun(ctx context.Context) {
	if err := a.Run(ctx); err != nil {
		panic("failed to run HTTP http_server: " + err.Error())
	}
}

// Run запускает HTTP сервер.
func (a *App) Run(ctx context.Context) error {
	const op = "httpapp.Run"

	log := a.log.With("op", op)

	log.InfoContext(ctx, "starting HTTP http_server")

	cfg := fiber.Config{
		ReadTimeout:   a.readTimeout,
		WriteTimeout:  a.writeTimeout,
		CaseSensitive: false,
	}

	fapp := fiber.New(cfg)

	// DISABLED FOR DEV
	//fapp.Use(recover.New())
	fapp.Use(middleware.Logging(a.log))

	webhookHandler := webhook.NewHandler(a.issueSvc, a.webhookSecret)

	baseRouter := fapp.Group("")

	webhookHandler.RegisterRoutes(baseRouter)

	fapp.Use(func(c *fiber.Ctx) error {
		return fiber.ErrBadRequest
	})

	a.mu.Lock()
	a.fapp = fapp
	a.mu.Unlock()

	addr := fmt.Sprintf("0.0.0.0:%d", a.port)

	return fapp.Listen(addr)
}

// Stop останавливает HTTP сервер.
// Нужно дожидаться завершения работы этого метода.
// Котнтекст должен быть с таймаутом, чтобы избежать
// зависания в случае проблем с остановкой сервера.
func (a *App) Stop(ctx context.Context) error {
	const op = "httpapp.Stop"

	log := a.log.With("op", op)

	log.Info("stopping HTTP http_server")

	a.mu.Lock()
	fapp := a.fapp
	a.mu.Unlock()

	if err := fapp.ShutdownWithContext(ctx); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
