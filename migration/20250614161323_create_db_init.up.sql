

CREATE TABLE IF NOT EXISTS public.users (
	id SERIAL PRIMARY KEY,
	username VARCHAR NOT NULL DEFAULT '',
	"password" VARCHAR NOT NULL,
	salt varchar DEFAULT ''::character varying NOT NULL,
	"status" VARCHAR NOT NULL DEFAULT 'active',
	"role" varchar DEFAULT 'employee'::character varying NOT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
	created_by int DEFAULT 0 NOT NULL,
	updated_by int DEFAULT 0 NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_users_status ON public.users ("status");
CREATE INDEX IF NOT EXISTS idx_users_username ON public.users (username);


CREATE TABLE IF NOT EXISTS public.employees (
	id serial NOT NULL,
	fullname varchar DEFAULT '' NOT NULL,
	salary float8 DEFAULT 0 NOT NULL,
	code varchar DEFAULT '' NOT NULL,
	created_at timestamptz DEFAULT now() NOT NULL,
	updated_at timestamptz DEFAULT now() NOT NULL,
	created_by int DEFAULT 0 NOT NULL,
	updated_by int DEFAULT 0 NOT NULL,
	CONSTRAINT employees_pk PRIMARY KEY (id)
);
CREATE INDEX IF NOT EXISTS employees_code_idx ON public.employees (code);


CREATE TABLE IF NOT EXISTS public.attendances (
	id serial4 NOT NULL,
	employee_id int4 DEFAULT 0 NOT NULL,
	"date" date NOT NULL,
	checked_in_at timestamptz DEFAULT now() NOT NULL,
	checked_out_at timestamptz NULL,
	created_at timestamptz DEFAULT now() NOT NULL,
	updated_at timestamptz DEFAULT now() NOT NULL,
	created_by int4 DEFAULT 0 NOT NULL,
	updated_by int4 DEFAULT 0 NOT NULL,
	CONSTRAINT attendances_pk PRIMARY KEY (id),
	CONSTRAINT attendances_unique UNIQUE (employee_id, date)
);
CREATE INDEX IF NOT EXISTS attendances_date_idx ON public.attendances USING btree (date);

CREATE TABLE IF NOT EXISTS public.overtimes (
	id serial4 NOT NULL,
	employee_id int4 DEFAULT 0 NOT NULL,
	"date" date NOT NULL,
	hours int4 DEFAULT 0 NOT NULL,
	"status" varchar DEFAULT 'validated'::character varying NOT NULL,
	created_at timestamptz DEFAULT now() NOT NULL,
	updated_at timestamptz DEFAULT now() NOT NULL,
	created_by int4 DEFAULT 0 NOT NULL,
	updated_by int4 DEFAULT 0 NOT NULL,
	CONSTRAINT overtime_pk PRIMARY KEY (id)
);
CREATE INDEX IF NOT EXISTS overtimes_date_idx ON public.overtimes USING btree (date);

CREATE TABLE public.reimbursements (
	id serial4 NOT NULL,
	employee_id int4 DEFAULT 0 NOT NULL,
	"date" date NOT NULL,
	amount float8 DEFAULT 0 NOT NULL,
	"description" text DEFAULT ''::text NOT NULL,
	created_at timestamptz DEFAULT now() NOT NULL,
	updated_at timestamptz DEFAULT now() NOT NULL,
	created_by int4 DEFAULT 0 NOT NULL,
	updated_by int4 DEFAULT 0 NOT NULL,
	payslip_id int4 DEFAULT 0 NOT NULL,
	status varchar DEFAULT 'validated'::character varying NOT NULL,
	CONSTRAINT reimbursements_pk PRIMARY KEY (id)
);
CREATE INDEX reimbursements_date_idx ON public.reimbursements USING btree (date);

CREATE TABLE IF NOT EXISTS public.payslips (
	id serial4 NOT NULL,
	employee_id int4 DEFAULT 0 NOT NULL,
	period_id int4 DEFAULT 0 NOT NULL,
	total_attendance_days int4 DEFAULT 0 NOT NULL,
	total_overtime_hours int4 DEFAULT 0 NOT NULL,
	total_attendance_salary float8 DEFAULT 0 NOT NULL,
	total_overtime_salary float8 DEFAULT 0 NOT NULL,
	total_reimbursement float8 DEFAULT 0 NOT NULL,
	total_take_home_pay float8 DEFAULT 0 NOT NULL,
	created_at timestamptz DEFAULT now() NOT NULL,
	updated_at timestamptz DEFAULT now() NOT NULL,
	created_by int4 DEFAULT 0 NOT NULL,
	updated_by int4 DEFAULT 0 NOT NULL,
	CONSTRAINT payslips_pk PRIMARY KEY (id),
	CONSTRAINT payslips_unique UNIQUE (employee_id, period_id)
);
CREATE INDEX IF NOT EXISTS payslips_employee_id_idx ON public.payslips USING btree (employee_id);


CREATE TABLE IF NOT EXISTS public.audit_logs (
	id serial4 NOT NULL,
	table_name varchar NOT NULL,
	"action" varchar NOT NULL,
	record_id int4 NOT NULL,
	old_data jsonb NULL,
	new_data jsonb NULL,
	changed_by int4 NOT NULL,
	changed_at timestamptz DEFAULT now() NOT NULL,
	ip_address varchar DEFAULT ''::character varying NOT NULL,
	request_id varchar DEFAULT ''::character varying NOT NULL,
	CONSTRAINT audit_logs_pkey PRIMARY KEY (id),
	CONSTRAINT audit_logs_unique UNIQUE (table_name, action, record_id)
);


CREATE TABLE public.attendances_period (
	id serial4 NOT NULL,
	period_start date DEFAULT now() NOT NULL,
	period_end date DEFAULT now() NOT NULL,
	is_payslip_generated bool DEFAULT false NULL,
	created_at timestamptz DEFAULT now() NOT NULL,
	updated_at timestamptz DEFAULT now() NOT NULL,
	created_by int4 DEFAULT 0 NOT NULL,
	updated_by int4 DEFAULT 0 NOT NULL,
	CONSTRAINT attendances_period_pk PRIMARY KEY (id)
);

INSERT INTO public.users (username,"password",status,created_at,updated_at,"role",salt,created_by,updated_by) VALUES
	 ('user1','XNJCb3D9+6V6YuXmeiwJaD2YuhWA1fxrYCwHEZa6p10=','active','2025-06-15 23:21:02.362','2025-06-15 23:21:02.362','employee','BIqwAg6TvxpX/3bc6KE5Ng==',0,0),
	 ('user2','iwk7EcGRKBOvhteqnjukgnrqzdNbkr8MkaivinbdBiQ=','active','2025-06-15 23:21:03.362','2025-06-15 23:21:03.362','employee','lO2whJ/dU6xRjhju+kYPpw==',0,0),
	 ('user3','cHbyEVYVUSEZR4BpJVfrU05//gdlZIsI9RczL4QHzU4=','active','2025-06-15 23:21:04.362','2025-06-15 23:21:04.362','employee','v0RzUswe4e3rNoM/kaHHWg==',0,0),
	 ('user4','xkvi2+l/y17mB3FVxULKaiG4KL/L01cM23Synv8raIE=','active','2025-06-15 23:21:05.362','2025-06-15 23:21:05.362','employee','ENpjbttXJDkGPjV4ablfmA==',0,0),
	 ('user5','SoFSj5xk79/bRx4D3BMMArShDr7iRFaKGpSJbLxJ+jM=','active','2025-06-15 23:21:06.362','2025-06-15 23:21:06.362','employee','9LR6/tDAgcWZYhnVEZzMhg==',0,0),
	 ('user6','jZias8HGMnuxwoCwPGfkOQDb+hw9o60Mvv4oxjND6QE=','active','2025-06-15 23:21:07.362','2025-06-15 23:21:07.362','employee','59rsoP2czo+uZ8jWeH5iGw==',0,0),
	 ('user7','M8lXgcxgEvn8pahVL5zAhtmbdd04V3NwaPs2WJ1eh4Q=','active','2025-06-15 23:21:08.362','2025-06-15 23:21:08.362','employee','g/kjgWGrBvPOUMFG4kc9Fw==',0,0),
	 ('user8','O2BwhsLCQqxHPsSdWASKqHzRyvZ0pG0d2gH8zURyee0=','active','2025-06-15 23:21:09.362','2025-06-15 23:21:09.362','employee','zR825p6FlvKUTugqDAPf6A==',0,0),
	 ('user9','VgfIZM9pER4vKESwfpvOfqqmkwquHBQrbssNn1s3DJE=','active','2025-06-15 23:21:10.362','2025-06-15 23:21:10.362','employee','CLTVNBiceXBDBcb50xWqJA==',0,0),
	 ('user10','Sg/4UaS2Eap1jdM/tt6mZ+I+ydCYRjyIJ6LBRLD6gME=','active','2025-06-15 23:21:11.362','2025-06-15 23:21:11.362','employee','ohj/oiES1AQsE8lIp3R6Ig==',0,0);
INSERT INTO public.users (username,"password",status,created_at,updated_at,"role",salt,created_by,updated_by) VALUES
	 ('user11','2wW6fNZMShTvJq1uFE12yO7Ssw4VyxMgPOckqkbcScA=','active','2025-06-15 23:21:12.362','2025-06-15 23:21:12.362','employee','QEnGX/yvrTqjawZJhRTXGg==',0,0),
	 ('user12','rPnVLclbCeLbkxVdh8mMM+04h83qpXUkNj3M3Q+A7JY=','active','2025-06-15 23:21:13.362','2025-06-15 23:21:13.362','employee','5+RJ35PLrjjk0dlQh0TNtw==',0,0),
	 ('user13','YyDXfr5cfadfrKUHKcljOaC4BDgX/iy7ZzEfWQ0Q4NQ=','active','2025-06-15 23:21:14.362','2025-06-15 23:21:14.362','employee','UY97yU0i7flhOOPLhXPJkA==',0,0),
	 ('user14','hn/AsCwvxtqlpbLDrMEXZh20FcSgX40+32r3mJfuIrc=','active','2025-06-15 23:21:15.362','2025-06-15 23:21:15.362','employee','5YqhlLh6OwXNLUZjgUqwzQ==',0,0),
	 ('user15','jO9vg6CKtOCTo2NY3hfecMa7PL/3TPy4FmpeZNsdnYE=','active','2025-06-15 23:21:16.362','2025-06-15 23:21:16.362','employee','LvUYykIKkDM8I3Nhk5iKEQ==',0,0),
	 ('user16','LTjmTMdU4Z+W7bnaWCuehXw0HE4QRMVBQus2W79fQ/Y=','active','2025-06-15 23:21:17.362','2025-06-15 23:21:17.362','employee','2A7jHIYGkClCYfc4dlZcqQ==',0,0),
	 ('user17','ANrfa6iglPAriCKb1jWp6p96VIyJccXcj0n/Rdekoy8=','active','2025-06-15 23:21:18.362','2025-06-15 23:21:18.362','employee','1Cdp/ecQvEqv+vnObvx29Q==',0,0),
	 ('user18','8lYZTf2RnqZqnQlONCtmOzq2T7CicaWXe2OcEuV7oJk=','active','2025-06-15 23:21:19.362','2025-06-15 23:21:19.362','employee','1fS4yGVKiVT9j0dCIJlr5g==',0,0),
	 ('user19','XJxPwESU/c9RgETEWprWJsIQnodiFSdIQPsRPhXTgj4=','active','2025-06-15 23:21:20.362','2025-06-15 23:21:20.362','employee','knjWIpTz75O2H/GGRUJCBg==',0,0),
	 ('user20','1X1CQ8qdVKvcN0vbDJVHmxbk+hx3aP1zowosAioXJEI=','active','2025-06-15 23:21:21.362','2025-06-15 23:21:21.362','employee','EoAmeA+pU8Az0sBfOoFJ3Q==',0,0);
INSERT INTO public.users (username,"password",status,created_at,updated_at,"role",salt,created_by,updated_by) VALUES
	 ('user21','NW0rA2AKCe1Hd9Fz80stbCmsAB438n8IEVf3zWPmQFA=','active','2025-06-15 23:21:22.362','2025-06-15 23:21:22.362','employee','5JfSmuBbG1WlM6fkchnJzQ==',0,0),
	 ('user22','iaCaBl7bbcqdiSlaIGjIQhZw9dHGLmlzmI4u1QIlvdc=','active','2025-06-15 23:21:23.362','2025-06-15 23:21:23.362','employee','wYcVe63xIys64dLptblsvg==',0,0),
	 ('user23','75dc7i9uiNoQLsv3G8ZRXVW0NJiTW8Nt2HCGWjFJvNY=','active','2025-06-15 23:21:24.362','2025-06-15 23:21:24.362','employee','eRi4vTRlo4UuLyswYN36lw==',0,0),
	 ('user24','SOBcE8EuiO7RyggnkP0MfQyzfibpn5dnmROkFNpUe0E=','active','2025-06-15 23:21:25.362','2025-06-15 23:21:25.362','employee','r7bavE1HhzelocvYwDcYiQ==',0,0),
	 ('user25','PY5bWEkC4LGlZ/vGoNiuwfbmSSsJ+b7khQJjWIdaQvA=','active','2025-06-15 23:21:26.362','2025-06-15 23:21:26.362','employee','G/IqFBWBrT3I+LIntHibqQ==',0,0),
	 ('user26','jgjReSSyZPnQhcg6IAtP4gSH0tRUsx2fnLF/KB4rPlM=','active','2025-06-15 23:21:27.362','2025-06-15 23:21:27.362','employee','eWX98yPJ3vcSqQODhThO/g==',0,0),
	 ('user27','KUPvv2ESW6HP8DFtu6WTyUu7YtIbTw15xpZLHXTF8O0=','active','2025-06-15 23:21:28.362','2025-06-15 23:21:28.362','employee','AITAe+RbJIAqg4ftNQtyag==',0,0),
	 ('user28','NKPnEM1zPALKimzXya0xdgQ0cCnWx9caKlsiiNkX+dA=','active','2025-06-15 23:21:29.362','2025-06-15 23:21:29.362','employee','p8kSMRaaKFwS19zZBqhYPA==',0,0),
	 ('user29','ISEfYoV+u3l3N2AllNE/9Z+uMGShGFEXsrrH+arv3CE=','active','2025-06-15 23:21:30.362','2025-06-15 23:21:30.362','employee','1zC+bMEU7AqIag2v+3bPOA==',0,0),
	 ('user30','GSCbD9KRv0ltM6k4FATT2VHrrW97sepr6UHHKOv4RWk=','active','2025-06-15 23:21:31.362','2025-06-15 23:21:31.362','employee','3h5/bv0lazl9mhvPwcdNzg==',0,0);
INSERT INTO public.users (username,"password",status,created_at,updated_at,"role",salt,created_by,updated_by) VALUES
	 ('user31','cVHavcgZUbXcLVrGQdKhAVyaKX6gg4L7SRem9qoICzg=','active','2025-06-15 23:21:32.362','2025-06-15 23:21:32.362','employee','TqKdISqxp1V5RZ+ouy4sqA==',0,0),
	 ('user32','ezV4i4aaFWuELleR+wndOClqcfHHmwJOtI7UVEOv3Xs=','active','2025-06-15 23:21:33.362','2025-06-15 23:21:33.362','employee','xakBS/uey1pIKqYjewdG8A==',0,0),
	 ('user33','/x+UIXGZkWjCXlkxAjtxfP4JcKL8TzM0qcBvCuLdlNc=','active','2025-06-15 23:21:34.362','2025-06-15 23:21:34.362','employee','n+8FaM50Jszp5KtUi/mP7Q==',0,0),
	 ('user34','O00N3y8votnaW9/05DrmUcKgIjGmI8omfxDe/QLh54Y=','active','2025-06-15 23:21:35.362','2025-06-15 23:21:35.362','employee','B0yMQy09O703+FnHEg2+zw==',0,0),
	 ('user35','CHAwAdp0Wl0oZw1G3jIVeu+C9ZaP1jPSz2ZtSNz4RIs=','active','2025-06-15 23:21:36.362','2025-06-15 23:21:36.362','employee','ECl3qo2hz/RL97G5Opz51A==',0,0),
	 ('user36','ynVsUE4wY8ge4AGxROXo0EwgOVcmGeO0jBNsD64fhXk=','active','2025-06-15 23:21:37.362','2025-06-15 23:21:37.362','employee','gVVUsjbJMJM/6zGd0RcdOQ==',0,0),
	 ('user37','myaMPGOE1rftJ3I+vlQgxiwvhOewXPPKWCqASZnNGAk=','active','2025-06-15 23:21:38.362','2025-06-15 23:21:38.362','employee','uk6rPOnohc/KexKhcqXVKA==',0,0),
	 ('user38','1drI1Ey1UMReH64byNm/Zo3ND0+G0b8wH7AN1VhMv00=','active','2025-06-15 23:21:39.362','2025-06-15 23:21:39.362','employee','kUPBFhKhx7HdwHV5yXXi6w==',0,0),
	 ('user39','ub5JlO6VXH//L1QBiN/btPCCSYTA36aW61Z5daHPxB0=','active','2025-06-15 23:21:40.362','2025-06-15 23:21:40.362','employee','5mvRORdjaUbROc2bxvNbdQ==',0,0),
	 ('user40','j8TzDoDkhLeHd6XIqQ3wE2Cjk9flQ+GTl8Q+W90RL4w=','active','2025-06-15 23:21:41.362','2025-06-15 23:21:41.362','employee','56JrvzDWYH0ayBKbCYKPWw==',0,0);
INSERT INTO public.users (username,"password",status,created_at,updated_at,"role",salt,created_by,updated_by) VALUES
	 ('user41','IhseF7fW0SwXsVf9nMy00aMqiGzQ9+4Iy7PzDX0tl34=','active','2025-06-15 23:21:42.362','2025-06-15 23:21:42.362','employee','ZqVdJh8QAGiX+eHuonzjqw==',0,0),
	 ('user42','L/nyDLnVKI5m4TzBd5XOCg8KHl/Vs43sVnpPa1Lx2Xw=','active','2025-06-15 23:21:43.362','2025-06-15 23:21:43.362','employee','rSWUi6gjhtdK5CTzBfHChQ==',0,0),
	 ('user43','lMKXS81fcQMyTnYrufhTY3SU8LJ1uAz56HGG18/FLSY=','active','2025-06-15 23:21:44.362','2025-06-15 23:21:44.362','employee','WafYBML7P2qraHu/c8kduw==',0,0),
	 ('user44','DXKflea8OlQLJuSngn9DQdqXHyX7tjYnkqUN6JxVwzU=','active','2025-06-15 23:21:45.362','2025-06-15 23:21:45.362','employee','9vTnnxSxYy1rF+ewkatQaA==',0,0),
	 ('user45','nOmcvgZOjka5W3G9f4dUpw392C9sZ15+V0LFcmn1ZPU=','active','2025-06-15 23:21:46.362','2025-06-15 23:21:46.362','employee','P+fR0EXcTCr0Fhi5tYTjjA==',0,0),
	 ('user46','H1P8RdOxZiaoIjMUf3IDpCJKmx+nbgFHmZrR1ztiSgY=','active','2025-06-15 23:21:47.362','2025-06-15 23:21:47.362','employee','ehTTyk4QPNw0/1hn07pbSQ==',0,0),
	 ('user47','kpTA/+SXVlnJwHUUVV+MufOQpkdz5LLsVDsVsBgTeU0=','active','2025-06-15 23:21:48.362','2025-06-15 23:21:48.362','employee','7mye7gyHhS0hWrnX/yAsyA==',0,0),
	 ('user48','wAMhPqmEbfZrfvDjzjpJDuhfii5dq2/mXbTuL6wXWaU=','active','2025-06-15 23:21:49.362','2025-06-15 23:21:49.362','employee','Owd1d+WZkHXCxg1lzUnD6w==',0,0),
	 ('user49','owX4iSREFcpKhl+to3psoOTy7U5QJTnLEXwaAIpr7jU=','active','2025-06-15 23:21:50.362','2025-06-15 23:21:50.362','employee','taJQ5WV9xxJ83MjkF/CwHQ==',0,0),
	 ('user50','yKt056p147v71lqB/0V4vqwc26311nw50zvisA1G7Qo=','active','2025-06-15 23:21:51.362','2025-06-15 23:21:51.362','employee','tVzdEoJuEO+/MVNca9PIgg==',0,0);
INSERT INTO public.users (username,"password",status,created_at,updated_at,"role",salt,created_by,updated_by) VALUES
	 ('user51','5ZWEodV4KNH6NIUMeoBHbufrJcVS8f+0unL8z2xqyFg=','active','2025-06-15 23:21:52.362','2025-06-15 23:21:52.362','employee','inQI0fTFQ5Rrl1pGCSXmWg==',0,0),
	 ('user52','fGq4FYV0LBE+L3wGpqAGaNguSVjRQav6QIud97WUc+A=','active','2025-06-15 23:21:53.362','2025-06-15 23:21:53.362','employee','Rf0iLTFyiCw+hIO1ttbBSQ==',0,0),
	 ('user53','lO3TfjbJc+nSFuJ0EgA5lbC8ZIfnDLQKZMo9Bxo2LQc=','active','2025-06-15 23:21:54.362','2025-06-15 23:21:54.362','employee','YDlZYua7BV8GOMUFKcKpgg==',0,0),
	 ('user54','YU1ivENPR8DmtssrmhkCaPqqA9u2orkyu4+cpzWF8dU=','active','2025-06-15 23:21:55.362','2025-06-15 23:21:55.362','employee','ea6Ik9nWSaIdBt0m5HPuvQ==',0,0),
	 ('user55','x7Lk8Rt/xl2xGDeYIOpnVfzDmrecOx5Ze1mAqiNnw4c=','active','2025-06-15 23:21:56.362','2025-06-15 23:21:56.362','employee','gvjDRVCoHD7SsffFNh8ycA==',0,0),
	 ('user56','mmqvB8kz89eTxJoTal5cx28i1vFjUOMqJCXwuf0Gusg=','active','2025-06-15 23:21:57.362','2025-06-15 23:21:57.362','employee','dMCLofB3iygHx7xhHC704g==',0,0),
	 ('user57','srtUMSTUOWEICirgc5APdM6xx5AjH0RHCFV1fKNsTcQ=','active','2025-06-15 23:21:58.362','2025-06-15 23:21:58.362','employee','rQNPN+SWFCUs0laphf1FdQ==',0,0),
	 ('user58','rouyga6dVHEM0dHctn2El4NCLBRc75ipQlJAfLAcY3c=','active','2025-06-15 23:21:59.362','2025-06-15 23:21:59.362','employee','qYDmyY4VRTadD8roSfPT4w==',0,0),
	 ('user59','x8ZLcMtid5vhVPMyiBMCUpX1yopzBaP166V1jrTdQCk=','active','2025-06-15 23:22:00.362','2025-06-15 23:22:00.362','employee','P5Kkv/XkYJ/cfP8T9GWIEA==',0,0),
	 ('user60','Ats9MmSv2tTQ05Ze0RZML5OYItpiM3jyltue+ZI/ulY=','active','2025-06-15 23:22:01.362','2025-06-15 23:22:01.362','employee','sjMmCTMrxHo4U0QQRVoZHw==',0,0);
INSERT INTO public.users (username,"password",status,created_at,updated_at,"role",salt,created_by,updated_by) VALUES
	 ('user61','5QC+yzOgYsyiaKYX3kXW5OBicm6ZPcOqpW0/WDd9p0w=','active','2025-06-15 23:22:02.362','2025-06-15 23:22:02.362','employee','YhhYC3FQDVOr/xH099ZQLQ==',0,0),
	 ('user62','ycBIl0ekYUxoInTlM5pob2vUMoYA+kclOLrBW83w/K4=','active','2025-06-15 23:22:03.362','2025-06-15 23:22:03.362','employee','CeOjPjdetC0wxQEJFl9nlg==',0,0),
	 ('user63','VAnLnRxB/mkOA5+JY98YcmOiWlKtDy3k8UlUtSQ/YLA=','active','2025-06-15 23:22:04.362','2025-06-15 23:22:04.362','employee','naQM5xQNrOTZuhtCwcGi1A==',0,0),
	 ('user64','Db6UhMiyvvE5RVkyIgKcQWyM3usKtCkGffg+6KD0gEg=','active','2025-06-15 23:22:05.362','2025-06-15 23:22:05.362','employee','Ty6yuQC8gjZ8zXcle5HmBw==',0,0),
	 ('user65','7D6KpVw0zVjc39ijBrxKqXkrx9fOAICtrdLHoehEYyI=','active','2025-06-15 23:22:06.362','2025-06-15 23:22:06.362','employee','MCen+AD+RAH33YmPzPBwEw==',0,0),
	 ('user66','T0N/tLoreHnsaw+gcPDLY0nm0uCygGgdCPJsmICDOGg=','active','2025-06-15 23:22:07.362','2025-06-15 23:22:07.362','employee','0YE5xlJvnRqB1B5gEcg0kA==',0,0),
	 ('user67','y5lM+JwkteU2dbT60d+y2FcsJSRUQ1kqi0pbnJIlQ50=','active','2025-06-15 23:22:08.362','2025-06-15 23:22:08.362','employee','dv8ElBdDHylMtg2cT1PwaQ==',0,0),
	 ('user68','J39wHhzPdpNHMr4NoSs+hDcak9OzELr6FNHeh3zAiEA=','active','2025-06-15 23:22:09.362','2025-06-15 23:22:09.362','employee','7+L92nPGoYTpsZx07S+CNA==',0,0),
	 ('user69','lEtoVTaQDqb2ekMgkgtFT3Z0douPmONJKJQU2aNwOF0=','active','2025-06-15 23:22:10.362','2025-06-15 23:22:10.362','employee','se81xKnGfvpP6ORv2CqqYg==',0,0),
	 ('user70','CeSnGzL7XM2WEd7c/bhesmkF9IxUESTti5nilqfTDZQ=','active','2025-06-15 23:22:11.362','2025-06-15 23:22:11.362','employee','gITNTbAW68cphDSQK6mUOA==',0,0);
INSERT INTO public.users (username,"password",status,created_at,updated_at,"role",salt,created_by,updated_by) VALUES
	 ('user71','qVRyY20TvD4AtiEjPicCQw0l8k+nWbi5j1r+GLTNmxI=','active','2025-06-15 23:22:12.362','2025-06-15 23:22:12.362','employee','A2YdOguPoTf37+V7GIWTDw==',0,0),
	 ('user72','kJeC/JkVTHJ3XgtNxm1JnHvgTD6vETPCVrZmqe16VZ0=','active','2025-06-15 23:22:13.362','2025-06-15 23:22:13.362','employee','xReIbHk7bxv/TeSFkWEUkQ==',0,0),
	 ('user73','8MvOH8BAlZZYqQNHhwZ31hk0R5TASDr3u1+LGumnpAw=','active','2025-06-15 23:22:14.362','2025-06-15 23:22:14.362','employee','EoiZR1sdRdGY7kNk9eznlA==',0,0),
	 ('user74','vGVH3y+8kpXPpLnQii9ZJ0WAstjIEF1kfKaNqyK48ls=','active','2025-06-15 23:22:15.362','2025-06-15 23:22:15.362','employee','cWKuFNxOR12u7JkiZPFMxw==',0,0),
	 ('user75','OK/VTksvRj2vQX7tbR2mMY5AQ2dYMLM8MMCN3D1iZkE=','active','2025-06-15 23:22:16.362','2025-06-15 23:22:16.362','employee','McPbKRcboYpk5RnxuCKfHg==',0,0),
	 ('user76','W7jZYD7714RyZ+JdpadNkKu7dh/NScqEL+VKpK0egCE=','active','2025-06-15 23:22:17.362','2025-06-15 23:22:17.362','employee','JJzXr0WTXrVLL7UW0h2FZw==',0,0),
	 ('user77','pNGSOYsLLC3opNTNLsYjGpPl/AGCoor3L8N1yp51gxI=','active','2025-06-15 23:22:18.362','2025-06-15 23:22:18.362','employee','TZ4bUIdzlnaifX36ILFvNg==',0,0),
	 ('user78','dHngLoazbTbt8Jddif3iwgtHeayiji5AkepYdrgkrOo=','active','2025-06-15 23:22:19.362','2025-06-15 23:22:19.362','employee','kxOvZxSdPwDLWedMu1Ebsw==',0,0),
	 ('user79','WZumcPPrs+Ew4PAe0+FvdjxBJtKeonN1Ojv5J/Jkg08=','active','2025-06-15 23:22:20.362','2025-06-15 23:22:20.362','employee','SPTxkAO+LdmjAYNyfFMK3A==',0,0),
	 ('user80','xdv88ccjkW+vwNhUduBWVfyM+PxU/TviFpzsc+MLfkg=','active','2025-06-15 23:22:21.362','2025-06-15 23:22:21.362','employee','rTJyLih4VDZH2uc+lj+dhg==',0,0);
INSERT INTO public.users (username,"password",status,created_at,updated_at,"role",salt,created_by,updated_by) VALUES
	 ('user81','rP+aOCm1azvGwCOcWlTW9mFcK5GSlYUBN7ynss2qKNI=','active','2025-06-15 23:22:22.362','2025-06-15 23:22:22.362','employee','hyXz4QSc85j1mS3aHc3eAA==',0,0),
	 ('user82','IzbtFl//zDIocxTFpi1pTK9yg6qPa3GB8rFlk7MZy80=','active','2025-06-15 23:22:23.362','2025-06-15 23:22:23.362','employee','nb4/R4Npln+JLOlvwlhW4A==',0,0),
	 ('user83','UGXBaRsBtLHeF+E9q3w/54a99+OSUS6jALUlWvcBYi0=','active','2025-06-15 23:22:24.362','2025-06-15 23:22:24.362','employee','ck4Of/FliCTjcjCyBF9KQA==',0,0),
	 ('user84','dAxpnk6eAI4Ug25b2QkQ3QeV+9hxdTsHUWTTscAu/PU=','active','2025-06-15 23:22:25.362','2025-06-15 23:22:25.362','employee','XNNQ58AUaH4M4l0ipZSpFQ==',0,0),
	 ('user85','A70MGyOOKUVdWGUmj/bt4rrH0R93uWSVZ8h9oXQlqek=','active','2025-06-15 23:22:26.362','2025-06-15 23:22:26.362','employee','GuLWqJlQUoXln7Pd10XdIw==',0,0),
	 ('user86','5GmKxsCNhqudA5qFlKFeHSQ0IdKD1B4julo37K2taRE=','active','2025-06-15 23:22:27.362','2025-06-15 23:22:27.362','employee','FGUPlcLZKQLJ5nLSBImHNA==',0,0),
	 ('user87','z5UYliktiGiD7Fa2BlTAzO3W/LI32jcJGZqyUQcUheI=','active','2025-06-15 23:22:28.362','2025-06-15 23:22:28.362','employee','oq8WoERWkhEIY4hbjuoL2w==',0,0),
	 ('user88','+Fd8M3bdf/o6ZV2+L0LmSezIPUoUHe1Szg4fq57QAU0=','active','2025-06-15 23:22:29.362','2025-06-15 23:22:29.362','employee','+goYDM69H5Dj+9kwZK8D2w==',0,0),
	 ('user89','/vv1f5r0vGILh6RAUAwyiVuc9isffdR9Rxhp4O6yZmk=','active','2025-06-15 23:22:30.362','2025-06-15 23:22:30.362','employee','fob9xqnMCdL44uUyfLVBng==',0,0),
	 ('user90','9Lkef+Grye5s8yzIse42bCD4GVyzi8D0SKBiDH4vcsI=','active','2025-06-15 23:22:31.362','2025-06-15 23:22:31.362','employee','JhUUh8aaRiBuzYdOxS5R5g==',0,0);
INSERT INTO public.users (username,"password",status,created_at,updated_at,"role",salt,created_by,updated_by) VALUES
	 ('user91','c2SNC3HYnJWiB9/+uiDZgpI/x5U/TbTp4vmhrIZdKvE=','active','2025-06-15 23:22:32.362','2025-06-15 23:22:32.362','employee','C/EWDG3P5MkEL1zdUqZdqA==',0,0),
	 ('user92','kM8541FuOgthymlPCWb10MEGGr+RpBmtYluykMPtCXk=','active','2025-06-15 23:22:33.362','2025-06-15 23:22:33.362','employee','0NTRaUuKebECY+50Aq7+lA==',0,0),
	 ('user93','MHcPaio5mlv5E2dvt44+zcfZLOaxrgY0vKHYYluwpD0=','active','2025-06-15 23:22:34.362','2025-06-15 23:22:34.362','employee','vN99Fs8T5MKa80U4n//2Bw==',0,0),
	 ('user94','LpvROEo4xayZg5Ml7p9DuF6hnydSeLrkodW639spiUM=','active','2025-06-15 23:22:35.362','2025-06-15 23:22:35.362','employee','TOrBkEzGslHzQWP7RWoh9w==',0,0),
	 ('user95','gc8kzyJgkgY1zwWLQzBi3YthMdV/YRuc4XtkfGdHwfE=','active','2025-06-15 23:22:36.362','2025-06-15 23:22:36.362','employee','oRdRdBfP5Su1WNCmU4cwjQ==',0,0),
	 ('user96','+xjIDbkdJejt2dNnwTuwCjExMkDydn6c7PiGpbDcTK8=','active','2025-06-15 23:22:37.362','2025-06-15 23:22:37.362','employee','qQdoInb+PGSc6BbyyHseHQ==',0,0),
	 ('user97','L446tfO9Pw/vitqwlPEgE1UUqqnjvxjRahGH22MuKJk=','active','2025-06-15 23:22:38.362','2025-06-15 23:22:38.362','employee','lUcRcki6tP6+4RiginawwA==',0,0),
	 ('user98','P/1PhI41FXc+J5o+ZcpbgrmC204Cv4U2bjU5fY0wI7k=','active','2025-06-15 23:22:39.362','2025-06-15 23:22:39.362','employee','ZgTulv+zTrVwxZnEEgAbDQ==',0,0),
	 ('user99','q1WTJNflBwwhxveEPvLuN1Bknggx5DI5SXzkZ8DP2Og=','active','2025-06-15 23:22:40.362','2025-06-15 23:22:40.362','employee','lnuc49XRyFP/mg2hThsdyA==',0,0),
	 ('user100','hAhEVTAAQ8gsfk+Np6NrhG3lOSXUaCuP6A7b3TQtDFM=','active','2025-06-15 23:22:41.362','2025-06-15 23:22:41.362','employee','DNq8+uuesKKa86d+l43F0w==',0,0);
INSERT INTO public.users (username,"password",status,created_at,updated_at,"role",salt,created_by,updated_by) VALUES
	 ('user1','XNJCb3D9+6V6YuXmeiwJaD2YuhWA1fxrYCwHEZa6p10=','active','2025-06-15 23:21:02.362','2025-06-15 23:21:02.362','employee','BIqwAg6TvxpX/3bc6KE5Ng==',0,0),
	 ('user2','iwk7EcGRKBOvhteqnjukgnrqzdNbkr8MkaivinbdBiQ=','active','2025-06-15 23:21:03.362','2025-06-15 23:21:03.362','employee','lO2whJ/dU6xRjhju+kYPpw==',0,0),
	 ('user3','cHbyEVYVUSEZR4BpJVfrU05//gdlZIsI9RczL4QHzU4=','active','2025-06-15 23:21:04.362','2025-06-15 23:21:04.362','employee','v0RzUswe4e3rNoM/kaHHWg==',0,0),
	 ('user4','xkvi2+l/y17mB3FVxULKaiG4KL/L01cM23Synv8raIE=','active','2025-06-15 23:21:05.362','2025-06-15 23:21:05.362','employee','ENpjbttXJDkGPjV4ablfmA==',0,0),
	 ('user5','SoFSj5xk79/bRx4D3BMMArShDr7iRFaKGpSJbLxJ+jM=','active','2025-06-15 23:21:06.362','2025-06-15 23:21:06.362','employee','9LR6/tDAgcWZYhnVEZzMhg==',0,0),
	 ('user6','jZias8HGMnuxwoCwPGfkOQDb+hw9o60Mvv4oxjND6QE=','active','2025-06-15 23:21:07.362','2025-06-15 23:21:07.362','employee','59rsoP2czo+uZ8jWeH5iGw==',0,0),
	 ('user7','M8lXgcxgEvn8pahVL5zAhtmbdd04V3NwaPs2WJ1eh4Q=','active','2025-06-15 23:21:08.362','2025-06-15 23:21:08.362','employee','g/kjgWGrBvPOUMFG4kc9Fw==',0,0),
	 ('user8','O2BwhsLCQqxHPsSdWASKqHzRyvZ0pG0d2gH8zURyee0=','active','2025-06-15 23:21:09.362','2025-06-15 23:21:09.362','employee','zR825p6FlvKUTugqDAPf6A==',0,0),
	 ('user9','VgfIZM9pER4vKESwfpvOfqqmkwquHBQrbssNn1s3DJE=','active','2025-06-15 23:21:10.362','2025-06-15 23:21:10.362','employee','CLTVNBiceXBDBcb50xWqJA==',0,0),
	 ('user10','Sg/4UaS2Eap1jdM/tt6mZ+I+ydCYRjyIJ6LBRLD6gME=','active','2025-06-15 23:21:11.362','2025-06-15 23:21:11.362','employee','ohj/oiES1AQsE8lIp3R6Ig==',0,0);
INSERT INTO public.users (username,"password",status,created_at,updated_at,"role",salt,created_by,updated_by) VALUES
	 ('user11','2wW6fNZMShTvJq1uFE12yO7Ssw4VyxMgPOckqkbcScA=','active','2025-06-15 23:21:12.362','2025-06-15 23:21:12.362','employee','QEnGX/yvrTqjawZJhRTXGg==',0,0),
	 ('user12','rPnVLclbCeLbkxVdh8mMM+04h83qpXUkNj3M3Q+A7JY=','active','2025-06-15 23:21:13.362','2025-06-15 23:21:13.362','employee','5+RJ35PLrjjk0dlQh0TNtw==',0,0),
	 ('user13','YyDXfr5cfadfrKUHKcljOaC4BDgX/iy7ZzEfWQ0Q4NQ=','active','2025-06-15 23:21:14.362','2025-06-15 23:21:14.362','employee','UY97yU0i7flhOOPLhXPJkA==',0,0),
	 ('user14','hn/AsCwvxtqlpbLDrMEXZh20FcSgX40+32r3mJfuIrc=','active','2025-06-15 23:21:15.362','2025-06-15 23:21:15.362','employee','5YqhlLh6OwXNLUZjgUqwzQ==',0,0),
	 ('user15','jO9vg6CKtOCTo2NY3hfecMa7PL/3TPy4FmpeZNsdnYE=','active','2025-06-15 23:21:16.362','2025-06-15 23:21:16.362','employee','LvUYykIKkDM8I3Nhk5iKEQ==',0,0),
	 ('user16','LTjmTMdU4Z+W7bnaWCuehXw0HE4QRMVBQus2W79fQ/Y=','active','2025-06-15 23:21:17.362','2025-06-15 23:21:17.362','employee','2A7jHIYGkClCYfc4dlZcqQ==',0,0),
	 ('user17','ANrfa6iglPAriCKb1jWp6p96VIyJccXcj0n/Rdekoy8=','active','2025-06-15 23:21:18.362','2025-06-15 23:21:18.362','employee','1Cdp/ecQvEqv+vnObvx29Q==',0,0),
	 ('user18','8lYZTf2RnqZqnQlONCtmOzq2T7CicaWXe2OcEuV7oJk=','active','2025-06-15 23:21:19.362','2025-06-15 23:21:19.362','employee','1fS4yGVKiVT9j0dCIJlr5g==',0,0),
	 ('user19','XJxPwESU/c9RgETEWprWJsIQnodiFSdIQPsRPhXTgj4=','active','2025-06-15 23:21:20.362','2025-06-15 23:21:20.362','employee','knjWIpTz75O2H/GGRUJCBg==',0,0),
	 ('user20','1X1CQ8qdVKvcN0vbDJVHmxbk+hx3aP1zowosAioXJEI=','active','2025-06-15 23:21:21.362','2025-06-15 23:21:21.362','employee','EoAmeA+pU8Az0sBfOoFJ3Q==',0,0);
INSERT INTO public.users (username,"password",status,created_at,updated_at,"role",salt,created_by,updated_by) VALUES
	 ('user21','NW0rA2AKCe1Hd9Fz80stbCmsAB438n8IEVf3zWPmQFA=','active','2025-06-15 23:21:22.362','2025-06-15 23:21:22.362','employee','5JfSmuBbG1WlM6fkchnJzQ==',0,0),
	 ('user22','iaCaBl7bbcqdiSlaIGjIQhZw9dHGLmlzmI4u1QIlvdc=','active','2025-06-15 23:21:23.362','2025-06-15 23:21:23.362','employee','wYcVe63xIys64dLptblsvg==',0,0),
	 ('user23','75dc7i9uiNoQLsv3G8ZRXVW0NJiTW8Nt2HCGWjFJvNY=','active','2025-06-15 23:21:24.362','2025-06-15 23:21:24.362','employee','eRi4vTRlo4UuLyswYN36lw==',0,0),
	 ('user24','SOBcE8EuiO7RyggnkP0MfQyzfibpn5dnmROkFNpUe0E=','active','2025-06-15 23:21:25.362','2025-06-15 23:21:25.362','employee','r7bavE1HhzelocvYwDcYiQ==',0,0),
	 ('user25','PY5bWEkC4LGlZ/vGoNiuwfbmSSsJ+b7khQJjWIdaQvA=','active','2025-06-15 23:21:26.362','2025-06-15 23:21:26.362','employee','G/IqFBWBrT3I+LIntHibqQ==',0,0),
	 ('user26','jgjReSSyZPnQhcg6IAtP4gSH0tRUsx2fnLF/KB4rPlM=','active','2025-06-15 23:21:27.362','2025-06-15 23:21:27.362','employee','eWX98yPJ3vcSqQODhThO/g==',0,0),
	 ('user27','KUPvv2ESW6HP8DFtu6WTyUu7YtIbTw15xpZLHXTF8O0=','active','2025-06-15 23:21:28.362','2025-06-15 23:21:28.362','employee','AITAe+RbJIAqg4ftNQtyag==',0,0),
	 ('user28','NKPnEM1zPALKimzXya0xdgQ0cCnWx9caKlsiiNkX+dA=','active','2025-06-15 23:21:29.362','2025-06-15 23:21:29.362','employee','p8kSMRaaKFwS19zZBqhYPA==',0,0),
	 ('user29','ISEfYoV+u3l3N2AllNE/9Z+uMGShGFEXsrrH+arv3CE=','active','2025-06-15 23:21:30.362','2025-06-15 23:21:30.362','employee','1zC+bMEU7AqIag2v+3bPOA==',0,0),
	 ('user30','GSCbD9KRv0ltM6k4FATT2VHrrW97sepr6UHHKOv4RWk=','active','2025-06-15 23:21:31.362','2025-06-15 23:21:31.362','employee','3h5/bv0lazl9mhvPwcdNzg==',0,0);
INSERT INTO public.users (username,"password",status,created_at,updated_at,"role",salt,created_by,updated_by) VALUES
	 ('user31','cVHavcgZUbXcLVrGQdKhAVyaKX6gg4L7SRem9qoICzg=','active','2025-06-15 23:21:32.362','2025-06-15 23:21:32.362','employee','TqKdISqxp1V5RZ+ouy4sqA==',0,0),
	 ('user32','ezV4i4aaFWuELleR+wndOClqcfHHmwJOtI7UVEOv3Xs=','active','2025-06-15 23:21:33.362','2025-06-15 23:21:33.362','employee','xakBS/uey1pIKqYjewdG8A==',0,0),
	 ('user33','/x+UIXGZkWjCXlkxAjtxfP4JcKL8TzM0qcBvCuLdlNc=','active','2025-06-15 23:21:34.362','2025-06-15 23:21:34.362','employee','n+8FaM50Jszp5KtUi/mP7Q==',0,0),
	 ('user34','O00N3y8votnaW9/05DrmUcKgIjGmI8omfxDe/QLh54Y=','active','2025-06-15 23:21:35.362','2025-06-15 23:21:35.362','employee','B0yMQy09O703+FnHEg2+zw==',0,0),
	 ('user35','CHAwAdp0Wl0oZw1G3jIVeu+C9ZaP1jPSz2ZtSNz4RIs=','active','2025-06-15 23:21:36.362','2025-06-15 23:21:36.362','employee','ECl3qo2hz/RL97G5Opz51A==',0,0),
	 ('user36','ynVsUE4wY8ge4AGxROXo0EwgOVcmGeO0jBNsD64fhXk=','active','2025-06-15 23:21:37.362','2025-06-15 23:21:37.362','employee','gVVUsjbJMJM/6zGd0RcdOQ==',0,0),
	 ('user37','myaMPGOE1rftJ3I+vlQgxiwvhOewXPPKWCqASZnNGAk=','active','2025-06-15 23:21:38.362','2025-06-15 23:21:38.362','employee','uk6rPOnohc/KexKhcqXVKA==',0,0),
	 ('user38','1drI1Ey1UMReH64byNm/Zo3ND0+G0b8wH7AN1VhMv00=','active','2025-06-15 23:21:39.362','2025-06-15 23:21:39.362','employee','kUPBFhKhx7HdwHV5yXXi6w==',0,0),
	 ('user39','ub5JlO6VXH//L1QBiN/btPCCSYTA36aW61Z5daHPxB0=','active','2025-06-15 23:21:40.362','2025-06-15 23:21:40.362','employee','5mvRORdjaUbROc2bxvNbdQ==',0,0),
	 ('user40','j8TzDoDkhLeHd6XIqQ3wE2Cjk9flQ+GTl8Q+W90RL4w=','active','2025-06-15 23:21:41.362','2025-06-15 23:21:41.362','employee','56JrvzDWYH0ayBKbCYKPWw==',0,0);
INSERT INTO public.users (username,"password",status,created_at,updated_at,"role",salt,created_by,updated_by) VALUES
	 ('user41','IhseF7fW0SwXsVf9nMy00aMqiGzQ9+4Iy7PzDX0tl34=','active','2025-06-15 23:21:42.362','2025-06-15 23:21:42.362','employee','ZqVdJh8QAGiX+eHuonzjqw==',0,0),
	 ('user42','L/nyDLnVKI5m4TzBd5XOCg8KHl/Vs43sVnpPa1Lx2Xw=','active','2025-06-15 23:21:43.362','2025-06-15 23:21:43.362','employee','rSWUi6gjhtdK5CTzBfHChQ==',0,0),
	 ('user43','lMKXS81fcQMyTnYrufhTY3SU8LJ1uAz56HGG18/FLSY=','active','2025-06-15 23:21:44.362','2025-06-15 23:21:44.362','employee','WafYBML7P2qraHu/c8kduw==',0,0),
	 ('user44','DXKflea8OlQLJuSngn9DQdqXHyX7tjYnkqUN6JxVwzU=','active','2025-06-15 23:21:45.362','2025-06-15 23:21:45.362','employee','9vTnnxSxYy1rF+ewkatQaA==',0,0),
	 ('user45','nOmcvgZOjka5W3G9f4dUpw392C9sZ15+V0LFcmn1ZPU=','active','2025-06-15 23:21:46.362','2025-06-15 23:21:46.362','employee','P+fR0EXcTCr0Fhi5tYTjjA==',0,0),
	 ('user46','H1P8RdOxZiaoIjMUf3IDpCJKmx+nbgFHmZrR1ztiSgY=','active','2025-06-15 23:21:47.362','2025-06-15 23:21:47.362','employee','ehTTyk4QPNw0/1hn07pbSQ==',0,0),
	 ('user47','kpTA/+SXVlnJwHUUVV+MufOQpkdz5LLsVDsVsBgTeU0=','active','2025-06-15 23:21:48.362','2025-06-15 23:21:48.362','employee','7mye7gyHhS0hWrnX/yAsyA==',0,0),
	 ('user48','wAMhPqmEbfZrfvDjzjpJDuhfii5dq2/mXbTuL6wXWaU=','active','2025-06-15 23:21:49.362','2025-06-15 23:21:49.362','employee','Owd1d+WZkHXCxg1lzUnD6w==',0,0),
	 ('user49','owX4iSREFcpKhl+to3psoOTy7U5QJTnLEXwaAIpr7jU=','active','2025-06-15 23:21:50.362','2025-06-15 23:21:50.362','employee','taJQ5WV9xxJ83MjkF/CwHQ==',0,0),
	 ('user50','yKt056p147v71lqB/0V4vqwc26311nw50zvisA1G7Qo=','active','2025-06-15 23:21:51.362','2025-06-15 23:21:51.362','employee','tVzdEoJuEO+/MVNca9PIgg==',0,0);
INSERT INTO public.users (username,"password",status,created_at,updated_at,"role",salt,created_by,updated_by) VALUES
	 ('user51','5ZWEodV4KNH6NIUMeoBHbufrJcVS8f+0unL8z2xqyFg=','active','2025-06-15 23:21:52.362','2025-06-15 23:21:52.362','employee','inQI0fTFQ5Rrl1pGCSXmWg==',0,0),
	 ('user52','fGq4FYV0LBE+L3wGpqAGaNguSVjRQav6QIud97WUc+A=','active','2025-06-15 23:21:53.362','2025-06-15 23:21:53.362','employee','Rf0iLTFyiCw+hIO1ttbBSQ==',0,0),
	 ('user53','lO3TfjbJc+nSFuJ0EgA5lbC8ZIfnDLQKZMo9Bxo2LQc=','active','2025-06-15 23:21:54.362','2025-06-15 23:21:54.362','employee','YDlZYua7BV8GOMUFKcKpgg==',0,0),
	 ('user54','YU1ivENPR8DmtssrmhkCaPqqA9u2orkyu4+cpzWF8dU=','active','2025-06-15 23:21:55.362','2025-06-15 23:21:55.362','employee','ea6Ik9nWSaIdBt0m5HPuvQ==',0,0),
	 ('user55','x7Lk8Rt/xl2xGDeYIOpnVfzDmrecOx5Ze1mAqiNnw4c=','active','2025-06-15 23:21:56.362','2025-06-15 23:21:56.362','employee','gvjDRVCoHD7SsffFNh8ycA==',0,0),
	 ('user56','mmqvB8kz89eTxJoTal5cx28i1vFjUOMqJCXwuf0Gusg=','active','2025-06-15 23:21:57.362','2025-06-15 23:21:57.362','employee','dMCLofB3iygHx7xhHC704g==',0,0),
	 ('user57','srtUMSTUOWEICirgc5APdM6xx5AjH0RHCFV1fKNsTcQ=','active','2025-06-15 23:21:58.362','2025-06-15 23:21:58.362','employee','rQNPN+SWFCUs0laphf1FdQ==',0,0),
	 ('user58','rouyga6dVHEM0dHctn2El4NCLBRc75ipQlJAfLAcY3c=','active','2025-06-15 23:21:59.362','2025-06-15 23:21:59.362','employee','qYDmyY4VRTadD8roSfPT4w==',0,0),
	 ('user59','x8ZLcMtid5vhVPMyiBMCUpX1yopzBaP166V1jrTdQCk=','active','2025-06-15 23:22:00.362','2025-06-15 23:22:00.362','employee','P5Kkv/XkYJ/cfP8T9GWIEA==',0,0),
	 ('user60','Ats9MmSv2tTQ05Ze0RZML5OYItpiM3jyltue+ZI/ulY=','active','2025-06-15 23:22:01.362','2025-06-15 23:22:01.362','employee','sjMmCTMrxHo4U0QQRVoZHw==',0,0);
INSERT INTO public.users (username,"password",status,created_at,updated_at,"role",salt,created_by,updated_by) VALUES
	 ('user61','5QC+yzOgYsyiaKYX3kXW5OBicm6ZPcOqpW0/WDd9p0w=','active','2025-06-15 23:22:02.362','2025-06-15 23:22:02.362','employee','YhhYC3FQDVOr/xH099ZQLQ==',0,0),
	 ('user62','ycBIl0ekYUxoInTlM5pob2vUMoYA+kclOLrBW83w/K4=','active','2025-06-15 23:22:03.362','2025-06-15 23:22:03.362','employee','CeOjPjdetC0wxQEJFl9nlg==',0,0),
	 ('user63','VAnLnRxB/mkOA5+JY98YcmOiWlKtDy3k8UlUtSQ/YLA=','active','2025-06-15 23:22:04.362','2025-06-15 23:22:04.362','employee','naQM5xQNrOTZuhtCwcGi1A==',0,0),
	 ('user64','Db6UhMiyvvE5RVkyIgKcQWyM3usKtCkGffg+6KD0gEg=','active','2025-06-15 23:22:05.362','2025-06-15 23:22:05.362','employee','Ty6yuQC8gjZ8zXcle5HmBw==',0,0),
	 ('user65','7D6KpVw0zVjc39ijBrxKqXkrx9fOAICtrdLHoehEYyI=','active','2025-06-15 23:22:06.362','2025-06-15 23:22:06.362','employee','MCen+AD+RAH33YmPzPBwEw==',0,0),
	 ('user66','T0N/tLoreHnsaw+gcPDLY0nm0uCygGgdCPJsmICDOGg=','active','2025-06-15 23:22:07.362','2025-06-15 23:22:07.362','employee','0YE5xlJvnRqB1B5gEcg0kA==',0,0),
	 ('user67','y5lM+JwkteU2dbT60d+y2FcsJSRUQ1kqi0pbnJIlQ50=','active','2025-06-15 23:22:08.362','2025-06-15 23:22:08.362','employee','dv8ElBdDHylMtg2cT1PwaQ==',0,0),
	 ('user68','J39wHhzPdpNHMr4NoSs+hDcak9OzELr6FNHeh3zAiEA=','active','2025-06-15 23:22:09.362','2025-06-15 23:22:09.362','employee','7+L92nPGoYTpsZx07S+CNA==',0,0),
	 ('user69','lEtoVTaQDqb2ekMgkgtFT3Z0douPmONJKJQU2aNwOF0=','active','2025-06-15 23:22:10.362','2025-06-15 23:22:10.362','employee','se81xKnGfvpP6ORv2CqqYg==',0,0),
	 ('user70','CeSnGzL7XM2WEd7c/bhesmkF9IxUESTti5nilqfTDZQ=','active','2025-06-15 23:22:11.362','2025-06-15 23:22:11.362','employee','gITNTbAW68cphDSQK6mUOA==',0,0);
INSERT INTO public.users (username,"password",status,created_at,updated_at,"role",salt,created_by,updated_by) VALUES
	 ('user71','qVRyY20TvD4AtiEjPicCQw0l8k+nWbi5j1r+GLTNmxI=','active','2025-06-15 23:22:12.362','2025-06-15 23:22:12.362','employee','A2YdOguPoTf37+V7GIWTDw==',0,0),
	 ('user72','kJeC/JkVTHJ3XgtNxm1JnHvgTD6vETPCVrZmqe16VZ0=','active','2025-06-15 23:22:13.362','2025-06-15 23:22:13.362','employee','xReIbHk7bxv/TeSFkWEUkQ==',0,0),
	 ('user73','8MvOH8BAlZZYqQNHhwZ31hk0R5TASDr3u1+LGumnpAw=','active','2025-06-15 23:22:14.362','2025-06-15 23:22:14.362','employee','EoiZR1sdRdGY7kNk9eznlA==',0,0),
	 ('user74','vGVH3y+8kpXPpLnQii9ZJ0WAstjIEF1kfKaNqyK48ls=','active','2025-06-15 23:22:15.362','2025-06-15 23:22:15.362','employee','cWKuFNxOR12u7JkiZPFMxw==',0,0),
	 ('user75','OK/VTksvRj2vQX7tbR2mMY5AQ2dYMLM8MMCN3D1iZkE=','active','2025-06-15 23:22:16.362','2025-06-15 23:22:16.362','employee','McPbKRcboYpk5RnxuCKfHg==',0,0),
	 ('user76','W7jZYD7714RyZ+JdpadNkKu7dh/NScqEL+VKpK0egCE=','active','2025-06-15 23:22:17.362','2025-06-15 23:22:17.362','employee','JJzXr0WTXrVLL7UW0h2FZw==',0,0),
	 ('user77','pNGSOYsLLC3opNTNLsYjGpPl/AGCoor3L8N1yp51gxI=','active','2025-06-15 23:22:18.362','2025-06-15 23:22:18.362','employee','TZ4bUIdzlnaifX36ILFvNg==',0,0),
	 ('user78','dHngLoazbTbt8Jddif3iwgtHeayiji5AkepYdrgkrOo=','active','2025-06-15 23:22:19.362','2025-06-15 23:22:19.362','employee','kxOvZxSdPwDLWedMu1Ebsw==',0,0),
	 ('user79','WZumcPPrs+Ew4PAe0+FvdjxBJtKeonN1Ojv5J/Jkg08=','active','2025-06-15 23:22:20.362','2025-06-15 23:22:20.362','employee','SPTxkAO+LdmjAYNyfFMK3A==',0,0),
	 ('user80','xdv88ccjkW+vwNhUduBWVfyM+PxU/TviFpzsc+MLfkg=','active','2025-06-15 23:22:21.362','2025-06-15 23:22:21.362','employee','rTJyLih4VDZH2uc+lj+dhg==',0,0);
INSERT INTO public.users (username,"password",status,created_at,updated_at,"role",salt,created_by,updated_by) VALUES
	 ('user81','rP+aOCm1azvGwCOcWlTW9mFcK5GSlYUBN7ynss2qKNI=','active','2025-06-15 23:22:22.362','2025-06-15 23:22:22.362','employee','hyXz4QSc85j1mS3aHc3eAA==',0,0),
	 ('user82','IzbtFl//zDIocxTFpi1pTK9yg6qPa3GB8rFlk7MZy80=','active','2025-06-15 23:22:23.362','2025-06-15 23:22:23.362','employee','nb4/R4Npln+JLOlvwlhW4A==',0,0),
	 ('user83','UGXBaRsBtLHeF+E9q3w/54a99+OSUS6jALUlWvcBYi0=','active','2025-06-15 23:22:24.362','2025-06-15 23:22:24.362','employee','ck4Of/FliCTjcjCyBF9KQA==',0,0),
	 ('user84','dAxpnk6eAI4Ug25b2QkQ3QeV+9hxdTsHUWTTscAu/PU=','active','2025-06-15 23:22:25.362','2025-06-15 23:22:25.362','employee','XNNQ58AUaH4M4l0ipZSpFQ==',0,0),
	 ('user85','A70MGyOOKUVdWGUmj/bt4rrH0R93uWSVZ8h9oXQlqek=','active','2025-06-15 23:22:26.362','2025-06-15 23:22:26.362','employee','GuLWqJlQUoXln7Pd10XdIw==',0,0),
	 ('user86','5GmKxsCNhqudA5qFlKFeHSQ0IdKD1B4julo37K2taRE=','active','2025-06-15 23:22:27.362','2025-06-15 23:22:27.362','employee','FGUPlcLZKQLJ5nLSBImHNA==',0,0),
	 ('user87','z5UYliktiGiD7Fa2BlTAzO3W/LI32jcJGZqyUQcUheI=','active','2025-06-15 23:22:28.362','2025-06-15 23:22:28.362','employee','oq8WoERWkhEIY4hbjuoL2w==',0,0),
	 ('user88','+Fd8M3bdf/o6ZV2+L0LmSezIPUoUHe1Szg4fq57QAU0=','active','2025-06-15 23:22:29.362','2025-06-15 23:22:29.362','employee','+goYDM69H5Dj+9kwZK8D2w==',0,0),
	 ('user89','/vv1f5r0vGILh6RAUAwyiVuc9isffdR9Rxhp4O6yZmk=','active','2025-06-15 23:22:30.362','2025-06-15 23:22:30.362','employee','fob9xqnMCdL44uUyfLVBng==',0,0),
	 ('user90','9Lkef+Grye5s8yzIse42bCD4GVyzi8D0SKBiDH4vcsI=','active','2025-06-15 23:22:31.362','2025-06-15 23:22:31.362','employee','JhUUh8aaRiBuzYdOxS5R5g==',0,0);
INSERT INTO public.users (username,"password",status,created_at,updated_at,"role",salt,created_by,updated_by) VALUES
	 ('user91','c2SNC3HYnJWiB9/+uiDZgpI/x5U/TbTp4vmhrIZdKvE=','active','2025-06-15 23:22:32.362','2025-06-15 23:22:32.362','employee','C/EWDG3P5MkEL1zdUqZdqA==',0,0),
	 ('user92','kM8541FuOgthymlPCWb10MEGGr+RpBmtYluykMPtCXk=','active','2025-06-15 23:22:33.362','2025-06-15 23:22:33.362','employee','0NTRaUuKebECY+50Aq7+lA==',0,0),
	 ('user93','MHcPaio5mlv5E2dvt44+zcfZLOaxrgY0vKHYYluwpD0=','active','2025-06-15 23:22:34.362','2025-06-15 23:22:34.362','employee','vN99Fs8T5MKa80U4n//2Bw==',0,0),
	 ('user94','LpvROEo4xayZg5Ml7p9DuF6hnydSeLrkodW639spiUM=','active','2025-06-15 23:22:35.362','2025-06-15 23:22:35.362','employee','TOrBkEzGslHzQWP7RWoh9w==',0,0),
	 ('user95','gc8kzyJgkgY1zwWLQzBi3YthMdV/YRuc4XtkfGdHwfE=','active','2025-06-15 23:22:36.362','2025-06-15 23:22:36.362','employee','oRdRdBfP5Su1WNCmU4cwjQ==',0,0),
	 ('user96','+xjIDbkdJejt2dNnwTuwCjExMkDydn6c7PiGpbDcTK8=','active','2025-06-15 23:22:37.362','2025-06-15 23:22:37.362','employee','qQdoInb+PGSc6BbyyHseHQ==',0,0),
	 ('user97','L446tfO9Pw/vitqwlPEgE1UUqqnjvxjRahGH22MuKJk=','active','2025-06-15 23:22:38.362','2025-06-15 23:22:38.362','employee','lUcRcki6tP6+4RiginawwA==',0,0),
	 ('user98','P/1PhI41FXc+J5o+ZcpbgrmC204Cv4U2bjU5fY0wI7k=','active','2025-06-15 23:22:39.362','2025-06-15 23:22:39.362','employee','ZgTulv+zTrVwxZnEEgAbDQ==',0,0),
	 ('user99','q1WTJNflBwwhxveEPvLuN1Bknggx5DI5SXzkZ8DP2Og=','active','2025-06-15 23:22:40.362','2025-06-15 23:22:40.362','employee','lnuc49XRyFP/mg2hThsdyA==',0,0),
	 ('user100','hAhEVTAAQ8gsfk+Np6NrhG3lOSXUaCuP6A7b3TQtDFM=','active','2025-06-15 23:22:41.362','2025-06-15 23:22:41.362','employee','DNq8+uuesKKa86d+l43F0w==',0,0);
INSERT INTO public.users (username,"password",status,created_at,updated_at,"role",salt,created_by,updated_by) VALUES
	 ('admin','Pep/e5hZObYronYQeNj3c2BWAiU7RAZ0uG0m7SFEjDM=','active','2025-06-15 15:21:01.362','2025-06-15 15:21:01.362','admin','rKxyNitzMfnfqmpWOsbezQ==',0,0);


INSERT INTO public.employees (fullname,salary,code,created_at,updated_at,created_by,updated_by,user_id) VALUES
	 ('Employee 101',90000000,'EMP0101','2025-06-15 23:22:41.362','2025-06-15 23:22:41.362',0,0,1),
	 ('Employee 1',40000000,'EMP0001','2025-06-15 23:21:02.362','2025-06-15 23:21:02.362',0,0,2),
	 ('Employee 2',50000000,'EMP0002','2025-06-15 23:21:03.362','2025-06-15 23:21:03.362',0,0,3),
	 ('Employee 3',40000000,'EMP0003','2025-06-15 23:21:04.362','2025-06-15 23:21:04.362',0,0,4),
	 ('Employee 4',90000000,'EMP0004','2025-06-15 23:21:05.362','2025-06-15 23:21:05.362',0,0,5),
	 ('Employee 5',60000000,'EMP0005','2025-06-15 23:21:06.362','2025-06-15 23:21:06.362',0,0,6),
	 ('Employee 6',70000000,'EMP0006','2025-06-15 23:21:07.362','2025-06-15 23:21:07.362',0,0,7),
	 ('Employee 7',80000000,'EMP0007','2025-06-15 23:21:08.362','2025-06-15 23:21:08.362',0,0,8),
	 ('Employee 8',70000000,'EMP0008','2025-06-15 23:21:09.362','2025-06-15 23:21:09.362',0,0,9),
	 ('Employee 9',60000000,'EMP0009','2025-06-15 23:21:10.362','2025-06-15 23:21:10.362',0,0,10);
INSERT INTO public.employees (fullname,salary,code,created_at,updated_at,created_by,updated_by,user_id) VALUES
	 ('Employee 10',30000000,'EMP0010','2025-06-15 23:21:11.362','2025-06-15 23:21:11.362',0,0,11),
	 ('Employee 11',50000000,'EMP0011','2025-06-15 23:21:12.362','2025-06-15 23:21:12.362',0,0,12),
	 ('Employee 12',70000000,'EMP0012','2025-06-15 23:21:13.362','2025-06-15 23:21:13.362',0,0,13),
	 ('Employee 13',30000000,'EMP0013','2025-06-15 23:21:14.362','2025-06-15 23:21:14.362',0,0,14),
	 ('Employee 14',70000000,'EMP0014','2025-06-15 23:21:15.362','2025-06-15 23:21:15.362',0,0,15),
	 ('Employee 15',80000000,'EMP0015','2025-06-15 23:21:16.362','2025-06-15 23:21:16.362',0,0,16),
	 ('Employee 16',80000000,'EMP0016','2025-06-15 23:21:17.362','2025-06-15 23:21:17.362',0,0,17),
	 ('Employee 17',80000000,'EMP0017','2025-06-15 23:21:18.362','2025-06-15 23:21:18.362',0,0,18),
	 ('Employee 18',30000000,'EMP0018','2025-06-15 23:21:19.362','2025-06-15 23:21:19.362',0,0,19),
	 ('Employee 19',30000000,'EMP0019','2025-06-15 23:21:20.362','2025-06-15 23:21:20.362',0,0,20);
INSERT INTO public.employees (fullname,salary,code,created_at,updated_at,created_by,updated_by,user_id) VALUES
	 ('Employee 20',30000000,'EMP0020','2025-06-15 23:21:21.362','2025-06-15 23:21:21.362',0,0,21),
	 ('Employee 21',70000000,'EMP0021','2025-06-15 23:21:22.362','2025-06-15 23:21:22.362',0,0,22),
	 ('Employee 22',80000000,'EMP0022','2025-06-15 23:21:23.362','2025-06-15 23:21:23.362',0,0,23),
	 ('Employee 23',60000000,'EMP0023','2025-06-15 23:21:24.362','2025-06-15 23:21:24.362',0,0,24),
	 ('Employee 24',90000000,'EMP0024','2025-06-15 23:21:25.362','2025-06-15 23:21:25.362',0,0,25),
	 ('Employee 25',80000000,'EMP0025','2025-06-15 23:21:26.362','2025-06-15 23:21:26.362',0,0,26),
	 ('Employee 26',60000000,'EMP0026','2025-06-15 23:21:27.362','2025-06-15 23:21:27.362',0,0,27),
	 ('Employee 27',50000000,'EMP0027','2025-06-15 23:21:28.362','2025-06-15 23:21:28.362',0,0,28),
	 ('Employee 28',40000000,'EMP0028','2025-06-15 23:21:29.362','2025-06-15 23:21:29.362',0,0,29),
	 ('Employee 29',60000000,'EMP0029','2025-06-15 23:21:30.362','2025-06-15 23:21:30.362',0,0,30);
INSERT INTO public.employees (fullname,salary,code,created_at,updated_at,created_by,updated_by,user_id) VALUES
	 ('Employee 30',40000000,'EMP0030','2025-06-15 23:21:31.362','2025-06-15 23:21:31.362',0,0,31),
	 ('Employee 31',50000000,'EMP0031','2025-06-15 23:21:32.362','2025-06-15 23:21:32.362',0,0,32),
	 ('Employee 32',70000000,'EMP0032','2025-06-15 23:21:33.362','2025-06-15 23:21:33.362',0,0,33),
	 ('Employee 33',30000000,'EMP0033','2025-06-15 23:21:34.362','2025-06-15 23:21:34.362',0,0,34),
	 ('Employee 34',40000000,'EMP0034','2025-06-15 23:21:35.362','2025-06-15 23:21:35.362',0,0,35),
	 ('Employee 35',30000000,'EMP0035','2025-06-15 23:21:36.362','2025-06-15 23:21:36.362',0,0,36),
	 ('Employee 36',50000000,'EMP0036','2025-06-15 23:21:37.362','2025-06-15 23:21:37.362',0,0,37),
	 ('Employee 37',30000000,'EMP0037','2025-06-15 23:21:38.362','2025-06-15 23:21:38.362',0,0,38),
	 ('Employee 38',40000000,'EMP0038','2025-06-15 23:21:39.362','2025-06-15 23:21:39.362',0,0,39),
	 ('Employee 39',80000000,'EMP0039','2025-06-15 23:21:40.362','2025-06-15 23:21:40.362',0,0,40);
INSERT INTO public.employees (fullname,salary,code,created_at,updated_at,created_by,updated_by,user_id) VALUES
	 ('Employee 40',90000000,'EMP0040','2025-06-15 23:21:41.362','2025-06-15 23:21:41.362',0,0,41),
	 ('Employee 41',70000000,'EMP0041','2025-06-15 23:21:42.362','2025-06-15 23:21:42.362',0,0,42),
	 ('Employee 42',80000000,'EMP0042','2025-06-15 23:21:43.362','2025-06-15 23:21:43.362',0,0,43),
	 ('Employee 43',70000000,'EMP0043','2025-06-15 23:21:44.362','2025-06-15 23:21:44.362',0,0,44),
	 ('Employee 44',30000000,'EMP0044','2025-06-15 23:21:45.362','2025-06-15 23:21:45.362',0,0,45),
	 ('Employee 45',40000000,'EMP0045','2025-06-15 23:21:46.362','2025-06-15 23:21:46.362',0,0,46),
	 ('Employee 46',80000000,'EMP0046','2025-06-15 23:21:47.362','2025-06-15 23:21:47.362',0,0,47),
	 ('Employee 47',80000000,'EMP0047','2025-06-15 23:21:48.362','2025-06-15 23:21:48.362',0,0,48),
	 ('Employee 48',50000000,'EMP0048','2025-06-15 23:21:49.362','2025-06-15 23:21:49.362',0,0,49),
	 ('Employee 49',70000000,'EMP0049','2025-06-15 23:21:50.362','2025-06-15 23:21:50.362',0,0,50);
INSERT INTO public.employees (fullname,salary,code,created_at,updated_at,created_by,updated_by,user_id) VALUES
	 ('Employee 50',70000000,'EMP0050','2025-06-15 23:21:51.362','2025-06-15 23:21:51.362',0,0,51),
	 ('Employee 51',60000000,'EMP0051','2025-06-15 23:21:52.362','2025-06-15 23:21:52.362',0,0,52),
	 ('Employee 52',30000000,'EMP0052','2025-06-15 23:21:53.362','2025-06-15 23:21:53.362',0,0,53),
	 ('Employee 53',60000000,'EMP0053','2025-06-15 23:21:54.362','2025-06-15 23:21:54.362',0,0,54),
	 ('Employee 54',30000000,'EMP0054','2025-06-15 23:21:55.362','2025-06-15 23:21:55.362',0,0,55),
	 ('Employee 55',30000000,'EMP0055','2025-06-15 23:21:56.362','2025-06-15 23:21:56.362',0,0,56),
	 ('Employee 56',60000000,'EMP0056','2025-06-15 23:21:57.362','2025-06-15 23:21:57.362',0,0,57),
	 ('Employee 57',70000000,'EMP0057','2025-06-15 23:21:58.362','2025-06-15 23:21:58.362',0,0,58),
	 ('Employee 58',40000000,'EMP0058','2025-06-15 23:21:59.362','2025-06-15 23:21:59.362',0,0,59),
	 ('Employee 59',30000000,'EMP0059','2025-06-15 23:22:00.362','2025-06-15 23:22:00.362',0,0,60);
INSERT INTO public.employees (fullname,salary,code,created_at,updated_at,created_by,updated_by,user_id) VALUES
	 ('Employee 60',70000000,'EMP0060','2025-06-15 23:22:01.362','2025-06-15 23:22:01.362',0,0,61),
	 ('Employee 61',60000000,'EMP0061','2025-06-15 23:22:02.362','2025-06-15 23:22:02.362',0,0,62),
	 ('Employee 62',50000000,'EMP0062','2025-06-15 23:22:03.362','2025-06-15 23:22:03.362',0,0,63),
	 ('Employee 63',70000000,'EMP0063','2025-06-15 23:22:04.362','2025-06-15 23:22:04.362',0,0,64),
	 ('Employee 64',90000000,'EMP0064','2025-06-15 23:22:05.362','2025-06-15 23:22:05.362',0,0,65),
	 ('Employee 65',80000000,'EMP0065','2025-06-15 23:22:06.362','2025-06-15 23:22:06.362',0,0,66),
	 ('Employee 66',40000000,'EMP0066','2025-06-15 23:22:07.362','2025-06-15 23:22:07.362',0,0,67),
	 ('Employee 67',60000000,'EMP0067','2025-06-15 23:22:08.362','2025-06-15 23:22:08.362',0,0,68),
	 ('Employee 68',60000000,'EMP0068','2025-06-15 23:22:09.362','2025-06-15 23:22:09.362',0,0,69),
	 ('Employee 69',70000000,'EMP0069','2025-06-15 23:22:10.362','2025-06-15 23:22:10.362',0,0,70);
INSERT INTO public.employees (fullname,salary,code,created_at,updated_at,created_by,updated_by,user_id) VALUES
	 ('Employee 70',80000000,'EMP0070','2025-06-15 23:22:11.362','2025-06-15 23:22:11.362',0,0,71),
	 ('Employee 71',60000000,'EMP0071','2025-06-15 23:22:12.362','2025-06-15 23:22:12.362',0,0,72),
	 ('Employee 72',90000000,'EMP0072','2025-06-15 23:22:13.362','2025-06-15 23:22:13.362',0,0,73),
	 ('Employee 73',80000000,'EMP0073','2025-06-15 23:22:14.362','2025-06-15 23:22:14.362',0,0,74),
	 ('Employee 74',80000000,'EMP0074','2025-06-15 23:22:15.362','2025-06-15 23:22:15.362',0,0,75),
	 ('Employee 75',70000000,'EMP0075','2025-06-15 23:22:16.362','2025-06-15 23:22:16.362',0,0,76),
	 ('Employee 76',90000000,'EMP0076','2025-06-15 23:22:17.362','2025-06-15 23:22:17.362',0,0,77),
	 ('Employee 77',80000000,'EMP0077','2025-06-15 23:22:18.362','2025-06-15 23:22:18.362',0,0,78),
	 ('Employee 78',50000000,'EMP0078','2025-06-15 23:22:19.362','2025-06-15 23:22:19.362',0,0,79),
	 ('Employee 79',30000000,'EMP0079','2025-06-15 23:22:20.362','2025-06-15 23:22:20.362',0,0,80);
INSERT INTO public.employees (fullname,salary,code,created_at,updated_at,created_by,updated_by,user_id) VALUES
	 ('Employee 80',50000000,'EMP0080','2025-06-15 23:22:21.362','2025-06-15 23:22:21.362',0,0,81),
	 ('Employee 81',80000000,'EMP0081','2025-06-15 23:22:22.362','2025-06-15 23:22:22.362',0,0,82),
	 ('Employee 82',30000000,'EMP0082','2025-06-15 23:22:23.362','2025-06-15 23:22:23.362',0,0,83),
	 ('Employee 83',50000000,'EMP0083','2025-06-15 23:22:24.362','2025-06-15 23:22:24.362',0,0,84),
	 ('Employee 84',50000000,'EMP0084','2025-06-15 23:22:25.362','2025-06-15 23:22:25.362',0,0,85),
	 ('Employee 85',50000000,'EMP0085','2025-06-15 23:22:26.362','2025-06-15 23:22:26.362',0,0,86),
	 ('Employee 86',30000000,'EMP0086','2025-06-15 23:22:27.362','2025-06-15 23:22:27.362',0,0,87),
	 ('Employee 87',60000000,'EMP0087','2025-06-15 23:22:28.362','2025-06-15 23:22:28.362',0,0,88),
	 ('Employee 88',70000000,'EMP0088','2025-06-15 23:22:29.362','2025-06-15 23:22:29.362',0,0,89),
	 ('Employee 89',70000000,'EMP0089','2025-06-15 23:22:30.362','2025-06-15 23:22:30.362',0,0,90);
INSERT INTO public.employees (fullname,salary,code,created_at,updated_at,created_by,updated_by,user_id) VALUES
	 ('Employee 90',80000000,'EMP0090','2025-06-15 23:22:31.362','2025-06-15 23:22:31.362',0,0,91),
	 ('Employee 91',40000000,'EMP0091','2025-06-15 23:22:32.362','2025-06-15 23:22:32.362',0,0,92),
	 ('Employee 92',50000000,'EMP0092','2025-06-15 23:22:33.362','2025-06-15 23:22:33.362',0,0,93),
	 ('Employee 93',50000000,'EMP0093','2025-06-15 23:22:34.362','2025-06-15 23:22:34.362',0,0,94),
	 ('Employee 94',30000000,'EMP0094','2025-06-15 23:22:35.362','2025-06-15 23:22:35.362',0,0,95),
	 ('Employee 95',70000000,'EMP0095','2025-06-15 23:22:36.362','2025-06-15 23:22:36.362',0,0,96),
	 ('Employee 96',40000000,'EMP0096','2025-06-15 23:22:37.362','2025-06-15 23:22:37.362',0,0,97),
	 ('Employee 97',70000000,'EMP0097','2025-06-15 23:22:38.362','2025-06-15 23:22:38.362',0,0,98),
	 ('Employee 98',40000000,'EMP0098','2025-06-15 23:22:39.362','2025-06-15 23:22:39.362',0,0,99),
	 ('Employee 99',40000000,'EMP0099','2025-06-15 23:22:40.362','2025-06-15 23:22:40.362',0,0,100);
INSERT INTO public.employees (fullname,salary,code,created_at,updated_at,created_by,updated_by,user_id) VALUES
	 ('Employee 100',90000000,'EMP0100','2025-06-15 23:22:41.362','2025-06-15 23:22:41.362',0,0,101);
