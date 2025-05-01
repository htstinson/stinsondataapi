package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/htstinson/stinsondataapi/api/internal/model"
)

// Customers
func (d *Database) GetProfile(ctx context.Context, id string) (*model.Profile, error) {
	fmt.Println("d GetProfile")

	var profile model.Profile

	err := d.DB.QueryRowContext(ctx,
		"SELECT id, parent_id, created_at, modified_at FROM profile WHERE id = $1",
		id,
	).Scan(&profile.Id, &profile.ParentId, &profile.CreatedAt, &profile.ModifiedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error getting customer: %w", err)
	}

	return &profile, nil
}

func (d *Database) GetProfileByParent(ctx context.Context, id string) (*model.Profile, error) {
	fmt.Println("d GetProfileByParent")

	var profile model.Profile

	err := d.DB.QueryRowContext(ctx,
		"SELECT id, parent_id, created_at, modified_at FROM profile WHERE parent_id = $1",
		id,
	).Scan(&profile.Id, &profile.ParentId, &profile.CreatedAt, &profile.ModifiedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error getting customer: %w", err)
	}

	return &profile, nil
}

func (d *Database) CreateProfile(ctx context.Context, schema_name string, parent_id string) (*model.Profile, error) {
	fmt.Println("d CreateProfile")

	profile := &model.Profile{
		Id:         uuid.New().String(),
		ParentId:   parent_id,
		CreatedAt:  time.Now(),
		ModifiedAt: time.Now(),
	}

	query := fmt.Sprintf(`INSERT INTO %s.profile (id, parent_id, created_at, modified_at) VALUES ($1, $2, $3, $4)`, schema_name)

	fmt.Println(query)

	_, err := d.DB.ExecContext(ctx, query, profile.Id, profile.ParentId, profile.CreatedAt, profile.ModifiedAt)
	if err != nil {
		return nil, fmt.Errorf("error creating profile: %w", err)
	}

	query = `UPDATE common.customers SET schema_name = $1 WHERE id = $2`

	fmt.Println(query)

	_, err = d.DB.ExecContext(ctx, query, schema_name, parent_id)
	if err != nil {
		return nil, fmt.Errorf("error updating parent: %w", err)
	}

	return profile, nil

}

func (d *Database) SelectProfiles(ctx context.Context, limit, offset int) ([]model.Profile, error) {

	fmt.Println("database.go SelectProfiles()")

	rows, err := d.DB.QueryContext(ctx,
		"SELECT id, parent_id, created_at, modified_at FROM profiles ORDER BY created_at LIMIT $1 OFFSET $2",
		limit, offset,
	)
	if err != nil {
		fmt.Println(err.Error())
		return nil, fmt.Errorf("error listing profiles: %w", err)
	}
	defer rows.Close()

	var profiles []model.Profile
	for rows.Next() {
		var profile model.Profile
		if err := rows.Scan(&profile.Id, &profile.ParentId, &profile.CreatedAt, profile.ModifiedAt); err != nil {
			fmt.Println(err.Error())
			return nil, fmt.Errorf("error scanning profile: %w", err)
		}

		profiles = append(profiles, profile)
	}
	return profiles, nil
}

func (d *Database) UpdateProfile(ctx context.Context, profile *model.Profile) error {
	fmt.Println("d UpdateProfile")

	query := `UPDATE profiles SET parent_id = $1 WHERE id = $2`

	_, err := d.DB.ExecContext(ctx, query, profile.ParentId, profile.Id)

	return err
}

func (d *Database) DeleteProfile(ctx context.Context, id string) error {
	fmt.Println("d DeleteProfile")

	query := `DELETE FROM profiles WHERE id = $1`

	_, err := d.DB.ExecContext(ctx, query, id)

	return err
}
