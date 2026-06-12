-- Maintenance scripts (manual use).
-- Run only when needed: psql -U dev -d local_db -f scripts/03_maintenance.sql
\i scripts/maintenance/cleanup_graph_text.sql