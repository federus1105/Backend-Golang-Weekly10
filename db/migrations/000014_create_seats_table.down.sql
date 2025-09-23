-- public.seats definition

-- Drop table

-- DROP TABLE public.seats;

CREATE TABLE public.seats (
	id int4 GENERATED ALWAYS AS IDENTITY( INCREMENT BY 1 MINVALUE 1 MAXVALUE 2147483647 START 1 CACHE 1 NO CYCLE) NOT NULL,
	codeseat varchar(255) NOT NULL,
	isstatus bool NULL,
	CONSTRAINT seats_pkey PRIMARY KEY (id)
);