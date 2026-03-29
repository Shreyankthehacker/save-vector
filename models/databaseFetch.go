package models

import (
	"os"

	"github.com/Shreyankthehacker/savector/config"
	"github.com/Shreyankthehacker/savector/proto_models"
	"github.com/Shreyankthehacker/savector/utils"
	"google.golang.org/protobuf/proto"
)






func FetchDatabase(Name string)(*DataBase,error){

	print("fetch is called")

meta_file_path:=Name+"/"+config.DATABASE_METAINFO_FILE
file,err:= os.ReadFile(meta_file_path)

if err!=nil{

if os.IsNotExist(err){
	utils.Logger.Fatal("No such database as ",Name,"exists")
}else{
	utils.Logger.Fatal("Failed to read file: ", err)
}
}
var db proto_models.DataBase
err = proto.Unmarshal(file,&db)
if err!=nil{
	utils.Logger.Fatalln("Failed fetching db0")
}



return &DataBase{

Name: db.Dbname,
Created_at: db.CreatedAt,
Updated_at: db.UpdatedAt,
Indexes:db.Indexes ,
},nil
}





