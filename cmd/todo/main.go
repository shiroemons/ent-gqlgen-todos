package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"entgo.io/contrib/entgql"
	"entgo.io/ent/dialect"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	_ "github.com/mattn/go-sqlite3"

	"github.com/shiroemons/ent-gqlgen-todos"
	"github.com/shiroemons/ent-gqlgen-todos/ent"
	"github.com/shiroemons/ent-gqlgen-todos/ent/migrate"
)

const defaultPort = "8080"

func graphqlHandler(ctx context.Context) gin.HandlerFunc {
	// Create ent.Client and run the schema migration.
	client, err := ent.Open(dialect.SQLite, "file:ent?mode=memory&cache=shared&_fk=1")
	if err != nil {
		log.Fatal("opening ent client", err)
	}

	if err = client.Schema.Create(
		ctx,
		migrate.WithGlobalUniqueID(true),
		migrate.WithDropIndex(true),
		migrate.WithDropColumn(true),
	); err != nil {
		log.Fatal("opening ent client", err)
	}

	h := handler.NewDefaultServer(todo.NewSchema(client))
	h.Use(entgql.Transactioner{TxOpener: client})

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

// Defining the Playground handler
func playgroundHandler() gin.HandlerFunc {
	h := playground.Handler("GraphQL playground", "/query")

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func GinContextToContextMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.WithValue(c.Request.Context(), "GinContextKey", c)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func main() {
	ctx := context.Background()

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	router := gin.Default()

	router.Use(GinContextToContextMiddleware())
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:" + port},
		AllowCredentials: true,
	}))

	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})
	router.POST("/query", graphqlHandler(ctx))
	router.GET("/", playgroundHandler())

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := srv.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			log.Printf("listen: %s\n", err)
		}
	}()

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be caught, so don't need to add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}
