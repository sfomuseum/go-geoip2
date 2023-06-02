package main

import (
	"encoding/json"
	"flag"
	"io"
	"log"
	"net"
	"os"

	"github.com/oschwald/geoip2-golang"
	"github.com/sfomuseum/go-geoip2/data"
)

func main() {

	as_json := flag.Bool("json", false, "Output data as a JSON array.")
	flag.Parse()

	db, err := geoip2.FromBytes(data.GeoLite2City)

	if err != nil {
		log.Fatal(err)
	}

	writers := []io.Writer{
		os.Stdout,
	}

	wr := io.MultiWriter(writers...)

	if *as_json {
		wr.Write([]byte(`[`))
	}
	
	for i, addr := range flag.Args() {

		ip := net.ParseIP(addr)

		record, err := db.City(ip)

		if err != nil {
			log.Fatalf("Failed to determine city for '%s', %v", addr, err)
		}
		
		if *as_json && i > 0 {
			wr.Write([]byte(`,`))
		}
		
		enc := json.NewEncoder(wr)
		err = enc.Encode(record)

		if err != nil {
			log.Fatalf("Failed to decode record for '%s', %v", addr, err)
		}
	}

	if *as_json {
		wr.Write([]byte(`]`))
	}
	
}
