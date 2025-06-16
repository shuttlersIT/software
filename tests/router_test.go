// tests/router_test.go
package tests

import (
	"software_management/routes"
	"testing"
)

func TestRouterNotNil(t *testing.T) {
	router := routes.RegisterRoutes()
	if router == nil {
		t.Fatal("Router returned by RegisterRoutes() is nil")
	}
}
