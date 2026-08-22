package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"green-compass-backend/pkg/security"
)

const (
	passwordHashCost = 12
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("Seed failed: %v", err)
	}
	log.Println("Seed completed successfully!")
}

func run() error {
	ctx := context.Background()

	dbURL := getEnv("DATABASE_URL", "postgres://greencompass:greencompass@localhost:5432/greencompass?sslmode=disable")

	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		return fmt.Errorf("connect to database: %w", err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		return fmt.Errorf("ping database: %w", err)
	}

	log.Println("Connected to database, starting seed...")

	tx, err := pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	userID, err := seedUsers(ctx, tx)
	if err != nil {
		return fmt.Errorf("seed users: %w", err)
	}
	log.Printf("Created user: %s", userID)

	adminID, err := seedAdminUser(ctx, tx)
	if err != nil {
		return fmt.Errorf("seed admin: %w", err)
	}
	log.Printf("Created admin: %s", adminID)

	placeIDs, err := seedPlaces(ctx, tx)
	if err != nil {
		return fmt.Errorf("seed places: %w", err)
	}
	log.Printf("Created %d places", len(placeIDs))

	err = seedUserSavedPlaces(ctx, tx, userID, placeIDs)
	if err != nil {
		return fmt.Errorf("seed user saved places: %w", err)
	}
	log.Println("Linked user to places")

	err = seedIndicatorDefinitions(ctx, tx)
	if err != nil {
		return fmt.Errorf("seed indicator definitions: %w", err)
	}
	log.Println("Created indicator definitions")

	indicatorIDs, err := getIndicatorIDs(ctx, tx)
	if err != nil {
		return fmt.Errorf("get indicator IDs: %w", err)
	}

	err = seedPlaceIndicators(ctx, tx, placeIDs, indicatorIDs)
	if err != nil {
		return fmt.Errorf("seed place indicators: %w", err)
	}
	log.Println("Created place indicators")

	err = seedPlaceContent(ctx, tx, placeIDs)
	if err != nil {
		return fmt.Errorf("seed place content: %w", err)
	}
	log.Println("Created place content (today/forecast/alert)")

	err = seedObservations(ctx, tx, userID, placeIDs)
	if err != nil {
		return fmt.Errorf("seed observations: %w", err)
	}
	log.Println("Created observations (community reports)")

	err = seedPlaceAssessments(ctx, tx, placeIDs)
	if err != nil {
		return fmt.Errorf("seed place assessments: %w", err)
	}
	log.Println("Created place assessments")

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func seedUsers(ctx context.Context, tx pgx.Tx) (uuid.UUID, error) {
	userID := uuid.New()
	passwordHash, err := hashPassword("password123")
	if err != nil {
		return uuid.Nil, err
	}

	// Check if user exists
	var existingID uuid.UUID
	err = tx.QueryRow(ctx, `SELECT id FROM users WHERE email = $1`, "demo@greencompass.app").Scan(&existingID)
	if err == nil {
		// User exists, update
		_, err = tx.Exec(ctx, `
			UPDATE users SET password_hash = $1, display_name = $2, language = $3, updated_at = now()
			WHERE email = $4
		`, passwordHash, "Demo User", "en", "demo@greencompass.app")
		if err != nil {
			return uuid.Nil, err
		}
		return existingID, nil
	}

	// Insert new user
	_, err = tx.Exec(ctx, `
		INSERT INTO users (id, phone_number, email, password_hash, display_name, language)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, userID, "+15551234567", "demo@greencompass.app", passwordHash, "Demo User", "en")

	if err != nil {
		return uuid.Nil, err
	}

	return userID, nil
}

func seedAdminUser(ctx context.Context, tx pgx.Tx) (uuid.UUID, error) {
	adminID := uuid.New()
	passwordHash, err := hashPassword("admin123")
	if err != nil {
		return uuid.Nil, err
	}

	var existingID uuid.UUID
	err = tx.QueryRow(ctx, `SELECT id FROM users WHERE email = $1`, "admin@greencompass.app").Scan(&existingID)
	if err == nil {
		_, err = tx.Exec(ctx, `
			UPDATE users SET password_hash = $1, display_name = $2, language = $3, is_platform_admin = $4, updated_at = now()
			WHERE email = $5
		`, passwordHash, "Platform Admin", "en", true, "admin@greencompass.app")
		if err != nil {
			return uuid.Nil, err
		}
		return existingID, nil
	}

	_, err = tx.Exec(ctx, `
		INSERT INTO users (id, phone_number, email, password_hash, display_name, language, is_platform_admin)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, adminID, "+15550000000", "admin@greencompass.app", passwordHash, "Platform Admin", "en", true)

	if err != nil {
		return uuid.Nil, err
	}

	return adminID, nil
}

func seedPlaces(ctx context.Context, tx pgx.Tx) ([]uuid.UUID, error) {
	places := []struct {
		id        uuid.UUID
		name      string
		placeType string
		lat       float64
		lon       float64
		extCode   string
	}{
		{uuid.New(), "Lower Valley", "community", -1.2921, 36.8219, "KE-LOWER-VALLEY"},
		{uuid.New(), "Riverside", "community", -1.3000, 36.8300, "KE-RIVERSIDE"},
		{uuid.New(), "Upper Highland", "community", -1.2800, 36.8100, "KE-UPPER-HIGHLAND"},
		{uuid.New(), "Greenfields", "community", -1.2700, 36.8000, "KE-GREENFIELDS"},
		{uuid.New(), "Lakeview", "community", -1.3100, 36.8400, "KE-LAKEVIEW"},
		{uuid.New(), "Hilltop", "ward", -1.2600, 36.7900, "KE-HILLTOP"},
	}

	placeIDs := make([]uuid.UUID, 0, len(places))
	for _, p := range places {
		_, err := tx.Exec(ctx, `
			INSERT INTO places (id, name, place_type, location, external_code)
			VALUES ($1, $2, $3, ST_SetSRID(ST_MakePoint($4, $5), 4326)::geography, $6)
			ON CONFLICT (id) DO NOTHING
		`, p.id, p.name, p.placeType, p.lon, p.lat, p.extCode)
		if err != nil {
			return nil, err
		}
		placeIDs = append(placeIDs, p.id)
	}

	return placeIDs, nil
}

func seedUserSavedPlaces(ctx context.Context, tx pgx.Tx, userID uuid.UUID, placeIDs []uuid.UUID) error {
	// First, set all existing saved places to not primary
	_, err := tx.Exec(ctx, `
		UPDATE user_saved_places SET is_primary = FALSE WHERE user_id = $1
	`, userID)
	if err != nil {
		return err
	}

	for i, placeID := range placeIDs {
		isPrimary := i == 0
		label := fmt.Sprintf("My %s", getPlaceName(ctx, tx, placeID))
		_, err := tx.Exec(ctx, `
			INSERT INTO user_saved_places (user_id, place_id, label, is_primary)
			VALUES ($1, $2, $3, $4)
			ON CONFLICT (user_id, place_id) DO UPDATE SET
				label = EXCLUDED.label,
				is_primary = EXCLUDED.is_primary
		`, userID, placeID, label, isPrimary)
		if err != nil {
			return err
		}
	}
	return nil
}

func getPlaceName(ctx context.Context, tx pgx.Tx, placeID uuid.UUID) string {
	var name string
	err := tx.QueryRow(ctx, "SELECT name FROM places WHERE id = $1", placeID).Scan(&name)
	if err != nil {
		return "Place"
	}
	return name
}

func seedIndicatorDefinitions(ctx context.Context, tx pgx.Tx) error {
	_, err := tx.Exec(ctx, `
		INSERT INTO indicator_definitions (id, code, display_name, category, source_variables, description)
		VALUES 
		(gen_random_uuid(), 'rain_intensity_trend', 'Rain Intensity Trend', 'weather', ARRAY['precipitation_mm'], 'Trend in precipitation over the observation period'),
		(gen_random_uuid(), 'rain_frequency', 'Rain Frequency', 'weather', ARRAY['precipitation_mm'], 'Number of rainy days in the observation period'),
		(gen_random_uuid(), 'temperature_trend', 'Temperature Trend', 'weather', ARRAY['temperature_celsius'], 'Mean temperature trend across the period'),
		(gen_random_uuid(), 'wind_speed_trend', 'Wind Speed Trend', 'weather', ARRAY['wind_speed_ms'], 'Wind speed trend across the period'),
		(gen_random_uuid(), 'soil_moisture_level', 'Soil Moisture Level', 'agriculture', ARRAY['soil_moisture_m3m3'], 'Current soil moisture as percentage of capacity'),
		(gen_random_uuid(), 'drought_stress', 'Drought Stress Index', 'agriculture', ARRAY['precipitation_mm', 'temperature_celsius', 'soil_moisture_m3m3'], 'Composite drought risk indicator'),
		(gen_random_uuid(), 'water_quality_concern', 'Water Quality Concern', 'water', ARRAY[]::TEXT[], 'Aggregate water quality assessment'),
		(gen_random_uuid(), 'air_quality_index', 'Air Quality Index (AQI)', 'air_quality', ARRAY[]::TEXT[], 'Composite air quality from PM2.5 and Ozone')
		ON CONFLICT (code) DO NOTHING
	`)
	return err
}

func getIndicatorIDs(ctx context.Context, tx pgx.Tx) (map[string]uuid.UUID, error) {
	rows, err := tx.Query(ctx, `SELECT id, code FROM indicator_definitions`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ids := make(map[string]uuid.UUID)
	for rows.Next() {
		var id uuid.UUID
		var code string
		if err := rows.Scan(&id, &code); err != nil {
			return nil, err
		}
		ids[code] = id
	}
	return ids, rows.Err()
}

func seedPlaceIndicators(ctx context.Context, tx pgx.Tx, placeIDs []uuid.UUID, indicatorIDs map[string]uuid.UUID) error {
	now := time.Now()
	periodStart := now.AddDate(0, 0, -7)
	periodEnd := now

	data := []struct {
		placeIdx       int
		indicatorCode  string
		value          float64
		unit           string
		trend          *string
		trendConf      *float64
		dataPoints     int
	}{
		{0, "rain_intensity_trend", 12.5, "mm/day", strPtr("increasing"), floatPtr(0.75), 7},
		{0, "rain_frequency", 4, "days/week", strPtr("stable"), floatPtr(0.6), 7},
		{0, "temperature_trend", 24.3, "°C", strPtr("stable"), floatPtr(0.5), 7},
		{0, "soil_moisture_level", 65, "%", strPtr("decreasing"), floatPtr(0.8), 7},
		{0, "drought_stress", 35, "index", strPtr("increasing"), floatPtr(0.7), 7},
		{0, "water_quality_concern", 2, "level", strPtr("stable"), floatPtr(0.5), 3},

		{1, "rain_intensity_trend", 8.2, "mm/day", strPtr("stable"), floatPtr(0.5), 7},
		{1, "rain_frequency", 3, "days/week", strPtr("stable"), floatPtr(0.4), 7},
		{1, "temperature_trend", 26.1, "°C", strPtr("increasing"), floatPtr(0.6), 7},
		{1, "soil_moisture_level", 58, "%", strPtr("stable"), floatPtr(0.5), 7},
		{1, "drought_stress", 42, "index", strPtr("stable"), floatPtr(0.5), 7},

		{2, "rain_intensity_trend", 15.8, "mm/day", strPtr("increasing"), floatPtr(0.85), 7},
		{2, "rain_frequency", 5, "days/week", strPtr("increasing"), floatPtr(0.8), 7},
		{2, "temperature_trend", 22.8, "°C", strPtr("stable"), floatPtr(0.4), 7},
		{2, "soil_moisture_level", 72, "%", strPtr("increasing"), floatPtr(0.7), 7},
		{2, "drought_stress", 28, "index", strPtr("decreasing"), floatPtr(0.75), 7},
		{2, "water_quality_concern", 1, "level", strPtr("stable"), floatPtr(0.4), 3},

		{3, "rain_intensity_trend", 6.1, "mm/day", strPtr("decreasing"), floatPtr(0.7), 7},
		{3, "rain_frequency", 2, "days/week", strPtr("decreasing"), floatPtr(0.65), 7},
		{3, "temperature_trend", 27.5, "°C", strPtr("increasing"), floatPtr(0.8), 7},
		{3, "soil_moisture_level", 45, "%", strPtr("decreasing"), floatPtr(0.85), 7},
		{3, "drought_stress", 55, "index", strPtr("increasing"), floatPtr(0.8), 7},

		{4, "rain_intensity_trend", 18.3, "mm/day", strPtr("increasing"), floatPtr(0.9), 7},
		{4, "rain_frequency", 6, "days/week", strPtr("increasing"), floatPtr(0.85), 7},
		{4, "water_quality_concern", 3, "level", strPtr("increasing"), floatPtr(0.75), 3},
		{4, "temperature_trend", 23.5, "°C", strPtr("stable"), floatPtr(0.5), 7},

		{5, "rain_intensity_trend", 9.7, "mm/day", strPtr("stable"), floatPtr(0.5), 7},
		{5, "rain_frequency", 3, "days/week", strPtr("stable"), floatPtr(0.5), 7},
		{5, "temperature_trend", 25.0, "°C", strPtr("increasing"), floatPtr(0.6), 7},
		{5, "soil_moisture_level", 62, "%", strPtr("stable"), floatPtr(0.5), 7},
		{5, "drought_stress", 38, "index", strPtr("stable"), floatPtr(0.5), 7},
	}

	for _, d := range data {
		indicatorID, ok := indicatorIDs[d.indicatorCode]
		if !ok {
			continue
		}
		placeID := placeIDs[d.placeIdx]

		_, err := tx.Exec(ctx, `
			INSERT INTO place_indicators (id, place_id, indicator_id, computed_at, period_start, period_end, value, unit, trend, trend_confidence, data_points_count)
			VALUES (gen_random_uuid(), $1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
			ON CONFLICT (place_id, indicator_id, period_start, period_end) DO UPDATE SET
				value = EXCLUDED.value,
				unit = EXCLUDED.unit,
				trend = EXCLUDED.trend,
				trend_confidence = EXCLUDED.trend_confidence,
				data_points_count = EXCLUDED.data_points_count,
				computed_at = EXCLUDED.computed_at
		`, placeID, indicatorID, now, periodStart, periodEnd, d.value, d.unit, d.trend, d.trendConf, d.dataPoints)
		if err != nil {
			return err
		}
	}
	return nil
}

func seedPlaceContent(ctx context.Context, tx pgx.Tx, placeIDs []uuid.UUID) error {
	now := time.Now()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	todayEnd := todayStart.Add(24 * time.Hour)
	forecastStart := todayEnd
	forecastEnd := forecastStart.Add(48 * time.Hour)
	alertStart := now.Add(-2 * time.Hour)
	alertEnd := now.Add(6 * time.Hour)

	content := []struct {
		placeIdx     int
		contentType  string
		headline     string
		bodyText     string
		callToAction *string
		periodStart  time.Time
		periodEnd    time.Time
		sourceIndIDs []string
	}{
		// Today content
		{0, "today", "Lower Valley: Light Rain Expected This Afternoon",
			"Scattered showers will develop after 2 PM, with rainfall totals of 5-10 mm expected through the evening. Temperatures will remain mild at 22-25°C. Soil moisture is adequate for most crops at 65%.",
			strPtr("Carry an umbrella if heading out this afternoon. Check drainage around homestead."),
			todayStart, todayEnd, []string{"rain_intensity_trend", "rain_frequency", "soil_moisture_level"}},

		{1, "today", "Riverside: Warm and Dry Conditions Continue",
			"High pressure keeps conditions dry and warm through the day. Temperatures reaching 28°C. No significant rainfall expected. Soil moisture declining at 58% - consider irrigation for vegetable gardens.",
			strPtr("Water vegetable gardens early morning. Monitor soil moisture daily."),
			todayStart, todayEnd, []string{"rain_intensity_trend", "temperature_trend", "soil_moisture_level"}},

		{2, "today", "Upper Highland: Heavy Rain Warning - Flood Risk Elevated",
			"Persistent heavy rainfall overnight has saturated soils. Additional 15-25 mm expected today. Stream levels rising - low-lying areas near the river may experience localized flooding. Drought stress decreasing as soils recharge.",
			strPtr("Avoid river crossings. Move livestock to higher ground. Check emergency supplies."),
			todayStart, todayEnd, []string{"rain_intensity_trend", "rain_frequency", "drought_stress", "water_quality_concern"}},

		{3, "today", "Greenfields: Hot and Dry - Drought Stress Building",
			"Temperatures climbing to 29°C with no rain in the 7-day forecast. Soil moisture at 45% and declining rapidly. Drought stress index rising to 55. Pasture conditions deteriorating.",
			strPtr("Prioritize water for high-value crops. Consider supplementary feed for livestock. Check water storage."),
			todayStart, todayEnd, []string{"rain_intensity_trend", "temperature_trend", "soil_moisture_level", "drought_stress"}},

		{4, "today", "Lakeview: Lake Levels Rising - Water Quality Advisory",
			"Continued heavy inflows have raised lake levels by 30 cm in 48 hours. Turbidity increasing - water treatment recommended for domestic use. Fishing conditions difficult due to high flows.",
			strPtr("Boil or treat lake water for drinking. Avoid swimming near inflows. Secure boats and lakeside equipment."),
			todayStart, todayEnd, []string{"rain_intensity_trend", "rain_frequency", "water_quality_concern"}},

		{5, "today", "Hilltop: Seasonal Conditions - Normal for This Time of Year",
			"Typical late-season weather with isolated afternoon showers possible. Temperatures 24-26°C. Soil moisture adequate at 62%. No significant concerns for agriculture or water supply.",
			strPtr("Continue normal farming activities. Good conditions for planting if soil workable."),
			todayStart, todayEnd, []string{"rain_intensity_trend", "temperature_trend", "soil_moisture_level"}},

		// Forecast content
		{0, "forecast", "Lower Valley: 48-Hour Outlook - Continued Showers",
			"On-and-off showers persist through tomorrow with another 10-15 mm accumulation. Cloudy conditions keep temperatures moderate. Good soil recharge for upcoming planting window.",
			strPtr("Prepare fields for planting once soils are workable. Maintain drainage channels."),
			forecastStart, forecastEnd, []string{"rain_intensity_trend", "rain_frequency"}},

		{1, "forecast", "Riverside: 48-Hour Outlook - Hot and Dry Persisting",
			"High pressure ridge strengthens - near-zero rain chance for next 48 hours. Temperatures 29-31°C. Evaporation rates high - soil moisture will drop below 50%.",
			strPtr("Increase irrigation frequency. Mulch garden beds to retain moisture. Plan activities for cooler hours."),
			forecastStart, forecastEnd, []string{"rain_intensity_trend", "temperature_trend"}},

		{2, "forecast", "Upper Highland: 48-Hour Outlook - Rain Easing, Flood Risk Remains",
			"Rain intensity decreasing but grounds remain saturated. Additional 5-10 mm possible. River levels peak tomorrow then slowly recede. Flood risk remains elevated through Thursday.",
			strPtr("Continue flood precautions. Monitor river gauges. Do not attempt flooded road crossings."),
			forecastStart, forecastEnd, []string{"rain_intensity_trend", "drought_stress"}},

		// Alert content
		{2, "alert", "FLOOD ALERT: Upper Highland - River Level Critical",
			"River gauge at Upper Highland Bridge reading 4.2m (flood stage 3.8m). Rapid rise over past 6 hours. Evacuation advisory for Riverside Road and Valley Floor settlements. Emergency services activated.",
			strPtr("EVACUATE if in advisory zone. Move to designated shelters. Call 119 for emergency assistance."),
			alertStart, alertEnd, []string{"rain_intensity_trend", "water_quality_concern"}},

		{3, "alert", "DROUGHT ADVISORY: Greenfields - Severe Soil Moisture Deficit",
			"Soil moisture at 45% - lowest recorded in 5 years. Drought stress index 55 and climbing. No significant rainfall in 14-day outlook. Crop failure risk high for non-irrigated maize.",
			strPtr("Activate drought contingency plans. Contact agricultural extension for support. Conserve water aggressively."),
			alertStart, alertEnd, []string{"soil_moisture_level", "drought_stress"}},
	}

	for _, c := range content {
		placeID := placeIDs[c.placeIdx]

		var sourceInds []uuid.UUID
		if len(c.sourceIndIDs) > 0 {
			rows, err := tx.Query(ctx, `SELECT id FROM indicator_definitions WHERE code = ANY($1)`, c.sourceIndIDs)
			if err != nil {
				return err
			}
			for rows.Next() {
				var id uuid.UUID
				if err := rows.Scan(&id); err != nil {
					rows.Close()
					return err
				}
				sourceInds = append(sourceInds, id)
			}
			rows.Close()
		}

		var callToAction *string
		if c.callToAction != nil {
			callToAction = c.callToAction
		}

		_, err := tx.Exec(ctx, `
			INSERT INTO place_content (id, place_id, generated_at, period_start, period_end, content_type, language, headline, body_text, call_to_action, source_indicators)
			VALUES (gen_random_uuid(), $1, $2, $3, $4, $5, 'en', $6, $7, $8, $9)
			ON CONFLICT (place_id, content_type, period_start, period_end, language) DO UPDATE SET
				headline = EXCLUDED.headline,
				body_text = EXCLUDED.body_text,
				call_to_action = EXCLUDED.call_to_action,
				source_indicators = EXCLUDED.source_indicators,
				generated_at = EXCLUDED.generated_at
		`, placeID, now, c.periodStart, c.periodEnd, c.contentType, c.headline, c.bodyText, callToAction, sourceInds)
		if err != nil {
			return err
		}
	}
	return nil
}

func seedObservations(ctx context.Context, tx pgx.Tx, userID uuid.UUID, placeIDs []uuid.UUID) error {
	now := time.Now()

	observations := []struct {
		reporterID uuid.UUID
		placeIdx   int
		lat        *float64
		lon        *float64
		category   string
		description string
		status     string
		createdAt  time.Time
	}{
		{userID, 0, floatPtr(-1.2921), floatPtr(36.8219), "flood",
			"Water rising fast near the footbridge on Main Road. The culvert appears blocked with debris. Water is now crossing the road surface - vehicles cannot pass.",
			"verified", now.Add(-3 * time.Hour)},
		{userID, 0, floatPtr(-1.2930), floatPtr(36.8220), "water_quality",
			"River water looks very muddy and has an unusual smell. Not normal for this time of year. Suspect upstream disturbance.",
			"pending", now.Add(-1 * time.Hour)},
		{userID, 1, floatPtr(-1.3005), floatPtr(36.8305), "drought",
			"Maize crop showing severe wilting despite irrigation. Leaves curling, bottom leaves drying. Soil cracks visible in unirrigated sections. Need urgent advice.",
			"pending", now.Add(-6 * time.Hour)},
		{userID, 2, floatPtr(-1.2805), floatPtr(36.8105), "flood",
			"Stream behind the school has overflowed its banks. Water approaching the classroom block foundations. Sandbags deployed but may not hold if rain continues.",
			"verified", now.Add(-12 * time.Hour)},
		{userID, 2, floatPtr(-1.2810), floatPtr(36.8110), "crop_damage",
			"Banana plantation - 15 mature plants toppled by saturated soil and wind overnight. Significant crop loss. Need assessment for insurance.",
			"verified", now.Add(-24 * time.Hour)},
		{userID, 3, floatPtr(-1.2705), floatPtr(36.8005), "drought",
			"Borehole yield dropped significantly - pump running 30 min longer to fill same tank. Neighbor reports similar. Water table dropping fast.",
			"pending", now.Add(-8 * time.Hour)},
		{userID, 4, floatPtr(-1.3105), floatPtr(36.8405), "water_quality",
			"Dead fish washing up on shore near the fishing pier. About 20 tilapia and catfish since yesterday. Water has greenish tinge. Reported to fisheries dept.",
			"verified", now.Add(-18 * time.Hour)},
		{userID, 4, floatPtr(-1.3110), floatPtr(36.8410), "flood",
			"Lake level up 40cm since Monday. Jetty underwater. Some lakeside homes have water in compounds. Community sandbagging this morning.",
			"pending", now.Add(-4 * time.Hour)},
		{userID, 5, floatPtr(-1.2605), floatPtr(36.7905), "air_quality",
			"Dust haze reducing visibility to <2km. Strong easterly winds. Respiratory issues reported at clinic - 12 cases today alone. Unusual for this season.",
			"pending", now.Add(-2 * time.Hour)},
	}

	for _, o := range observations {
		var location interface{}
		if o.lat != nil && o.lon != nil {
			location = fmt.Sprintf("SRID=4326;POINT(%f %f)", *o.lon, *o.lat)
		}

		_, err := tx.Exec(ctx, `
			INSERT INTO observations (id, reporter_id, place_id, location, category, description, status, created_at, updated_at)
			VALUES (gen_random_uuid(), $1, $2, 
				CASE WHEN $3::text <> '' THEN $3::geography ELSE NULL END,
				$4, $5, $6, $7, $7)
		`, o.reporterID, placeIDs[o.placeIdx], location, o.category, o.description, o.status, o.createdAt)
		if err != nil {
			return err
		}
	}
	return nil
}

func seedPlaceAssessments(ctx context.Context, tx pgx.Tx, placeIDs []uuid.UUID) error {
	now := time.Now()
	periodStart := now.AddDate(0, 0, -7)
	periodEnd := now

	assessments := []struct {
		placeIdx      int
		urgency       int
		confidence    int
		summary       string
		affectedGroups []string
	}{
		{0, 45, 70, "Moderate rainfall increasing. Good soil recharge for agriculture. Monitor drainage in low areas.", []string{"farmers", "residents"}},
		{1, 30, 65, "Dry and warm conditions persisting. Irrigation demand increasing. Low flood risk.", []string{"farmers", "gardeners"}},
		{2, 85, 90, "HIGH URGENCY: Active flooding. River levels critical. Evacuation advisory in effect for vulnerable zones. Immediate action required.", []string{"all_residents", "emergency_services", "livestock_owners"}},
		{3, 70, 80, "Drought stress escalating rapidly. Soil moisture critically low. Crop failure risk high without intervention. Water conservation critical.", []string{"farmers", "livestock_owners", "water_users"}},
		{4, 60, 75, "Lake levels rising. Water quality declining. Flood risk for shoreline properties. Fishery impacts likely.", []string{"lakeside_residents", "fishers", "water_users"}},
		{5, 25, 60, "Normal seasonal conditions. No significant hazards. Good conditions for field preparation.", []string{"farmers", "residents"}},
	}

	for _, a := range assessments {
		placeID := placeIDs[a.placeIdx]

		indicatorRows, err := tx.Query(ctx, `SELECT indicator_id FROM place_indicators WHERE place_id = $1`, placeID)
		if err != nil {
			return err
		}
		var applicableIndicators []uuid.UUID
		for indicatorRows.Next() {
			var id uuid.UUID
			if err := indicatorRows.Scan(&id); err != nil {
				indicatorRows.Close()
				return err
			}
			applicableIndicators = append(applicableIndicators, id)
		}
		indicatorRows.Close()

		_, err = tx.Exec(ctx, `
			INSERT INTO place_assessments (id, place_id, assessed_at, period_start, period_end, urgency_score, confidence_score, applicable_indicators, affected_groups, assessment_summary)
			VALUES (gen_random_uuid(), $1, $2, $3, $4, $5, $6, $7, $8, $9)
			ON CONFLICT (place_id, period_start, period_end) DO UPDATE SET
				urgency_score = EXCLUDED.urgency_score,
				confidence_score = EXCLUDED.confidence_score,
				applicable_indicators = EXCLUDED.applicable_indicators,
				affected_groups = EXCLUDED.affected_groups,
				assessment_summary = EXCLUDED.assessment_summary,
				assessed_at = EXCLUDED.assessed_at
		`, placeID, now, periodStart, periodEnd, a.urgency, a.confidence, applicableIndicators, a.affectedGroups, a.summary)
		if err != nil {
			return err
		}
	}
	return nil
}

func hashPassword(password string) (string, error) {
	return security.HashPassword(password)
}

func strPtr(s string) *string {
	return &s
}

func floatPtr(f float64) *float64 {
	return &f
}