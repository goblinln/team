package config

import (
	"fmt"
)

// AppInfo holds basic information for this server.
type AppInfo struct {
	Name      string
	Port      int
	LoginType LoginType
}

// Addr returns address for this service.
func (a *AppInfo) Addr() string {
	return fmt.Sprintf(":%d", a.Port)
}
