package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/doug-martin/goqu/v9"
	"github.com/jmoiron/sqlx/types"
	"github.com/raystack/frontier/billing/plan"
	"github.com/raystack/frontier/pkg/db"
)

type Plan struct {
	ID string `db:"id"`

	Name        string  `db:"name"`
	Title       *string `db:"title"`
	Description *string `db:"description"`
	Interval    *string `db:"interval"`

	State    string             `db:"state"`
	Metadata types.NullJSONText `db:"metadata"`

	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}

func (c Plan) transform() (plan.Plan, error) {
	var unmarshalledMetadata map[string]any
	if c.Metadata.Valid {
		if err := c.Metadata.Unmarshal(&unmarshalledMetadata); err != nil {
			return plan.Plan{}, err
		}
	}
	planTitle := ""
	if c.Title != nil {
		planTitle = *c.Title
	}
	planDescription := ""
	if c.Description != nil {
		planDescription = *c.Description
	}
	recurringInterval := ""
	if c.Interval != nil {
		recurringInterval = *c.Interval
	}
	return plan.Plan{
		ID:          c.ID,
		Name:        c.Name,
		Title:       planTitle,
		Description: planDescription,
		Interval:    recurringInterval,
		State:       c.State,
		Metadata:    unmarshalledMetadata,
		CreatedAt:   c.CreatedAt,
		UpdatedAt:   c.UpdatedAt,
		DeletedAt:   c.DeletedAt,
	}, nil
}

type BillingPlanRepository struct {
	dbc *db.Client
}

func NewBillingPlanRepository(dbc *db.Client) *BillingPlanRepository {
	return &BillingPlanRepository{
		dbc: dbc,
	}
}

func (r BillingPlanRepository) Create(ctx context.Context, toCreate plan.Plan) (plan.Plan, error) {
	if toCreate.Metadata == nil {
		toCreate.Metadata = make(map[string]any)
	}
	marshaledMetadata, err := json.Marshal(toCreate.Metadata)
	if err != nil {
		return plan.Plan{}, err
	}
	if toCreate.ID == "" {
		toCreate.ID = uuid.New().String()
	}

	query, params, err := dialect.Insert(TABLE_BILLING_PLANS).Rows(
		goqu.Record{
			"id":          toCreate.ID,
			"name":        toCreate.Name,
			"title":       toCreate.Title,
			"description": toCreate.Description,
			"interval":    toCreate.Interval,
			"state":       toCreate.State,
			"metadata":    marshaledMetadata,
			"updated_at":  goqu.L("now()"),
		}).Returning(&Plan{}).ToSQL()
	if err != nil {
		return plan.Plan{}, fmt.Errorf("%w: %s", parseErr, err)
	}

	var planModel Plan
	if err = r.dbc.WithTimeout(ctx, TABLE_BILLING_PLANS, "Create", func(ctx context.Context) error {
		return r.dbc.QueryRowxContext(ctx, query, params...).StructScan(&planModel)
	}); err != nil {
		return plan.Plan{}, fmt.Errorf("%w: %s", dbErr, err)
	}

	return planModel.transform()
}

func (r BillingPlanRepository) GetByID(ctx context.Context, id string) (plan.Plan, error) {
	stmt := dialect.Select().From(TABLE_BILLING_PLANS).Where(goqu.Ex{
		"id": id,
	})
	query, params, err := stmt.ToSQL()
	if err != nil {
		return plan.Plan{}, fmt.Errorf("%w: %s", parseErr, err)
	}

	var planModel Plan
	if err = r.dbc.WithTimeout(ctx, TABLE_BILLING_PLANS, "GetByID", func(ctx context.Context) error {
		return r.dbc.QueryRowxContext(ctx, query, params...).StructScan(&planModel)
	}); err != nil {
		err = checkPostgresError(err)
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return plan.Plan{}, plan.ErrNotFound
		}
		return plan.Plan{}, fmt.Errorf("%w: %s", dbErr, err)
	}

	return planModel.transform()
}

func (r BillingPlanRepository) GetByName(ctx context.Context, name string) (plan.Plan, error) {
	stmt := dialect.Select().From(TABLE_BILLING_PLANS).Where(goqu.Ex{
		"name": name,
	})
	query, params, err := stmt.ToSQL()
	if err != nil {
		return plan.Plan{}, fmt.Errorf("%w: %s", parseErr, err)
	}

	var planModel Plan
	if err = r.dbc.WithTimeout(ctx, TABLE_BILLING_PLANS, "GetByName", func(ctx context.Context) error {
		return r.dbc.QueryRowxContext(ctx, query, params...).StructScan(&planModel)
	}); err != nil {
		err = checkPostgresError(err)
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return plan.Plan{}, plan.ErrNotFound
		}
		return plan.Plan{}, fmt.Errorf("%w: %s", dbErr, err)
	}

	return planModel.transform()
}

func (r BillingPlanRepository) UpdateByName(ctx context.Context, toUpdate plan.Plan) (plan.Plan, error) {
	if strings.TrimSpace(toUpdate.Name) == "" {
		return plan.Plan{}, plan.ErrInvalidName
	}
	marshaledMetadata, err := json.Marshal(toUpdate.Metadata)
	if err != nil {
		return plan.Plan{}, fmt.Errorf("%w: %s", parseErr, err)
	}

	query, params, err := dialect.Update(TABLE_BILLING_PLANS).Set(
		goqu.Record{
			"title":       toUpdate.Title,
			"description": toUpdate.Description,
			"metadata":    marshaledMetadata,
			"updated_at":  goqu.L("now()"),
		}).Where(goqu.Ex{
		"name": toUpdate.Name,
	}).Returning(&Plan{}).ToSQL()
	if err != nil {
		return plan.Plan{}, fmt.Errorf("%w: %s", queryErr, err)
	}

	var planModel Plan
	if err = r.dbc.WithTimeout(ctx, TABLE_BILLING_PLANS, "UpdateByName", func(ctx context.Context) error {
		return r.dbc.QueryRowxContext(ctx, query, params...).StructScan(&planModel)
	}); err != nil {
		err = checkPostgresError(err)
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return plan.Plan{}, plan.ErrNotFound
		case errors.Is(err, ErrInvalidTextRepresentation):
			return plan.Plan{}, plan.ErrInvalidUUID
		default:
			return plan.Plan{}, fmt.Errorf("%s: %w", txnErr, err)
		}
	}

	return planModel.transform()
}

func (r BillingPlanRepository) List(ctx context.Context, filter plan.Filter) ([]plan.Plan, error) {
	stmt := dialect.Select().From(TABLE_BILLING_PLANS)
	query, params, err := stmt.ToSQL()
	if err != nil {
		return nil, fmt.Errorf("%w: %s", parseErr, err)
	}

	var planModels []Plan
	if err = r.dbc.WithTimeout(ctx, TABLE_BILLING_PLANS, "List", func(ctx context.Context) error {
		return r.dbc.SelectContext(ctx, &planModels, query, params...)
	}); err != nil {
		return nil, fmt.Errorf("%w: %s", dbErr, err)
	}

	var plans []plan.Plan
	for _, planModel := range planModels {
		plan, err := planModel.transform()
		if err != nil {
			return nil, err
		}
		plans = append(plans, plan)
	}
	return plans, nil
}