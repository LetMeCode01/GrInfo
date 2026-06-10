BEGIN;

-- Safety snapshot before restore.
CREATE TABLE IF NOT EXISTS grinfo_questions_backup_before_csv_restore AS TABLE grinfo_questions WITH NO DATA;
TRUNCATE grinfo_questions_backup_before_csv_restore;
INSERT INTO grinfo_questions_backup_before_csv_restore SELECT * FROM grinfo_questions;

CREATE TABLE IF NOT EXISTS grinfo_question_options_backup_before_csv_restore AS TABLE grinfo_question_options WITH NO DATA;
TRUNCATE grinfo_question_options_backup_before_csv_restore;
INSERT INTO grinfo_question_options_backup_before_csv_restore SELECT * FROM grinfo_question_options;

CREATE TEMP TABLE tmp_import_questions (
  id BIGINT,
  category TEXT,
  difficulty TEXT,
  elo_rating INTEGER,
  enunt TEXT,
  explicatie_raspuns TEXT,
  options TEXT,
  correct_indexes TEXT,
  graph_data TEXT
) ON COMMIT DROP;

COPY tmp_import_questions (id, category, difficulty, elo_rating, enunt, explicatie_raspuns, options, correct_indexes, graph_data)
FROM '/tmp/grinfo_all_questions_repaired.csv' WITH (FORMAT csv, HEADER true);

-- Keep exactly the 60-question dataset from CSV.
DELETE FROM grinfo_question_options qo
USING grinfo_questions q
JOIN grinfo_categories c ON c.id = q.category_id
WHERE qo.question_id = q.id
  AND c.slug IN ('orientate', 'neorientate')
  AND q.id NOT IN (SELECT id FROM tmp_import_questions);

DELETE FROM grinfo_questions q
USING grinfo_categories c
WHERE c.id = q.category_id
  AND c.slug IN ('orientate', 'neorientate')
  AND q.id NOT IN (SELECT id FROM tmp_import_questions);

-- Upsert questions from CSV.
INSERT INTO grinfo_questions (id, category_id, difficulty, elo_rating, enunt, explicatie_raspuns, graph_data, is_active)
SELECT
  t.id,
  c.id,
  CASE WHEN t.difficulty IN ('usoara','medie','grea') THEN t.difficulty ELSE 'medie' END,
  COALESCE(t.elo_rating, 1000),
  t.enunt,
  t.explicatie_raspuns,
  CASE
    WHEN t.graph_data IS NULL OR btrim(t.graph_data) = '' OR btrim(t.graph_data) = '{}' THEN '{}'::jsonb
    ELSE t.graph_data::jsonb
  END,
  TRUE
FROM tmp_import_questions t
JOIN grinfo_categories c ON c.slug = t.category
ON CONFLICT (id) DO UPDATE
SET category_id = EXCLUDED.category_id,
    difficulty = EXCLUDED.difficulty,
    elo_rating = EXCLUDED.elo_rating,
    enunt = EXCLUDED.enunt,
    explicatie_raspuns = EXCLUDED.explicatie_raspuns,
    graph_data = EXCLUDED.graph_data,
    is_active = TRUE;

-- Replace options for imported question IDs.
DELETE FROM grinfo_question_options
WHERE question_id IN (SELECT id FROM tmp_import_questions);

INSERT INTO grinfo_question_options (question_id, option_index, option_text, is_correct)
SELECT
  t.id,
  split_part(opt.item, ':', 1)::int AS option_index,
  regexp_replace(opt.item, '^[0-9]+:', '') AS option_text,
  split_part(opt.item, ':', 1)::int = ANY(
    CASE
      WHEN t.correct_indexes IS NULL OR btrim(t.correct_indexes) = '' THEN ARRAY[]::int[]
      ELSE string_to_array(replace(t.correct_indexes, ' ', ''), ',')::int[]
    END
  ) AS is_correct
FROM tmp_import_questions t
CROSS JOIN LATERAL regexp_split_to_table(t.options, '\s\|\s') AS opt(item);

-- Keep only active imported questions in GrInfo categories.
UPDATE grinfo_questions q
SET is_active = FALSE
FROM grinfo_categories c
WHERE c.id = q.category_id
  AND c.slug IN ('orientate', 'neorientate')
  AND q.id NOT IN (SELECT id FROM tmp_import_questions);

-- Ensure imported questions are active.
UPDATE grinfo_questions
SET is_active = TRUE
WHERE id IN (SELECT id FROM tmp_import_questions);

-- Fix sequences after explicit IDs.
SELECT setval(pg_get_serial_sequence('grinfo_questions', 'id'), COALESCE((SELECT MAX(id) FROM grinfo_questions), 1), true);
SELECT setval(pg_get_serial_sequence('grinfo_question_options', 'id'), COALESCE((SELECT MAX(id) FROM grinfo_question_options), 1), true);

COMMIT;

-- Validation
SELECT COUNT(*) AS active_orientate_neorientate
FROM grinfo_questions q
JOIN grinfo_categories c ON c.id = q.category_id
WHERE c.slug IN ('orientate','neorientate') AND q.is_active = TRUE;
