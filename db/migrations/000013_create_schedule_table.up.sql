-- public.schedule definition

-- Drop table

-- DROP TABLE public.schedule;

CREATE TABLE public.schedule (
	id int4 GENERATED ALWAYS AS IDENTITY( INCREMENT BY 1 MINVALUE 1 MAXVALUE 2147483647 START 1 CACHE 1 NO CYCLE) NOT NULL,
	id_movie int4 NOT NULL,
	"date" date NOT NULL,
	id_cinema int4 NULL,
	id_time int4 NULL,
	id_location int4 NULL,
	CONSTRAINT schedule_pkey PRIMARY KEY (id)
);


-- public.schedule foreign keys

ALTER TABLE public.schedule ADD CONSTRAINT schedule_id_cinema_fkey FOREIGN KEY (id_cinema) REFERENCES public.cinema(id);
ALTER TABLE public.schedule ADD CONSTRAINT schedule_id_location_fkey FOREIGN KEY (id_location) REFERENCES public."location"(id);
ALTER TABLE public.schedule ADD CONSTRAINT schedule_id_movie_fkey FOREIGN KEY (id_movie) REFERENCES public.movies(id) ON DELETE CASCADE;
ALTER TABLE public.schedule ADD CONSTRAINT schedule_id_time_fkey FOREIGN KEY (id_time) REFERENCES public."time"(id);