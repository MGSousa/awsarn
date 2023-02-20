package main

import (
	"flag"
	"net/http"
	"strings"

	"github.com/MGSousa/awsarn/awsarn"
	log "github.com/sirupsen/logrus"
)

var (
	arn string
)

// command that implements the main executable.
func main() {
	flag.StringVar(&arn, "arn", "", "ARN string to validate")
	flag.Parse()

	if strings.Trim(arn, " ") == "" {
		log.Fatalln("ARN not specified!")
	}

	if ok := awsarn.Valid(arn, http.DefaultClient); ok {
		log.Infoln("ARN is valid")
		return
	}
	log.Fatalln("ARN is not valid")
}
