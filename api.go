package jsonpatch

import (
	"github.com/evanphx/json-patch/json"
)

// GetAPI returns the json API instance being used
func GetAPI() json.API {
	return json.GetAPI()
}

// SetAPI changes the json API instance being used
func SetAPI(newApi json.API) {
	json.SetAPI(newApi)
}

// ResetAPI sets back to using standard json API
func ResetAPI() {
	json.ResetAPI()
}
