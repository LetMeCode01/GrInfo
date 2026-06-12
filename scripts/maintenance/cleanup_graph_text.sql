BEGIN;

UPDATE grinfo_questions
SET enunt = '[Neorientate 5] Cum se numește un graf neorientat în care s-a eliminat cel puțin o muchie dintr-un graf complet K_n, dar fără a elimina noduri?'
WHERE id = 93;

UPDATE grinfo_question_options
SET option_text = 'm^2'
WHERE question_id = 92 AND option_index = 2;

UPDATE grinfo_question_options
SET option_text = 'Ciclu (C3)'
WHERE question_id = 98 AND option_index = 1;

COMMIT;
