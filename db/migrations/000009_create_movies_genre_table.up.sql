-- public.movies_genre definition

-- Drop table

-- DROP TABLE public.movies_genre;

CREATE TABLE public.movies_genre (
	id_genre int4 NOT NULL,
	id_movies int4 NOT NULL
);


-- public.movies_genre foreign keys

ALTER TABLE public.movies_genre ADD CONSTRAINT movies_genre_id_genre_fkey FOREIGN KEY (id_genre) REFERENCES public.genres(id);
ALTER TABLE public.movies_genre ADD CONSTRAINT movies_genre_id_movies_fkey FOREIGN KEY (id_movies) REFERENCES public.movies(id) ON DELETE CASCADE;