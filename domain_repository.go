package hivdomainstatus

import (
	"database/sql"
	_ "github.com/lib/pq"
)

type DomainRepositoryInterface interface {
	Persist(m *Domain) (result sql.Result, err error)
}

type DomainRepository struct {
	DomainRepositoryInterface
	db         *sql.DB
	TABLE_NAME string
	FIELDS     string
}

func NewDomainRepository(db *sql.DB) (repo *DomainRepository) {
	repo = new(DomainRepository)
	repo.db = db
	repo.TABLE_NAME = "domain"
	repo.FIELDS = "id, name"
	return
}

func (repo *DomainRepository) Persist(m *Domain) (result sql.Result, err error) {
	result, err = repo.db.Exec("INSERT INTO "+repo.TABLE_NAME+" "+
		"("+repo.FIELDS+") "+
		"VALUES($1, $2)",
		m.Id,
		m.Name)
	return
}

func (repo *DomainRepository) rowsToResult(rows *sql.Rows) (result []*Domain, err error) {
	result = make([]*Domain, 0)
	for rows.Next() {
		var m = new(Domain)
		err = rows.Scan(&m.Id,
			&m.Name)
		if err != nil {
			return
		}
		result = append(result, m)
	}
	err = rows.Err()
	return
}

func (repo *DomainRepository) FindAll() (result []*Domain, err error) {
	rows, err := repo.db.Query("SELECT " + repo.FIELDS + " FROM " + repo.TABLE_NAME)
	if err != nil {
		return
	}
	defer rows.Close()
	result, err = repo.rowsToResult(rows)
	return
}