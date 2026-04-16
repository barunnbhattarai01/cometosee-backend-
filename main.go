package main

import (
	"cometosee/config"
	"cometosee/di"
	"cometosee/firebase"
	"cometosee/intailizer"
	"cometosee/middleware"
	"cometosee/routes"
	"encoding/json"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/mukezhz/pay-np/esewa"
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
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
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
	//routing
	routes.AuthRoutes(gor, authcontroller)
	routes.PostRoutes(gor, postcontroller)
	routes.WebscoketRoutes(gor, wscontroller)
	routes.SetupUserfilterroutes(gor, userfilterctrl)

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

func initiateHandler(w http.ResponseWriter, r *http.Request) {
	payload := &esewa.EsewaPayload{
		Amount:           "100",
		TaxAmount:        "0",
		ProductCode:      "EPAYTEST", // eSewa sandbox merchant code
		TotalAmount:      "100",
		TransactionUUID:  "txn-001",
		SuccessURL:       "http://localhost:8080/esewa/verify",
		FailureURL:       "http://localhost:8080/esewa/failure",
		SignedFieldNames: "total_amount,transaction_uuid,product_code",
	}

	client, err := esewa.New("8gBm/:&EnhH.1/q", payload) // sandbox secret
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	sig, err := client.GenerateSignature()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	payload.Signature = sig
	json.NewEncoder(w).Encode(payload)
}

func verifyHandler(w http.ResponseWriter, r *http.Request) {
	data := r.URL.Query().Get("data")
	if data == "" {
		http.Error(w, "missing data param", 400)
		return
	}

	client, _ := esewa.New("8gBm/:&EnhH.1/q", &esewa.EsewaPayload{})
	err := client.VerifySignature(data)
	if err != nil {
		http.Error(w, "invalid signature: "+err.Error(), 400)
		return
	}

	w.Write([]byte("payment verified ok"))
}
