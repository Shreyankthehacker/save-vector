package main

import (
	"fmt"

	"github.com/Shreyankthehacker/savector/models"
	_ "github.com/Shreyankthehacker/savector/utils"

	_ "fmt"

	_ "github.com/Shreyankthehacker/savector/models"
)

func main(){


db,_:= models.CreateDatabase("Shreyankdb2")
db.CreateIndex("Third Index",10)
idx,_:=db.FetchIndex("Third Index")
idx.InsertVector("vec1", []float32{0.12, 0.87, 0.33, 0.45, 0.91, 0.02, 0.76, 0.54, 0.29, 0.68})
idx.InsertVector("vec2", []float32{0.55, 0.14, 0.92, 0.31, 0.73, 0.88, 0.09, 0.47, 0.66, 0.25})
idx.InsertVector("vec3", []float32{0.99, 0.21, 0.44, 0.78, 0.11, 0.63, 0.35, 0.82, 0.57, 0.04})
idx.InsertVector("vec4", []float32{0.18, 0.69, 0.27, 0.53, 0.84, 0.39, 0.71, 0.06, 0.95, 0.48})
fmt.Print(idx.FetchVectorByInternalId(2))
fmt.Print(idx.FetchInternalIdByExternalId("vec1"))
fmt.Print(idx.FetchInternalIdByExternalId("vec3"))


}