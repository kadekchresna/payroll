

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