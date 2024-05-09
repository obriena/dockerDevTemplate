package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	chiprometheus "github.com/766b/chi-prometheus"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/obriena/dockerdevtemplate/domain"
	"github.com/obriena/dockerdevtemplate/infra"
	"github.com/obriena/dockerdevtemplate/post"
	postrepo "github.com/obriena/dockerdevtemplate/post/repository/mysqlrepo"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/client_golang/prometheus/push"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var Db *gorm.DB

var (
	queryAllTime = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "task_event_query_all_duration",
		Help: "Time it took to complete a taks",
	})
	queryAllCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "task_event_query_all_total",
		Help: "How many tasks have been processed",
	},
		[]string{"task"}).WithLabelValues("task")
)

func main() {
	log.Println("Starting services")
	context := context.Background()

	service := newService()
	if err := service.Init(context); err != nil {
		log.Println("Service did not initialize.", err)
		return
	}
	r := chi.NewRouter()

	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(chiprometheus.NewMiddleware("service_name"))

	r.Use(middleware.Timeout(60 * time.Second))

	r.Handle("/metrics", promhttp.Handler())

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		message := make(map[string]string)
		message["message"] = "Hi"
		elapsed := time.Since(start)

		log.Println("Elapsed Time: ")
		log.Println(elapsed.Milliseconds())
		RespondJSON(w, r, message)
	})

	r.Route("/posts", func(r chi.Router) {
		/*
			Chi has routes that can have pagination built in somethin like this
			r.With(paginate).Get("/path", ...)
		*/
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			posts, err := service.postInteractor.ReadAll(context)
			if err != nil {
				log.Println("Error retrieving posts: ", err)
			} else {
				RespondJSON(w, r, posts)
			}
			elapsed := time.Since(start)

			queryAllTime.Set(float64(elapsed.Milliseconds()))
			queryAllCounter.Add(1)
			pushProcessingDuration(queryAllTime)
			pushProcessingCount(queryAllCounter)
		})
	})

	http.ListenAndServe(":80", r)

	log.Println("Service ending")
}

type Service struct {
	db             *gorm.DB
	postInteractor domain.PostInteractor
}

func newService() *Service {
	return &Service{}
}

func (s *Service) Init(ctx context.Context) error {
	config := infra.GetConfig()
	db, err := connectDB(config)
	if err != nil {
		return err
	}
	s.db = db
	aRepo := postrepo.NewRepository(db)
	s.postInteractor = post.NewInteractor(aRepo)
	return nil
}

func pushProcessingDuration(processingTime prometheus.Gauge) {
	if err := push.New("http://host.docker.internal:9091", "task_event_process_duration").
		Collector(processingTime).
		Grouping("db", "event-service").
		Push(); err != nil {
		log.Println("Could not push completion time to Pushgateway:", err)
	}
}

func pushProcessingCount(processedCounter prometheus.Counter) {
	if err := push.New("http://host.docker.internal:9091", "task_event_processing_total").
		Collector(processedCounter).
		Grouping("db", "event-service").
		Push(); err != nil {
		fmt.Println("Could not push tasks processed to Pushgateway:", err)
	}
}
func RespondJSON(w http.ResponseWriter, r *http.Request, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	b, _ := json.Marshal(data)

	w.Write(b)
}

func connectDB(config *infra.Config) (*gorm.DB, error) {
	var err error
	dsn := config.DbUser + ":" + config.DbPassword + "@tcp" + "(" + config.DbHost + ":" + config.DbPort + ")/" + config.DbName + "?" + "parseTime=true&loc=Local"
	log.Println("@tcp(" + config.DbHost + ":" + config.DbPort + ")/" + config.DbName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Printf("Error connecting to database : error=%v", err)
		return nil, err
	}

	return db, nil
}
