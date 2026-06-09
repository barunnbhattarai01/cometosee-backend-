package main

import (
	"cometosee/config"
	"cometosee/controller"
	"cometosee/di"
	"cometosee/firebase"
	"cometosee/intailizer"
	"cometosee/middleware"
	"cometosee/routes"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/cors"
)

type Api struct {
	Addr string
}

func init() {
	intailizer.Loadenv()
	config.InitCache()
	intailizer.DatabaseConnection()
	intailizer.Syncdatabase()
	firebase.Initailize()
}

func main() {
	defer intailizer.DB.Close()
	config.SetupGarbageCollector()

	go func() {
		//barrrrun for localhost it ok but for production it nottt ok to expose pprof
		log.Printf("pprof running in :6060")
		http.ListenAndServe("localhost:6060", nil)
	}()

	Port := ":" + os.Getenv("PORT")

	handler := &Api{Addr: Port}

	gor := mux.NewRouter()

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})
	//recovery->ratelimit->cors
	recoverhandler := middleware.PanicrecoveryMiddleware(gor)
	ratelimithandler := middleware.RateLimitMiddleware(recoverhandler)

	corshandler := c.Handler(ratelimithandler)

	srv := &http.Server{
		Addr:              handler.Addr,
		Handler:           corshandler,
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       120 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
	}

	//di
	authcontroller := di.SetupAuthcontroller()
	postcontroller := di.SetupPostController()
	wscontroller := di.SetupMessage()
	userfilterctrl := di.UserfilterDI()
	subcriptiondi := di.SubcribntionDi()
	connectiondi := di.SetupConnectionController()
	agoradi, _ := di.SetupAgoraDI()
	userdetailinfodi := di.InitializeUserDetailInfo()
	discoverydi := di.Setupdiscoveryfeed()
	//routing
	routes.AuthRoutes(gor, authcontroller)
	routes.PostRoutes(gor, postcontroller)
	routes.WebscoketRoutes(gor, wscontroller)
	routes.SetupUserfilterroutes(gor, userfilterctrl)
	routes.ConnectionRoutes(gor, connectiondi)
	routes.SubcribtionRoute(gor, subcriptiondi)
	routes.SetupAgoraRoutes(gor, agoradi)
	routes.InitializeUserDetailInfoRoutes(gor, userdetailinfodi)
	routes.Discoveryroutes(gor, discoverydi)
	gor.HandleFunc("/esewa/initiate", controller.InitiateHandler).Methods("POST")
	gor.HandleFunc("/esewa/verify", controller.VerifyHandler).Methods("GET")
	gor.HandleFunc("/esewa/failure", controller.FailureHandler).Methods("GET")
	//promethheus mertics
	//had to put it on differenr server for standard code
	go func() {
		log.Printf("promethus mertics running on:2112")
		http.Handle("/metrics", promhttp.Handler())
		log.Fatal(http.ListenAndServe(":2112", nil))
	}()

	if err := srv.ListenAndServe(); err != nil {
		log.Fatal("error in server")
	}
}
