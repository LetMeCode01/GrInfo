-- Backup target tables, delete neorientate usoara except ids 27,28, insert provided 8 questions
BEGIN;

-- Backup current questions and options
CREATE TABLE IF NOT EXISTS grinfo_questions_backup_manual AS TABLE grinfo_questions WITH NO DATA;
INSERT INTO grinfo_questions_backup_manual SELECT * FROM grinfo_questions;

CREATE TABLE IF NOT EXISTS grinfo_question_options_backup_manual AS TABLE grinfo_question_options WITH NO DATA;
INSERT INTO grinfo_question_options_backup_manual SELECT * FROM grinfo_question_options;

-- Delete neorientate usoara except ids 27,28
DELETE FROM grinfo_questions
WHERE category_id = (SELECT id FROM grinfo_categories WHERE slug='neorientate')
  AND difficulty = 'usoara'
  AND id NOT IN (27,28);

-- Insert the 8 new Neorientate easy questions (IDs will be assigned automatically)

-- 1
WITH ins AS (
 INSERT INTO grinfo_questions (category_id,difficulty,elo_rating,enunt,explicatie_raspuns,graph_data,is_active)
 VALUES ((SELECT id FROM grinfo_categories WHERE slug='neorientate'),'usoara',1000,
 'În timpul unei parcurgeri în lățime (BFS) a unui graf neorientat, ce structură de date se folosește în mod ideal pentru gestionarea nodurilor descoperite?',
 'Parcurgerea BFS folosește o coadă (FIFO) pentru a procesa nodurile în ordinea strictă a distanței față de nodul de start.', '{}'::jsonb, TRUE)
 RETURNING id
)
INSERT INTO grinfo_question_options (question_id, option_index, option_text, is_correct)
SELECT id, ord, opt, (ord = 1)
FROM ins, (VALUES (0,'Stivă (Stack)'),(1,'Coadă (Queue)'),(2,'Listă simplu înlănțuită fără pointeri'),(3,'Arbore binar de căutare')) AS v(ord,opt);

-- 2
WITH ins AS (
 INSERT INTO grinfo_questions (category_id,difficulty,elo_rating,enunt,explicatie_raspuns,graph_data,is_active)
 VALUES ((SELECT id FROM grinfo_categories WHERE slug='neorientate'),'usoara',1000,
 'Dacă un graf neorientat este reprezentat prin liste de adiacență, spațiul total de memorie ocupat de pointerii/elementele listelor este proporțional cu:',
 'Fiecare muchie (i, j) este memorată de două ori: j apare în lista lui i, iar i apare în lista lui j. Deci numărul total de elemente din liste este 2m.', '{}'::jsonb, TRUE)
 RETURNING id
)
INSERT INTO grinfo_question_options (question_id, option_index, option_text, is_correct)
SELECT id, ord, opt, (ord = 1)
FROM ins, (VALUES (0,'n^2'),(1,'2*m'),(2,'$m^2$'),(3,'n* m')) AS v(ord,opt);

-- 3
WITH ins AS (
 INSERT INTO grinfo_questions (category_id,difficulty,elo_rating,enunt,explicatie_raspuns,graph_data,is_active)
 VALUES ((SELECT id FROM grinfo_categories WHERE slug='neorientate'),'usoara',1000,
 'Cum se numește un graf neorientat în care s-a eliminat cel puțin o muchie dintr-un graf complet $K_n$, dar fără a elimina noduri?',
 'Un graf parțial al unui graf G se obține păstrând toate nodurile și eliminând o parte din muchii.', '{}'::jsonb, TRUE)
 RETURNING id
)
INSERT INTO grinfo_question_options (question_id, option_index, option_text, is_correct)
SELECT id, ord, opt, (ord = 0)
FROM ins, (VALUES (0,'Graf parțial'),(1,'Subgraf'),(2,'Graf indus'),(3,'Graf complementar')) AS v(ord,opt);

-- 4
WITH ins AS (
 INSERT INTO grinfo_questions (category_id,difficulty,elo_rating,enunt,explicatie_raspuns,graph_data,is_active)
 VALUES ((SELECT id FROM grinfo_categories WHERE slug='neorientate'),'usoara',1000,
 'Câte muchii are graful alăturat?', 'Graful conține exact două conexiuni definite în JSON: între 1-2 și 2-3.',
 '{"nodes":[{"id":1},{"id":2},{"id":3}],"edges":[{"from":1,"to":2},{"from":2,"to":3}]}'::jsonb, TRUE)
 RETURNING id
)
INSERT INTO grinfo_question_options (question_id, option_index, option_text, is_correct)
SELECT id, ord, opt, (ord = 1)
FROM ins, (VALUES (0,'1'),(1,'2'),(2,'3'),(3,'4')) AS v(ord,opt);

-- 5
WITH ins AS (
 INSERT INTO grinfo_questions (category_id,difficulty,elo_rating,enunt,explicatie_raspuns,graph_data,is_active)
 VALUES ((SELECT id FROM grinfo_categories WHERE slug='neorientate'),'usoara',1000,
 'Care este gradul nodului 2 în graful dat?', 'Nodul 2 are 3 muchii incidente în JSON, care îl conectează direct de nodurile 1, 3 și 4.',
 '{"nodes":[{"id":1},{"id":2},{"id":3},{"id":4}],"edges":[{"from":1,"to":2},{"from":2,"to":3},{"from":2,"to":4}]}'::jsonb, TRUE)
 RETURNING id
)
INSERT INTO grinfo_question_options (question_id, option_index, option_text, is_correct)
SELECT id, ord, opt, (ord = 2)
FROM ins, (VALUES (0,'1'),(1,'2'),(2,'3'),(3,'4')) AS v(ord,opt);

-- 6
WITH ins AS (
 INSERT INTO grinfo_questions (category_id,difficulty,elo_rating,enunt,explicatie_raspuns,graph_data,is_active)
 VALUES ((SELECT id FROM grinfo_categories WHERE slug='neorientate'),'usoara',1000,
 'Câte noduri distincte are graful descris?', 'În array-ul nodes sunt declarate explicit 4 noduri (1, 2, 3 și 4), chiar dacă unele nu au muchii.',
 '{"nodes":[{"id":1},{"id":2},{"id":3},{"id":4}],"edges":[{"from":1,"to":2},{"from":2,"to":3}]}'::jsonb, TRUE)
 RETURNING id
)
INSERT INTO grinfo_question_options (question_id, option_index, option_text, is_correct)
SELECT id, ord, opt, (ord = 2)
FROM ins, (VALUES (0,'2'),(1,'3'),(2,'4'),(3,'5')) AS v(ord,opt);

-- 7
WITH ins AS (
 INSERT INTO grinfo_questions (category_id,difficulty,elo_rating,enunt,explicatie_raspuns,graph_data,is_active)
 VALUES ((SELECT id FROM grinfo_categories WHERE slug='neorientate'),'usoara',1000,
 'Este graful alăturat un graf conex?', 'Nu este conex. Muchiile formează două grupuri complet separate: {1,2} și {3,4}, fără nicio legătură între ele.',
 '{"nodes":[{"id":1},{"id":2},{"id":3},{"id":4}],"edges":[{"from":1,"to":2},{"from":3,"to":4}]}'::jsonb, TRUE)
 RETURNING id
)
INSERT INTO grinfo_question_options (question_id, option_index, option_text, is_correct)
SELECT id, ord, opt, (ord = 1)
FROM ins, (VALUES (0,'Da, este complet conex'),(1,'Nu, are componente izolate'),(2,'Nu se poate determina din JSON'),(3,'Da, deoarece nu are noduri cu gradul 0')) AS v(ord,opt);

-- 8
WITH ins AS (
 INSERT INTO grinfo_questions (category_id,difficulty,elo_rating,enunt,explicatie_raspuns,graph_data,is_active)
 VALUES ((SELECT id FROM grinfo_categories WHERE slug='neorientate'),'usoara',1047,
 'Ce formă geometrică / structură simplă închisă descrie graful ilustrat?', 'Muchiile (1,2), (2,3) și (3,1) închid o buclă perfectă de 3 noduri, formând un ciclu simplu.',
 '{"nodes":[{"id":1},{"id":2},{"id":3}],"edges":[{"from":1,"to":2},{"from":2,"to":3},{"from":3,"to":1}]}'::jsonb, TRUE)
 RETURNING id
)
INSERT INTO grinfo_question_options (question_id, option_index, option_text, is_correct)
SELECT id, ord, opt, (ord = 1)
FROM ins, (VALUES (0,'Arbore liniar'),(1,'Ciclu ($C_3$)'),(2,'Graf stea'),(3,'Graf bipartit complet')) AS v(ord,opt);

COMMIT;

-- Note: If you prefer marking inactive instead of DELETE, let me know.
