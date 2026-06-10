BEGIN;
-- Backup current orientate grele questions and options
CREATE TABLE IF NOT EXISTS grinfo_questions_backup_orientate_grea AS TABLE grinfo_questions WITH NO DATA;
INSERT INTO grinfo_questions_backup_orientate_grea SELECT * FROM grinfo_questions WHERE category_id=1 AND difficulty='grea';
CREATE TABLE IF NOT EXISTS grinfo_question_options_backup_orientate_grea AS TABLE grinfo_question_options WITH NO DATA;
INSERT INTO grinfo_question_options_backup_orientate_grea SELECT * FROM grinfo_question_options WHERE question_id IN (SELECT id FROM grinfo_questions WHERE category_id=1 AND difficulty='grea');

-- Update tags and set ELOs for kept IDs 51 and 54 -> Orientate 21,22
UPDATE grinfo_questions SET enunt = '[Orientate 21] ' || regexp_replace(enunt, '^(([\\[]Orientate[^\\]]*[\\]]\\s*)+)', ''), elo_rating = 1650 WHERE id = 51;
UPDATE grinfo_questions SET enunt = '[Orientate 22] ' || regexp_replace(enunt, '^(([\\[]Orientate[^\\]]*[\\]]\\s*)+)', ''), elo_rating = 1677 WHERE id = 54;

-- Delete other orientate grele questions
DELETE FROM grinfo_questions WHERE category_id=1 AND difficulty='grea' AND id NOT IN (51,54);

-- Insert new orientate grele questions [Orientate 23..30] with ELOs 1704..1893 (step 27)
-- Q1 -> Orientate 23 (elo 1704)
INSERT INTO grinfo_questions (category_id,difficulty,elo_rating,enunt,explicatie_raspuns,graph_data,is_active) VALUES (1,'grea',1704,$q$[Orientate 23] Ce algoritm determină drumurile minime de la un singur nod sursă către toate celelalte noduri într-un graf orientat ce poate conține și costuri negative, dar fără circuite negative?$q$,$q$Bellman-Ford rulează în timp O(n * m) și, spre deosebire de Dijkstra, poate procesa corect arce cu ponderi negative.$q$,'{}'::jsonb,true);
INSERT INTO grinfo_question_options (question_id, option_index, option_text, is_correct) VALUES (currval(pg_get_serial_sequence('grinfo_questions','id')),0,$q$Algoritmul lui Dijkstra$q$,false),(currval(pg_get_serial_sequence('grinfo_questions','id')),1,$q$Algoritmul Bellman-Ford$q$,true),(currval(pg_get_serial_sequence('grinfo_questions','id')),2,$q$Algoritmul lui Prim$q$,false),(currval(pg_get_serial_sequence('grinfo_questions','id')),3,$q$Algoritmul Kruskal$q$,false);

-- Q2 -> Orientate 24 (elo 1731)
INSERT INTO grinfo_questions (category_id,difficulty,elo_rating,enunt,explicatie_raspuns,graph_data,is_active) VALUES (1,'grea',1731,$q$[Orientate 24] Ce metodă algoritmică eficientă (bazată pe două parcurgeri DFS, dintre care a doua se realizează pe graful transpus) este folosită pentru a găsi componentele tare conexe?$q$,$q$Algoritmul lui Kosaraju folosește graful inversat (transpus) și timpii de finalizare din DFS pentru a izola componentele tare conexe în timp liniar O(n+m).$q$,'{}'::jsonb,true);
INSERT INTO grinfo_question_options (question_id, option_index, option_text, is_correct) VALUES (currval(pg_get_serial_sequence('grinfo_questions','id')),0,$q$Algoritmul lui Kosaraju$q$,true),(currval(pg_get_serial_sequence('grinfo_questions','id')),1,$q$Algoritmul lui Kruskal$q$,false),(currval(pg_get_serial_sequence('grinfo_questions','id')),2,$q$Algoritmul lui Lee$q$,false),(currval(pg_get_serial_sequence('grinfo_questions','id')),3,$q$Algoritmul Edmonds-Karp$q$,false);

-- Q3 -> Orientate 25 (elo 1758)
INSERT INTO grinfo_questions (category_id,difficulty,elo_rating,enunt,explicatie_raspuns,graph_data,is_active) VALUES (1,'grea',1758,$q$[Orientate 25] Dacă un graf orientat are n noduri și nu are circuite (este DAG), numărul maxim posibil de arce pe care le poate avea este:$q$,$q$Într-o sortare topologică, arcele pot merge doar de la stânga la dreapta. Numărul maxim de arce apare când fiecare nod are legături spre toate nodurile din dreapta lui, adică exact ca în graful complet: n(n-1)/2.$q$,'{}'::jsonb,true);
INSERT INTO grinfo_question_options (question_id, option_index, option_text, is_correct) VALUES (currval(pg_get_serial_sequence('grinfo_questions','id')),0,$q$n-1$q$,false),(currval(pg_get_serial_sequence('grinfo_questions','id')),1,$q$n(n-1)/2$q$,true),(currval(pg_get_serial_sequence('grinfo_questions','id')),2,$q$n^2$q$,false),(currval(pg_get_serial_sequence('grinfo_questions','id')),3,$q$2 * m$q$,false);

-- Q4 -> Orientate 26 (elo 1785)
INSERT INTO grinfo_questions (category_id,difficulty,elo_rating,enunt,explicatie_raspuns,graph_data,is_active) VALUES (1,'grea',1785,$q$[Orientate 26] Care este costul drumului minim de la nodul 1 la nodul 4 în graful ponderat de mai jos?$q$,$q$Traseul direct sau prin alte noduri oferă costuri variate: calea 1 -> 2 -> 4 oferă costul 2 + 5 = 7, care este mai mic decât calea 1 -> 3 -> 4 (1 + 8 = 9).$q$,$j${"nodes":[{"id":1},{"id":2},{"id":3},{"id":4}],"edges":[{"from":1,"to":2,"label":"2"},{"from":2,"to":4,"label":"5"},{"from":1,"to":3,"label":"1"},{"from":3,"to":4,"label":"8"}]}$j$::jsonb,true);
INSERT INTO grinfo_question_options (question_id, option_index, option_text, is_correct) VALUES (currval(pg_get_serial_sequence('grinfo_questions','id')),0,$q$5$q$,false),(currval(pg_get_serial_sequence('grinfo_questions','id')),1,$q$7$q$,true),(currval(pg_get_serial_sequence('grinfo_questions','id')),2,$q$8$q$,false),(currval(pg_get_serial_sequence('grinfo_questions','id')),3,$q$9$q$,false);

-- Q5 -> Orientate 27 (elo 1812)
INSERT INTO grinfo_questions (category_id,difficulty,elo_rating,enunt,explicatie_raspuns,graph_data,is_active) VALUES (1,'grea',1812,$q$[Orientate 27] Câte componente tare conexe (SCC) se pot izola în graful orientat de mai jos?$q$,$q$Nodurile {1, 2} pot fi accesate reciproc (formează o componentă). Nodurile {3} și {4} nu au circuite de întoarcere către restul grafului, deci fiecare constituie o componentă maximală proprie. Total: 3 componente.$q$,$j${"nodes":[{"id":1},{"id":2},{"id":3},{"id":4}],"edges":[{"from":1,"to":2},{"from":2,"to":1},{"from":3,"to":4}]}$j$::jsonb,true);
INSERT INTO grinfo_question_options (question_id, option_index, option_text, is_correct) VALUES (currval(pg_get_serial_sequence('grinfo_questions','id')),0,$q$1$q$,false),(currval(pg_get_serial_sequence('grinfo_questions','id')),1,$q$2$q$,true),(currval(pg_get_serial_sequence('grinfo_questions','id')),2,$q$3$q$,false),(currval(pg_get_serial_sequence('grinfo_questions','id')),3,$q$4$q$,false);

-- Q6 -> Orientate 28 (elo 1839)
INSERT INTO grinfo_questions (category_id,difficulty,elo_rating,enunt,explicatie_raspuns,graph_data,is_active) VALUES (1,'grea',1839,$q$[Orientate 28] Care este singura secvență liniară corectă rezultată în urma unei sortări topologice a acestui DAG?$q$,$q$Dacă arcele merg de la 1 la 2 și de la 2 la 3, nodul 1 trebuie să apară obligatoriu primul, urmat de 2, iar la final 3.$q$,$j${"nodes":[{"id":1},{"id":2},{"id":3}],"edges":[{"from":1,"to":2},{"from":2,"to":3}]}$j$::jsonb,true);
INSERT INTO grinfo_question_options (question_id, option_index, option_text, is_correct) VALUES (currval(pg_get_serial_sequence('grinfo_questions','id')),0,$q$3, 2, 1$q$,false),(currval(pg_get_serial_sequence('grinfo_questions','id')),1,$q$1, 2, 3$q$,true),(currval(pg_get_serial_sequence('grinfo_questions','id')),2,$q$2, 1, 3$q$,false),(currval(pg_get_serial_sequence('grinfo_questions','id')),3,$q$1, 3, 2$q$,false);

-- Q7 -> Orientate 29 (elo 1866)
INSERT INTO grinfo_questions (category_id,difficulty,elo_rating,enunt,explicatie_raspuns,graph_data,is_active) VALUES (1,'grea',1866,$q$[Orientate 29] Există vreo cale validă de acces (drum orientat) care să pornească din nodul 4 și să ajungă la nodul 1?$q$,$q$Nu există nicio posibilitate. Lanțul de arce progresează unidirecțional în direcția 1 -> 2 -> 3 -> 4. Nu există arce de întoarcere.$q$,$j${"nodes":[{"id":1},{"id":2},{"id":3},{"id":4}],"edges":[{"from":1,"to":2},{"from":2,"to":3},{"from":3,"to":4}]}$j$::jsonb,true);
INSERT INTO grinfo_question_options (question_id, option_index, option_text, is_correct) VALUES (currval(pg_get_serial_sequence('grinfo_questions','id')),0,$q$Da, prin arcul direct$q$,false),(currval(pg_get_serial_sequence('grinfo_questions','id')),1,$q$Nu, deoarece toate arcele merg în sens invers$q$,true),(currval(pg_get_serial_sequence('grinfo_questions','id')),2,$q$Doar dacă costul este minim$q$,false),(currval(pg_get_serial_sequence('grinfo_questions','id')),3,$q$Nu se poate determina$q$,false);

-- Q8 -> Orientate 30 (elo 1893)
INSERT INTO grinfo_questions (category_id,difficulty,elo_rating,enunt,explicatie_raspuns,graph_data,is_active) VALUES (1,'grea',1893,$q$[Orientate 30] Care este valoarea sumei gradelor exterioare (sum d^+(v)) pentru toate nodurile din graful de mai jos?$q$,$q$Numărăm arcele care ies din fiecare nod: nodul 1 are 2 arce de ieșire, nodul 2 are un arc de ieșire, nodul 3 are 0. Suma lor este 2 + 1 + 0 = 3.$q$,$j${"nodes":[{"id":1},{"id":2},{"id":3}],"edges":[{"from":1,"to":2},{"from":1,"to":3},{"from":2,"to":3}]}$j$::jsonb,true);
INSERT INTO grinfo_question_options (question_id, option_index, option_text, is_correct) VALUES (currval(pg_get_serial_sequence('grinfo_questions','id')),0,$q$2$q$,false),(currval(pg_get_serial_sequence('grinfo_questions','id')),1,$q$3$q$,true),(currval(pg_get_serial_sequence('grinfo_questions','id')),2,$q$4$q$,false),(currval(pg_get_serial_sequence('grinfo_questions','id')),3,$q$5$q$,false);

COMMIT;
