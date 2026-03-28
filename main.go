package main

import (
	"github.com/Shreyankthehacker/savector/models"
	_ "github.com/Shreyankthehacker/savector/utils"

	_ "fmt"

	_ "github.com/Shreyankthehacker/savector/models"
)

func main(){


db,_:= models.CreateDatabase("Shreyank vector")
idx,_:= db.CreateIndex("First Index",10)
idx.InsertVector("firstvectorbitch",[]float32{1,2,3,4,5,6,7,8,9,0});

}