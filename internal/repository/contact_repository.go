package repository

import (
	"context"

	"golang-clean-architecture/ent"
	"golang-clean-architecture/ent/contact"
	"golang-clean-architecture/internal/dto"

	"github.com/sirupsen/logrus"
)

type ContactRepository struct {
	Log *logrus.Logger
}

func NewContactRepository(log *logrus.Logger) *ContactRepository {
	return &ContactRepository{Log: log}
}

func (r *ContactRepository) FindByIdAndUserId(ctx context.Context, client *ent.Client, id string, userId string) (*ent.Contact, error) {
	c, err := client.Contact.Query().
		Where(contact.ID(id), contact.UserID(userId)).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return c, nil
}

func (r *ContactRepository) Search(ctx context.Context, client *ent.Client, request *dto.SearchContactRequest) ([]*ent.Contact, int64, error) {
	if request.Page < 1 {
		request.Page = 1
	}
	if request.Size < 1 {
		request.Size = 10
	}

	query := client.Contact.Query().Where(contact.UserID(request.UserId))

	if request.Name != "" {
		query = query.Where(
			contact.Or(
				contact.FirstNameContainsFold(request.Name),
				contact.LastNameContainsFold(request.Name),
			),
		)
	}
	if request.Phone != "" {
		query = query.Where(contact.PhoneContainsFold(request.Phone))
	}
	if request.Email != "" {
		query = query.Where(contact.EmailContainsFold(request.Email))
	}

	total, err := query.Count(ctx)
	if err != nil {
		r.Log.WithError(err).Error("Failed to count contacts")
		return nil, 0, err
	}

	offset := (request.Page - 1) * request.Size
	contacts, err := query.Limit(request.Size).Offset(offset).All(ctx)
	if err != nil {
		r.Log.WithError(err).Error("Failed to search contacts")
		return nil, 0, err
	}

	return contacts, int64(total), nil
}

func (r *ContactRepository) Create(ctx context.Context, client *ent.Client, c *ent.Contact) error {
	builder := client.Contact.Create().
		SetID(c.ID).
		SetFirstName(c.FirstName).
		SetUserID(c.UserID).
		SetCreatedAt(c.CreatedAt).
		SetUpdatedAt(c.UpdatedAt)

	if c.LastName != nil {
		builder = builder.SetNillableLastName(c.LastName)
	}
	if c.Email != nil {
		builder = builder.SetNillableEmail(c.Email)
	}
	if c.Phone != nil {
		builder = builder.SetNillablePhone(c.Phone)
	}

	_, err := builder.Save(ctx)
	if err != nil {
		r.Log.WithError(err).WithField("contact_id", c.ID).Error("Failed to create contact")
		return err
	}
	r.Log.WithField("contact_id", c.ID).Debug("Contact created successfully")
	return nil
}

func (r *ContactRepository) Update(ctx context.Context, client *ent.Client, c *ent.Contact) error {
	builder := client.Contact.UpdateOneID(c.ID).
		SetFirstName(c.FirstName).
		SetUpdatedAt(c.UpdatedAt)

	if c.LastName != nil {
		builder = builder.SetLastName(*c.LastName)
	} else {
		builder = builder.ClearLastName()
	}
	if c.Email != nil {
		builder = builder.SetEmail(*c.Email)
	} else {
		builder = builder.ClearEmail()
	}
	if c.Phone != nil {
		builder = builder.SetPhone(*c.Phone)
	} else {
		builder = builder.ClearPhone()
	}

	_, err := builder.Save(ctx)
	if err != nil {
		r.Log.WithError(err).WithField("contact_id", c.ID).Error("Failed to update contact")
		return err
	}
	r.Log.WithField("contact_id", c.ID).Debug("Contact updated successfully")
	return nil
}

func (r *ContactRepository) Delete(ctx context.Context, client *ent.Client, id string) error {
	err := client.Contact.DeleteOneID(id).Exec(ctx)
	if err != nil {
		r.Log.WithError(err).WithField("contact_id", id).Error("Failed to delete contact")
		return err
	}
	r.Log.WithField("contact_id", id).Debug("Contact deleted successfully")
	return nil
}
