package storer

import (
	"context"
	"errors"

	"github.com/flukis/pagination-comparation/domain"
	"github.com/flukis/pagination-comparation/utils/pagination"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

type Product interface {
	fetch(ctx context.Context, query string, args ...interface{}) (res []domain.Product, err error)
	GetPaginationWithPage(ctx context.Context, page, limit int, sortColumn, sortOrder string) (res []domain.Product, nextpage, total int, err error)
	GetPaginationWithCursor(ctx context.Context, cursor string, limit int) (res []domain.Product, nextcursor string, total int, err error)
}

type productStorer struct {
	db *pgxpool.Pool
}

// fetch implements Product.
func (p *productStorer) fetch(ctx context.Context, query string, args ...interface{}) (res []domain.Product, err error) {
	rows, err := p.db.Query(ctx, query, args...)
	if err != nil {
		log.Err(err)
		return nil, err
	}
	defer func() {
		errRow := rows.Conn().Close(ctx)
		if errRow != nil {
			log.Err(errRow).Msg("error when close connection defer")
		}
	}()
	res = make([]domain.Product, 0)
	for rows.Next() {
		t := domain.Product{}
		if err = rows.Scan(
			&t.ID,
			&t.Name,
			&t.Description,
			&t.Amount,
			&t.CreatedAt,
		); err != nil {
			log.Err(err)
			return nil, err
		}
		res = append(res, t)
	}
	return res, nil
}

// GetPaginationWithCursor implements Product.
func (p *productStorer) GetPaginationWithCursor(ctx context.Context, cursor string, limit int) (res []domain.Product, nextcursor string, total int, err error) {
	queryGetTotal := `SELECT COUNT(id) as cnt FROM products`
	row := p.db.QueryRow(ctx, queryGetTotal)
	if err := row.Scan(&total); err != nil {
		log.Err(err)
		return nil, nextcursor, total, err
	}
	if total == 0 {
		err = errors.New("data is empty")
		return nil, nextcursor, total, err
	}
	log.Debug().Int("count", total).Msg("found product items")

	query := `
		SELECT
			id,
			name,
			description,
			amount,
			created_at
		FROM products
		WHERE created_at > $1
		ORDER BY created_at
		LIMIT $2
	`
	decodesCursor, err := pagination.DecodeCursor(cursor)
	if err != nil && cursor != "" {
		return nil, nextcursor, total, err
	}
	res, err = p.fetch(ctx, query, decodesCursor, limit)
	if err != nil {
		return nil, nextcursor, total, err
	}
	if len(res) == int(limit) {
		nextcursor = pagination.EncodeCursor(res[len(res)-1].CreatedAt)
	}
	return
}

// GetPaginationWithPage implements Product.
func (p *productStorer) GetPaginationWithPage(ctx context.Context, page, limit int, sortColumn, sortOrder string) (res []domain.Product, nextpage, total int, err error) {
	if sortOrder != "ASC" && sortOrder != "DESC" {
		err = errors.New("sort order must either ASC or DESC")
		return nil, nextpage, total, err
	}
	queryGetTotal := `SELECT COUNT(id) as cnt FROM products`
	row := p.db.QueryRow(ctx, queryGetTotal)
	if err := row.Scan(&total); err != nil {
		log.Err(err)
		return nil, nextpage, total, err
	}
	if total == 0 {
		err = errors.New("data is empty")
		return nil, nextpage, total, err
	}
	log.Debug().Int("count", total).Msg("found product items")

	var query string

	if sortOrder == "DESC" {
		query = `
			SELECT
				id,
				name,
				description,
				amount,
				created_at
			FROM products
			ORDER BY $1 DESC
			LIMIT $2 OFFSET $3
		`
	} else {
		query = `
			SELECT
				id,
				name,
				description,
				amount,
				created_at
			FROM products
			ORDER BY $1 ASC
			LIMIT $2 OFFSET $3
		`
	}
	offset := (page - 1) * limit
	res, err = p.fetch(ctx, query, sortColumn, limit, offset)
	if err != nil {
		return nil, nextpage, total, err
	}
	if len(res) == limit {
		nextpage = page + 1
	}
	return
}

func NewProductStorer(db *pgxpool.Pool) Product {
	return &productStorer{db}
}
