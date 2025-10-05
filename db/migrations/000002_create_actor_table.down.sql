-- public.actor definition

-- Drop table

-- DROP TABLE public.actor;

CREATE TABLE actor (
	id int4 GENERATED ALWAYS AS IDENTITY( INCREMENT BY 1 MINVALUE 1 MAXVALUE 2147483647 START 1 CACHE 1 NO CYCLE) NOT NULL,
	name varchar(255) NOT NULL,
	CONSTRAINT actor_pkey PRIMARY KEY (id)
);