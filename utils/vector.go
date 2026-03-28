package utils

import (
	"bytes"
	"encoding/binary"
	"errors"
)



func Flatten(data interface{})([]float32,error){


switch v :=data.(type){
case []float32:return v,nil
case [][]float32 :
	var result []float32
	for _,row := range v{
		result = append(result , row...)
	}
	return result,nil
}
return nil, errors.New("cannot flatten: only []float32 or [][]float32 supported")
}


func FloatToByte(data []float32)([]byte , error){

buffer:= new(bytes.Buffer)

err:= binary.Write(buffer,binary.LittleEndian,data)
if err!=nil{
	Logger.Fatal("Issue in converting vector to byte")
	return  nil,err
}

return buffer.Bytes(),nil


}