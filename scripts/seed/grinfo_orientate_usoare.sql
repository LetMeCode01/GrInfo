BEGIN;
-- Backup current orientate usoara questions and options
CREATE TABLE IF NOT EXISTS grinfo_questions_backup_orientate_usoara AS TABLE grinfo_questions WITH NO DATA;
INSERT INTO grinfo_questions_backup_orientate_usoara SELECT * FROM grinfo_questions WHERE category_id=1 AND difficulty='usoara';
CREATE TABLE IF NOT EXISTS grinfo_question_options_backup_orientate_usoara AS TABLE grinfo_question_options WITH NO DATA;
INSERT INTO grinfo_question_options_backup_orientate_usoara SELECT * FROM grinfo_question_options WHERE question_id IN (SELECT id FROM grinfo_questions WHERE category_id=1 AND difficulty='usoara');
-- Update tags and elo for kept IDs 3 and 4
UPDATE grinfo_questions SET enunt = '[Orientate 1] ' || regexp_replace(enunt, '^([\[]Orientate[^\]]*\]\\s*)+', ''), elo_rating = 900 WHERE id = 3;
UPDATE grinfo_questions SET enunt = '[Orientate 2] ' || regexp_replace(enunt, '^([\[]Orientate[^\]]*\]\\s*)+', ''), elo_rating = 921 WHERE id = 4;
-- Delete other orientate usoara questions
DELETE FROM grinfo_questions WHERE category_id=1 AND difficulty='usoara' AND id NOT IN (3,4);

-- Insert new orientate usoara questions [Orientate 3..10] with cleaned text
-- Q1 -> Orientate 3 (elo 942)
INSERT INTO grinfo_questions (category_id,difficulty,elo_rating,enunt,explicatie_raspuns,graph_data,is_active) VALUES (1,'usoara',942,$q$[Orientate 3] Cum se calculează gradul total al unui nod într-un graf orientat?$q$,$q$Gradul total reprezintă numărul total de arce conectate la nod, adică cele care intră plus cele care ies.$q$,'{}'::jsonb,true);
INSERT INTO grinfo_question_options (question_id, option_index, option_text, is_correct) VALUES (currval(pg_get_serial_sequence('grinfo_questions','id')),0,$q$Este egală cu gradul interior (in-degree)$q$,false),(currval(pg_get_serial_sequence('grinfo_questions','id')),1,$q$Este suma dintre gradul interior și gradul exterior (d^-(v) + d^+(v))$q$,true),(currval(pg_get_serial_sequence('grinfo_questions','id')),2,$q$Este diferența lor$q$,false),(currval(pg_get_serial_sequence('grinfo_questions','id')),3,$q$Este egal cu numărul total de noduri din graf$q$,false);

-- Q2 -> Orientate 4 (elo 963)
INSERT INTO grinfo_questions (category_id,difficulty,elo_rating,enunt,explicatie_raspuns,graph_data,is_active) VALUES (1,'usoara',963,$q$[Orientate 4] Într-o matrice de adiacență a unui graf orientat fără bucle, elementele de pe diagonala principală sunt întotdeauna:$q$,$q$Diagonala reprezintă arcele de la un nod la el însuși (i -> i). Deoarece graful nu are bucle, toate valorile de pe diagonală sunt 0.$q$,'{}'::jsonb,true);
INSERT INTO grinfo_question_options (question_id, option_index, option_text, is_correct) VALUES (currval(pg_get_serial_sequence('grinfo_questions','id')),0,$q$Egale cu 1$q$,false),(currval(pg_get_serial_sequence('grinfo_questions','id')),1,$q$Egale cu 0$q$,true),(currval(pg_get_serial_sequence('grinfo_questions','id')),2,$q$Egale cu costul maxim$q$,false),(currval(pg_get_serial_sequence('grinfo_questions','id')),3,$q$Valori negative$q$,false);

-- Q3 -> Orientate 5 (elo 984)
INSERT INTO grinfo_questions (category_id,difficulty,elo_rating,enunt,explicatie_raspuns,graph_data,is_active) VALUES (1,'usoara',984,$q$[Orientate 5] Care este numărul maxim de arce într-un graf orientat cu n noduri care nu conține bucle?$q$,$q$Între oricare două noduri distincte i și j pot exista maximum două arce: un arc de la i la j și unul de la j la i. Totalul este n(n-1).$q$,'{}'::jsonb,true);
INSERT INTO grinfo_question_options (question_id, option_index, option_text, is_correct) VALUES (currval(pg_get_serial_sequence('grinfo_questions','id')),0,$q$n(n-1)/2$q$,false),(currval(pg_get_serial_sequence('grinfo_questions','id')),1,$q$n(n-1)$q$,true),(currval(pg_get_serial_sequence('grinfo_questions','id')),2,$q$n^2$q$,false),(currval(pg_get_serial_sequence('grinfo_questions','id')),3,$q$n-1$q$,false);

-- Q4 -> Orientate 6 (elo 1005)
INSERT INTO grinfo_questions (category_id,difficulty,elo_rating,enunt,explicatie_raspuns,graph_data,is_active) VALUES (1,'usoara',1005,$q$[Orientate 6] Care este gradul interior (in-degree) al nodului 2 în graful orientat alăturat?$q$,$q$Există un singur arc direcționat către nodul 2, cel care pornește din nodul 1.$q$,$j${"nodes":[{"id":1},{"id":2}],"edges":[{"from":1,"to":2}]}$j$::jsonb,true);
INSERT INTO grinfo_question_options (question_id, option_index, option_text, is_correct) VALUES (currval(pg_get_serial_sequence('grinfo_questions','id')),0,$q$0$q$,false),(currval(pg_get_serial_sequence('grinfo_questions','id')),1,$q$1$q$,true),(currval(pg_get_serial_sequence('grinfo_questions','id')),2,$q$2$q$,false),(currval(pg_get_serial_sequence('grinfo_questions','id')),3,$q$3$q$,false);

-- Q5 -> Orientate 7 (elo 1026)
INSERT INTO grinfo_questions (category_id,difficulty,elo_rating,enunt,explicatie_raspuns,graph_data,is_active) VALUES (1,'usoara',1026,$q$[Orientate 7] Care este costul (ponderea) arcului care pleacă direct din nodul 1?$q$,$q$Arcul definit de la sursa 1 la destinația 2 are setat un label cu valoarea explicită "5".$q$,$j${"nodes":[{"id":1},{"id":2},{"id":3}],"edges":[{"from":1,"to":2,"label":"5"},{"from":2,"to":3,"label":"3"}]}$j$::jsonb,true);
INSERT INTO grinfo_question_options (question_id, option_index, option_text, is_correct) VALUES (currval(pg_get_serial_sequence('grinfo_questions','id')),0,$q$3$q$,false),(currval(pg_get_serial_sequence('grinfo_questions','id')),1,$q$5$q$,true),(currval(pg_get_serial_sequence('grinfo_questions','id')),2,$q$8$q$,false),(currval(pg_get_serial_sequence('grinfo_questions','id')),3,$q$Nu are cost$q$,false);

-- Q6 -> Orientate 8 (elo 1047)
INSERT INTO grinfo_questions (category_id,difficulty,elo_rating,enunt,explicatie_raspuns,graph_data,is_active) VALUES (1,'usoara',1047,$q$[Orientate 8] Este posibil să navigăm pe arce urmând sensul lor de la nodul 1 la nodul 3?$q$,$q$Da, sensurile arcelor permit formarea drumului orientat valid 1 -> 2 -> 3.$q$,$j${"nodes":[{"id":1},{"id":2},{"id":3}],"edges":[{"from":1,"to":2},{"from":2,"to":3}]}$j$::jsonb,true);
INSERT INTO grinfo_question_options (question_id, option_index, option_text, is_correct) VALUES (currval(pg_get_serial_sequence('grinfo_questions','id')),0,$q$Da, prin drumul direct$q$,true),(currval(pg_get_serial_sequence('grinfo_questions','id')),1,$q$Nu, sensul arcelor blochează trecerea$q$,false),(currval(pg_get_serial_sequence('grinfo_questions','id')),2,$q$Doar dacă graful se transformă în graf neorientat$q$,false),(currval(pg_get_serial_sequence('grinfo_questions','id')),3,$q$Nu se poate determina$q$,false);

-- Q7 -> Orientate 9 (elo 1068)
INSERT INTO grinfo_questions (category_id,difficulty,elo_rating,enunt,explicatie_raspuns,graph_data,is_active) VALUES (1,'usoara',1068,$q$[Orientate 9] Care este gradul exterior (out-degree) al nodului 1 în structura dată?$q$,$q$Din nodul 1 pornesc (ies) exact două arce: unul orientat spre nodul 2 și altul spre nodul 3.$q$,$j${"nodes":[{"id":1},{"id":2},{"id":3}],"edges":[{"from":1,"to":2},{"from":1,"to":3}]}$j$::jsonb,true);
INSERT INTO grinfo_question_options (question_id, option_index, option_text, is_correct) VALUES (currval(pg_get_serial_sequence('grinfo_questions','id')),0,$q$0$q$,false),(currval(pg_get_serial_sequence('grinfo_questions','id')),1,$q$1$q$,false),(currval(pg_get_serial_sequence('grinfo_questions','id')),2,$q$2$q$,true),(currval(pg_get_serial_sequence('grinfo_questions','id')),3,$q$3$q$,false);

-- Q8 -> Orientate 10 (elo 1089)
INSERT INTO grinfo_questions (category_id,difficulty,elo_rating,enunt,explicatie_raspuns,graph_data,is_active) VALUES (1,'usoara',1089,$q$[Orientate 10] Verificând prezența sau absența buclelor închise în sensul săgeților, este acest graf un DAG?$q$,$q$Da. Drumul merge doar înainte (1 -> 2 și 2 -> 3). Nu există nicio cale de întoarcere, deci nu avem circuite. Este un DAG valid.$q$,$j${"nodes":[{"id":1},{"id":2},{"id":3}],"edges":[{"from":1,"to":2},{"from":2,"to":3}]}$j$::jsonb,true);
INSERT INTO grinfo_question_options (question_id, option_index, option_text, is_correct) VALUES (currval(pg_get_serial_sequence('grinfo_questions','id')),0,$q$Da, este un graf orientat aciclic$q$,true),(currval(pg_get_serial_sequence('grinfo_questions','id')),1,$q$Nu, conține circuite$q$,false),(currval(pg_get_serial_sequence('grinfo_questions','id')),2,$q$Este un graf neorientat$q$,false),(currval(pg_get_serial_sequence('grinfo_questions','id')),3,$q$Nu se poate defini ca DAG din lipsa costurilor$q$,false);

COMMIT;
