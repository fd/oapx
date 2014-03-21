--
-- Access Tokens
--
CREATE TABLE access_tokens (
  id               BIGSERIAL   NOT NULL,
  identity_id      BIGINT      NOT NULL,
  application_id   BIGINT      NOT NULL,
  access_token     VARCHAR(64) NOT NULL,
  refresh_token    VARCHAR(64),
  expires_in       INTEGER     NOT NULL,
  created_at       TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP,

  scope        VARCHAR(1024) NOT NULL,
  redirect_uri VARCHAR(1024) NOT NULL,

  PRIMARY KEY (id),
  FOREIGN KEY (identity_id)    REFERENCES identities (id) ON DELETE CASCADE ON UPDATE CASCADE,
  FOREIGN KEY (application_id) REFERENCES applications (id) ON DELETE CASCADE ON UPDATE CASCADE,
  UNIQUE      (access_token),
  UNIQUE      (refresh_token)
);
