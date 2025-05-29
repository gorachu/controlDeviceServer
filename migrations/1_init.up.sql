-- Таблица шкафов (shields)
CREATE TABLE IF NOT EXISTS shields (
    shield_id INTEGER PRIMARY KEY AUTOINCREMENT,
    number_in_list INTEGER NOT NULL UNIQUE,
    shipping_date DATE,
    customer TEXT,
    inspector TEXT,
    install_address TEXT,
    has_sim INTEGER NOT NULL DEFAULT 0,
    phone_num TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
    );

-- Контроллеры
CREATE TABLE IF NOT EXISTS controllers (
    controller_id INTEGER PRIMARY KEY AUTOINCREMENT,
    number_in_list INTEGER NOT NULL UNIQUE,
    imei TEXT NOT NULL UNIQUE,
    type TEXT,
    firmware TEXT,
    inspector TEXT,
    comment TEXT,
    in_shield INTEGER DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Input-модули
CREATE TABLE IF NOT EXISTS inputmodules (
    inputmodule_id INTEGER PRIMARY KEY AUTOINCREMENT,
    number_in_list INTEGER NOT NULL UNIQUE,
    type TEXT,
    firmware TEXT,
    inspector TEXT,
    comment TEXT,
    in_shield INTEGER DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- LCD модули
CREATE TABLE IF NOT EXISTS lcds (
    lcd_id INTEGER PRIMARY KEY AUTOINCREMENT,
    number_in_list INTEGER NOT NULL UNIQUE,
    type TEXT,
    firmware TEXT,
    inspector TEXT,
    comment TEXT,
    in_shield INTEGER DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Анализаторы тока
CREATE TABLE IF NOT EXISTS current_analyzers (
    current_analyzer_id INTEGER PRIMARY KEY AUTOINCREMENT,
    number_in_list INTEGER NOT NULL UNIQUE,
    type TEXT,
    firmware TEXT,
    inspector TEXT,
    comment TEXT,
    in_shield INTEGER DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);