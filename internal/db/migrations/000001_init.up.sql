CREATE TABLE songs (
    id BIGSERIAL PRIMARY KEY,
    group_title VARCHAR(255) NOT NULL,
    song_title VARCHAR(255) NOT NULL,
    release_date DATE NOT NULL,
    song_text TEXT NOT NULL,
    link VARCHAR(255) NOT NULL
);