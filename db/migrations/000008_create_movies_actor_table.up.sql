-- public.movies_actor definition

-- Drop table

-- DROP TABLE public.movies_actor;

CREATE TABLE movies_actor (
	id_movie int4 NOT NULL,
	id_actor int4 NOT NULL
);


-- public.movies_actor foreign keys

ALTER TABLE movies_actor ADD CONSTRAINT movies_actor_id_actor_fkey FOREIGN KEY (id_actor) REFERENCES actor(id);
ALTER TABLE movies_actor ADD CONSTRAINT movies_actor_id_movie_fkey FOREIGN KEY (id_movie) REFERENCES movies(id) ON DELETE CASCADE;