package hivdomainstatus

import (
	"encoding/json"
	"net/http"
)

type EntryPoint struct {
	JsonLDContext string       `json:"@context"`
	Domains       *JsonLDTypedModel `json:"domains"`
}

type EntryPointController struct {
}

func (c *EntryPointController) EntryPointHandler(w http.ResponseWriter, r *http.Request, routeParams []string) {
	if r.Method != "GET" {
		w.WriteHeader(400)
		return
	}
	entryPoint := new(EntryPoint)
	entryPoint.JsonLDContext = "http://jsonld.click4life.hiv/EntryPoint"
	entryPoint.Domains = new(JsonLDTypedModel)
	entryPoint.Domains.JsonLDContext = "http://jsonld.click4life.hiv/List"
	entryPoint.Domains.JsonLDType = "http://jsonld.click4life.hiv/Domain"
	entryPoint.Domains.JsonLDId = "/domain"
	w.Header().Add("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	encoder.Encode(entryPoint)
}
