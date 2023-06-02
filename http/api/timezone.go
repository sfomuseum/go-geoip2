package api

import (
	"net"
	"net/http"

	"github.com/aaronland/go-http-sanitize"
	"github.com/oschwald/geoip2-golang"
)

// TimeZoneHandler returns a `http.Handler` instance that will return the timezone for a given IP address.
// Addresses are inferred from a "?address={IPADDRESS}" parameter or from the IP address of the requestor
// in that order.
func TimeZoneHandler(db *geoip2.Reader) http.Handler {

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		addr := req.RemoteAddr

		q_addr, err := sanitize.GetString(req, "address")

		if err != nil {
			http.Error(rsp, "Invalid address parameter", http.StatusBadRequest)
			return
		}

		if q_addr != "" {
			addr = q_addr
		}

		ip := net.ParseIP(addr)

		record, err := db.City(ip)

		if err != nil {
			http.Error(rsp, err.Error(), http.StatusInternalServerError)
			return
		}

		rsp.Header().Set("Content-type", "text/plain")
		rsp.Write([]byte(record.Location.TimeZone))
		return
	}

	return http.HandlerFunc(fn)
}
