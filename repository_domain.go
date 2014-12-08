package hivdomainstatus

import (
	"database/sql"
	_ "github.com/lib/pq"
	"fmt"
)

type DomainRepositoryInterface interface {
	Persist(domain *Domain) (err error)
	Remove(domain *Domain) (err error)
	FindAll() (domains []*Domain, err error)
	FindPaginated(numitems int, offsetKey string) (domains []*Domain, err error)
	Stats() (count int, maxKey string, err error)
	FindById(id int64) (domain *Domain, err error)
	FindByName(name string) (domain *Domain, err error)
}

type DomainRepository struct {
	DomainRepositoryInterface
	db         *sql.DB
	TABLE_NAME string
	FIELDS     string
	OFFSET_FIELD  string
	CREATED_FIELD string
}

func NewDomainRepository(db *sql.DB) (repo *DomainRepository) {
	repo = new(DomainRepository)
	repo.db = db
	repo.TABLE_NAME = "domain"
	repo.FIELDS = "name, valid"
	repo.OFFSET_FIELD = "id"
	repo.CREATED_FIELD = "created"
	return
}

func (repo *DomainRepository) Persist(domain *Domain) (err error) {
	if domain.Id > 0 {
		_, err = repo.db.Exec("UPDATE "+repo.TABLE_NAME+" "+
			"SET valid = $1 WHERE id = $2",domain.Valid, domain.Id)
	} else {
		err = repo.db.QueryRow("INSERT INTO "+repo.TABLE_NAME+" "+
			"("+repo.FIELDS+") "+
			"VALUES($1, $2) RETURNING id, created",
			domain.Name, domain.Valid).Scan(&domain.Id, &domain.Created)
	}
	return
}

func (repo *DomainRepository) Remove(domain *Domain) (err error) {
	_, err = repo.db.Exec("DELETE FROM "+repo.TABLE_NAME+" "+
		"WHERE " + repo.OFFSET_FIELD + " = $1",
		domain.Id)
	return
}

func (repo *DomainRepository) rowsToResult(rows *sql.Rows) (domains []*Domain, err error) {
	domains = make([]*Domain, 0)
	for rows.Next() {
		var domain = new(Domain)
		err = rows.Scan(&domain.Id, &domain.Name, &domain.Valid, &domain.Created)
		if err != nil {
			return
		}
		domains = append(domains, domain)
	}
	err = rows.Err()
	return
}

func (repo *DomainRepository) FindAll() (domains []*Domain, err error) {
	rows, err := repo.db.Query("SELECT " + repo.OFFSET_FIELD + "," + repo.FIELDS+","+repo.CREATED_FIELD + " FROM " + repo.TABLE_NAME)
	if err != nil {
		return
	}
	defer rows.Close()
	domains, err = repo.rowsToResult(rows)
	return
}

func (repo *DomainRepository) FindPaginated(numitems int, offsetKey string) (domains []*Domain, err error) {
	var rows *sql.Rows
	if len(offsetKey) > 0 {
		rows, err = repo.db.Query("SELECT "+repo.OFFSET_FIELD + "," + repo.FIELDS+","+repo.CREATED_FIELD+" "+"FROM "+repo.TABLE_NAME+" WHERE "+repo.OFFSET_FIELD+" > $1 ORDER BY "+repo.OFFSET_FIELD+" ASC LIMIT $2", offsetKey, numitems)
	} else {
		rows, err = repo.db.Query("SELECT "+repo.OFFSET_FIELD + "," + repo.FIELDS+","+repo.CREATED_FIELD+" "+"FROM "+repo.TABLE_NAME+" ORDER BY "+repo.OFFSET_FIELD+" ASC LIMIT $1", numitems)
	}
	if err != nil {
		return
	}
	defer rows.Close()
	domains, err = repo.rowsToResult(rows)
	return
}

func (repo *DomainRepository) Stats() (count int, maxKey string, err error) {
	var maxKeyInt sql.NullInt64
	err = repo.db.QueryRow("SELECT COUNT("+repo.OFFSET_FIELD+"), MAX("+repo.OFFSET_FIELD+") FROM "+repo.TABLE_NAME).Scan(&count, &maxKeyInt)
	if maxKeyInt.Valid {
		// If table is empty MAX(id) is null
		maxKey = fmt.Sprintf("%d", maxKeyInt.Int64)
	}
	return
}

func (repo *DomainRepository) FindById(id int64) (domain *Domain, err error) {
	domain = new(Domain)
	err = repo.db.QueryRow("SELECT " + repo.OFFSET_FIELD + "," + repo.FIELDS+","+repo.CREATED_FIELD + " FROM " + repo.TABLE_NAME + " WHERE " + repo.OFFSET_FIELD + " = $1", id).Scan(&domain.Id, &domain.Name, &domain.Valid, &domain.Created)
	return
}

func (repo *DomainRepository) FindByName(name string) (domain *Domain, err error) {
	domain = new(Domain)
	err = repo.db.QueryRow("SELECT " + repo.OFFSET_FIELD + "," + repo.FIELDS+","+repo.CREATED_FIELD + " FROM " + repo.TABLE_NAME + " WHERE name = $1", name).Scan(&domain.Id, &domain.Name, &domain.Valid, &domain.Created)
	return
}