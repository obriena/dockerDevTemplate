package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	chiprometheus "github.com/766b/chi-prometheus"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/client_golang/prometheus/push"
)

var (
	processingTime = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "task_event_process_duration",
		Help: "Time it took to complete a taks",
	})
	processedCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "task_event_processing_total",
		Help: "How many tasks have been processed",
	},
		[]string{"task"}).WithLabelValues("task")
)

func main() {
	log.Println("Starting services")

	processingTime = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "task_event_process_duration",
		Help: "Time it took to complete a taks",
	})
	processedCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "task_event_processing_total",
		Help: "How many tasks have been processed",
	},
		[]string{"task"}).WithLabelValues("task")

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
		processingTime.Set(float64(elapsed.Milliseconds()))
		processedCounter.Add(1)
		log.Println("Elapsed Time: ")
		log.Println(elapsed.Milliseconds())
		pushProcessingDuration(processingTime)
		pushProcessingCount(processedCounter)
		RespondJSON(w, r, message)
	})

	// RESTy routes for "articles" resource
	r.Route("/articles", func(r chi.Router) {
		// r.With(paginate).Get("/", listArticles)                           // GET /articles
		// r.With(paginate).Get("/{month}-{day}-{year}", listArticlesByDate) // GET /articles/01-16-2017

		// r.Post("/", createArticle)                                        // POST /articles
		// r.Get("/search", searchArticles)                                  // GET /articles/search

		// Regexp url parameters:
		// r.Get("/{articleSlug:[a-z-]+}", getArticleBySlug)                // GET /articles/home-is-toronto

		// Subrouters:
		r.Route("/{articleID}", func(r chi.Router) {

			// r.Use(ArticleCtx)
			r.Get("/", getArticle) // GET /articles/123
			// r.Put("/", updateArticle)                                       // PUT /articles/123
			// r.Delete("/", deleteArticle)                                    // DELETE /articles/123
		})
	})

	http.ListenAndServe(":80", r)

	log.Println("Service ending")
}

func getArticle(w http.ResponseWriter, r *http.Request) {
	//ctx := r.Context()
	article := Article{ID: "0", UserID: 0, Title: "Best Article", Slug: "uniqueId"}
	// if !ok {
	// 	http.Error(w, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
	// 	return
	// }
	//w.Write([]byte(fmt.Sprintf("title:%s", article.Title)))
	RespondJSON(w, r, article)
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

type Article struct {
	ID     string `json:"id"`
	UserID int64  `json:"user_id"` // the author
	Title  string `json:"title"`
	Slug   string `json:"slug"`
}
