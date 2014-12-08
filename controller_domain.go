package hivdomainstatus

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

type DomainController struct {
	domainRepo      DomainRepositoryInterface
	domainCheckRepo DomainCheckRepositoryInterface
}

func (c *DomainController) ListingHandler(w http.ResponseWriter, r *http.Request, routeParams []string) {
	if r.Method == "POST" {
		c.createItem(w, r, routeParams)
		return
	}
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
	items, findErr := c.domainRepo.FindPaginated(itemsPerPage, offsetKey)
	if findErr != nil {
		HttpProblem(w, http.StatusInternalServerError, findErr.Error())
		return
	}

	total, maxKey, statsErr := c.domainRepo.Stats()
	if statsErr != nil {
		HttpProblem(w, http.StatusInternalServerError, statsErr.Error())
		return
	}

	list := new(DomainListModel)
	list.Total = total
	list.JsonLDContext = "http://jsonld.click4life.hiv/List"
	list.JsonLDType = "http://jsonld.click4life.hiv/Domain"
	list.JsonLDId = getHttpHost(r)
	list.Items = make([]*DomainModel, len(items))

	for i, item := range items {
		e := transformEntity(item, getHttpHost(r)+"/domain/%d")
		list.Items[i] = e
	}

	w.Header().Add("Content-Type", "application/json")
	// Add nwext link
	if len(items) > 0 {
		last := list.Items[len(items)-1]
		w.Header().Add("Link", fmt.Sprintf(`<%s/domain?offsetKey=%s>; rel="next"`, getHttpHost(r), last.Id))
	} else {
		w.Header().Add("Link", fmt.Sprintf(`<%s/domain?offsetKey=%s>; rel="next"`, getHttpHost(r), maxKey))
	}
	encoder := json.NewEncoder(w)
	encoder.Encode(list)
}

func transformEntity(e *Domain, route string) (m *DomainModel) {
	m = new(DomainModel)
	m.JsonLDContext = "http://jsonld.click4life.hiv/Domain"
	m.JsonLDId = fmt.Sprintf(route, e.Id)
	m.Id = fmt.Sprintf("%d", e.Id)
	m.Name = e.Name
	m.Check = new(JsonLDTypedModel)
	m.Check.JsonLDContext = "http://jsonld.click4life.hiv/DomainCheck"
	m.Check.JsonLDId = m.JsonLDId + "/check"
	return
}

func getHttpHost(r *http.Request) string {
	proto := "https"
	if r.TLS == nil {
		proto = "http"
	}
	return fmt.Sprintf("%s://%s", proto, r.Host)
}

func (c *DomainController) createItem(w http.ResponseWriter, r *http.Request, routeParams []string) {
	if r.Header.Get("Content-Type") != "application/json" {
		HttpProblem(w, http.StatusBadRequest, "Expected application/json got "+r.Header.Get("Content-Type"))
		return
	}
	b, err := ioutil.ReadAll(r.Body)
	r.Body.Close()
	if err != nil {
		HttpProblem(w, http.StatusBadRequest, "Failed to read body: "+err.Error())
		return
	}

	var m DomainModel
	unmarshalErr := json.Unmarshal(b, &m)
	if unmarshalErr != nil {
		HttpProblem(w, http.StatusBadRequest, "Failed to parse request: "+bytes.NewBuffer(b).String())
		return
	}
	domain := new(Domain)
	domain.Name = m.Name
	err = c.domainRepo.Persist(domain)
	if err != nil {
		HttpProblem(w, http.StatusBadRequest, "Failed to create domain: "+err.Error())
		return
	}
	m = *transformEntity(domain, getHttpHost(r)+"/domain/%d")
	w.Header().Add("Location", m.JsonLDId)
	w.WriteHeader(http.StatusCreated)
}

func (c *DomainController) ItemHandler(w http.ResponseWriter, r *http.Request, routeParams []string) {
	id, err := strconv.ParseInt(routeParams[1], 0, 64)
	if err != nil {
		HttpProblem(w, http.StatusBadRequest, "Invalid id: "+routeParams[1])
		return
	}
	domain, findErr := c.domainRepo.FindById(id)
	if findErr != nil {
		HttpProblem(w, http.StatusNotFound, "Domain not found: "+routeParams[1])
		return
	}

	if r.Method == "DELETE" {
		c.domainRepo.Remove(domain)
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	m := transformEntity(domain, getHttpHost(r)+"/domain/%d")
	encoder := json.NewEncoder(w)
	encoder.Encode(m)
}

func transformCheckEntity(check *DomainCheck, route string) (m *DomainCheckModel) {
	m = new(DomainCheckModel)
	m.JsonLDContext = "http://jsonld.click4life.hiv/DomainCheck"
	// m.JsonLDId = fmt.Sprintf(route, check.Id)
	m.JsonLDId = route
	m.Id = fmt.Sprintf("%d", check.Id)
	m.Domain = check.Domain
	m.DnsOK = check.DnsOK
	m.Addresses = check.Addresses
	m.URL = check.URL
	m.StatusCode = check.StatusCode
	m.ScriptPresent = check.ScriptPresent
	m.IframeTarget = check.IframeTarget
	m.IframeTargetOk = check.IframeTargetOk
	m.Valid = check.Valid
	return
}

func (c *DomainController) CheckHandler(w http.ResponseWriter, r *http.Request, routeParams []string) {
	id, err := strconv.ParseInt(routeParams[1], 0, 64)
	if err != nil {
		HttpProblem(w, http.StatusBadRequest, "Invalid id: "+routeParams[1])
		return
	}
	domain, findErr := c.domainRepo.FindById(id)
	if findErr != nil {
		HttpProblem(w, http.StatusNotFound, "Domain not found: "+routeParams[1])
		return
	}
	domainCheck, checkErr := c.domainCheckRepo.FindLatestByDomain(domain.Name)
	if checkErr != nil {
		HttpProblem(w, http.StatusNotFound, "Domain check not found: "+routeParams[1])
		return
	}
	w.Header().Add("Content-Type", "application/json")
	m := transformCheckEntity(domainCheck, fmt.Sprintf(getHttpHost(r)+"/domain/%d/check", domain.Id))
	encoder := json.NewEncoder(w)
	encodeErr := encoder.Encode(m)
	if encodeErr != nil {
		log.Fatalln(encodeErr.Error())
	}
}
