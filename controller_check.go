package hivdomainstatus

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type DomainCheckController struct {
	domainCheckRepo DomainCheckRepositoryInterface
}

func (c *DomainCheckController) ListingHandler(w http.ResponseWriter, r *http.Request, routeParams []string) {
	if r.Method != "GET" {
		HttpProblem(w, http.StatusBadRequest, "Method not allow: "+r.Method)
		return
	}
	formErr := r.ParseForm()
	if formErr != nil {
		HttpProblem(w, http.StatusInternalServerError, formErr.Error())
		return
	}

	itemsPerPage := 100
	offsetKey := r.Form.Get("offsetKey")
	items, findErr := c.domainCheckRepo.FindPaginated(itemsPerPage, offsetKey)
	if findErr != nil {
		HttpProblem(w, http.StatusInternalServerError, findErr.Error())
		return
	}

	total, maxKey, statsErr := c.domainCheckRepo.Stats()
	if statsErr != nil {
		HttpProblem(w, http.StatusInternalServerError, statsErr.Error())
		return
	}

	list := new(DomainCheckListModel)
	list.Total = total
	list.JsonLDContext = "http://jsonld.click4life.hiv/List"
	list.JsonLDType = "http://jsonld.click4life.hiv/DomainCheck"
	list.JsonLDId = getHttpHost(r)
	list.Items = make([]*DomainCheckModel, len(items))

	for i, item := range items {
		e := transformCheckEntity(item, getHttpHost(r)+"/check/%d")
		list.Items[i] = e
	}

	w.Header().Add("Content-Type", "application/json")
	// Add nwext link
	if len(items) > 0 {
		last := list.Items[len(items)-1]
		w.Header().Add("Link", fmt.Sprintf(`<%s/check?offsetKey=%s>; rel="next"`, getHttpHost(r), last.Id))
	} else {
		w.Header().Add("Link", fmt.Sprintf(`<%s/check?offsetKey=%s>; rel="next"`, getHttpHost(r), maxKey))
	}
	encoder := json.NewEncoder(w)
	encoder.Encode(list)
}
