package decorate

import "testing"

var original string = `{ "id": "1024",  "name": "linux", "spec":{"memory":"4G", "cpu": "intel"}}`


func TestReplacePatch(t *testing.T){
	t.Log(original)
	jp := JsonPatch{}
	jp.AddReplacePatch("/id", "2048")
	jp.AddReplacePatch("/name", "macOS")
	jp.AddReplacePatch("/spec/memory", "8G")
	jp.AddReplacePatch("/spec/cpu", "amd")
	result, err := jp.ApplyPatch(original)
	if err != nil{
		t.Error(err.Error())
	}else{
		t.Log(result)
	}
}