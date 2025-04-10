-- Создание таблицы пользователей
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(64) NOT NULL UNIQUE,
    email VARCHAR(128) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    avatar_url VARCHAR(255) DEFAULT NULL,
    role VARCHAR(20) NOT NULL DEFAULT 'user',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
    );

-- Создание таблицы жанров
CREATE TABLE IF NOT EXISTS genres (
    id SERIAL PRIMARY KEY,
    name VARCHAR(64) NOT NULL UNIQUE
    );

-- Создание таблицы манги
CREATE TABLE IF NOT EXISTS manga (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    alter_title VARCHAR(255) DEFAULT NULL,
    description TEXT,
    cover_url VARCHAR(255) DEFAULT NULL,
    year INTEGER,
    status VARCHAR(20) NOT NULL DEFAULT 'ongoing',
    author VARCHAR(255) NOT NULL,
    artist VARCHAR(255) DEFAULT NULL,
    rating DECIMAL(3,2) DEFAULT 0.0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
    );

-- Таблица связи манги и жанров
CREATE TABLE IF NOT EXISTS manga_genres (
    manga_id INTEGER NOT NULL REFERENCES manga(id) ON DELETE CASCADE,
    genre_id INTEGER NOT NULL REFERENCES genres(id) ON DELETE CASCADE,
    PRIMARY KEY (manga_id, genre_id)
    );

-- Таблица глав
CREATE TABLE IF NOT EXISTS chapters (
    id SERIAL PRIMARY KEY,
    manga_id INTEGER NOT NULL REFERENCES manga(id) ON DELETE CASCADE,
    number DECIMAL(5,2) NOT NULL,
    title VARCHAR(255) NOT NULL,
    page_count INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (manga_id, number)
    );

-- Таблица страниц
CREATE TABLE IF NOT EXISTS pages (
    id SERIAL PRIMARY KEY,
    chapter_id INTEGER NOT NULL REFERENCES chapters(id) ON DELETE CASCADE,
    number INTEGER NOT NULL,
    image_url VARCHAR(255) NOT NULL,
    UNIQUE (chapter_id, number)
    );

-- Таблица закладок
CREATE TABLE IF NOT EXISTS bookmarks (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    manga_id INTEGER NOT NULL REFERENCES manga(id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (user_id, manga_id)
    );

-- Таблица истории чтения
CREATE TABLE IF NOT EXISTS read_history (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    manga_id INTEGER NOT NULL REFERENCES manga(id) ON DELETE CASCADE,
    chapter_id INTEGER NOT NULL REFERENCES chapters(id) ON DELETE CASCADE,
    page INTEGER NOT NULL DEFAULT 1,
    read_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (user_id, manga_id, chapter_id)
    );

-- Индексы для оптимизации запросов
CREATE INDEX idx_manga_title ON manga(title);
CREATE INDEX idx_manga_status ON manga(status);
CREATE INDEX idx_manga_year ON manga(year);
CREATE INDEX idx_manga_rating ON manga(rating);
CREATE INDEX idx_chapters_manga_id ON chapters(manga_id);
CREATE INDEX idx_pages_chapter_id ON pages(chapter_id);
CREATE INDEX idx_bookmarks_user_id ON bookmarks(user_id);
CREATE INDEX idx_read_history_user_id ON read_history(user_id);
CREATE INDEX idx_read_history_manga_id ON read_history(manga_id);

-- Вставка начальных жанров
INSERT INTO genres (name) VALUES
                              ('Боевик'),
                              ('Комедия'),
                              ('Драма'),
                              ('Фэнтези'),
                              ('Ужасы'),
                              ('Романтика'),
                              ('Научная фантастика'),
                              ('Сёнэн'),
                              ('Сёдзё'),
                              ('Психология'),
                              ('Приключения'),
                              ('Повседневность'),
                              ('Спорт'),
                              ('Детектив'),
                              ('Исторический')
    ON CONFLICT (name) DO NOTHING;

-- Добавление функции для обновления timestamp
CREATE OR REPLACE FUNCTION update_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Добавление триггеров для автоматического обновления updated_at
CREATE TRIGGER update_users_timestamp
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION update_timestamp();

CREATE TRIGGER update_manga_timestamp
    BEFORE UPDATE ON manga
    FOR EACH ROW
    EXECUTE FUNCTION update_timestamp();

CREATE TRIGGER update_chapters_timestamp
    BEFORE UPDATE ON chapters
    FOR EACH ROW
    EXECUTE FUNCTION update_timestamp();