-- public.movies definition

-- Drop table

-- DROP TABLE public.movies;

CREATE TABLE movies (
	id int4 GENERATED ALWAYS AS IDENTITY( INCREMENT BY 1 MINVALUE 1 MAXVALUE 2147483647 START 1 CACHE 1 NO CYCLE) NOT NULL,
	image varchar(255) NULL,
	backdrop varchar(255) NULL,
	title varchar(255) NOT NULL,
	release_date date NOT NULL,
	duration varchar(255) NOT NULL,
	id_director int4 NOT NULL,
	synopsis text NOT NULL,
	rating float8 NOT NULL,
	is_deleted bool DEFAULT false NULL,
	CONSTRAINT movies_pkey PRIMARY KEY (id)
);


-- public.movies foreign keys

ALTER TABLE public.movies ADD CONSTRAINT movies_id_director_fkey FOREIGN KEY (id_director) REFERENCES public.director(id);