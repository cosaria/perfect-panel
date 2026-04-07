package tool

import (
	"net"

	"github.com/oschwald/geoip2-golang"
	"github.com/perfect-panel/server/config"
)

type GeoIPCityReader interface {
	City(ipAddress net.IP) (*geoip2.City, error)
}

type Deps struct {
	Restart func() error
	Config  *config.Config
	GeoIPDB GeoIPCityReader
}
