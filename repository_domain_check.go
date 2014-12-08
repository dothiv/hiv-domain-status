package hivdomainstatus

import (
	"database/sql"
	_ "github.com/lib/pq"
	"fmt"
)

type DomainCheckRepositoryInterface interface {
	Persist(result *DomainCheck) (err error)
	Remove(domain *Domain) (err error)
	FindAll() (results []*DomainCheck, err error)
	FindByDomain(domain string) (result *DomainCheck, err error)
	FindPaginated(numitems int, offsetKey string) (domains []*Domain, err error)
	Stats() (count int, maxKey string, err error)
}

type DomainCheckRepository struct {
	DomainCheckRepositoryInterface
	db         *sql.DB
	TABLE_NAME string
	ID_FIELD  string
	FIELDS     string
	CREATED_FIELD string
}

func NewDomainCheckRepository(db *sql.DB) (repo *DomainCheckRepository) {
	repo = new(DomainCheckRepository)
	repo.db = db
	repo.TABLE_NAME = "domain_check"
	repo.FIELDS = "domain, url, status_code, script_present, iframe_target, iframe_target_ok, valid"
	repo.ID_FIELD = "id"
	repo.CREATED_FIELD = "created"
	return
}

func (repo *DomainCheckRepository) Persist(result *DomainCheck) (err error) {
	if result.Id > 0 {
		_, err = repo.db.Exec("UPDATE "+repo.TABLE_NAME+" "+
			"SET domain = $1, url = $2, status_code = $3, script_present = $4, iframe_target = $5, iframe_target_ok = $6, valid = $7 WHERE id = $8",
			result.Domain, result.URL, result.StatusCode, result.ScriptPresent, result.IframeTarget, result.IframeTargetOk, result.Valid, result.Id)
	} else {
		err = repo.db.QueryRow("INSERT INTO "+repo.TABLE_NAME+" "+
			"("+repo.FIELDS+") "+
			"VALUES($1, $2, $3, $4, $5, $6, $7) RETURNING id, created",
			result.Domain, result.URL, result.StatusCode, result.ScriptPresent, result.IframeTarget, result.IframeTargetOk, result.Valid).Scan(&result.Id, &result.Created)
	}
	return
}

func (repo *DomainCheckRepository) Remove(result *DomainCheck) (err error) {
	_, err = repo.db.Exec("DELETE FROM "+repo.TABLE_NAME+" "+
		"WHERE " + repo.ID_FIELD + " = $1",
		result.Id)
	return
}

func (repo *DomainCheckRepository) rowsToResult(rows *sql.Rows) (results []*DomainCheck, err error) {
	results = make([]*DomainCheck, 0)
	for rows.Next() {
		var result = new(DomainCheck)
		err = rows.Scan(&result.Id, &result.Domain, &result.URL, &result.StatusCode, &result.ScriptPresent, &result.IframeTarget, &result.IframeTargetOk, &result.Valid, &result.Created)
		if err != nil {
			return
		}
		results = append(results, result)
	}
	err = rows.Err()
	return
}

func (repo *DomainCheckRepository) FindAll() (results []*DomainCheck, err error) {
	rows, err := repo.db.Query("SELECT " + repo.ID_FIELD + "," + repo.FIELDS + "," + repo.CREATED_FIELD+" FROM " + repo.TABLE_NAME)
	if err != nil {
		return
	}
	defer rows.Close()
	results, err = repo.rowsToResult(rows)
	return
}

func (repo *DomainCheckRepository) FindPaginated(numitems int, offsetKey string) (results []*DomainCheck, err error) {
	var rows *sql.Rows
	if len(offsetKey) > 0 {
		rows, err = repo.db.Query("SELECT "+repo.ID_FIELD + "," + repo.FIELDS+"," + repo.CREATED_FIELD+" "+"FROM "+repo.TABLE_NAME+" WHERE "+repo.ID_FIELD+" > $1 ORDER BY "+repo.ID_FIELD+" ASC LIMIT $2", offsetKey, numitems)
	} else {
		rows, err = repo.db.Query("SELECT "+repo.ID_FIELD + "," + repo.FIELDS+"," + repo.CREATED_FIELD+" "+"FROM "+repo.TABLE_NAME+" ORDER BY "+repo.ID_FIELD+" ASC LIMIT $1", numitems)
	}
	if err != nil {
		return
	}
	defer rows.Close()
	results, err = repo.rowsToResult(rows)
	return
}

func (repo *DomainCheckRepository) Stats() (count int, maxKey string, err error) {
	var maxKeyInt sql.NullInt64
	err = repo.db.QueryRow("SELECT COUNT("+repo.ID_FIELD+"), MAX("+repo.ID_FIELD+") FROM "+repo.TABLE_NAME).Scan(&count, &maxKeyInt)
	if maxKeyInt.Valid {
		// If table is empty MAX(id) is null
		maxKey = fmt.Sprintf("%d", maxKeyInt.Int64)
	}
	return
}

func (repo *DomainCheckRepository) FindById(id int64) (result *DomainCheck, err error) {
	result = new(DomainCheck)
	err = repo.db.QueryRow("SELECT " + repo.ID_FIELD + "," + repo.FIELDS+","+repo.CREATED_FIELD + " FROM " + repo.TABLE_NAME + " WHERE " + repo.ID_FIELD + " = $1", id).Scan(&result.Id, &result.Domain, &result.URL, &result.StatusCode, &result.ScriptPresent, &result.IframeTarget, &result.IframeTargetOk, &result.Valid, &result.Created)
	return
}

func (repo *DomainCheckRepository) FindByDomain(domain string) (result *DomainCheck, err error) {
	result = new(DomainCheck)
	err = repo.db.QueryRow("SELECT " + repo.ID_FIELD + "," + repo.FIELDS+","+repo.CREATED_FIELD + " FROM " + repo.TABLE_NAME + " WHERE domain = $1", domain).Scan(&result.Id, &result.Domain, &result.URL, &result.StatusCode, &result.ScriptPresent, &result.IframeTarget, &result.IframeTargetOk, &result.Valid, &result.Created)
	return
}