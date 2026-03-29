package main

import (
	"github.com/Shreyankthehacker/savector/models"
	_ "github.com/Shreyankthehacker/savector/utils"

	_ "fmt"

	_ "github.com/Shreyankthehacker/savector/models"
)

func main(){


db,_:= models.CreateDatabase("Shreyank vector")
db.CreateIndex("Third Index",20)
idx,_:=db.FetchIndex("Third Index")
idx.InsertVector("random vector2",[]float32{1,2,3,4,5,6,7,8,9,0})


}