
CREATE TABLE songs_table (
  id SERIAL PRIMARY KEY,
  group VARCHAR(255) NOT NULL,
  song VARCHAR(255) NOT NULL, --name
  release_date DATE NOT NULL,
  text TEXT NOT NULL,
  link TEXT DEFAULT NULL
  );