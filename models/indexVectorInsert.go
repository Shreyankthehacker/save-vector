package models
import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"

	"github.com/Shreyankthehacker/savector/config"
	"github.com/Shreyankthehacker/savector/utils"
)



func (index *Index)InsertVector(external_id string , data interface{})error{

// flatten if not flat 
// add raw to the raw file 
fmt.Print("Inserting vector")
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
file,err:=os.OpenFile(index.path+"/"+config.INDEX_VECTOR_FILE,os.O_WRONLY,0644)

if err!=nil{
	utils.Logger.Fatal("Couldnt open orr find raw file")
	return  err
}
defer file.Close()

offset:= index.count*index.dim*4  // 4 bcz each number takes 4 bytes
print("Writing after ",offset)
n,err := file.WriteAt(buffered_vector,int64(offset))
if err!=nil{
	utils.Logger.Fatal("Could not write to the raw vector ");
	return err
}else{
	fmt.Print("Wrote ",n," bytes of data ")
}
index.count++;
fmt.Println(index.count,"has been increased")



/// after raw is written we increment ctr by 1 bcz to access then external id is mapped with the count 
 
// for map we'll use [length(external_id)][external_id] format keep it seq coz we have internal_id 0,1,2,3.. no need to waste space for that 
meta_file_path := index.path+"/"+config.INDEX_METADATA_FILE

old_data,err:=os.ReadFile(meta_file_path)

if err!=nil{
	fmt.Println("Couldnt open the meta data file")
}


length_of_meta:=int32(binary.LittleEndian.Uint32(old_data[:4]))
modified_data:=index.getMetaDataBytes()
new_length_of_meta:=len(modified_data)

final_buffer:= new(bytes.Buffer)
binary.Write(final_buffer,binary.LittleEndian,int32(new_length_of_meta))
final_buffer.Write(modified_data)
final_buffer.Write(old_data[int32(4)+length_of_meta:])

meta_file,err:=os.OpenFile(index.path+"/"+config.INDEX_METADATA_FILE,os.O_RDWR,0644)



if err!=nil{
	utils.Logger.Fatal("Couldnt open orr find map file")
	return  err
}




// length_eid := byte(len(external_id))
length_eid := uint32(len(external_id))
eid_byte  := []byte(external_id)
binary.Write(final_buffer,binary.LittleEndian,length_eid)
binary.Write(final_buffer,binary.LittleEndian,eid_byte)


// map done and raw vectors can be inserted
meta_file.Write(final_buffer.Bytes())


defer meta_file.Close()

return nil



}

