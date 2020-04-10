-- From: https://dba.stackexchange.com/a/117661
--       https://stackoverflow.com/a/28849656/9911189

-- Create database
CREATE DATABASE contradb;
REVOKE ALL ON DATABASE contradb FROM public;

-- Create users for new schemas
CREATE USER contracore_mgr WITH ENCRYPTED PASSWORD 'uTiXe3oYJDv9Z4Ef';
CREATE USER contracore_usr WITH ENCRYPTED PASSWORD 'EvPvkro59Jb7RK3o';
CREATE USER contracore_ro  WITH ENCRYPTED PASSWORD 'G2e2e6frXA8ytod5';
CREATE USER contraweb_mgr  WITH ENCRYPTED PASSWORD '95WaR3vBbC23RqBv';
CREATE USER contraweb_usr  WITH ENCRYPTED PASSWORD 'U475jBKZfK3xhbVZ';
CREATE USER contraweb_ro   WITH ENCRYPTED PASSWORD 'oUkDtezN2Z2kF83c';

GRANT contracore_usr TO contracore_mgr;
GRANT contracore_ro  TO contracore_usr;
GRANT contraweb_usr  TO contraweb_mgr;
GRANT contraweb_ro   TO contraweb_usr;

GRANT CONNECT ON DATABASE contradb TO contracore_ro; -- others inherit
GRANT CONNECT ON DATABASE contradb TO contraweb_ro; -- others inherit

\connect contradb

-- Create ContraCore schema
CREATE SCHEMA contracore AUTHORIZATION contracore_mgr;

SET search_path = contracore;

-- These are not inheritable
ALTER ROLE contracore_mgr IN DATABASE contradb SET search_path = contracore;
ALTER ROLE contracore_usr IN DATABASE contradb SET search_path = contracore;
ALTER ROLE contracore_ro  IN DATABASE contradb SET search_path = contracore;

GRANT CREATE ON SCHEMA contracore TO contracore_mgr;
GRANT USAGE  ON SCHEMA contracore TO contracore_ro ; -- contra_usr inherits

-- Set default privileges
-- -> Read only
ALTER DEFAULT PRIVILEGES FOR ROLE contracore_mgr GRANT SELECT ON TABLES TO contracore_ro;

-- -> Read only for sequences
--    Not recommended; this is read-write because users with USAGE can use nextval()
-- ALTER DEFAULT PRIVILEGES FOR ROLE contra_mgr GRANT USAGE ON SEQUENCES TO contra_ro;

-- -> Read/write
ALTER DEFAULT PRIVILEGES FOR ROLE contracore_mgr GRANT INSERT, UPDATE, DELETE, TRUNCATE ON TABLES TO contracore_usr;

-- -> Read/write for sequences
ALTER DEFAULT PRIVILEGES FOR ROLE contracore_mgr GRANT USAGE, SELECT, UPDATE ON SEQUENCES TO contracore_usr;



-- Create ContraWeb schema
CREATE SCHEMA contraweb AUTHORIZATION contraweb_mgr;

SET search_path = contraweb;

-- These are not inheritable
ALTER ROLE contraweb_mgr IN DATABASE contradb SET search_path = contraweb;
ALTER ROLE contraweb_usr IN DATABASE contradb SET search_path = contraweb;
ALTER ROLE contraweb_ro  IN DATABASE contradb SET search_path = contraweb;

GRANT CREATE ON SCHEMA contraweb TO contraweb_mgr;
GRANT USAGE  ON SCHEMA contraweb TO contraweb_ro ; -- contra_usr inherits

-- Set default privileges
-- -> Read only
ALTER DEFAULT PRIVILEGES FOR ROLE contraweb_mgr GRANT SELECT ON TABLES TO contraweb_ro;

-- -> Read only for sequences
--    Not recommended; this is read-write because users with USAGE can use nextval()
-- ALTER DEFAULT PRIVILEGES FOR ROLE contra_mgr GRANT USAGE ON SEQUENCES TO contra_ro;

-- -> Read/write
ALTER DEFAULT PRIVILEGES FOR ROLE contraweb_mgr GRANT INSERT, UPDATE, DELETE, TRUNCATE ON TABLES TO contraweb_usr;

-- -> Read/write for sequences
ALTER DEFAULT PRIVILEGES FOR ROLE contraweb_mgr GRANT USAGE, SELECT, UPDATE ON SEQUENCES TO contraweb_usr;
