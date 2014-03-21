--
-- Accounts
--
CREATE TABLE accounts (
  id          BIGSERIAL    NOT NULL,
  identity_id BIGINT       NOT NULL,
  -- provider_id BIGINT    NOT NULL,
  remote_id   VARCHAR(128) NOT NULL,
  name        VARCHAR(256) NOT NULL,
  email       VARCHAR(512) NOT NULL,
  picture     VARCHAR(512) NOT NULL,
  raw_profile TEXT         NOT NULL,
  raw_token   TEXT         NOT NULL,

  PRIMARY KEY (id),
  FOREIGN KEY (identity_id) REFERENCES identities (id) ON DELETE CASCADE ON UPDATE CASCADE,
  -- FOREIGN KEY (provider_id) REFERENCES providers (id) ON DELETE CASCADE ON UPDATE CASCADE,
  UNIQUE      (remote_id),

  CONSTRAINT remote_id_length  CHECK (char_length(remote_id) > 0),
  CONSTRAINT name_length       CHECK (char_length(name) > 0),
  CONSTRAINT email_length      CHECK (char_length(email) > 0)
);
