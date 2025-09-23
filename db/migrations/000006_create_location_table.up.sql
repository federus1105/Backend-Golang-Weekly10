-- public."location" definition

-- Drop table

-- DROP TABLE public."location";

CREATE TABLE public."location" (
	id int4 GENERATED ALWAYS AS IDENTITY( INCREMENT BY 1 MINVALUE 1 MAXVALUE 2147483647 START 1 CACHE 1 NO CYCLE) NOT NULL,
	"name" varchar(255) NOT NULL,
	CONSTRAINT location_pkey PRIMARY KEY (id)
);