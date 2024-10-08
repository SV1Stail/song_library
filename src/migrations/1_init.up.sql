CREATE TABLE songs_table (
  id SERIAL PRIMARY KEY,
  band VARCHAR(255) NOT NULL,
  song VARCHAR(255) NOT NULL,
  release_date VARCHAR(255) NOT NULL,
  text TEXT NOT NULL,
  link TEXT DEFAULT NULL
);


INSERT INTO songs_table (band, song, release_date, text, link) VALUES
('Queen', 'Bohemian Rhapsody', '1975', 'Is this the real life

nsomethingsong

Inareyou ready

  I dont hear you\n\nnready?\n', 'https://example.com/bohemian_rhapsody'),
('The Beatles', 'Hey Jude', '1968', 'Hey Jude, don''t make it bad 

work hard\nwork hard 

go push it

du hast 

duhast\nduhast

hast', 'https://example.com/hey_jude');



