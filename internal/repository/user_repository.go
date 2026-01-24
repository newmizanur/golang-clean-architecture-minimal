package repository

import (
	"context"
	"database/sql"
	"golang-clean-architecture/internal/apperror"
	dbmodel "golang-clean-architecture/internal/entity/db/model"
	t "golang-clean-architecture/internal/entity/db/table"

	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/sirupsen/logrus"
)

type UserRepository struct {
	DB  *sql.DB
	Log *logrus.Logger
}

func NewUserRepository(db *sql.DB, log *logrus.Logger) *UserRepository {
	return &UserRepository{
		DB:  db,
		Log: log,
	}
}

func (r *UserRepository) CountById(ctx context.Context, tx *sql.Tx, id string) (int64, error) {
	var result struct {
		Total int64
	}
	stmt := mysql.SELECT(mysql.COUNT(t.Users.ID).AS("total")).
		FROM(t.Users).
		WHERE(t.Users.ID.EQ(mysql.String(id)))
	db := qrm.Queryable(r.DB)
	if tx != nil {
		db = tx
		stmt = stmt.FOR(mysql.UPDATE())
	}
	err := stmt.QueryContext(ctx, db, &result)
	return result.Total, err
}

func (r *UserRepository) FindById(ctx context.Context, tx *sql.Tx, id string) (*dbmodel.Users, error) {
	stmt := mysql.SELECT(
		t.Users.ID,
		t.Users.Name,
		t.Users.Password,
		t.Users.CreatedAt,
		t.Users.UpdatedAt,
	).
		FROM(t.Users).
		WHERE(t.Users.ID.EQ(mysql.String(id))).
		LIMIT(1)
	var db qrm.Queryable = r.DB
	if tx != nil {
		db = tx
		stmt = stmt.FOR(mysql.UPDATE())
	}
	r.Log.Print(stmt.DebugSql())
	user := new(dbmodel.Users)
	if err := stmt.QueryContext(ctx, db, user); err != nil {
		r.Log.Error(err)
		if apperror.IsNoRows(err) {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) Create(ctx context.Context, tx *sql.Tx, user *dbmodel.Users) error {
	stmt := t.Users.INSERT(
		t.Users.ID,
		t.Users.Password,
		t.Users.Name,
		t.Users.CreatedAt,
		t.Users.UpdatedAt,
	).MODEL(user)
	db := qrm.Executable(r.DB)
	if tx != nil {
		db = tx
	}
	_, err := stmt.ExecContext(ctx, db)
	return err
}

func (r *UserRepository) Update(ctx context.Context, tx *sql.Tx, user *dbmodel.Users) error {
	stmt := t.Users.UPDATE(
		t.Users.Password,
		t.Users.Name,
		t.Users.CreatedAt,
		t.Users.UpdatedAt,
	).
		SET(
			t.Users.Password.SET(mysql.String(user.Password)),
			t.Users.Name.SET(mysql.String(user.Name)),
			t.Users.UpdatedAt.SET(mysql.Int(user.UpdatedAt)),
		).
		WHERE(t.Users.ID.EQ(mysql.String(user.ID)))
	db := qrm.Executable(r.DB)
	if tx != nil {
		db = tx
	}
	_, err := stmt.ExecContext(ctx, db)
	return err
}

func (r *UserRepository) Delete(ctx context.Context, tx *sql.Tx, user *dbmodel.Users) error {
	stmt := t.Users.DELETE().
		WHERE(t.Users.ID.EQ(mysql.String(user.ID)))
	db := qrm.Executable(r.DB)
	if tx != nil {
		db = tx
	}
	_, err := stmt.ExecContext(ctx, db)
	return err
}
