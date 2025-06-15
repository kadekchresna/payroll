

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
	created_at timestamptz DEFAULT now() NOT NULL,
	updated_at timestamptz DEFAULT now() NOT NULL,
	created_by int4 DEFAULT 0 NOT NULL,
	updated_by int4 DEFAULT 0 NOT NULL,
	CONSTRAINT attendances_pk PRIMARY KEY (id),
	CONSTRAINT attendances_unique UNIQUE (employee_id, date)
);
CREATE INDEX IF NOT EXISTS attendances_date_idx ON public.attendances USING btree (date);


CREATE TABLE IF NOT EXISTS public.overtimes (
	id int4 DEFAULT nextval('overtime_id_seq'::regclass) NOT NULL,
	employee_id int4 DEFAULT 0 NOT NULL,
	"date" date NOT NULL,
	hours int4 DEFAULT 0 NOT NULL,
	created_at timestamptz DEFAULT now() NOT NULL,
	updated_at timestamptz DEFAULT now() NOT NULL,
	created_by int4 DEFAULT 0 NOT NULL,
	updated_by int4 DEFAULT 0 NOT NULL,
	CONSTRAINT overtime_pk PRIMARY KEY (id)
);
CREATE INDEX IF NOT EXISTS overtimes_date_idx ON public.overtimes USING btree (date);


CREATE TABLE IF NOT EXISTS public.reimbursements (
	id serial4 NOT NULL,
	employee_id int4 DEFAULT 0 NOT NULL,
	"date" date NOT NULL,
	amount float8 DEFAULT 0 NOT NULL,
	"description" text DEFAULT ''::text NOT NULL,
	payslip_id int4 DEFAULT 0 NOT NULL,
	created_at timestamptz DEFAULT now() NOT NULL,
	updated_at timestamptz DEFAULT now() NOT NULL,
	created_by int4 DEFAULT 0 NOT NULL,
	updated_by int4 DEFAULT 0 NOT NULL,
	CONSTRAINT reimbursements_pk PRIMARY KEY (id)
);
CREATE INDEX IF NOT EXISTS reimbursements_date_idx ON public.reimbursements USING btree (date);



CREATE TABLE public.payslips (
	id serial4 NOT NULL,
	employee_id int4 DEFAULT 0 NOT NULL,
	period_start date NOT NULL,
	period_end date NOT NULL,
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
	CONSTRAINT payslips_unique UNIQUE (employee_id, period_start, period_end)
);
CREATE INDEX payslips_employee_id_idx ON public.payslips USING btree (employee_id);


CREATE TABLE IF NOT EXISTS audit_logs (
    id SERIAL PRIMARY KEY,
    table_name VARCHAR NOT NULL,
    "action" VARCHAR NOT NULL,
    record_id INT NOT NULL,
    old_data JSONB,
    new_data JSONB,
    changed_by INT NOT NULL, 
    changed_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    ip_address VARCHAR NOT NULL DEFAULT '',
    request_id VARCHAR NOT NULL DEFAULT ''
);
