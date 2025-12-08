package main

import (
	"cometosee/controller"
	"cometosee/firebase"
	"cometosee/intailizer"
	"cometosee/middleware"
	"context"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type Api struct {
	Addr string
}

func init() {
	intailizer.Loadenv()
	intailizer.DatabaseConnection()
	intailizer.Syncdatabase()
	firebase.Initailize()
}

func main() {
	defer intailizer.DB.Close()

	Port := ":" + os.Getenv("PORT")

	handler := &Api{Addr: Port}

	gor := mux.NewRouter()

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})
	//recovery->ratelimit->cors
	recoverhandler := middleware.PanicrecoveryMiddleware(gor)
	ratelimithandler := middleware.RateLimitMiddleware(recoverhandler)

	corshandler := c.Handler(ratelimithandler)

	srv := &http.Server{
		Addr:    handler.Addr,
		Handler: corshandler,
	}
	manager := controller.NewManger(context.Background())
	//routing
	gor.HandleFunc("/signup", controller.Signup).Methods("POST")
	gor.HandleFunc("/login", controller.Login).Methods("POST")
	gor.HandleFunc("/post", controller.UploadPost).Methods("POST")
	gor.HandleFunc("/post/like", controller.LikePost).Methods("POST")
	gor.HandleFunc("/post/comment", controller.CommentPost).Methods("POST")
	gor.HandleFunc("/ws/manager", manager.ServeWS).Methods("GET")

	if err := srv.ListenAndServe(); err != nil {
		log.Fatal("error in server")
	}
}
