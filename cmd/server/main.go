package main

import (
	"bfm-example/internal/adapter/handler"
	"bfm-example/internal/adapter/transport"
	"bfm-example/internal/config"
	"bfm-example/internal/domain/repository"
	"bfm-example/internal/domain/usecase"
	"bfm-example/internal/infrastructure/httpclient"
	"bfm-example/pkg/headerforward"
	"bfm-example/pkg/validate"
	"context"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"go.uber.org/fx"
)

// --- Modules ---

var moduleInfra = fx.Module("infra",
	fx.Provide(
		config.Load,
		headerforward.NewClient,
	),
)

var moduleOrder = fx.Module("order",
	fx.Provide(
		// Repositories
		fx.Annotate(
			httpclient.NewUserRepo,
			fx.As(new(repository.UserRepository)),
		),
		fx.Annotate(
			httpclient.NewOrderRepo,
			fx.As(new(repository.OrderRepository)),
		),
		// Validator
		validate.New,
		// Use cases
		usecase.NewProcessOrder,
		// Handlers
		handler.NewProcessOrderHandler,
	),
)

var moduleMCP = fx.Module("mcp",
	fx.Provide(transport.RegisterTools),
)

// --- Run server lifecycle ---

func runServer(lc fx.Lifecycle, cfg config.Config, server *mcp.Server) {
	srv := &http.Server{
		Addr: cfg.Port,
		Handler: mcp.NewStreamableHTTPHandler(func(r *http.Request) *mcp.Server {
			return server
		}, nil),
	}

	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			ln, err := net.Listen("tcp", srv.Addr)
			if err != nil {
				return fmt.Errorf("listen on %s: %w", srv.Addr, err)
			}
			go func() {
				log.Printf("Starting BFF Server on %s", ln.Addr().String())
				if err := srv.Serve(ln); err != nil && err != http.ErrServerClosed {
					log.Printf("BFF server exited: %v", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Println("Shutting down server...")
			return srv.Shutdown(ctx)
		},
	})
}

// --- Main Function ---

func main() {
	app := fx.New(
		moduleInfra,
		moduleOrder,
		moduleMCP,
		fx.Invoke(runServer),
	)
	app.Run()
}
