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
















func (idx *Index) LoadAllVectorFromLayer(layer int) ([][]float32, error) {
	
	
	
	
	count := idx.count
	if count == 0 {
		return [][]float32{}, nil
	}

	var result [][]float32
	for i := 0; i < count; i++ {
		v, err := idx.FetchVectorByInternalId(i)
		if err != nil {
			return nil, err
		}
		result = append(result, v.data)
	}
	return result, nil
}