package utils

import (
	"log"
	"os"
)


var Logger = log.New(os.Stdout,"savector",log.LstdFlags)