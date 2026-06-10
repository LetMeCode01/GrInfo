BEGIN;

-- Keep IDs 76 and 77 as [Neorientate 21] and [Neorientate 22]
UPDATE grinfo_questions
SET enunt = '[Neorientate 21] ' || regexp_replace(enunt, '^(\[Neorientate[^\]]*\]\s*)+', '')
WHERE id = 76;

UPDATE grinfo_questions
SET enunt = '[Neorientate 22] ' || regexp_replace(enunt, '^(\[Neorientate[^\]]*\]\s*)+', '')
WHERE id = 77;

-- Remove the 13 older heavy leftovers
DELETE FROM grinfo_question_options WHERE question_id IN (78,79,80,81,82,83,84,85,86,87,88,89,90);
DELETE FROM grinfo_questions WHERE id IN (78,79,80,81,82,83,84,85,86,87,88,89,90);

-- New heavy questions [Neorientate 23..30]
INSERT INTO grinfo_questions (category_id, difficulty, elo_rating, enunt, explicatie_raspuns, graph_data, is_active)
VALUES ((SELECT id FROM grinfo_categories WHERE slug='neorientate'), 'grea', 1704,
$q$[Neorientate 23] Care este numărul maxim de muchii pe care le poate avea un graf neorientat cu n noduri pentru a rămâne deconectat, în cel mai nefavorabil caz?$q$,
$q$Dacă izolăm un singur nod, celelalte n - 1 noduri pot forma un graf complet cu maximum (n - 1)(n - 2) / 2 muchii, fără ca graful total să devină conex.$q$,
'{}'::jsonb, TRUE)
RETURNING id;
INSERT INTO grinfo_question_options (question_id, option_index, option_text, is_correct)
VALUES (currval(pg_get_serial_sequence('grinfo_questions','id')),0,'n - 1',FALSE),(currval(pg_get_serial_sequence('grinfo_questions','id')),1,'(n - 1)(n - 2) / 2',TRUE),(currval(pg_get_serial_sequence('grinfo_questions','id')),2,'n(n - 1) / 2',FALSE),(currval(pg_get_serial_sequence('grinfo_questions','id')),3,'n',FALSE);

INSERT INTO grinfo_questions (category_id, difficulty, elo_rating, enunt, explicatie_raspuns, graph_data, is_active)
VALUES ((SELECT id FROM grinfo_categories WHERE slug='neorientate'), 'grea', 1731,
$q$[Neorientate 24] Ce algoritm folosește o strategie greedy pentru a construi un arbore parțial de cost minim, alegând succesiv muchiile cu costul cel mai mic și evitând ciclurile?$q$,
$q$Algoritmul lui Kruskal sortează muchiile crescător după cost și le adaugă pe rând dacă nu creează cicluri, folosind seturi disjuncte.$q$,
'{}'::jsonb, TRUE)
RETURNING id;
INSERT INTO grinfo_question_options (question_id, option_index, option_text, is_correct)
VALUES (currval(pg_get_serial_sequence('grinfo_questions','id')),0,'Algoritmul lui Dijkstra',FALSE),(currval(pg_get_serial_sequence('grinfo_questions','id')),1,'Algoritmul lui Kruskal',TRUE),(currval(pg_get_serial_sequence('grinfo_questions','id')),2,'Algoritmul lui Prim',FALSE),(currval(pg_get_serial_sequence('grinfo_questions','id')),3,'Algoritmul Roy-Floyd',FALSE);

INSERT INTO grinfo_questions (category_id, difficulty, elo_rating, enunt, explicatie_raspuns, graph_data, is_active)
VALUES ((SELECT id FROM grinfo_categories WHERE slug='neorientate'), 'grea', 1758,
$q$[Neorientate 25] Într-un graf Hamiltonian neorientat cu n noduri, un ciclu Hamiltonian parcurge:$q$,
$q$Ciclul Hamiltonian trece prin toate nodurile exact o singură dată, cu excepția primului nod, care coincide cu ultimul.$q$,
'{}'::jsonb, TRUE)
RETURNING id;
INSERT INTO grinfo_question_options (question_id, option_index, option_text, is_correct)
VALUES (currval(pg_get_serial_sequence('grinfo_questions','id')),0,'Toate muchiile grafului o singură dată',FALSE),(currval(pg_get_serial_sequence('grinfo_questions','id')),1,'Toate nodurile grafului exact o singură dată, cu excepția primului care coincide cu ultimul',TRUE),(currval(pg_get_serial_sequence('grinfo_questions','id')),2,'Doar nodurile cu grad impar',FALSE),(currval(pg_get_serial_sequence('grinfo_questions','id')),3,'Fiecare componentă conexă de două ori',FALSE);

INSERT INTO grinfo_questions (category_id, difficulty, elo_rating, enunt, explicatie_raspuns, graph_data, is_active)
VALUES ((SELECT id FROM grinfo_categories WHERE slug='neorientate'), 'grea', 1785,
$q$[Neorientate 26] Care din nodurile grafului de mai jos reprezintă un punct de articulare?$q$,
$q$Nodul 2 este punctul de articulare. Dacă este eliminat, nodul 1 rămâne separat de grupul format din {3,4}.$q$,
$g${"nodes":[{"id":1},{"id":2},{"id":3},{"id":4}],"edges":[{"from":1,"to":2},{"from":2,"to":3},{"from":3,"to":4},{"from":2,"to":4}]}$g$::jsonb, TRUE)
RETURNING id;
INSERT INTO grinfo_question_options (question_id, option_index, option_text, is_correct)
VALUES (currval(pg_get_serial_sequence('grinfo_questions','id')),0,'Nodul 1',FALSE),(currval(pg_get_serial_sequence('grinfo_questions','id')),1,'Nodul 2',TRUE),(currval(pg_get_serial_sequence('grinfo_questions','id')),2,'Nodul 3',FALSE),(currval(pg_get_serial_sequence('grinfo_questions','id')),3,'Nodul 4',FALSE);

INSERT INTO grinfo_questions (category_id, difficulty, elo_rating, enunt, explicatie_raspuns, graph_data, is_active)
VALUES ((SELECT id FROM grinfo_categories WHERE slug='neorientate'), 'grea', 1812,
$q$[Neorientate 27] Care este diametrul grafului din imagine, adică cea mai mare valoare dintre toate distanțele minime?$q$,
$q$Cea mai mare distanță minimă este între nodurile 1 și 4. Drumul minim necesită 3 muchii: 1 -> 2 -> 3 -> 4.$q$,
$g${"nodes":[{"id":1},{"id":2},{"id":3},{"id":4}],"edges":[{"from":1,"to":2},{"from":2,"to":3},{"from":3,"to":4}]}$g$::jsonb, TRUE)
RETURNING id;
INSERT INTO grinfo_question_options (question_id, option_index, option_text, is_correct)
VALUES (currval(pg_get_serial_sequence('grinfo_questions','id')),0,'1',FALSE),(currval(pg_get_serial_sequence('grinfo_questions','id')),1,'2',FALSE),(currval(pg_get_serial_sequence('grinfo_questions','id')),2,'3',TRUE),(currval(pg_get_serial_sequence('grinfo_questions','id')),3,'4',FALSE);

INSERT INTO grinfo_questions (category_id, difficulty, elo_rating, enunt, explicatie_raspuns, graph_data, is_active)
VALUES ((SELECT id FROM grinfo_categories WHERE slug='neorientate'), 'grea', 1839,
$q$[Neorientate 28] Câte subgrafuri de tip triunghi, adică clici de dimensiune 3 sau K_3, conține această structură?$q$,
$q$Se observă două triunghiuri: {1,2,3} și {3,4,5}. Nodul 3 este comun, dar fiecare triunghi este distinct.$q$,
$g${"nodes":[{"id":1},{"id":2},{"id":3},{"id":4},{"id":5}],"edges":[{"from":1,"to":2},{"from":2,"to":3},{"from":3,"to":1},{"from":3,"to":4},{"from":4,"to":5},{"from":5,"to":3}]}$g$::jsonb, TRUE)
RETURNING id;
INSERT INTO grinfo_question_options (question_id, option_index, option_text, is_correct)
VALUES (currval(pg_get_serial_sequence('grinfo_questions','id')),0,'0',FALSE),(currval(pg_get_serial_sequence('grinfo_questions','id')),1,'1',FALSE),(currval(pg_get_serial_sequence('grinfo_questions','id')),2,'2',TRUE),(currval(pg_get_serial_sequence('grinfo_questions','id')),3,'3',FALSE);

INSERT INTO grinfo_questions (category_id, difficulty, elo_rating, enunt, explicatie_raspuns, graph_data, is_active)
VALUES ((SELECT id FROM grinfo_categories WHERE slug='neorientate'), 'grea', 1866,
$q$[Neorientate 29] Care este numărul minim de muchii ce trebuie adăugate pentru a transforma graful dat într-unul 2-conex, fără puncte de articulare?$q$,
$q$Graful este un lanț simplu 1 - 2 - 3. Dacă adăugăm o muchie între nodurile terminale 1 și 3, obținem un singur ciclu și eliminăm punctele de articulare.$q$,
$g${"nodes":[{"id":1},{"id":2},{"id":3}],"edges":[{"from":1,"to":2},{"from":2,"to":3}]}$g$::jsonb, TRUE)
RETURNING id;
INSERT INTO grinfo_question_options (question_id, option_index, option_text, is_correct)
VALUES (currval(pg_get_serial_sequence('grinfo_questions','id')),0,'0',FALSE),(currval(pg_get_serial_sequence('grinfo_questions','id')),1,'1',TRUE),(currval(pg_get_serial_sequence('grinfo_questions','id')),2,'2',FALSE),(currval(pg_get_serial_sequence('grinfo_questions','id')),3,'3',FALSE);

INSERT INTO grinfo_questions (category_id, difficulty, elo_rating, enunt, explicatie_raspuns, graph_data, is_active)
VALUES ((SELECT id FROM grinfo_categories WHERE slug='neorientate'), 'grea', 1893,
$q$[Neorientate 30] Analizând gradele nodurilor din graful de mai jos, ce putem spune despre existența unui circuit Eulerian?$q$,
$q$Toate cele 4 noduri au gradul 2, deci au grad par. Cum graful este și conex, el conține un circuit Eulerian.$q$,
$g${"nodes":[{"id":1},{"id":2},{"id":3},{"id":4}],"edges":[{"from":1,"to":2},{"from":2,"to":3},{"from":3,"to":4},{"from":4,"to":1}]}$g$::jsonb, TRUE)
RETURNING id;
INSERT INTO grinfo_question_options (question_id, option_index, option_text, is_correct)
VALUES (currval(pg_get_serial_sequence('grinfo_questions','id')),0,'Conține un circuit Eulerian',TRUE),(currval(pg_get_serial_sequence('grinfo_questions','id')),1,'Nu conține circuit Eulerian deoarece are noduri cu grad impar',FALSE),(currval(pg_get_serial_sequence('grinfo_questions','id')),2,'Conține doar drum Eulerian, nu și circuit',FALSE),(currval(pg_get_serial_sequence('grinfo_questions','id')),3,'Este un graf complet deconectat',FALSE);

COMMIT;

SELECT id, enunt, elo_rating
FROM grinfo_questions q
JOIN grinfo_categories c ON c.id = q.category_id
WHERE c.slug = 'neorientate' AND q.difficulty = 'grea'
ORDER BY id;
