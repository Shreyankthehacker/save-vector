package models

import (
	"fmt"
	"os"


	"github.com/Shreyankthehacker/savector/config"
	"github.com/Shreyankthehacker/savector/utils"
)



func (db* DataBase)CreateIndex(name string,dim int) (*Index , error){


if dim<0{
	return nil,fmt.Errorf("invalid dimension")
}

//  raw.db for raw vector , index.db for hnsw storage , map.db for mapping bw external and internal id
index_path := db.Name+"/"+name
info,err:=os.Stat(index_path)
if err==nil && info.IsDir(){
	utils.Logger.Printf("Such vector already exist or there might be some issue in creating index")
	return db.FetchIndex(name)
}



err = os.Mkdir(index_path,0755)
if err!=nil{
	utils.Logger.Print("Index creation failed")
	return  nil,err
}
utils.Logger.Println("Index created successfully")



err = os.WriteFile(index_path+"/"+config.INDEX_HNSW_FILE,[]byte{},0644)
if err!=nil{
utils.Logger.Printf("failed to create %s  file",config.INDEX_HNSW_FILE)
}


err = os.WriteFile(index_path+"/"+config.INDEX_VECTOR_FILE,[]byte{},0644)
if err!=nil{
utils.Logger.Printf("failed to create %s  file",config.INDEX_VECTOR_FILE)
}

index:= &Index{
	Name : name ,
	data :[]float32{},
	dim: dim,
	count :0,
	path : index_path,

}
index_bytes := index.getMetaDataBytes()
err = os.WriteFile(index_path+"/"+config.INDEX_METADATA_FILE,index_bytes,0644)
if err!=nil{
utils.Logger.Printf("failed to create %s  file",config.INDEX_METADATA_FILE)
}

utils.Logger.Println("Created all the necessaryy files")




db.Indexes = append(db.Indexes, index.Name)
err = db.AddDatabaseDetailstometa()



if err!=nil{
	utils.Logger.Fatal("Error writing to the db")
}

utils.Logger.Printf("Index %s   created successfully",name)

return index,nil

}



// Index.db is basically a graph of numbers 1,2,3,4,5.......