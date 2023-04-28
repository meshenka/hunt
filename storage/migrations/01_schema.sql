CREATE TABLE account (
    id uuid NOT NULL PRIMARY KEY,
    handle text NOT NULL, -- displayed username
    description text DEFAULT '' NOT NULL, -- a nice introduction text
    created_at timestamptz DEFAULT now() NOT NULL,
    updated_at timestamptz DEFAULT now() NOT NULL,
    enabled boolean DEFAULT false NOT NULL, -- is account enabled, aka account has been 2FAed
    deleted boolean DEFAULT false NOT NULL,
    deleted_at timestamptz,
    roles jsonb DEFAULT '[]' NOT NULL
);

CREATE INDEX account_handle ON account (handle);
CREATE INDEX account_roles ON account USING GIN (roles);
CREATE INDEX account_enabled ON account (enabled);
CREATE INDEX account_deleted ON account (deleted);

COMMENT ON TABLE "account" IS 'users of the service';
COMMENT ON COLUMN "account"."id" IS 'uniq identifier of a user';
COMMENT ON COLUMN "account"."handle" IS 'display name of the user, user is free to update it from app';
COMMENT ON COLUMN "account"."description" IS 'a nice introduction text';
COMMENT ON COLUMN "account"."enabled" IS 'true if account is able to use features of the service';
COMMENT ON COLUMN "account"."deleted" IS 'soft delete flag';
COMMENT ON COLUMN "account"."roles" IS 'a list of users permissions';

CREATE TABLE company (
    id uuid NOT NULL PRIMARY KEY,
    name text NOT NULL, 
    description text DEFAULT '' NOT NULL, 
    address text DEFAULT '' NOT NULL,
    created_at timestamptz DEFAULT now() NOT NULL,
    updated_at timestamptz DEFAULT now() NOT NULL,
    author_id uuid NOT NULL,
    size_low int DEFAULT 1 NOT NULL,
    size_high int DEFAULT 100 NOT NULL,
    category text DEFAULT '' NOT NULL,
    link text DEFAULT '' NOT NULL,
    CONSTRAINT fk_company_author FOREIGN KEY (author_id) REFERENCES account(id)
);
CREATE INDEX company_name ON company (name);
CREATE INDEX company_category ON company (category);

CREATE TABLE opportunity (
    id uuid NOT NULL PRIMARY KEY,
    title text NOT NULL, 
    description text DEFAULT '' NOT NULL, 
    created_at timestamptz DEFAULT now() NOT NULL,
    updated_at timestamptz DEFAULT now() NOT NULL,
    author_id uuid NOT NULL,
    category text DEFAULT '' NOT NULL, -- backend, frontend, fullstack
    link text DEFAULT '' NOT NULL,
    salary_low int DEFAULT 22000 NOT NULL,
    salary_high int DEFAULT 30000 NOT NULL,
    location text DEFAULT '' NOT NULL,
    location_long float,
    location_lat float,
    kind text DEFAULT 'cdi' NOT NULL, -- cdi cdd part time
    work_format text DEFAULT '' NOT NULL, -- remote, hybrid, on site 
    contact text DEFAULT '' NOT NULL, -- who do i talk to
    CONSTRAINT fk_opportunity_author FOREIGN KEY (author_id) REFERENCES account(id)
);
CREATE INDEX opportunity_title ON opportunity (title);
CREATE INDEX opportunity_category ON opportunity (category);
CREATE INDEX opportunity_salary ON opportunity USING BRIN (salary_low, salary_high);
CREATE INDEX opportunity_location ON opportunity (location);
CREATE INDEX opportunity_workformat ON opportunity (work_format);

CREATE TABLE opportunity_note (
    id uuid NOT NULL PRIMARY KEY,
    description text DEFAULT '' NOT NULL, 
    created_at timestamptz DEFAULT now() NOT NULL,
    updated_at timestamptz DEFAULT now() NOT NULL,
    author_id uuid NOT NULL,
    opportunity_id uuid NOT NULL,
    tags text DEFAULT '' NOT NULL, -- backend, frontend, fullstack
    CONSTRAINT fk_opportunity_note_author FOREIGN KEY (author_id) REFERENCES account(id),
    CONSTRAINT fk_opportunity_note_opportuniry FOREIGN KEY (opportunity_id) REFERENCES opportunity(id)
);

CREATE INDEX opportunity_note_created ON opportunity_note USING BRIN(created_at);
CREATE INDEX opportunity_note_tags ON opportunity_note (tags);

CREATE TABLE company_note (
    id uuid NOT NULL PRIMARY KEY,
    description text DEFAULT '' NOT NULL, 
    created_at timestamptz DEFAULT now() NOT NULL,
    updated_at timestamptz DEFAULT now() NOT NULL,
    author_id uuid NOT NULL,
    company_id uuid NOT NULL,
    tags text DEFAULT '' NOT NULL, -- backend, frontend, fullstack
    CONSTRAINT fk_company_note_author FOREIGN KEY (author_id) REFERENCES account(id),
    CONSTRAINT fk_company_note_company FOREIGN KEY (company_id) REFERENCES company(id)
);
CREATE INDEX company_note_created ON company_note USING BRIN(created_at);
CREATE INDEX company_note_tags ON company_note (tags);
