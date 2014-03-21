--
-- Authorizations
--
CREATE TABLE authorizations (
  id             BIGSERIAL   NOT NULL,
  identity_id    BIGINT      NOT NULL,
  application_id BIGINT      NOT NULL,
  code           VARCHAR(64) NOT NULL,
  expires_in     INTEGER     NOT NULL,
  created_at     TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP,

  state        VARCHAR(512)  NOT NULL,
  scope        VARCHAR(1024) NOT NULL,
  redirect_uri VARCHAR(1024) NOT NULL,

  PRIMARY KEY (id),
  FOREIGN KEY (identity_id)    REFERENCES identities (id) ON DELETE CASCADE ON UPDATE CASCADE,
  FOREIGN KEY (application_id) REFERENCES applications (id) ON DELETE CASCADE ON UPDATE CASCADE,
  UNIQUE      (code)
);
