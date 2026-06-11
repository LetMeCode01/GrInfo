BEGIN;

-- Permanent hotfix for existing DB rows:
-- [Neorientate 20] correct answer must be option_index = 2 (value "3").
UPDATE grinfo_question_options qo
SET is_correct = (qo.option_index = 2)
WHERE qo.question_id IN (
  SELECT q.id
  FROM grinfo_questions q
  JOIN grinfo_categories c ON c.id = q.category_id
  WHERE c.slug = 'neorientate'
    AND q.enunt LIKE '[Neorientate 20]%'
);

COMMIT;
