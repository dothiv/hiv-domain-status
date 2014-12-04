package hivdomainstatus

import (
	"database/sql"
	_ "github.com/lib/pq"
	"fmt"
)

type DomainRepositoryInterface interface {
	Persist(m *Domain) (result sql.Result, err error)
	FindAll() (result []*Domain, err error)
	FindPaginated(numitems int, offsetKey string) (result []*Domain, err error)
	Stats() (count int, maxKey string, err error)
}

type DomainRepository struct {
	DomainRepositoryInterface
	db         *sql.DB
	TABLE_NAME string
	FIELDS     string
	OFFSET_FIELD  string
}

func NewDomainRepository(db *sql.DB) (repo *DomainRepository) {
	repo = new(DomainRepository)
	repo.db = db
	repo.TABLE_NAME = "domain"
	repo.FIELDS = "name"
	repo.OFFSET_FIELD = "id"
	return
}

func (repo *DomainRepository) Persist(m *Domain) (result sql.Result, err error) {
	result, err = repo.db.Exec("INSERT INTO "+repo.TABLE_NAME+" "+
		"("+repo.FIELDS+") "+
		"VALUES($1)",
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
	rows, err := repo.db.Query("SELECT " + repo.OFFSET_FIELD + "," + repo.FIELDS + " FROM " + repo.TABLE_NAME)
	if err != nil {
		return
	}
	defer rows.Close()
	result, err = repo.rowsToResult(rows)
	return
}

func (repo *DomainRepository) FindPaginated(numitems int, offsetKey string) (result []*Domain, err error) {
	var rows *sql.Rows
	if len(offsetKey) > 0 {
		rows, err = repo.db.Query("SELECT "+repo.OFFSET_FIELD + "," + repo.FIELDS+" "+"FROM "+repo.TABLE_NAME+" WHERE "+repo.OFFSET_FIELD+" > $1 ORDER BY "+repo.OFFSET_FIELD+" ASC LIMIT $2", offsetKey, numitems)
	} else {
		rows, err = repo.db.Query("SELECT "+repo.OFFSET_FIELD + "," + repo.FIELDS+" "+"FROM "+repo.TABLE_NAME+" ORDER BY "+repo.OFFSET_FIELD+" ASC LIMIT $1", numitems)
	}
	if err != nil {
		return
	}
	defer rows.Close()
	result, err = repo.rowsToResult(rows)
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