package metrics

import (
	"context"
	"fmt"

	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
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
	e.HideBanner = true

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

func (s *Server) Stop(ctx context.Context) error {
	return s.e.Shutdown(ctx)
}
