package models

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/Shreyankthehacker/savector/config"
	
)



func getIndexFromMetaDataByte(data  []byte)(*Index , error){
fmt.Print("metadata decode is called")
parts := bytes.Split(data, []byte(config.DELIMITER_INDEX_METADATA_FILE)) 
if (len(parts)!=5){
	print("....lol")
	return nil,fmt.Errorf("Some internal issue with decoding meta , maybe corrupted")
}
name:=string(parts[0])
path:=string(parts[1])
dim:=int32(binary.LittleEndian.Uint32(parts[2]))
count:=int32(binary.LittleEndian.Uint32(parts[3]))
fmt.Printf("Name: %s, Path: %s, Dim: %d, Count: %d\n", name, path, dim, count)
return &Index{
	path: path,
	Name: name,
	dim: int(dim),
	count: int(count),
	data: []float32{},
},nil
}




func (index *Index)getMetaDataBytes()[]byte{
	payload:=new(bytes.Buffer)
	payload.Write([]byte(index.Name))
	payload.Write([]byte(config.DELIMITER_INDEX_METADATA_FILE))
	payload.Write([]byte(index.path))
	payload.Write([]byte(config.DELIMITER_INDEX_METADATA_FILE))
	binary.Write(payload,binary.LittleEndian,int32(index.dim))
	payload.Write([]byte(config.DELIMITER_INDEX_METADATA_FILE))
	binary.Write(payload,binary.LittleEndian,int32(index.count))
	payload.Write([]byte(config.DELIMITER_INDEX_METADATA_FILE))
	
	final_length := len(payload.Bytes())
	final_buffer:=new(bytes.Buffer)
	binary.Write(final_buffer,binary.LittleEndian,int32(final_length))
	final_buffer.Write(payload.Bytes())
	return final_buffer.Bytes()
}




