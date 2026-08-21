DROP INDEX IF EXISTS place_content_type_idx;
DROP INDEX IF EXISTS place_content_place_generated_idx;
DROP TABLE IF EXISTS place_content;

DROP INDEX IF EXISTS place_assessments_place_assessed_idx;
DROP TABLE IF EXISTS place_assessments;

DROP TRIGGER IF EXISTS indicator_applicability_rules_set_updated_at ON indicator_applicability_rules;
DROP TABLE IF EXISTS indicator_applicability_rules;

DROP INDEX IF EXISTS place_indicators_period_idx;
DROP INDEX IF EXISTS place_indicators_place_computed_idx;
DROP TABLE IF EXISTS place_indicators;

DROP TRIGGER IF EXISTS indicator_definitions_set_updated_at ON indicator_definitions;
DROP TABLE IF EXISTS indicator_definitions;
