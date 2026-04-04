package models

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"

	"github.com/Shreyankthehacker/savector/config"
)


type Vector struct{
	data []float32
}






func (idx *Index)FetchVectorByInternalId(id int) (*Vector,error){
	raw_vector_file_path:=idx.path+"/"+config.INDEX_VECTOR_FILE
	file,err:= os.OpenFile(raw_vector_file_path,os.O_RDONLY,0646)
	if err!=nil{
		fmt.Println("Some issue in opening raw vector file")
		return nil,err
	}
	defer file.Close()
	// buffer := make([]byte, 4)

	// _, err = io.ReadFull(file, buffer)
	// if err != nil {
	// 	return nil, err
	// }

	//actually inside 
	
	// offset:=int64(binary.LittleEndian.Uint32(buffer))
	vsize:=int64(idx.dim*4)
	offset:=int64(id*int(vsize))
	vector_buffer:=make([]byte,vsize)
	file.Seek(offset-4,io.SeekCurrent)
	_, err = io.ReadFull(file, vector_buffer)
	if err != nil {
		return nil, err
	}


	data:=make([]float32,idx.dim)

	err=binary.Read(bytes.NewReader(vector_buffer),binary.LittleEndian,data)

	if err!=nil{
		fmt.Print("some issue in decoding the vector")
		return nil,err
	}
	return &Vector{
		data:data,
	},nil


}

func (idx *Index)FetchInternalIdByExternalId(external_id string)(int , error){

	map_vector_file_path:=idx.path+"/"+config.INDEX_METADATA_FILE
	file,err:= os.OpenFile(map_vector_file_path,os.O_RDONLY,0646)
	if err!=nil{
		fmt.Println("Some issue in opening raw vector file")
		return -1,err
	}
	defer file.Close()
	buffer := make([]byte, 4)

	_, err = io.ReadFull(file, buffer)
	if err != nil {
		return -1, err
	}

	
	offset:=int64(binary.LittleEndian.Uint32(buffer))
	external_id_bytes := []byte(external_id)
	file.Seek(offset,io.SeekCurrent)
	// skipped the meta data for the index now iterating the external id map 

	for ctr := range(idx.count){
		file.Read(buffer)
		length_to_skip:=int32(binary.LittleEndian.Uint32(buffer))
		if(length_to_skip==int32((len(external_id_bytes)))){
			temp_buffer:= make([]byte,length_to_skip)
			file.Read(temp_buffer)
			if(bytes.Equal(temp_buffer,external_id_bytes)){
				return ctr,nil
			}//else dont need to do anything coz we already moved ptr
		}else {
			file.Seek(int64(length_to_skip),io.SeekCurrent)
		}
	}


	return -1,nil
}




func (idx *Index) FetchVectorByExternalId(id string)(*Vector , error){


internal_id := 0;
return idx.FetchVectorByInternalId(internal_id)



}