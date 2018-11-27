START TRANSACTION;

CREATE TABLE users (
  id            int unsigned auto_increment NOT NULL,
  name          varchar(32)                 NOT NULL         DEFAULT '',
  email         varchar(120)                NOT NULL,
  status        ENUM ('active', 'disabled') NOT NULL         DEFAULT 'active',
  password_hash binary(60)                  NOT NULL,
  created_at    timestamp                   NOT NULL         DEFAULT CURRENT_TIMESTAMP,
  updated_at    timestamp                   NOT NULL         DEFAULT CURRENT_TIMESTAMP
  ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  UNIQUE (email)
)
  ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci;

COMMIT;