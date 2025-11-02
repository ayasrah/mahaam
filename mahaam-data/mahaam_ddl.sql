DROP TABLE IF EXISTS app.tasks;
DROP TABLE IF EXISTS app.plan_members;
DROP TABLE IF EXISTS app.suggested_emails;
DROP TABLE IF EXISTS app.devices;
DROP TABLE IF EXISTS app.plans;
DROP TABLE IF EXISTS app.users;
--
DROP TABLE IF EXISTS monitor.health;
DROP TABLE IF EXISTS monitor.log;
DROP TABLE IF EXISTS monitor.traffic;
--
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
--

CREATE TABLE app.users (
	id uuid NOT NULL DEFAULT uuid_generate_v4 (),
	email varchar(255) NULL,
	name varchar(50) NULL,
	created_at timestamptz NOT NULL,
	updated_at timestamptz NULL,
	CONSTRAINT users_pkey PRIMARY KEY (id)
);
CREATE UNIQUE INDEX users_unique_index_email ON app.users (email);
--

CREATE TABLE app.devices (
	id uuid NOT NULL,
	user_id uuid NOT NULL,
	platform TEXT NULL,
	fingerprint varchar(255) NOT NULL,
	info varchar(255) NOT NULL,
	created_at timestamptz NOT NULL,
	updated_at timestamptz NULL,
	CONSTRAINT devices_pkey PRIMARY KEY (id),
	CONSTRAINT devices_user_id_fkey FOREIGN KEY (user_id) REFERENCES app.users (id) ON DELETE CASCADE
);
--

CREATE TABLE app.suggested_emails (
	id uuid NOT NULL,
	user_id uuid NOT NULL,
	email varchar(255) NULL,
	created_at timestamptz NOT NULL,
	CONSTRAINT suggested_emails_pkey PRIMARY KEY (id),
	CONSTRAINT suggested_emails_user_id_fkey FOREIGN KEY (user_id) REFERENCES app.users (id) ON DELETE CASCADE
);
CREATE UNIQUE INDEX suggested_emails_unique_index_user_id_email ON app.suggested_emails (user_id, email);
--

CREATE TABLE app.plans (
	id uuid NOT NULL DEFAULT uuid_generate_v4 (),
	user_id uuid NOT NULL,
	type VARCHAR(50) NOT NULL,
	status varchar(50) NOT NULL,
	title varchar(100) NULL,
	starts date NULL,
	ends date NULL,
	done_percent varchar(10) NOT NULL,
	sort_order int4 NOT NULL,
	created_at timestamptz NOT NULL,
	updated_at timestamptz NULL,
	CONSTRAINT plans_pk PRIMARY KEY (id),
	CONSTRAINT plans_user_id_fkey FOREIGN KEY (user_id) REFERENCES app.users (id) ON DELETE CASCADE
);
CREATE INDEX plans_index_type ON app.plans (type);
--

CREATE TABLE app.plan_members (
	plan_id uuid NOT NULL,
	user_id uuid NOT NULL,
	created_at timestamptz NOT NULL,
	CONSTRAINT plan_members_pkey PRIMARY KEY (plan_id, user_id),
	CONSTRAINT plan_members_plan_id_fkey FOREIGN KEY (plan_id) REFERENCES app.plans (id) ON DELETE CASCADE,
	CONSTRAINT plan_members_user_id_fkey FOREIGN KEY (user_id) REFERENCES app.users (id) ON DELETE CASCADE
);
--

CREATE TABLE app.tasks (
	id uuid NOT NULL DEFAULT uuid_generate_v4 (),
	plan_id uuid NOT NULL,
	title varchar(255) NOT NULL,
	done bool NOT NULL,
	sort_order int4 NOT NULL,
	created_at timestamptz NOT NULL,
	updated_at timestamptz NULL,
	CONSTRAINT tasks_pkey PRIMARY KEY (id),
	CONSTRAINT tasks_fkey FOREIGN KEY (plan_id) REFERENCES app.plans (id) ON DELETE CASCADE
);
--

CREATE TABLE monitor.logs (
	id uuid NOT NULL DEFAULT uuid_generate_v4 (),
	traffic_id uuid NULL,
	"type" varchar(50) NULL,
	message text NOT NULL,
	node_ip varchar(20) NOT NULL,
	created_at timestamptz NOT NULL,
	CONSTRAINT logs_pkey PRIMARY KEY (id)
);
--

CREATE TABLE monitor.health (
	id uuid NOT NULL DEFAULT uuid_generate_v4 (),
	api_name varchar(50) NOT NULL,
	api_version varchar(20) NOT NULL,
	node_ip varchar(20) NOT NULL,
	node_name varchar(50) NOT NULL,
	env_name varchar(10) NOT NULL,
	started_at timestamptz NOT NULL,
	pulsed_at timestamptz NULL,
	stopped_at timestamptz NULL,
	CONSTRAINT health_pkey PRIMARY KEY (id)
);
--

CREATE TABLE monitor.traffic (
	id uuid NOT NULL DEFAULT uuid_generate_v4 (),
	health_id uuid NOT NULL,
	method varchar(20) NOT NULL,
	path varchar(255) NOT NULL,
	code int2 NULL,
	elapsed int8 NULL,
	headers text NULL,
	request text NULL,
	response text NULL,
	created_at timestamptz NOT NULL,
	CONSTRAINT traffic_pkey PRIMARY KEY (id)
);