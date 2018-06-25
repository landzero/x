package main

import (
	"flag"
	"landzero.net/x/log"
	"net"
	"os"
	"path/filepath"

	"landzero.net/x/os/minit"
)

var sock string

func main() {
	// parse flags
	flag.StringVar(&sock, "L", "/var/run/minit/minit.sock", "socket file to listen")
	flag.Parse()
	// try remove existing sock file
	os.Remove(sock)
	// try create parrent directory
	os.MkdirAll(filepath.Dir(sock), os.FileMode(0755))
	// listen sock file
	var err error
	var l net.Listener
	if l, err = net.Listen("unix", sock); err != nil {
		log.Println("Failed to listen", sock, err)
		return
	}
	log.Println("Listening on", sock)
	// the listen loop
	minit.Serve(l)
}

func printHelp() {
	println("Minit")
	println("Usage:")
	flag.PrintDefaults()
}
