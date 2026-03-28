package models

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"

	"github.com/Shreyankthehacker/savector/utils"
)

type Index struct{
	path string
	Name string 
	data []float32  // one flat vector all over 
	dim int 
	count int 
}


// raw vector storage id*


func (db* DataBase)CreateIndex(name string,dim int) (*Index , error){


if dim<0{
	return nil,fmt.Errorf("invalid dimension")
}

//  raw.db for raw vector , index.db for hnsw storage , map.db for mapping bw external and internal id
index_path := db.Name+"/"+name
info,err:=os.Stat(index_path)
if err==nil && info.IsDir(){
	utils.Logger.Fatal("Such vector already exist or there might be some issue in creating index")
	return nil,err
}

err = os.Mkdir(index_path,0755)
if err!=nil{
	utils.Logger.Print("Index creation failed")
	return  nil,err
}

utils.Logger.Println("Index created successfully")

err = os.WriteFile(index_path+"/"+INDEX_VECTOR_STORAGE_FILE,[]byte{},0644)
if err!=nil{
utils.Logger.Printf("failed to create %s  file",INDEX_VECTOR_STORAGE_FILE)
}


err = os.WriteFile(index_path+"/"+INDEX_VECTOR_FILE,[]byte{},0644)
if err!=nil{
utils.Logger.Printf("failed to create %s  file",INDEX_VECTOR_FILE)
}

err = os.WriteFile(index_path+"/"+INDEX_METADATA_FILE,[]byte{},0644)
if err!=nil{
utils.Logger.Printf("failed to create %s  file",INDEX_METADATA_FILE)
}


utils.Logger.Println("Created all the necessaryy files")



index:= &Index{
	Name : name ,
	data :[]float32{},
	dim: dim,
	count :0,
	path : index_path,

}


// TODO In index metadata file for first few bytes we'll keep index related data , then we will keep map 


db.Indexes = append(db.Indexes, index.Name)

err = db.AddDatabaseDetailstometa()



if err!=nil{
	utils.Logger.Fatal("Error writing to the db")
}

utils.Logger.Printf("Index %s   created successfully",name)

return index,nil

}




func (index *Index)InsertVector(external_id string , data interface{})error{

// flatten if not flat 
// add raw to the raw file 
flatten_vector,err := utils.Flatten(data)
if err!=nil{
	fmt.Print("Cant flatten the vector pplease make it compatible")
	return err
}

if len(flatten_vector)!=index.dim{
	utils.Logger.Fatal("Uncompatible vector size please match the correct dimension......")
	return err
}

// writing to the raw file 
buffered_vector,err:=utils.FloatToByte(flatten_vector)
if err!=nil{
	utils.Logger.Fatalln("Vector creation failed")
	return err
}
file,err:=os.OpenFile(index.path+"/"+INDEX_VECTOR_FILE,os.O_WRONLY,0644)

if err!=nil{
	utils.Logger.Fatal("Couldnt open orr find raw file")
	return  err
}
defer file.Close()

offset:= index.count*index.dim
n,err := file.WriteAt(buffered_vector,int64(offset))
if err!=nil{
	utils.Logger.Fatal("Could not write to the raw vector ");
	return err
}else{
	fmt.Print("Wrote ",n," bytes of data ")
}
index.count++;

/// after raw is written we increment ctr by 1 bcz to access then external id is mapped with the count 
 


// for map we'll use [length(external_id)][external_id] format keep it seq coz we have internal_id 0,1,2,3.. no need to waste space for that 


meta_file,err:=os.OpenFile(index.path+"/"+INDEX_METADATA_FILE,os.O_APPEND|os.O_WRONLY,0644)

if err!=nil{
	utils.Logger.Fatal("Couldnt open orr find map file")
	return  err
}

length_eid := byte(len(external_id))
eid_byte  := []byte(external_id)
buffer:= new(bytes.Buffer)
binary.Write(buffer,binary.LittleEndian,length_eid)
binary.Write(buffer,binary.LittleEndian,eid_byte)


// map done and raw vectors can be inserted
meta_file.Write(buffer.Bytes())






defer meta_file.Close()

return nil



}