package indicators

import (
	"context"
	"log/slog"
	"math"
	"sort"
	"time"

	"github.com/google/uuid"
	"green-compass-backend/internal/normalization"
)

type indicatorRepo interface {
	ListDefinitions(ctx context.Context) ([]Definition, error)
	Store(ctx context.Context, ind *Indicator) error
	GetLatest(ctx context.Context, placeID, indicatorID uuid.UUID) (*Indicator, error)
	GetForPeriod(ctx context.Context, placeID, indicatorID uuid.UUID, periodStart, periodEnd time.Time) (*Indicator, error)
}

type NormalizationRepository interface {
	ListCanonicalObservations(ctx context.Context, placeID, sourceID uuid.UUID, from, to time.Time) ([]normalization.CanonicalObservation, error)
}

type Service struct {
	repo          indicatorRepo
	normRepo      NormalizationRepository
	logger        *slog.Logger
}

func NewService(repo indicatorRepo, normRepo NormalizationRepository, logger *slog.Logger) *Service {
	return &Service{
		repo:     repo,
		normRepo: normRepo,
		logger:   logger,
	}
}

// ComputeForPlace computes all indicators for a place over a given period
func (s *Service) ComputeForPlace(ctx context.Context, placeID uuid.UUID, periodStart, periodEnd time.Time) ([]Indicator, error) {
	defs, err := s.repo.ListDefinitions(ctx)
	if err != nil {
		return nil, err
	}

	var results []Indicator
	for _, def := range defs {
		ind, err := s.compute(ctx, placeID, def, periodStart, periodEnd)
		if err != nil {
			s.logger.Warn("failed to compute indicator", "place_id", placeID, "indicator", def.Code, "err", err)
			continue
		}
		if err := s.repo.Store(ctx, ind); err != nil {
			s.logger.Error("failed to store indicator", "place_id", placeID, "indicator", def.Code, "err", err)
			continue
		}
		results = append(results, *ind)
	}

	return results, nil
}

// compute derives a single indicator value from normalized observations
func (s *Service) compute(ctx context.Context, placeID uuid.UUID, def Definition, periodStart, periodEnd time.Time) (*Indicator, error) {
	// For now, compute against all sources; future: check which source provides this variable
	var allObservations []normalization.CanonicalObservation
	sourceIDs := []uuid.UUID{
		uuid.Nil, // Placeholder; would iterate through enabled sources in production
	}

	for _, sourceID := range sourceIDs {
		// Skip if sourceID is nil; in production, iterate over enabled sources
		if sourceID == uuid.Nil {
			continue
		}
		obs, err := s.normRepo.ListCanonicalObservations(ctx, placeID, sourceID, periodStart, periodEnd)
		if err != nil {
			continue
		}
		allObservations = append(allObservations, obs...)
	}

	ind := &Indicator{
		ID:          uuid.New(),
		PlaceID:     placeID,
		IndicatorID: def.ID,
		ComputedAt:  time.Now().UTC(),
		PeriodStart: periodStart,
		PeriodEnd:   periodEnd,
		Metadata:    make(map[string]interface{}),
	}

	switch def.Code {
	case "rain_intensity_trend":
		s.computeRainIntensityTrend(allObservations, ind)
	case "rain_frequency":
		s.computeRainFrequency(allObservations, ind)
	case "temperature_trend":
		s.computeTemperatureTrend(allObservations, ind)
	case "wind_speed_trend":
		s.computeWindSpeedTrend(allObservations, ind)
	case "soil_moisture_level":
		s.computeSoilMoistureLevel(allObservations, ind)
	case "drought_stress":
		s.computeDroughtStress(allObservations, ind)
	default:
		// Stub indicators
		ind.Value = 0
		ind.Unit = "unknown"
		ind.Metadata["stubbed"] = true
	}

	return ind, nil
}

func (s *Service) computeRainIntensityTrend(obs []normalization.CanonicalObservation, ind *Indicator) {
	prec := filterByVariable(obs, normalization.VarPrecipitation)
	if len(prec) < 2 {
		ind.Value = 0
		ind.Unit = "mm/day"
		ind.Trend = nil
		ind.DataPointsCount = len(prec)
		return
	}

	sort.Slice(prec, func(i, j int) bool {
		return prec[i].ObservedAt.Before(prec[j].ObservedAt)
	})

	half := len(prec) / 2
	firstHalfMean := mean(prec[:half])
	secondHalfMean := mean(prec[half:])

	ind.Value = secondHalfMean
	ind.Unit = "mm/day"
	ind.DataPointsCount = len(prec)

	if math.Abs(secondHalfMean-firstHalfMean) < 0.5 {
		stable := "stable"
		ind.Trend = &stable
		ind.TrendConfidence = floatPtr(0.7)
	} else if secondHalfMean > firstHalfMean {
		incr := "increasing"
		ind.Trend = &incr
		ind.TrendConfidence = floatPtr(0.8)
	} else {
		decr := "decreasing"
		ind.Trend = &decr
		ind.TrendConfidence = floatPtr(0.8)
	}

	ind.Metadata["first_half_mean"] = firstHalfMean
	ind.Metadata["second_half_mean"] = secondHalfMean
}

func (s *Service) computeRainFrequency(obs []normalization.CanonicalObservation, ind *Indicator) {
	prec := filterByVariable(obs, normalization.VarPrecipitation)
	rainyDays := 0
	for _, o := range prec {
		if o.Value > 0.5 {
			rainyDays++
		}
	}

	ind.Value = float64(rainyDays)
	ind.Unit = "days"
	ind.DataPointsCount = len(prec)
	ind.Metadata["total_observations"] = len(prec)
}

func (s *Service) computeTemperatureTrend(obs []normalization.CanonicalObservation, ind *Indicator) {
	temps := filterByVariable(obs, normalization.VarTemperature)
	if len(temps) < 2 {
		ind.Value = 0
		ind.Unit = normalization.UnitCelsius
		ind.DataPointsCount = len(temps)
		return
	}

	sort.Slice(temps, func(i, j int) bool {
		return temps[i].ObservedAt.Before(temps[j].ObservedAt)
	})

	half := len(temps) / 2
	firstHalfMean := mean(temps[:half])
	secondHalfMean := mean(temps[half:])

	ind.Value = secondHalfMean
	ind.Unit = normalization.UnitCelsius
	ind.DataPointsCount = len(temps)

	if math.Abs(secondHalfMean-firstHalfMean) < 1.0 {
		stable := "stable"
		ind.Trend = &stable
		ind.TrendConfidence = floatPtr(0.75)
	} else if secondHalfMean > firstHalfMean {
		incr := "increasing"
		ind.Trend = &incr
		ind.TrendConfidence = floatPtr(0.8)
	} else {
		decr := "decreasing"
		ind.Trend = &decr
		ind.TrendConfidence = floatPtr(0.8)
	}

	ind.Metadata["first_half_mean"] = firstHalfMean
	ind.Metadata["second_half_mean"] = secondHalfMean
}

func (s *Service) computeWindSpeedTrend(obs []normalization.CanonicalObservation, ind *Indicator) {
	winds := filterByVariable(obs, normalization.VarWindSpeed)
	if len(winds) == 0 {
		ind.Value = 0
		ind.Unit = normalization.UnitMetersPerSecond
		ind.DataPointsCount = 0
		return
	}

	ind.Value = mean(winds)
	ind.Unit = normalization.UnitMetersPerSecond
	ind.DataPointsCount = len(winds)
	ind.Trend = nil
	ind.TrendConfidence = nil
}

func (s *Service) computeSoilMoistureLevel(obs []normalization.CanonicalObservation, ind *Indicator) {
	moisture := filterByVariable(obs, normalization.VarSoilMoisture)
	if len(moisture) == 0 {
		ind.Value = 0
		ind.Unit = "percent"
		ind.DataPointsCount = 0
		return
	}

	ind.Value = mean(moisture) * 100
	ind.Unit = "percent"
	ind.DataPointsCount = len(moisture)
	ind.Trend = nil
}

func (s *Service) computeDroughtStress(obs []normalization.CanonicalObservation, ind *Indicator) {
	prec := filterByVariable(obs, normalization.VarPrecipitation)
	temps := filterByVariable(obs, normalization.VarTemperature)
	moisture := filterByVariable(obs, normalization.VarSoilMoisture)

	if len(prec) == 0 || len(temps) == 0 || len(moisture) == 0 {
		ind.Value = 0
		ind.Unit = "index"
		ind.DataPointsCount = 0
		return
	}

	precipMean := mean(prec)
	tempMean := mean(temps)
	moistureMean := mean(moisture) * 100

	stress := (50.0 - precipMean) / 50.0 * (tempMean - 20.0) / 10.0 * (50.0 - moistureMean) / 50.0 * 100.0
	if stress < 0 {
		stress = 0
	}
	if stress > 100 {
		stress = 100
	}

	ind.Value = stress
	ind.Unit = "index"
	ind.DataPointsCount = (len(prec) + len(temps) + len(moisture)) / 3
	ind.Metadata["precipitation_mean"] = precipMean
	ind.Metadata["temperature_mean"] = tempMean
	ind.Metadata["moisture_mean"] = moistureMean
}

// Helper functions
func filterByVariable(obs []normalization.CanonicalObservation, variable string) []normalization.CanonicalObservation {
	var filtered []normalization.CanonicalObservation
	for _, o := range obs {
		if o.Variable == variable && o.Valid {
			filtered = append(filtered, o)
		}
	}
	return filtered
}

func mean(obs []normalization.CanonicalObservation) float64 {
	if len(obs) == 0 {
		return 0
	}
	sum := 0.0
	for _, o := range obs {
		sum += o.Value
	}
	return sum / float64(len(obs))
}

func floatPtr(v float64) *float64 {
	return &v
}
