package models


import (
	"os"
"time"
	"github.com/Shreyankthehacker/savector/config"
	"github.com/Shreyankthehacker/savector/proto_models"
	"github.com/Shreyankthehacker/savector/utils"
	"google.golang.org/protobuf/proto"
)



func CreateDatabase(name string)(*DataBase,error){

info,err:= os.Stat(name)

if err==nil && info.IsDir(){
	utils.Logger.Println("Such database already exist please create a new one")
	return FetchDatabase(name)
}


err= os.Mkdir(name,0755)
if err!=nil{
utils.Logger.Println("Database creation failed ->",err)
return nil,err
}
utils.Logger.Println("Database folder created successfully")

err=os.WriteFile(name+"/"+config.DATABASE_METAINFO_FILE,[]byte{},0644)

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

meta_file_path:=d.Name+"/"+config.DATABASE_METAINFO_FILE
err = os.WriteFile(meta_file_path,data,0644)

if err!=nil{
	utils.Logger.Fatalln("Couldnt write into the meta.db")
	return  err
}
utils.Logger.Println("Successfully written to the meta file")


return nil


}