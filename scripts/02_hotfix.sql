-- One-time hotfix scripts.
-- Run only when needed: psql -U dev -d local_db -f scripts/02_hotfix.sql
\i scripts/hotfix/fix_neorientate20_correct_option.sql
\i scripts/hotfix/grinfo_neorientate_replace.sql