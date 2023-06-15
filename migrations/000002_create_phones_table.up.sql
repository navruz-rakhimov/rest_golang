CREATE TABLE IF NOT EXISTS phones (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    phone TEXT NOT NULL,
    description TEXT,
    is_fax INTEGER NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users (id)
)