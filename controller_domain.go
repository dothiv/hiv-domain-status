package hivdomainstatus

import (
	"encoding/json"
	"net/http"
	"os"
	"fmt"
)

type DomainController struct {
	repo DomainRepositoryInterface
}

func (c *DomainController) ListingHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(400)
		return
	}
	formErr := r.ParseForm()
	if formErr != nil {
		w.WriteHeader(500)
		os.Stderr.WriteString(formErr.Error() + "\n")
		return
	}

	itemsPerPage := 100
	offsetKey := r.Form.Get("offsetKey")
	items, findErr := c.repo.FindPaginated(itemsPerPage, offsetKey)
	if findErr != nil {
		w.WriteHeader(500)
		os.Stderr.WriteString(findErr.Error() + "\n")
		return
	}

	total, maxKey, statsErr := c.repo.Stats()
	if statsErr != nil {
		w.WriteHeader(500)
		os.Stderr.WriteString(statsErr.Error() + "\n")
		return
	}

	list := new(DomainListModel)
	list.Total = total
	list.JsonLDContext = "http://jsonld.click4life.hiv/List"
	list.JsonLDType = "http://jsonld.click4life.hiv/Domain"
	list.JsonLDId = "/domain"
	list.Items = make([]*DomainModel, len(items))

	proto := "https"
	if (r.TLS == nil) {
		proto = "http"
	}

	for i := range items {
		e := new(DomainModel)
		e.JsonLDContext = "http://jsonld.click4life.hiv/Domain"
		e.JsonLDId = fmt.Sprintf("%s://%s%s/%d", proto, r.Host, r.URL, items[i].Id)
		e.Id = fmt.Sprintf("%d", items[i].Id)
		e.Name = items[i].Name
		list.Items[i] = e
	}

	w.Header().Add("Content-Type", "application/json")
	// Add nwext link
	if len(items) > 0 {
		last := list.Items[len(items)-1]
		w.Header().Add("Link", `</domain?offsetKey=`+last.Id+`>; rel="next"`)
	} else {
		w.Header().Add("Link", `</domain?offsetKey=`+maxKey+`>; rel="next"`)
	}
	encoder := json.NewEncoder(w)
	encoder.Encode(list)
}
