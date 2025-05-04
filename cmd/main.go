package main

import(
	"fmt"
	"net/http"
	"time"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"habit21/internal/config"
	"habit21/internal/storage"
	"habit21/internal/handlers"
)

func main(){
	cfg, err := config.MustLoad()
	if err!=nil {
		//logger stuf
		fmt.Errorf("failed at loading config files", err.Error())
	}
	conn, err := storage.New(cfg)
	defer conn.DB.Close()
	if err!=nil {
		fmt.Println(err.Error())
	}
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
        AllowedOrigins: []string{"*"},
        AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
        AllowedHeaders: []string{"Content-Type"},
        AllowCredentials: false,
    }))
	trash := &handlers.Handler{Storage: conn}
	_ = trash
	r.Get("/", trash.Show)
	r.Get("/api/heatmap", trash.ShowTable)
	r.Put("/{id}", trash.Update)
	r.Post("/", trash.Create)
	r.Post("/api/login", trash.LoginUser)
	r.Post("/api/register", trash.AddUser)
	r.Delete("/{id}", trash.Delete)

	go func() {
		ticker := time.NewTicker(24*time.Hour) // negizi kate bar bul zherde
		var err error
		defer ticker.Stop()
		for {
			<-ticker.C
			err = conn.UpdateDefault()
			_ = err
		}
	}()
	http.ListenAndServe(fmt.Sprintf(":%s",cfg.Http.Port), r)
}