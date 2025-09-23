CREATE TABLE public.account (
	image varchar(255) NULL,
	firstname varchar(255) NULL,
	lastname varchar(255) NULL,
	phonenumber varchar(20) NULL,
	update_at timestamp NULL,
	point int4 DEFAULT 0 NULL,
	user_id int4 NOT NULL,
	CONSTRAINT account_pkey PRIMARY KEY (user_id)
);