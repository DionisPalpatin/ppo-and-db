-- Триггер на автовыдачу прав при регистрастрации или изменении данных пользователя
CREATE OR REPLACE FUNCTION update_user_role()
    RETURNS TRIGGER AS $$
BEGIN
    IF NEW.role = 0 THEN
        -- Роль читателя (reader)
        EXECUTE format('ALTER ROLE user_%s SET ROLE Reader', NEW.id);
    ELSIF NEW.role = 1 THEN
        -- Роль автора (author)
        EXECUTE format('ALTER ROLE user_%s SET ROLE Author', NEW.id);
    ELSIF NEW.role = 2 THEN
        -- Роль администратора (admin)
        EXECUTE format('ALTER ROLE user_%s SET ROLE Admin', NEW.id);
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;


DROP TRIGGER IF EXISTS update_user_role_trigger ON main.users;
CREATE TRIGGER update_user_role_trigger
    AFTER INSERT OR UPDATE ON main.users
    FOR EACH ROW
EXECUTE PROCEDURE update_user_role();