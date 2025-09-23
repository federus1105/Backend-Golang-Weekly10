-- public.order_seat definition

-- Drop table

-- DROP TABLE public.order_seat;

CREATE TABLE public.order_seat (
	id_order int4 NOT NULL,
	id_seats int4 NOT NULL
);


-- public.order_seat foreign keys

ALTER TABLE public.order_seat ADD CONSTRAINT order_seat_id_order_fkey FOREIGN KEY (id_order) REFERENCES public.orders(id);
ALTER TABLE public.order_seat ADD CONSTRAINT order_seat_id_seats_fkey FOREIGN KEY (id_seats) REFERENCES public.seats(id);