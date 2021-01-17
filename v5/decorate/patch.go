package decorate

import (
"encoding/json"
jsonPatch "github.com/evanphx/json-patch/v5"
)


const (
	OpReplace string = "replace"
)

type ReplacePatch struct{
	Op string `json:"op"`
	Path string `json:"path"`
	Value string `json:"value"`
}

type JsonPatch struct{
	Patches []interface{}
}

func (jp *JsonPatch) AddReplacePatch(path, value string)  {
	if len(jp.Patches) == 0{
		jp.Patches = make([]interface{}, 0)
	}
	rp := ReplacePatch{Op:OpReplace, Path: path, Value: value}
	jp.Patches = append(jp.Patches, rp)
}


func (jp *JsonPatch) ApplyPatch( originalJson string) (string, error){
	patches, err := json.Marshal(jp.Patches)
	if err != nil{
		return  "", err
	}
	decodedPatches, err := jsonPatch.DecodePatch(patches)
	if err != nil {
		return  "", err
	}
	result, err := decodedPatches.Apply([]byte(originalJson))
	if err != nil{
		return "", err
	}
	return string(result), nil
}