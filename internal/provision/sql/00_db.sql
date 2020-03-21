-- From: https://dba.stackexchange.com/a/117661
--       https://stackoverflow.com/a/28849656/9911189

-- Create database
CREATE DATABASE contradb;
REVOKE ALL ON DATABASE contradb FROM public;

-- Create users for new schema
CREATE USER contra_mgr WITH ENCRYPTED PASSWORD 'uTiXe3oYJDv9Z4Ef';
CREATE USER contra_usr WITH ENCRYPTED PASSWORD 'EvPvkro59Jb7RK3o';
CREATE USER contra_ro  WITH ENCRYPTED PASSWORD 'G2e2e6frXA8ytod5';

GRANT contra_usr TO contra_mgr;
GRANT contra_ro  TO contra_usr;

GRANT CONNECT ON DATABASE contradb TO contra_ro; -- others inherit

-- Create new schema
\connect contradb

CREATE SCHEMA contra AUTHORIZATION contra_mgr;

SET search_path = contra;

-- These are not inheritable
ALTER ROLE contra_mgr IN DATABASE contradb SET search_path = contra;
ALTER ROLE contra_usr IN DATABASE contradb SET search_path = contra;
ALTER ROLE contra_ro  IN DATABASE contradb SET search_path = contra;

GRANT CREATE ON SCHEMA contra TO contra_mgr;
GRANT USAGE  ON SCHEMA contra TO contra_ro ; -- contra_usr inherits

-- Set default privileges
-- -> Read only
ALTER DEFAULT PRIVILEGES FOR ROLE contra_mgr GRANT SELECT ON TABLES TO contra_ro;

-- -> Read only for sequences
--    Not recommended; this is read-write because users with USAGE can use nextval()
-- ALTER DEFAULT PRIVILEGES FOR ROLE contra_mgr GRANT USAGE ON SEQUENCES TO contra_ro;

-- -> Read/write
ALTER DEFAULT PRIVILEGES FOR ROLE contra_mgr GRANT INSERT, UPDATE, DELETE, TRUNCATE ON TABLES TO contra_usr;

-- -> Read/write for sequences
ALTER DEFAULT PRIVILEGES FOR ROLE contra_mgr GRANT USAGE, SELECT, UPDATE ON SEQUENCES TO contra_usr;
