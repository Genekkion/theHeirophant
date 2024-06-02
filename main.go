package main

import "flag"

func main() {
	portFlag := flag.String("port", "30303", "port for server to run on")

	flag.Parse()

	server := initServer(*portFlag)
	server.run()
}
