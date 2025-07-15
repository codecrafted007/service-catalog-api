package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/codecrafted007/service-catalog-api/internal/logger"
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

func (ss *sqliteStore) ListServices(ctx context.Context, filter, sort string, page, limit int) ([]model.Service, error) {
	offset := (page - 1) * limit

	var queryBuilder strings.Builder
	queryBuilder.WriteString(`
		SELECT s.id, s.name, s.description, s.created_at,
		GROUP_CONCAT(v.version, ',') AS versions
		FROM services s
		LEFT JOIN versions v ON s.id = v.service_id`)

	args := make([]interface{}, 0)
	conditions := make([]string, 0)

	if filter != "" {
		conditions = append(conditions, "s.name LIKE ? OR s.description LIKE ?")
		filterValue := fmt.Sprintf("%%%s%%", filter)
		args = append(args, filterValue, filterValue)
	}

	if len(conditions) > 0 {
		queryBuilder.WriteString(" WHERE ")
		queryBuilder.WriteString(strings.Join(conditions, " AND "))
	}

	queryBuilder.WriteString(" GROUP BY s.id")

	if sort != "" {
		sortBy := "s.name"
		if strings.ToLower(sort) == "createdat" {
			sortBy = "s.created_at"
		}
		queryBuilder.WriteString(" ORDER BY " + sortBy)
	}

	queryBuilder.WriteString(" LIMIT ? OFFSET ?")
	args = append(args, limit, offset)
	logger.L().Infow("Executing query", "query", queryBuilder.String(), "args", args)
	rows, err := ss.db.QueryContext(ctx, queryBuilder.String(), args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var services []model.Service

	for rows.Next() {
		var svc model.Service
		var versionsCSV sql.NullString

		err := rows.Scan(&svc.ID, &svc.Name, &svc.Description, &svc.CreatedAt, &versionsCSV)
		if err != nil {
			return nil, err
		}
		if versionsCSV.Valid {
			svc.Versions = strings.Split(versionsCSV.String, ",")
		} else {
			svc.Versions = []string{}
		}
		services = append(services, svc)
	}

	return services, nil
}

func (ss *sqliteStore) GetServiceById(ctx context.Context, serviceId int) (*model.Service, error) {
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
			svc.Versions = append(svc.Versions, verStr.String)
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

func (s *sqliteStore) CreateService(ctx context.Context, service *model.Service) (int64, error) {
	result, err := s.db.Exec(`
		INSERT INTO services (name, description, created_at)
		VALUES (?, ?, CURRENT_TIMESTAMP)
	`, service.Name, service.Description)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (s *sqliteStore) UpdateService(ctx context.Context, id int, service *model.Service) error {
	result, err := s.db.Exec(`
		UPDATE services
		SET name = ?, description = ?
		WHERE id = ?
	`, service.Name, service.Description, id)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (s *sqliteStore) DeleteService(ctx context.Context, id int) error {
	_, err := s.db.Exec(`
		DELETE FROM services
		WHERE id = ?
	`, id)
	return err
}

func (s *sqliteStore) CreateVersion(ctx context.Context, v *model.Version) (int64, error) {
	result, err := s.db.ExecContext(ctx, `
		INSERT INTO versions (service_id, version, changelog, created_at)
		VALUES (?, ?, ?, ?)
	`, v.ServiceID, v.Version, v.Changelog, v.CreatedAt)
	if err != nil {
		return 0, err
	}
	lastInsertID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return lastInsertID, nil
}

func (s *sqliteStore) GetVersionsByServiceID(ctx context.Context, serviceID int64) ([]*model.Version, error) {
	var versions []*model.Version
	err := s.db.SelectContext(ctx, &versions, `
		SELECT id, service_id, version, changelog, created_at
		FROM versions
		WHERE service_id = ?
		ORDER BY created_at DESC`, serviceID)
	return versions, err
}

func (s *sqliteStore) GetVersionByID(ctx context.Context, versionID int64) (*model.Version, error) {
	var version model.Version
	err := s.db.GetContext(ctx, &version, `
		SELECT id, service_id, version, changelog, created_at
		FROM versions
		WHERE id = ?
	`, versionID)
	if err != nil {
		return nil, err
	}
	return &version, nil
}

func (s *sqliteStore) DeleteVersionByID(ctx context.Context, versionID int64) (bool, error) {
	result, err := s.db.ExecContext(ctx, "DELETE FROM versions WHERE id = ?", versionID)
	if err != nil {
		return false, err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return false, err
	}
	return rows > 0, nil
}
