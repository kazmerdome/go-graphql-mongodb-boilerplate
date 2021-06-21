package utility

import (
	"fmt"
	"log"
)

func ShowLogo() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	fmt.Println()
	fmt.Println()
	fmt.Println("\033[0;97m kazmerdome - go-graphql-mongodb-boilerplate ")
	fmt.Println("\033[0;34m----------------------------------------------")
}

func LoggerFormat() string {
	return "\033[0;34mlatency: \033[0;97m${latency_human} " +
		"\033[0;34mremote_ip: \033[0;97m${remote_ip} " +
		"\033[0;34mmethod: \033[0;97m${method} " +
		"\033[0;34mstatus: \033[0;97m${status} " +
		"\033[0;34mdate: \033[0;97m${time_rfc3339} " +
		"\033[0;34mpath: \033[0;97m${path} " +
		"\033[0;34mquery: \033[0;97m${query} " +
		"\033[0;34muri: \033[0;97m${uri} " +
		"\n"
}
