package models

import (
	"os"
	"time"

	"github.com/Shreyankthehacker/savector/proto_models"
	"github.com/Shreyankthehacker/savector/utils"
	"google.golang.org/protobuf/proto"
	_ "google.golang.org/protobuf/proto"
	
)





type DataBase struct{
Name string 
Created_at string
Updated_at string 
Indexes []string
}





func CreateDatabase(name string)(*DataBase,error){

info,err:= os.Stat(name)

if err==nil && info.IsDir(){
	utils.Logger.Fatalln("Such database already exist please create a new one")
	return nil,nil
}


err= os.Mkdir(name,0755)
if err!=nil{
utils.Logger.Println("Database creation failed ->",err)
return nil,err
}
utils.Logger.Println("Database folder created successfully")

err=os.WriteFile(name+"/"+DATABASE_METAINFO_FILE,[]byte{},0644)

if err!=nil{
	utils.Logger.Println("failed to create meta info file",err)
	return nil,err
}
utils.Logger.Println("Meta info file created as well")


db:= &DataBase{
	Name: name,
	Created_at: time.Now().UTC().Format(time.RFC3339),
Updated_at: time.Now().UTC().Format(time.RFC3339),
	Indexes: nil,
}

db.AddDatabaseDetailstometa()
utils.Logger.Printf("Database %s   created successfully",name)

return db,nil

}


func(d *DataBase) AddDatabaseDetailstometa()error{



databaseModel:= &proto_models.DataBase{
	Dbname: d.Name,
	CreatedAt: d.Created_at,
	UpdatedAt: d.Updated_at,
	Indexes: d.Indexes,
}

data,err:= proto.Marshal(databaseModel)

if err!=nil{
	utils.Logger.Fatal("Got a issue while marshalling the database meta file")
	return err
}

meta_file_path:=d.Name+"/"+DATABASE_METAINFO_FILE
err = os.WriteFile(meta_file_path,data,0644)

if err!=nil{
	utils.Logger.Fatalln("Couldnt write into the meta.db")
	return  err
}
utils.Logger.Println("Successfully written to the meta file")


return nil


}


func FetchDatabase(Name string)(*DataBase,error){

meta_file_path:=Name+"/"+DATABASE_METAINFO_FILE
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





