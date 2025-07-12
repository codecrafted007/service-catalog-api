package sqlite

import (
	"database/sql"
	"strings"

	"github.com/codecrafted007/service-catalog-api/internal/storage"
	"github.com/codecrafted007/service-catalog-api/model"
	"github.com/jmoiron/sqlx"
)

type sqliteStore struct {
	db *sqlx.DB
}

func New(path string) (storage.Storage, error) {
	db, err := sqlx.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}
	return &sqliteStore{db: db}, nil
}

func (ss *sqliteStore) DB() *sqlx.DB {
	return ss.db
}

func (ss *sqliteStore) ListServices(filter, sort string, page, limit int) ([]model.Service, error) {
	offset := (page - 1) * limit
	var queryBuilder strings.Builder
	queryBuilder.WriteString("select * from services")

	args := map[string]interface{}{
		"limit":  limit,
		"offset": offset,
	}

	queryClauses := []string{}

	if filter != "" {
		queryClauses = append(queryClauses, "name like :filter or description like :filter")
		args["filter"] = "%" + filter + "%"
	}

	if len(queryClauses) > 0 {
		queryBuilder.WriteString(" WHERE " + strings.Join(queryClauses, " AND "))
	}

	if sort != "" {
		sortBy := "name"
		if strings.ToLower(sort) == "createdAt" {
			sortBy = "createdAt"
		}
		queryBuilder.WriteString(" ORDER BY " + sortBy)
	}

	queryBuilder.WriteString(" limit :limit offset :offset")

	stmt, err := ss.db.PrepareNamed(queryBuilder.String())
	if err != nil {
		return nil, err
	}
	var services []model.Service
	err = stmt.Select(&services, args)
	return services, err
}

func (ss *sqliteStore) GetServiceById(serviceId int) (*model.Service, error) {
	rows, err := ss.db.Queryx(`
		SELECT 
			s.id AS service_id, s.name, s.description, s.created_at AS service_created_at,
			v.id AS version_id, v.version, v.created_at AS version_created_at
		FROM services s
		LEFT JOIN versions v ON s.id = v.service_id
		WHERE s.id = ?
	`, serviceId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var svc *model.Service
	for rows.Next() {
		var (
			sid                        int
			name, desc                 string
			svcCreatedAt, verCreatedAt sql.NullTime
			verID                      sql.NullInt64
			verStr                     sql.NullString
		)
		err := rows.Scan(&sid, &name, &desc, &svcCreatedAt, &verID, &verStr, &verCreatedAt)
		if err != nil {
			return nil, err
		}

		if svc == nil {
			svc = &model.Service{
				ID:          sid,
				Name:        name,
				Description: desc,
				CreatedAt:   svcCreatedAt.Time,
			}
		}

		if verID.Valid && verStr.Valid && verCreatedAt.Valid {
			svc.Versions = append(svc.Versions, model.Version{
				ID:        int(verID.Int64),
				Version:   verStr.String,
				CreatedAt: verCreatedAt.Time,
			})
		}
	}

	if svc == nil {
		return nil, sql.ErrNoRows
	}

	return svc, nil
}

func (s *sqliteStore) IsValidAPIKey(key string) bool {
	var exists bool
	err := s.db.Get(&exists, `SELECT EXISTS(SELECT 1 FROM api_keys WHERE key = ?)`, key)
	return err == nil && exists
}

func (s *sqliteStore) CreateService(service *model.Service) (int64, error) {
	result, err := s.db.Exec(`
		INSERT INTO services (name, description, created_at)
		VALUES (?, ?, CURRENT_TIMESTAMP)
	`, service.Name, service.Description)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (s *sqliteStore) UpdateService(id int, service *model.Service) error {
	_, err := s.db.Exec(`
		UPDATE services
		SET name = ?, description = ?
		WHERE id = ?
	`, service.Name, service.Description, id)
	return err
}

func (s *sqliteStore) DeleteService(id int) error {
	_, err := s.db.Exec(`
		DELETE FROM services
		WHERE id = ?
	`, id)
	return err
}
