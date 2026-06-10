BEGIN;

-- Tag existing IDs 41 and 42 and set ELOs
UPDATE grinfo_questions SET enunt = '[Neorientate 11] ' || enunt, elo_rating = 1300 WHERE id = 41;
UPDATE grinfo_questions SET enunt = '[Neorientate 12] ' || enunt, elo_rating = 1327 WHERE id = 42;

-- Insert new medium questions (Neorientate 13..20) using currval to attach options

INSERT INTO grinfo_questions (category_id, difficulty, elo_rating, enunt, explicatie_raspuns, graph_data, is_active)
VALUES ((SELECT id FROM grinfo_categories WHERE slug='neorientate'),'medie',1354,$q$[Neorientate 13] În urma parcurgerii în adâncime (DFS) a unui graf neorientat conex, muchiile folosite pentru a descoperi noduri noi formează:$q$,$q$Parcurgerea oricărui graf conex prin DFS extrage o structură aciclică și conexă care atinge toate nodurile, numită arbore de parcurgere.$q$, $j${}$j$::jsonb, TRUE);
INSERT INTO grinfo_question_options (question_id, option_index, option_text, is_correct)
VALUES (currval(pg_get_serial_sequence('grinfo_questions','id')),0,'Un circuit',FALSE),(currval(pg_get_serial_sequence('grinfo_questions','id')),1,'Un graf bipartit',FALSE),(currval(pg_get_serial_sequence('grinfo_questions','id')),2,'Un arbore parțial de parcurgere',TRUE),(currval(pg_get_serial_sequence('grinfo_questions','id')),3,'O clică completă',FALSE);

INSERT INTO grinfo_questions (category_id, difficulty, elo_rating, enunt, explicatie_raspuns, graph_data, is_active)
VALUES ((SELECT id FROM grinfo_categories WHERE slug='neorientate'),'medie',1381,$q$[Neorientate 14] Care este numărul maxim de componente conexe pe care le poate avea un graf neorientat cu n noduri?$q$,$q$Cazul extrem este graful nul (fără nicio muchie). Fiecare nod reprezintă o componentă conexă de sine stătătoare, deci avem n componente.$q$,$j${}$j$::jsonb, TRUE);
INSERT INTO grinfo_question_options (question_id, option_index, option_text, is_correct)
VALUES (currval(pg_get_serial_sequence('grinfo_questions','id')),0,'1',FALSE),(currval(pg_get_serial_sequence('grinfo_questions','id')),1,'n',TRUE),(currval(pg_get_serial_sequence('grinfo_questions','id')),2,'n-1',FALSE),(currval(pg_get_serial_sequence('grinfo_questions','id')),3,'0',FALSE);

INSERT INTO grinfo_questions (category_id, difficulty, elo_rating, enunt, explicatie_raspuns, graph_data, is_active)
VALUES ((SELECT id FROM grinfo_categories WHERE slug='neorientate'),'medie',1408,$q$[Neorientate 15] Cum se obține un subgraf al unui graf G prin definiție?$q$,$q$Subgraful se generează prin selectarea unei submulțimi de noduri și păstrarea exclusivă a acelor muchii originale care le conectează între ele.$q$,$j${}$j$::jsonb, TRUE);
INSERT INTO grinfo_question_options (question_id, option_index, option_text, is_correct)
VALUES (currval(pg_get_serial_sequence('grinfo_questions','id')),0,'Eliminând doar muchii, păstrând toate nodurile',FALSE),(currval(pg_get_serial_sequence('grinfo_questions','id')),1,'Păstrând o submulțime de noduri și doar muchiile din G care au ambele extremități în acea submulțime',TRUE),(currval(pg_get_serial_sequence('grinfo_questions','id')),2,'Adăugând noduri noi',FALSE),(currval(pg_get_serial_sequence('grinfo_questions','id')),3,'Inversând legăturile',FALSE);

INSERT INTO grinfo_questions (category_id, difficulty, elo_rating, enunt, explicatie_raspuns, graph_data, is_active)
VALUES ((SELECT id FROM grinfo_categories WHERE slug='neorientate'),'medie',1435,$q$[Neorientate 16] Câte cicluri simple conține graful definit mai jos?$q$,$q$Nodurile 1-2-3 formează un singur ciclu închis. Muchia (3,4) este doar o ramură exterioară care nu închide nicio buclă.$q$,$j${"nodes":[{"id":1},{"id":2},{"id":3},{"id":4}],"edges":[{"from":1,"to":2},{"from":2,"to":3},{"from":3,"to":1},{"from":3,"to":4}]}$j$::jsonb, TRUE);
INSERT INTO grinfo_question_options (question_id, option_index, option_text, is_correct)
VALUES (currval(pg_get_serial_sequence('grinfo_questions','id')),0,'0',FALSE),(currval(pg_get_serial_sequence('grinfo_questions','id')),1,'1',TRUE),(currval(pg_get_serial_sequence('grinfo_questions','id')),2,'2',FALSE),(currval(pg_get_serial_sequence('grinfo_questions','id')),3,'3',FALSE);

INSERT INTO grinfo_questions (category_id, difficulty, elo_rating, enunt, explicatie_raspuns, graph_data, is_active)
VALUES ((SELECT id FROM grinfo_categories WHERE slug='neorientate'),'medie',1462,$q$[Neorientate 17] Care este numărul minim de muchii care trebuie eliminate pentru a transforma acest graf într-un arbore parțial?$q$,$q$Graful are 3 noduri și 3 muchii. Un arbore cu 3 noduri are nevoie de exact 2 muchii (n-1). Eliminând oricare o muchie din ciclu, obținem un arbore.$q$,$j${"nodes":[{"id":1},{"id":2},{"id":3}],"edges":[{"from":1,"to":2},{"from":2,"to":3},{"from":3,"to":1}]}$j$::jsonb, TRUE);
INSERT INTO grinfo_question_options (question_id, option_index, option_text, is_correct)
VALUES (currval(pg_get_serial_sequence('grinfo_questions','id')),0,'0',FALSE),(currval(pg_get_serial_sequence('grinfo_questions','id')),1,'1',TRUE),(currval(pg_get_serial_sequence('grinfo_questions','id')),2,'2',FALSE),(currval(pg_get_serial_sequence('grinfo_questions','id')),3,'3',FALSE);

INSERT INTO grinfo_questions (category_id, difficulty, elo_rating, enunt, explicatie_raspuns, graph_data, is_active)
VALUES ((SELECT id FROM grinfo_categories WHERE slug='neorientate'),'medie',1489,$q$[Neorientate 18] Care este lungimea (numărul de muchii) a celui mai scurt drum de la nodul 1 la nodul 4?$q$,$q$Singura cale de conectare în acest graf liniar este 1 -> 2 -> 3 -> 4, ceea ce înseamnă exact 3 muchii parcurse.$q$,$j${"nodes":[{"id":1},{"id":2},{"id":3},{"id":4}],"edges":[{"from":1,"to":2},{"from":2,"to":3},{"from":3,"to":4}]}$j$::jsonb, TRUE);
INSERT INTO grinfo_question_options (question_id, option_index, option_text, is_correct)
VALUES (currval(pg_get_serial_sequence('grinfo_questions','id')),0,'1',FALSE),(currval(pg_get_serial_sequence('grinfo_questions','id')),1,'2',TRUE),(currval(pg_get_serial_sequence('grinfo_questions','id')),2,'3',FALSE),(currval(pg_get_serial_sequence('grinfo_questions','id')),3,'4',FALSE);

INSERT INTO grinfo_questions (category_id, difficulty, elo_rating, enunt, explicatie_raspuns, graph_data, is_active)
VALUES ((SELECT id FROM grinfo_categories WHERE slug='neorientate'),'medie',1516,$q$[Neorientate 19] Câte componente conexe distincte are acest graf?$q$,$q$Sunt 3 componente conexe: grupul format din {1,2}, grupul format din {3,4} și nodul complet izolat {5}.$q$,$j${"nodes":[{"id":1},{"id":2},{"id":3},{"id":4},{"id":5}],"edges":[{"from":1,"to":2},{"from":3,"to":4}]}$j$::jsonb, TRUE);
INSERT INTO grinfo_question_options (question_id, option_index, option_text, is_correct)
VALUES (currval(pg_get_serial_sequence('grinfo_questions','id')),0,'1',FALSE),(currval(pg_get_serial_sequence('grinfo_questions','id')),1,'2',FALSE),(currval(pg_get_serial_sequence('grinfo_questions','id')),2,'3',FALSE),(currval(pg_get_serial_sequence('grinfo_questions','id')),3,'4',TRUE);

INSERT INTO grinfo_questions (category_id, difficulty, elo_rating, enunt, explicatie_raspuns, graph_data, is_active)
VALUES ((SELECT id FROM grinfo_categories WHERE slug='neorientate'),'medie',1543,$q$[Neorientate 20] Care este valoarea gradului maxim (Delta(G)) întâlnit în graful de mai jos?$q$,$q$Nodul 1 este nod central și are 3 muchii care pleacă din el. Toate celelalte noduri au gradul 1.$q$,$j${"nodes":[{"id":1},{"id":2},{"id":3},{"id":4}],"edges":[{"from":1,"to":2},{"from":1,"to":3},{"from":1,"to":4}]}$j$::jsonb, TRUE);
INSERT INTO grinfo_question_options (question_id, option_index, option_text, is_correct)
VALUES (currval(pg_get_serial_sequence('grinfo_questions','id')),0,'1',FALSE),(currval(pg_get_serial_sequence('grinfo_questions','id')),1,'2',TRUE),(currval(pg_get_serial_sequence('grinfo_questions','id')),2,'3',FALSE),(currval(pg_get_serial_sequence('grinfo_questions','id')),3,'4',FALSE);

COMMIT;

-- Show inserted/updated Neorientate questions 11..20
SELECT id, enunt, elo_rating FROM grinfo_questions WHERE enunt LIKE '[Neorientate %' ORDER BY id;
