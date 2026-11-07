package metrics

import (
	"fmt"

	echoprometheus "github.com/labstack/echo-prometheus"
	"github.com/labstack/echo/v5"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	ServiceStatus = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name:      "service_status",
		Namespace: "prompage",
		Help:      "Thet status of the service. 1 = up, 0 = down",
	}, []string{"name", "group"})
	ServiceUptime = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name:      "service_uptime",
		Namespace: "prompage",
		Help:      "The uptime percentage of the service",
	}, []string{"name", "group"})
)

type Server struct {
	e    *echo.Echo
	port int
}

func NewServer(port int) *Server {
	e := echo.New()

	e.GET("/metrics", echoprometheus.NewHandler())

	return &Server{
		e:    e,
		port: port,
	}
}

func (s *Server) Init() {
	prometheus.MustRegister(ServiceStatus)
	prometheus.MustRegister(ServiceUptime)
}

func (s *Server) Start() error {
	return s.e.Start(fmt.Sprintf(":%d", s.port))
}
