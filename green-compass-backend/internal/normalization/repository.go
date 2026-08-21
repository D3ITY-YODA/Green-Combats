package normalization

import (
	"context"

	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrNotFound = errors.New("normalized observation not found")
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) Store(ctx context.Context, obs []CanonicalObservation) error {
	if len(obs) == 0 {
		return nil
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	query := `
		INSERT INTO normalized_observations (
			id, source_key, dataset_key, topic_key, variable, geometry,
			place_id, value, text_value, unit, observed_at, valid_from, valid_until,
			retrieved_at, is_forecast, quality_status, source_version, raw_asset_ref, license
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19)
		ON CONFLICT DO NOTHING
	`

	for _, o := range obs {
		var geom interface{}
		if o.Geometry != nil {
			geom = o.Geometry
		}

		var topicKey string
		if o.TopicKey != "" {
			topicKey = o.TopicKey
		} else {
			topicKey = inferTopicKey(o.Variable)
		}

		_, err := tx.Exec(ctx, query,
			o.ID, o.SourceCode, "", topicKey, o.Variable, geom,
			o.PlaceID, o.Value, nil, o.Unit, o.ObservedAt, nil, nil,
			time.Now().UTC(), false, "accepted", "", "", "",
		)
		if err != nil {
			return fmt.Errorf("insert normalized observation: %w", err)
		}
	}

	return tx.Commit(ctx)
}

func (r *Repository) ListByPlace(ctx context.Context, placeID uuid.UUID, from, to time.Time) ([]CanonicalObservation, error) {
	query := `
		SELECT id, source_key, dataset_key, topic_key, variable,
		       ST_AsText(geometry) as geometry_wkt,
		       place_id, value, text_value, unit, observed_at, valid_from, valid_until,
		       retrieved_at, is_forecast, quality_status, source_version, raw_asset_ref, license
		FROM normalized_observations
		WHERE place_id = $1 AND observed_at >= $2 AND observed_at <= $3
		ORDER BY observed_at ASC
	`

	rows, err := r.pool.Query(ctx, query, placeID, from, to)
	if err != nil {
		return nil, fmt.Errorf("query normalized observations: %w", err)
	}
	defer rows.Close()

	var results []CanonicalObservation
	for rows.Next() {
		var o CanonicalObservation
		var geomWkt *string
		err := rows.Scan(
			&o.ID, &o.SourceCode, &o.DatasetKey, &o.TopicKey, &o.Variable,
			&geomWkt, &o.PlaceID, &o.Value, &o.TextValue, &o.Unit,
			&o.ObservedAt, &o.ValidFrom, &o.ValidUntil,
			&o.RetrievedAt, &o.IsForecast, &o.QualityStatus, &o.SourceVersion, &o.RawAssetRef, &o.License,
		)
		if err != nil {
			return nil, fmt.Errorf("scan observation: %w", err)
		}
		results = append(results, o)
	}

	return results, rows.Err()
}

func (r *Repository) ListBySourceAndPlace(ctx context.Context, sourceKey, datasetKey string, placeID uuid.UUID, from, to time.Time) ([]CanonicalObservation, error) {
	query := `
		SELECT id, source_key, dataset_key, topic_key, variable,
		       ST_AsText(geometry) as geometry_wkt,
		       place_id, value, text_value, unit, observed_at, valid_from, valid_until,
		       retrieved_at, is_forecast, quality_status, source_version, raw_asset_ref, license
		FROM normalized_observations
		WHERE source_key = $1 AND dataset_key = $2 AND place_id = $3 AND observed_at >= $4 AND observed_at <= $5
		ORDER BY observed_at ASC
	`

	rows, err := r.pool.Query(ctx, query, sourceKey, datasetKey, placeID, from, to)
	if err != nil {
		return nil, fmt.Errorf("query normalized observations: %w", err)
	}
	defer rows.Close()

	var results []CanonicalObservation
	for rows.Next() {
		var o CanonicalObservation
		var geomWkt *string
		err := rows.Scan(
			&o.ID, &o.SourceCode, &o.DatasetKey, &o.TopicKey, &o.Variable,
			&geomWkt, &o.PlaceID, &o.Value, &o.TextValue, &o.Unit,
			&o.ObservedAt, &o.ValidFrom, &o.ValidUntil,
			&o.RetrievedAt, &o.IsForecast, &o.QualityStatus, &o.SourceVersion, &o.RawAssetRef, &o.License,
		)
		if err != nil {
			return nil, fmt.Errorf("scan observation: %w", err)
		}
		results = append(results, o)
	}

	return results, rows.Err()
}

func (r *Repository) ListByTopicAndPlace(ctx context.Context, topicKey string, placeID uuid.UUID, from, to time.Time) ([]CanonicalObservation, error) {
	query := `
		SELECT id, source_key, dataset_key, topic_key, variable,
		       ST_AsText(geometry) as geometry_wkt,
		       place_id, value, text_value, unit, observed_at, valid_from, valid_until,
		       retrieved_at, is_forecast, quality_status, source_version, raw_asset_ref, license
		FROM normalized_observations
		WHERE topic_key = $1 AND place_id = $2 AND observed_at >= $3 AND observed_at <= $4
		ORDER BY observed_at ASC
	`

	rows, err := r.pool.Query(ctx, query, topicKey, placeID, from, to)
	if err != nil {
		return nil, fmt.Errorf("query normalized observations: %w", err)
	}
	defer rows.Close()

	var results []CanonicalObservation
	for rows.Next() {
		var o CanonicalObservation
		var geomWkt *string
		err := rows.Scan(
			&o.ID, &o.SourceCode, &o.DatasetKey, &o.TopicKey, &o.Variable,
			&geomWkt, &o.PlaceID, &o.Value, &o.TextValue, &o.Unit,
			&o.ObservedAt, &o.ValidFrom, &o.ValidUntil,
			&o.RetrievedAt, &o.IsForecast, &o.QualityStatus, &o.SourceVersion, &o.RawAssetRef, &o.License,
		)
		if err != nil {
			return nil, fmt.Errorf("scan observation: %w", err)
		}
		results = append(results, o)
	}

	return results, rows.Err()
}

func inferTopicKey(variable string) string {
	switch variable {
	case "temperature_celsius", "precipitation_mm", "relative_humidity_percent", "wind_speed_ms", "weather_code":
		return "local_outlook"
	case "soil_moisture_m3m3", "solar_radiation_mj":
		return "food_and_agriculture"
	default:
		return "local_outlook"
	}
}