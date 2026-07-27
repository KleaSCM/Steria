// Author: KleaSCM
// Email: KleaSCM@gmail.com
// Name of the file: main.go
// Description: Entry point for the Steria multi-user file browser web server. Uses the internal/web package for all web logic.

package main

import "steria/internal/web"

func main() {
	web.StartServer(":8080")
}
