package repository

import (
	"context"
	"database/sql"
	"golang-clean-architecture/internal/apperror"
	dbmodel "golang-clean-architecture/internal/entity/db/model"
	t "golang-clean-architecture/internal/entity/db/table"
	"golang-clean-architecture/internal/model"

	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/sirupsen/logrus"
)

type ContactRepository struct {
	DB  *sql.DB
	Log *logrus.Logger
}

func NewContactRepository(db *sql.DB, log *logrus.Logger) *ContactRepository {
	return &ContactRepository{
		DB:  db,
		Log: log,
	}
}

func (r *ContactRepository) FindByIdAndUserId(ctx context.Context, tx *sql.Tx, id string, userId string) (*dbmodel.Contacts, error) {
	stmt := mysql.SELECT(t.Contacts.AllColumns).
		FROM(t.Contacts).
		WHERE(
			t.Contacts.ID.EQ(mysql.String(id)).
				AND(t.Contacts.UserID.EQ(mysql.String(userId))),
		).
		LIMIT(1)
	db := qrm.Queryable(r.DB)
	if tx != nil {
		db = tx
		stmt = stmt.FOR(mysql.UPDATE())
	}
	contact := new(dbmodel.Contacts)
	if err := stmt.QueryContext(ctx, db, contact); err != nil {
		if apperror.IsNoRows(err) {
			return nil, nil
		}
		return nil, err
	}
	return contact, nil
}

func (r *ContactRepository) Search(ctx context.Context, tx *sql.Tx, request *model.SearchContactRequest) ([]dbmodel.Contacts, int64, error) {
	var contacts []dbmodel.Contacts
	filter := r.filterContact(request)
	offset := (request.Page - 1) * request.Size

	stmt := mysql.SELECT(t.Contacts.AllColumns).
		FROM(t.Contacts).
		WHERE(filter).
		LIMIT(int64(request.Size)).
		OFFSET(int64(offset))
	db := qrm.Queryable(r.DB)
	if tx != nil {
		db = tx
		stmt = stmt.FOR(mysql.UPDATE())
	}
	if err := stmt.QueryContext(ctx, db, &contacts); err != nil {
		return nil, 0, err
	}

	var result struct {
		Total int64
	}
	countStmt := mysql.SELECT(mysql.COUNT(t.Contacts.ID).AS("total")).
		FROM(t.Contacts).
		WHERE(filter)
	if tx != nil {
		countStmt = countStmt.FOR(mysql.UPDATE())
	}
	if err := countStmt.QueryContext(ctx, db, &result); err != nil {
		return nil, 0, err
	}

	return contacts, result.Total, nil
}

func (r *ContactRepository) filterContact(request *model.SearchContactRequest) mysql.BoolExpression {
	condition := t.Contacts.UserID.EQ(mysql.String(request.UserId))

	if name := request.Name; name != "" {
		pattern := "%" + name + "%"
		nameCondition := t.Contacts.FirstName.LIKE(mysql.String(pattern)).
			OR(t.Contacts.LastName.LIKE(mysql.String(pattern)))
		condition = condition.AND(nameCondition)
	}

	if phone := request.Phone; phone != "" {
		pattern := "%" + phone + "%"
		condition = condition.AND(t.Contacts.Phone.LIKE(mysql.String(pattern)))
	}

	if email := request.Email; email != "" {
		pattern := "%" + email + "%"
		condition = condition.AND(t.Contacts.Email.LIKE(mysql.String(pattern)))
	}

	return condition
}

func (r *ContactRepository) Create(ctx context.Context, tx *sql.Tx, contact *dbmodel.Contacts) error {
	stmt := t.Contacts.INSERT(
		t.Contacts.ID,
		t.Contacts.FirstName,
		t.Contacts.LastName,
		t.Contacts.Email,
		t.Contacts.Phone,
		t.Contacts.UserID,
		t.Contacts.CreatedAt,
		t.Contacts.UpdatedAt,
	).MODEL(contact)
	db := qrm.Executable(r.DB)
	if tx != nil {
		db = tx
	}
	_, err := stmt.ExecContext(ctx, db)
	return err
}

func (r *ContactRepository) Update(ctx context.Context, tx *sql.Tx, contact *dbmodel.Contacts) error {
	stmt := t.Contacts.UPDATE(
		t.Contacts.FirstName,
		t.Contacts.LastName,
		t.Contacts.Email,
		t.Contacts.Phone,
		t.Contacts.UserID,
		t.Contacts.CreatedAt,
		t.Contacts.UpdatedAt,
	).
		SET(
			t.Contacts.FirstName.SET(mysql.String(contact.FirstName)),
			t.Contacts.LastName.SET(stringExprOrNull(contact.LastName)),
			t.Contacts.Email.SET(stringExprOrNull(contact.Email)),
			t.Contacts.Phone.SET(stringExprOrNull(contact.Phone)),
			t.Contacts.UpdatedAt.SET(mysql.Int(contact.UpdatedAt)),
		).
		WHERE(t.Contacts.ID.EQ(mysql.String(contact.ID)))
	db := qrm.Executable(r.DB)
	if tx != nil {
		db = tx
	}
	_, err := stmt.ExecContext(ctx, db)
	return err
}

func (r *ContactRepository) Delete(ctx context.Context, tx *sql.Tx, contact *dbmodel.Contacts) error {
	stmt := t.Contacts.DELETE().
		WHERE(t.Contacts.ID.EQ(mysql.String(contact.ID)))
	db := qrm.Executable(r.DB)
	if tx != nil {
		db = tx
	}
	_, err := stmt.ExecContext(ctx, db)
	return err
}
