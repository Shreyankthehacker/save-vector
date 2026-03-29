package models

import (

	"github.com/Shreyankthehacker/savector/utils"
)


func(db *DataBase)FetchIndex(name string)(*Index ,error ){


index_path := db.Name+"/"+name 

index_buffer,err:=utils.GetIndexBuffer(index_path)

if err!=nil{
utils.Logger.Println("failed to load the buffer properly")
return nil,err
}

return getIndexFromMetaDataByte(index_buffer)

}
