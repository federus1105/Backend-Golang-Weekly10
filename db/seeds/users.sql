INSERT INTO public.users (email,"password","role") VALUES
	 ('firdaus.example@gmail.com','$argon2id$v=19$m=65536,t=2,p=4$UYfOmhLHm0HzQ+y+JzQ5cw$+UQfB6M5eZJVPPrIAwmaQLykAmg0OjkGsA0ioLleMY0','Admin'::public."user_role"),
	 ('agus@example.com','$argon2id$v=19$m=65536,t=2,p=4$zNhxx0dm9Jk+ykMtGZ/uAg$I8skCF8IT9Zh36ooGOfYJ44hpg3WYSm7SMhIegR+d5k','User'::public."user_role"),
	 ('user@gmail.cp,','$argon2id$v=19$m=65536,t=2,p=4$4nIYP32YE39uQ1KQbODJOw$ER/vNFYtMouTjcn/SlFcnOvtajtGFHmQGjjTG7GIc30','User'::public."user_role");
