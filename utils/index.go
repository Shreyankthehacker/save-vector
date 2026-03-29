package utils

import (
	"fmt"
	"os"
	"github.com/Shreyankthehacker/savector/config"
	"encoding/binary"
)

func GetIndexBuffer(index_path string)([]byte, error){

	

fmt.Print("reached get index bufferr",index_path)
_, err := os.Stat(index_path + "/" + config.INDEX_METADATA_FILE)
if os.IsNotExist(err) {
    fmt.Println("File does not exist")
	return nil,err
}
file,err:=os.OpenFile(index_path+"/"+config.INDEX_METADATA_FILE,os.O_RDONLY,0644)
if err != nil {
    Logger.Println("Couldnt open the index metadata file")
	return nil,err
}
defer file.Close()

buffer:= make([]byte,4)
_,err=file.Read(buffer)
if err != nil {
    Logger.Println("Problem reading the metadata length")
	return nil,err
}
val := int(binary.LittleEndian.Uint32(buffer))

index_buffer:=make([]byte,val)
_,err=file.Read(index_buffer)
if err != nil {
    Logger.Println("Problem reading the metadata value")
	return nil,err
}
fmt.Print("left get idx")
return index_buffer,nil

}



