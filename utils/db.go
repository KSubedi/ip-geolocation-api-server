package utils

import (
	"fmt"

	"github.com/oschwald/geoip2-golang"
)

var DB *geoip2.Reader

func NewDB() {
	database, err := geoip2.Open("GeoLite2-City.mmdb")
	if err != nil {
		fmt.Print(err)
	}

	DB = database
}
