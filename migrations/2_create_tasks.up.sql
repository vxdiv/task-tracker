START TRANSACTION;

CREATE TABLE tasks
(
  id          int unsigned auto_increment                                                                NOT NULL,
  name        varchar(255)                                                                               NOT NULL,
  description text                                                                                       NOT NULL DEFAULT '',
  type        ENUM ('improvement', 'future', 'bug')                                                      NOT NULL DEFAULT 'improvement',
  owner_id    int unsigned                                                                               NOT NULL,
  assigned_id int unsigned                                                                               NOT NULL,
  status      ENUM ('open', 'in_progress', 'resolve', 'close', 'hold', 'reopen')                         NOT NULL DEFAULT 'open',
  due_date    timestamp                                                                                  NOT NULL,
  resolution  ENUM ('done', 'fixed', 'duplicate', 'incomplete', 'cannot_reproduce', 'do_not_need_to_do') NOT NULL DEFAULT 'done',
  priority    ENUM ('trivial', 'major', 'critical', 'asap')                                              NOT NULL DEFAULT 'major',
  created_at  timestamp                                                                                  NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at  timestamp                                                                                  NOT NULL DEFAULT CURRENT_TIMESTAMP
    ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (id)
)
  ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci;

COMMIT;
