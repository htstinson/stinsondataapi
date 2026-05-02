package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/htstinson/stinsondataapi/api/internal/model"
)

/*
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
*/

func (d *Database) GetProfile(ctx context.Context, subscriber *model.Subscriber) (*model.Profile, error) {
	fmt.Println("d GetProfile")

	fmt.Println("subscriber.Id", &subscriber.Id)

	var profile model.Profile

	query := fmt.Sprintf(`SELECT id, parent_id, created_at, modified_at, legal_name, phone, fax, website, email, linkedin, facebook, instagram, x, youtube, pinterest, google_business, yelp, glassdoor, github, nextdoor, bizapedia FROM %s.profile WHERE parent_id = $1`, subscriber.Schema_Name)

	err := d.DB.QueryRowContext(ctx,
		query,
		&subscriber.Id,
	).Scan(&profile.Id, &profile.ParentId, &profile.CreatedAt, &profile.ModifiedAt,
		&profile.Legal_Name, &profile.Phone, &profile.Fax, &profile.Email, &profile.Website, &profile.LinkedIn, &profile.Facebook,
		&profile.Instagram, &profile.X, &profile.YouTube, &profile.Pinterest, &profile.GoogleBusiness,
		&profile.Yelp, &profile.GlassDoor, &profile.Github, &profile.NextDoor, &profile.Bizapedia)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		fmt.Println(err.Error())
		return nil, fmt.Errorf("error getting profile: %w", err)
	}

	return &profile, nil
}

func (d *Database) CreateProfile(ctx context.Context, subscriber model.Subscriber, profile model.Profile) (*model.Profile, error) {
	fmt.Println("d CreateProfile")

	profile.Id = uuid.New().String()
	//profile.ParentId = subscriber.Id
	profile.CreatedAt = time.Now()
	profile.ModifiedAt = time.Now()

	query := fmt.Sprintf(`INSERT INTO %s.profile (id, parent_id, created_at, modified_at, 
		legal_name, phone, fax, email, website, linkedin, facebook, instagram, x, youtube, pinterest, 
		google_business, yelp, glassdoor, github, nextdoor, bizapedia
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21)`, subscriber.Schema_Name)

	_, err := d.DB.ExecContext(ctx, query, profile.Id, profile.ParentId, profile.CreatedAt, profile.ModifiedAt,
		profile.Legal_Name, profile.Phone, profile.Fax, profile.Email, profile.Website, profile.LinkedIn,
		profile.Facebook, profile.Instagram, profile.X, profile.YouTube, profile.Pinterest,
		profile.GoogleBusiness, profile.Yelp, profile.GlassDoor, profile.Github, profile.NextDoor, profile.Bizapedia)
	if err != nil {
		return nil, fmt.Errorf("error creating profile: %w", err)
	}

	query = `UPDATE common.subscribers SET schema_name = $1 WHERE id = $2`

	_, err = d.DB.ExecContext(ctx, query, subscriber.Schema_Name, subscriber.Id)
	if err != nil {
		return nil, fmt.Errorf("error updating parent: %w", err)
	}

	return &profile, nil

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
