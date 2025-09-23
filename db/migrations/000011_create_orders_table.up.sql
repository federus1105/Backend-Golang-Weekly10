-- public.orders definition

-- Drop table

-- DROP TABLE public.orders;

CREATE TABLE public.orders (
	id int4 GENERATED ALWAYS AS IDENTITY( INCREMENT BY 1 MINVALUE 1 MAXVALUE 2147483647 START 1 CACHE 1 NO CYCLE) NOT NULL,
	id_schedule int4 NOT NULL,
	id_user int4 NULL,
	id_payment_method int4 NULL,
	total numeric(10, 2) NULL,
	fullname varchar(255) NOT NULL,
	email varchar(255) NOT NULL,
	phone_number varchar(20) NOT NULL,
	created_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
	paid bool NOT NULL,
	CONSTRAINT orders_email_key UNIQUE (email),
	CONSTRAINT orders_pkey PRIMARY KEY (id)
);


-- public.orders foreign keys

ALTER TABLE public.orders ADD CONSTRAINT orders_id_payment_method_fkey FOREIGN KEY (id_payment_method) REFERENCES public.payment_method(id);
ALTER TABLE public.orders ADD CONSTRAINT orders_id_payment_method_fkey1 FOREIGN KEY (id_payment_method) REFERENCES public.payment_method(id);
ALTER TABLE public.orders ADD CONSTRAINT orders_id_schedule_fkey FOREIGN KEY (id_schedule) REFERENCES public.schedule(id);
ALTER TABLE public.orders ADD CONSTRAINT orders_id_user_fkey FOREIGN KEY (id_user) REFERENCES public.users(id);