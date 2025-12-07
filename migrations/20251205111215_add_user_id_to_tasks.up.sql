-- Добавляем колонку user_id и foreign key 
-- (в поле task теперь будет поле, которое указывает на то, какому пользователю пренадлежит данная таска.
-- у одного пользователя может быть несколько тасок)
ALTER TABLE task_structs ADD COLUMN user_id INTEGER;

ALTER TABLE task_structs
ADD CONSTRAINT fk_tasks_user_id 
FOREIGN KEY (user_id) REFERENCES user_structs(id) 
ON DELETE CASCADE;

CREATE INDEX idx_tasks_user_id ON task_structs(user_id);