-- Откат: удаляем constraint и колонку
ALTER TABLE task_structs DROP CONSTRAINT fk_tasks_user_id;
DROP INDEX IF EXISTS idx_tasks_user_id;
ALTER TABLE task_structs DROP COLUMN user_id;