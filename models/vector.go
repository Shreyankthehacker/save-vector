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
	buffer := make([]byte, 4)

	_, err = io.ReadFull(file, buffer)
	if err != nil {
		return nil, err
	}
	
	offset:=int64(binary.LittleEndian.Uint32(buffer))
	vsize:=int64(idx.dim*4)
	offset+=int64(id*int(vsize))
	vector_buffer:=make([]byte,vsize)
	file.Seek(offset,io.SeekCurrent)
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

func FetchInternalIdByExternalId(external_id string)(int , error){
	return 0,nil
}

func (idx *Index) FetchVectorByExternalId(id string)(*Vector , error){


internal_id := 0;
return idx.FetchVectorByInternalId(internal_id)



}