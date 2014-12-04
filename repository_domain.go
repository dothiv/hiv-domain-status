package hivdomainstatus

import (
	"database/sql"
	_ "github.com/lib/pq"
	"fmt"
)

type DomainRepositoryInterface interface {
	Persist(m *Domain) (err error)
	Remove(m *Domain) (err error)
	FindAll() (result []*Domain, err error)
	FindPaginated(numitems int, offsetKey string) (result []*Domain, err error)
	Stats() (count int, maxKey string, err error)
	FindById(id int64) (d *Domain, err error)
	FindByName(name string) (d *Domain, err error)
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
	repo.FIELDS = "name, valid"
	repo.OFFSET_FIELD = "id"
	return
}

func (repo *DomainRepository) Persist(m *Domain) (err error) {
	if m.Id > 0 {
		_, err = repo.db.Exec("UPDATE "+repo.TABLE_NAME+" "+
			"SET valid = $1 WHERE id = $2",m.Valid, m.Id)
	} else {
		err = repo.db.QueryRow("INSERT INTO "+repo.TABLE_NAME+" "+
			"("+repo.FIELDS+") "+
			"VALUES($1, $2) RETURNING id",
			m.Name, m.Valid).Scan(&m.Id)		
	}
	return
}

func (repo *DomainRepository) Remove(m *Domain) (err error) {
	_, err = repo.db.Exec("DELETE FROM "+repo.TABLE_NAME+" "+
		"WHERE " + repo.OFFSET_FIELD + " = $1",
		m.Id)
	print(fmt.Sprintf("<%d>", m.Id))
	return
}

func (repo *DomainRepository) rowsToResult(rows *sql.Rows) (result []*Domain, err error) {
	result = make([]*Domain, 0)
	for rows.Next() {
		var m = new(Domain)
		err = rows.Scan(&m.Id, &m.Name, &m.Valid)
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

func (repo *DomainRepository) FindById(id int64) (d *Domain, err error) {
	d = new(Domain)
	err = repo.db.QueryRow("SELECT " + repo.OFFSET_FIELD + "," + repo.FIELDS + " FROM " + repo.TABLE_NAME + " WHERE " + repo.OFFSET_FIELD + " = $1", id).Scan(&d.Id, &d.Name, &d.Valid)
	return
}

func (repo *DomainRepository) FindByName(name string) (d *Domain, err error) {
	d = new(Domain)
	err = repo.db.QueryRow("SELECT " + repo.OFFSET_FIELD + "," + repo.FIELDS + " FROM " + repo.TABLE_NAME + " WHERE name = $1", name).Scan(&d.Id, &d.Name, &d.Valid)
	return
}