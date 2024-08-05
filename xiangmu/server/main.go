package main

import "flag"

func main() {
	ip := flag.String("ip", "127.0.0.1", "specify server ip")
	port := flag.Int("port", 8888, "specify server prot")
	flag.Parse()

	server := NewServer(*ip, *port)
	server.Start()
}
