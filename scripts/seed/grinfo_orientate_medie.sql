BEGIN;
-- Backup current orientate medie questions and options
CREATE TABLE IF NOT EXISTS grinfo_questions_backup_orientate_medie AS TABLE grinfo_questions WITH NO DATA;
INSERT INTO grinfo_questions_backup_orientate_medie SELECT * FROM grinfo_questions WHERE category_id=1 AND difficulty='medie';
CREATE TABLE IF NOT EXISTS grinfo_question_options_backup_orientate_medie AS TABLE grinfo_question_options WITH NO DATA;
INSERT INTO grinfo_question_options_backup_orientate_medie SELECT * FROM grinfo_question_options WHERE question_id IN (SELECT id FROM grinfo_questions WHERE category_id=1 AND difficulty='medie');

-- Update tags and set ELOs for kept IDs 18 and 19 -> Orientate 11,12
UPDATE grinfo_questions SET enunt = '[Orientate 11] ' || regexp_replace(enunt, '^((\\[Orientate[^\\]]*\\]\\s*)+)', ''), elo_rating = 1300 WHERE id = 18;
UPDATE grinfo_questions SET enunt = '[Orientate 12] ' || regexp_replace(enunt, '^((\\[Orientate[^\\]]*\\]\\s*)+)', ''), elo_rating = 1327 WHERE id = 19;

-- Delete other orientate medie questions
DELETE FROM grinfo_questions WHERE category_id=1 AND difficulty='medie' AND id NOT IN (18,19);

-- Insert new orientate medie questions [Orientate 13..20] with ELOs 1354..1543 (step 27)
-- Q1 -> Orientate 13 (elo 1354)
INSERT INTO grinfo_questions (category_id,difficulty,elo_rating,enunt,explicatie_raspuns,graph_data,is_active) VALUES (1,'medie',1354,$q$[Orientate 13] Ce algoritm clasic folosește programarea dinamică pentru a determina matricea drumurilor (existența drumurilor între oricare două noduri) într-un graf orientat?$q$,$q$Algoritmul Roy-Warshall construiește în mod eficient încheierea tranzitivă a unui graf, determinând dacă există sau nu drumuri între noduri.$q$,'{}'::jsonb,true);
INSERT INTO grinfo_question_options (question_id, option_index, option_text, is_correct) VALUES (currval(pg_get_serial_sequence('grinfo_questions','id')),0,$q$Algoritmul lui Dijkstra$q$,false),(currval(pg_get_serial_sequence('grinfo_questions','id')),1,$q$Algoritmul Roy-Warshall$q$,true),(currval(pg_get_serial_sequence('grinfo_questions','id')),2,$q$Algoritmul lui Prim$q$,false),(currval(pg_get_serial_sequence('grinfo_questions','id')),3,$q$Kruskal$q$,false);

-- Q2 -> Orientate 14 (elo 1381)
INSERT INTO grinfo_questions (category_id,difficulty,elo_rating,enunt,explicatie_raspuns,graph_data,is_active) VALUES (1,'medie',1381,$q$[Orientate 14] Dacă un graf orientat conține un arc de la nodul i la nodul j, cum afectează acest lucru simetria matricei de adiacență?$q$,$q$Simetria perfectă apare doar la grafurile neorientate. La cele orientate, relațiile de adiacență au sens unic.$q$,'{}'::jsonb,true);
INSERT INTO grinfo_question_options (question_id, option_index, option_text, is_correct) VALUES (currval(pg_get_serial_sequence('grinfo_questions','id')),0,$q$Matricea rămâne perfect simetrică$q$,false),(currval(pg_get_serial_sequence('grinfo_questions','id')),1,$q$Matricea nu mai este garantat simetrică, deoarece arcul (i,j) nu implică existența arcului (j,i)$q$,true),(currval(pg_get_serial_sequence('grinfo_questions','id')),2,$q$Toate elementele de deasupra diagonalei devin zero$q$,false),(currval(pg_get_serial_sequence('grinfo_questions','id')),3,$q$Matricea se inversează$q$,false);

-- Q3 -> Orientate 15 (elo 1408)
INSERT INTO grinfo_questions (category_id,difficulty,elo_rating,enunt,explicatie_raspuns,graph_data,is_active) VALUES (1,'medie',1408,$q$[Orientate 15] Într-un graf orientat tare conex (strongly connected):$q$,$q$Aceasta este definiția conexiunii tari pe grafuri orientate: accesibilitate reciprocă deplină între oricare două noduri.$q$,'{}'::jsonb,true);
INSERT INTO grinfo_question_options (question_id, option_index, option_text, is_correct) VALUES (currval(pg_get_serial_sequence('grinfo_questions','id')),0,$q$Nu există niciun circuit$q$,false),(currval(pg_get_serial_sequence('grinfo_questions','id')),1,$q$Din orice nod i se poate ajunge în orice nod j urmând sensul arcelor$q$,true),(currval(pg_get_serial_sequence('grinfo_questions','id')),2,$q$Toate nodurile au gradul interior egal cu 0$q$,false),(currval(pg_get_serial_sequence('grinfo_questions','id')),3,$q$Numărul de arce este egal cu n-1$q$,false);

-- Q4 -> Orientate 16 (elo 1435)
INSERT INTO grinfo_questions (category_id,difficulty,elo_rating,enunt,explicatie_raspuns,graph_data,is_active) VALUES (1,'medie',1435,$q$[Orientate 16] Care este costul total cumulat al drumului parcurs în sensul arcelor 1 -> 2 -> 3?$q$,$q$Adunăm costurile arcelor componente: arcul (1,2) costă 10, iar arcul (2,3) costă 5. Totalul este 10 + 5 = 15.$q$,$j${"nodes":[{"id":1},{"id":2},{"id":3}],"edges":[{"from":1,"to":2,"label":"10"},{"from":2,"to":3,"label":"5"}]}$j$::jsonb,true);
INSERT INTO grinfo_question_options (question_id, option_index, option_text, is_correct) VALUES (currval(pg_get_serial_sequence('grinfo_questions','id')),0,$q$5$q$,false),(currval(pg_get_serial_sequence('grinfo_questions','id')),1,$q$10$q$,false),(currval(pg_get_serial_sequence('grinfo_questions','id')),2,$q$15$q$,true),(currval(pg_get_serial_sequence('grinfo_questions','id')),3,$q$50$q$,false);

-- Q5 -> Orientate 17 (elo 1462)
INSERT INTO grinfo_questions (category_id,difficulty,elo_rating,enunt,explicatie_raspuns,graph_data,is_active) VALUES (1,'medie',1462,$q$[Orientate 17] Câte circuite orientate distincte se găsesc în graful următor?$q$,$q$Există o singură buclă închisă care respectă sensul arcelor: traseul 1 -> 2 -> 3 -> 1.$q$,$j${"nodes":[{"id":1},{"id":2},{"id":3}],"edges":[{"from":1,"to":2},{"from":2,"to":3},{"from":3,"to":1}]}$j$::jsonb,true);
INSERT INTO grinfo_question_options (question_id, option_index, option_text, is_correct) VALUES (currval(pg_get_serial_sequence('grinfo_questions','id')),0,$q$0$q$,false),(currval(pg_get_serial_sequence('grinfo_questions','id')),1,$q$1$q$,true),(currval(pg_get_serial_sequence('grinfo_questions','id')),2,$q$2$q$,false),(currval(pg_get_serial_sequence('grinfo_questions','id')),3,$q$3$q$,false);

-- Q6 -> Orientate 18 (elo 1489)
INSERT INTO grinfo_questions (category_id,difficulty,elo_rating,enunt,explicatie_raspuns,graph_data,is_active) VALUES (1,'medie',1489,$q$[Orientate 18] Ce valoare numerică va fi stocată în celula A[1][2] a matricei de adiacență ponderate a grafului?$q$,$q$Celula A[1][2] salvează direct ponderea arcului existent de la nodul 1 la nodul 2, care este 7.$q$,$j${"nodes":[{"id":1},{"id":2}],"edges":[{"from":1,"to":2,"label":"7"}]}$j$::jsonb,true);
INSERT INTO grinfo_question_options (question_id, option_index, option_text, is_correct) VALUES (currval(pg_get_serial_sequence('grinfo_questions','id')),0,$q$0$q$,false),(currval(pg_get_serial_sequence('grinfo_questions','id')),1,$q$1$q$,false),(currval(pg_get_serial_sequence('grinfo_questions','id')),2,$q$7$q$,true),(currval(pg_get_serial_sequence('grinfo_questions','id')),3,$q$Infinity$q$,false);

-- Q7 -> Orientate 19 (elo 1516)
INSERT INTO grinfo_questions (category_id,difficulty,elo_rating,enunt,explicatie_raspuns,graph_data,is_active) VALUES (1,'medie',1516,$q$[Orientate 19] Este acest graf orientat unul tare conex?$q$,$q$Da, este tare conex. Putem ajunge de la 1 la 2 direct, iar de la 2 la 1 înapoi tot direct.$q$,$j${"nodes":[{"id":1},{"id":2}],"edges":[{"from":1,"to":2},{"from":2,"to":1}]}$j$::jsonb,true);
INSERT INTO grinfo_question_options (question_id, option_index, option_text, is_correct) VALUES (currval(pg_get_serial_sequence('grinfo_questions','id')),0,$q$Da, deoarece există drum de la oricare nod la oricare altul$q$,true),(currval(pg_get_serial_sequence('grinfo_questions','id')),1,$q$Nu, deoarece arcele nu au costuri$q$,false),(currval(pg_get_serial_sequence('grinfo_questions','id')),2,$q$Nu, este doar un graf slab conex$q$,false),(currval(pg_get_serial_sequence('grinfo_questions','id')),3,$q$Este un DAG$q$,false);

-- Q8 -> Orientate 20 (elo 1543)
INSERT INTO grinfo_questions (category_id,difficulty,elo_rating,enunt,explicatie_raspuns,graph_data,is_active) VALUES (1,'medie',1543,$q$[Orientate 20] Care este costul celui mai scurt drum (cost minim ponderat) de la nodul 1 la nodul 3?$q$,$q$Deși există arc direct de la 1 la 3 cu costul 10, ruta ocolitoare 1 -> 2 -> 3 costă doar 2 + 3 = 5, fiind cea optimă.$q$,$j${"nodes":[{"id":1},{"id":2},{"id":3}],"edges":[{"from":1,"to":3,"label":"10"},{"from":1,"to":2,"label":"2"},{"from":2,"to":3,"label":"3"}]}$j$::jsonb,true);
INSERT INTO grinfo_question_options (question_id, option_index, option_text, is_correct) VALUES (currval(pg_get_serial_sequence('grinfo_questions','id')),0,$q$2$q$,false),(currval(pg_get_serial_sequence('grinfo_questions','id')),1,$q$3$q$,false),(currval(pg_get_serial_sequence('grinfo_questions','id')),2,$q$5$q$,true),(currval(pg_get_serial_sequence('grinfo_questions','id')),3,$q$10$q$,false);

COMMIT;
