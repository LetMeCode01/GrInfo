-- Script: grinfo_rebuild_questions.sql
-- Purpose: Backup grinfo questions/options, mark inactive extras, insert new Neorientate easy questions,
--          and recalculate elo ratings with step=21 for 30 active questions per category.
-- IMPORTANT: This script marks old questions inactive instead of deleting them. It does NOT renumber IDs.

BEGIN;

-- 1) Backup current tables (one-time, safe)
CREATE TABLE IF NOT EXISTS grinfo_questions_backup AS TABLE grinfo_questions WITH NO DATA;
INSERT INTO grinfo_questions_backup SELECT * FROM grinfo_questions;

CREATE TABLE IF NOT EXISTS grinfo_question_options_backup AS TABLE grinfo_question_options WITH NO DATA;
INSERT INTO grinfo_question_options_backup SELECT * FROM grinfo_question_options;

-- Convenience: get category ids
SELECT id INTO TEMP TABLE tmp_categories FROM grinfo_categories WHERE slug IN ('neorientate','orientate');

-- Helper: function to activate a limited number of existing questions per category
-- We'll: 1) deactivate all in category, 2) activate specified ones (by enunt tag or id), 3) activate oldest remaining to reach 30.

-- ===== Neorientate =====
-- Deactivate all questions in category
UPDATE grinfo_questions
SET is_active = FALSE
WHERE category_id = (SELECT id FROM grinfo_categories WHERE slug = 'neorientate');

-- Activate the seed questions the user requested (keep these)
UPDATE grinfo_questions
SET is_active = TRUE
WHERE category_id = (SELECT id FROM grinfo_categories WHERE slug = 'neorientate')
  AND (
    enunt ILIKE '%[Neorientate 2]%' OR enunt ILIKE '%[Neorientate 3]%'
    OR id IN (16,17,31,36) -- medium:16,17 hard:31,36 as requested
  );

-- Activate additional existing rows (by id order) until we have 30 active questions
WITH needed AS (
  SELECT GREATEST(0, 30 - COUNT(*)) AS n
  FROM grinfo_questions
  WHERE category_id = (SELECT id FROM grinfo_categories WHERE slug = 'neorientate') AND is_active = TRUE
), to_pick AS (
  SELECT id FROM grinfo_questions
  WHERE category_id = (SELECT id FROM grinfo_categories WHERE slug = 'neorientate') AND is_active = FALSE
  ORDER BY id
  LIMIT (SELECT n FROM needed)
)
UPDATE grinfo_questions SET is_active = TRUE WHERE id IN (SELECT id FROM to_pick);

-- ===== Insert new Neorientate easy questions provided by user =====
-- We'll insert them as difficulty = 'usoara' and graph_data exactly as provided.

-- Question 1: BFS queue
WITH ins AS (
  INSERT INTO grinfo_questions (category_id, difficulty, elo_rating, enunt, explicatie_raspuns, graph_data, is_active)
  VALUES ((SELECT id FROM grinfo_categories WHERE slug='neorientate'), 'usoara', 1000,
    'În timpul unei parcurgeri în lățime (BFS) a unui graf neorientat, ce structură de date se folosește în mod ideal pentru gestionarea nodurilor descoperite?',
    'Parcurgerea BFS folosește o coadă (FIFO) pentru a procesa nodurile în ordinea strictă a distanței față de nodul de start.',
    '{}'::jsonb, TRUE)
  RETURNING id
)
INSERT INTO grinfo_question_options (question_id, option_index, option_text, is_correct)
SELECT id, ord, opt, (ord = 1)
FROM ins, (VALUES (0,'Stivă (Stack)'),(1,'Coadă (Queue)'),(2,'Listă simplu înlănțuită fără pointeri'),(3,'Arbore binar de căutare')) AS v(ord,opt);

-- Question 2: lists of adjacency space 2*m
WITH ins AS (
  INSERT INTO grinfo_questions (category_id, difficulty, elo_rating, enunt, explicatie_raspuns, graph_data, is_active)
  VALUES ((SELECT id FROM grinfo_categories WHERE slug='neorientate'), 'usoara', 1000,
    'Dacă un graf neorientat este reprezentat prin liste de adiacență, spațiul total de memorie ocupat de pointerii/elementele listelor este proporțional cu:',
    'Fiecare muchie (i, j) este memorată de două ori: j apare în lista lui i, iar i apare în lista lui j. Deci numărul total de elemente din liste este 2m.',
    '{}'::jsonb, TRUE)
  RETURNING id
)
INSERT INTO grinfo_question_options (question_id, option_index, option_text, is_correct)
SELECT id, ord, opt, (ord = 1)
FROM ins, (VALUES (0,'n^2'),(1,'2*m'),(2,'$m^2$'),(3,'n* m')) AS v(ord,opt);

-- Question 3: graf partial
WITH ins AS (
  INSERT INTO grinfo_questions (category_id, difficulty, elo_rating, enunt, explicatie_raspuns, graph_data, is_active)
  VALUES ((SELECT id FROM grinfo_categories WHERE slug='neorientate'), 'usoara', 1000,
    'Cum se numește un graf neorientat în care s-a eliminat cel puțin o muchie dintr-un graf complet $K_n$, dar fără a elimina noduri?',
    'Un graf parțial al unui graf G se obține păstrând toate nodurile și eliminând o parte din muchii.',
    '{}'::jsonb, TRUE)
  RETURNING id
)
INSERT INTO grinfo_question_options (question_id, option_index, option_text, is_correct)
SELECT id, ord, opt, (ord = 0)
FROM ins, (VALUES (0,'Graf parțial'),(1,'Subgraf'),(2,'Graf indus'),(3,'Graf complementar')) AS v(ord,opt);

-- Question 4: Câte muchii are graful alăturat? (2 edges)
WITH ins AS (
  INSERT INTO grinfo_questions (category_id, difficulty, elo_rating, enunt, explicatie_raspuns, graph_data, is_active)
  VALUES ((SELECT id FROM grinfo_categories WHERE slug='neorientate'), 'usoara', 1000,
    'Câte muchii are graful alăturat?',
    'Graful conține exact două conexiuni definite în JSON: între 1-2 și 2-3.',
    '{"nodes":[{"id":1},{"id":2},{"id":3}],"edges":[{"from":1,"to":2},{"from":2,"to":3}]}'::jsonb, TRUE)
  RETURNING id
)
INSERT INTO grinfo_question_options (question_id, option_index, option_text, is_correct)
SELECT id, ord, opt, (ord = 1)
FROM ins, (VALUES (0,'1'),(1,'2'),(2,'3'),(3,'4')) AS v(ord,opt);

-- Question 5: degree of node 2 is 3 (index 2)
WITH ins AS (
  INSERT INTO grinfo_questions (category_id, difficulty, elo_rating, enunt, explicatie_raspuns, graph_data, is_active)
  VALUES ((SELECT id FROM grinfo_categories WHERE slug='neorientate'), 'usoara', 1000,
    'Care este gradul nodului 2 în graful dat?',
    'Nodul 2 are 3 muchii incidente în JSON, care îl conectează direct de nodurile 1, 3 și 4.',
    '{"nodes":[{"id":1},{"id":2},{"id":3},{"id":4}],"edges":[{"from":1,"to":2},{"from":2,"to":3},{"from":2,"to":4}]}'::jsonb, TRUE)
  RETURNING id
)
INSERT INTO grinfo_question_options (question_id, option_index, option_text, is_correct)
SELECT id, ord, opt, (ord = 2)
FROM ins, (VALUES (0,'1'),(1,'2'),(2,'3'),(3,'4')) AS v(ord,opt);

-- Question 6: Câte noduri distincte? (4 nodes -> answer index 2)
WITH ins AS (
  INSERT INTO grinfo_questions (category_id, difficulty, elo_rating, enunt, explicatie_raspuns, graph_data, is_active)
  VALUES ((SELECT id FROM grinfo_categories WHERE slug='neorientate'), 'usoara', 1000,
    'Câte noduri distincte are graful descris?',
    'În array-ul nodes sunt declarate explicit 4 noduri (1, 2, 3 și 4), chiar dacă unele nu au muchii.',
    '{"nodes":[{"id":1},{"id":2},{"id":3},{"id":4}],"edges":[{"from":1,"to":2},{"from":2,"to":3}]}'::jsonb, TRUE)
  RETURNING id
)
INSERT INTO grinfo_question_options (question_id, option_index, option_text, is_correct)
SELECT id, ord, opt, (ord = 2)
FROM ins, (VALUES (0,'2'),(1,'3'),(2,'4'),(3,'5')) AS v(ord,opt);

-- Question 7: Is graph connected? (answer index 1)
WITH ins AS (
  INSERT INTO grinfo_questions (category_id, difficulty, elo_rating, enunt, explicatie_raspuns, graph_data, is_active)
  VALUES ((SELECT id FROM grinfo_categories WHERE slug='neorientate'), 'usoara', 1000,
    'Este graful alăturat un graf conex?',
    'Nu este conex. Muchiile formează două grupuri complet separate: {1,2} și {3,4}, fără nicio legătură între ele.',
    '{"nodes":[{"id":1},{"id":2},{"id":3},{"id":4}],"edges":[{"from":1,"to":2},{"from":3,"to":4}]}'::jsonb, TRUE)
  RETURNING id
)
INSERT INTO grinfo_question_options (question_id, option_index, option_text, is_correct)
SELECT id, ord, opt, (ord = 1)
FROM ins, (VALUES (0,'Da, este complet conex'),(1,'Nu, are componente izolate'),(2,'Nu se poate determina din JSON'),(3,'Da, deoarece nu are noduri cu gradul 0')) AS v(ord,opt);

-- Question 8: triangle cycle C3 (answer index 1)
WITH ins AS (
  INSERT INTO grinfo_questions (category_id, difficulty, elo_rating, enunt, explicatie_raspuns, graph_data, is_active)
  VALUES ((SELECT id FROM grinfo_categories WHERE slug='neorientate'), 'usoara', 1047,
    'Ce formă geometrică / structură simplă închisă descrie graful ilustrat?',
    'Muchiile (1,2), (2,3) și (3,1) închid o buclă perfectă de 3 noduri, formând un ciclu simplu.',
    '{"nodes":[{"id":1},{"id":2},{"id":3}],"edges":[{"from":1,"to":2},{"from":2,"to":3},{"from":3,"to":1}]}'::jsonb, TRUE)
  RETURNING id
)
INSERT INTO grinfo_question_options (question_id, option_index, option_text, is_correct)
SELECT id, ord, opt, (ord = 1)
FROM ins, (VALUES (0,'Arbore liniar'),(1,'Ciclu ($C_3$)'),(2,'Graf stea'),(3,'Graf bipartit complet')) AS v(ord,opt);

-- After inserts, ensure we still have exactly 30 active questions: if less, activate more existing; if more, deactivate newest beyond 30 by id.
-- Trim or expand to exactly 30 active per category (neorientate)
WITH active_list AS (
  SELECT id FROM grinfo_questions
  WHERE category_id = (SELECT id FROM grinfo_categories WHERE slug='neorientate') AND is_active = TRUE
  ORDER BY CASE difficulty WHEN 'usoara' THEN 1 WHEN 'medie' THEN 2 WHEN 'grea' THEN 3 END, id
), cnt AS (
  SELECT COUNT(*) AS c FROM active_list
)
-- If more than 30, deactivate those after first 30
UPDATE grinfo_questions
SET is_active = FALSE
WHERE id IN (
  SELECT id FROM (
    SELECT id, ROW_NUMBER() OVER (ORDER BY CASE difficulty WHEN 'usoara' THEN 1 WHEN 'medie' THEN 2 WHEN 'grea' THEN 3 END, id) AS rn
    FROM grinfo_questions
    WHERE category_id = (SELECT id FROM grinfo_categories WHERE slug='neorientate') AND is_active = TRUE
  ) t WHERE t.rn > 30
);

-- If less than 30, activate more (oldest inactive) to reach 30
WITH needed AS (
  SELECT GREATEST(0, 30 - (SELECT COUNT(*) FROM grinfo_questions WHERE category_id=(SELECT id FROM grinfo_categories WHERE slug='neorientate') AND is_active=TRUE)) AS n
), to_pick AS (
  SELECT id FROM grinfo_questions
  WHERE category_id = (SELECT id FROM grinfo_categories WHERE slug='neorientate') AND is_active = FALSE
  ORDER BY id
  LIMIT (SELECT n FROM needed)
)
UPDATE grinfo_questions SET is_active = TRUE WHERE id IN (SELECT id FROM to_pick);

-- ===== Orientate: apply same pattern =====
UPDATE grinfo_questions
SET is_active = FALSE
WHERE category_id = (SELECT id FROM grinfo_categories WHERE slug = 'orientate');

-- Activate seed ones (tags and ids as requested)
UPDATE grinfo_questions
SET is_active = TRUE
WHERE category_id = (SELECT id FROM grinfo_categories WHERE slug = 'orientate')
  AND (
    enunt ILIKE '%[Orientate 3]%' OR enunt ILIKE '%[Orientate 4]%'
    OR id IN (18,19,26,29)
  );

-- Activate additional existing rows until 30
WITH needed AS (
  SELECT GREATEST(0, 30 - COUNT(*)) AS n
  FROM grinfo_questions
  WHERE category_id = (SELECT id FROM grinfo_categories WHERE slug = 'orientate') AND is_active = TRUE
), to_pick AS (
  SELECT id FROM grinfo_questions
  WHERE category_id = (SELECT id FROM grinfo_categories WHERE slug = 'orientate') AND is_active = FALSE
  ORDER BY id
  LIMIT (SELECT n FROM needed)
)
UPDATE grinfo_questions SET is_active = TRUE WHERE id IN (SELECT id FROM to_pick);

-- Trim any excess to exactly 30 (deactivate those beyond first 30 in ordering)
UPDATE grinfo_questions
SET is_active = FALSE
WHERE id IN (
  SELECT id FROM (
    SELECT id, ROW_NUMBER() OVER (ORDER BY CASE difficulty WHEN 'usoara' THEN 1 WHEN 'medie' THEN 2 WHEN 'grea' THEN 3 END, id) AS rn
    FROM grinfo_questions
    WHERE category_id = (SELECT id FROM grinfo_categories WHERE slug='orientate') AND is_active = TRUE
  ) t WHERE t.rn > 30
);

-- If less than 30, activate more to reach 30
WITH needed AS (
  SELECT GREATEST(0, 30 - (SELECT COUNT(*) FROM grinfo_questions WHERE category_id=(SELECT id FROM grinfo_categories WHERE slug='orientate') AND is_active=TRUE)) AS n
), to_pick AS (
  SELECT id FROM grinfo_questions
  WHERE category_id = (SELECT id FROM grinfo_categories WHERE slug = 'orientate') AND is_active = FALSE
  ORDER BY id
  LIMIT (SELECT n FROM needed)
)
UPDATE grinfo_questions SET is_active = TRUE WHERE id IN (SELECT id FROM to_pick);

-- ===== Recalculate Elo ratings =====
-- We'll use base = 912 and step = 21 as requested. Ordering: difficulty (usoara, medie, grea) then id.

-- Neorientate
WITH ordered AS (
  SELECT id, ROW_NUMBER() OVER (ORDER BY CASE difficulty WHEN 'usoara' THEN 1 WHEN 'medie' THEN 2 WHEN 'grea' THEN 3 END, id) AS rn
  FROM grinfo_questions
  WHERE category_id = (SELECT id FROM grinfo_categories WHERE slug='neorientate') AND is_active = TRUE
)
UPDATE grinfo_questions q
SET elo_rating = 912 + (o.rn - 1) * 21
FROM ordered o
WHERE q.id = o.id;

-- Orientate
WITH ordered AS (
  SELECT id, ROW_NUMBER() OVER (ORDER BY CASE difficulty WHEN 'usoara' THEN 1 WHEN 'medie' THEN 2 WHEN 'grea' THEN 3 END, id) AS rn
  FROM grinfo_questions
  WHERE category_id = (SELECT id FROM grinfo_categories WHERE slug='orientate') AND is_active = TRUE
)
UPDATE grinfo_questions q
SET elo_rating = 912 + (o.rn - 1) * 21
FROM ordered o
WHERE q.id = o.id;

COMMIT;

-- NOTE: This script does NOT renumber primary keys to 1..60. Renumbering is destructive because of
-- foreign keys (grinfo_session_answers, sessions, logs). If you want renumbering, reply to confirm
-- and I will produce and run a controlled truncate+reinsert flow that preserves referential integrity
-- or migrates old session references.
