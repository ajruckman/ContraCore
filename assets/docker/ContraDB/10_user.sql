\c contradb contraweb_mgr

CREATE TABLE IF NOT EXISTS contraweb.user
(
    username VARCHAR(16) NOT NULL,

    salt     VARCHAR(32) NOT NULL,
    password VARCHAR(64) NOT NULL,
    role     VARCHAR(13) NOT NULL,

    CONSTRAINT user_pk PRIMARY KEY (username),

    CONSTRAINT user_role_chk CHECK (role IN ('restricted', 'privileged', 'administrator'))
);

CREATE TABLE IF NOT EXISTS contraweb.user_session
(
    username     VARCHAR(16) NOT NULL,
    token        VARCHAR(32) NOT NULL,
    created_at   TIMESTAMP   NOT NULL,
    refreshed_at TIMESTAMP   NOT NULL,

    CONSTRAINT user_session_pk PRIMARY KEY (username),
    CONSTRAINT user_session_user_id_fk FOREIGN KEY (username) REFERENCES contraweb.user (username)
);
