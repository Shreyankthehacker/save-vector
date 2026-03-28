package main

import (
	"github.com/Shreyankthehacker/savector/models"
	_ "github.com/Shreyankthehacker/savector/utils"

	_ "fmt"

	_ "github.com/Shreyankthehacker/savector/models"
)

func main(){


db,_:= models.CreateDatabase("Shreyank vector")
db.FetchIndex("Third Index")
}