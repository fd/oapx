--
-- Applications
--
CREATE TABLE applications (
  id       BIGSERIAL    NOT NULL,
  owner_id BIGINT       NOT NULL,
  name     VARCHAR(128) NOT NULL,

  client_id     VARCHAR(16) NOT NULL,
  client_secret VARCHAR(32) NOT NULL,
  redirect_uri  VARCHAR(512) NOT NULL,

  PRIMARY KEY (id),
  FOREIGN KEY (owner_id) REFERENCES identities (id) ON DELETE RESTRICT ON UPDATE CASCADE,
  UNIQUE      (client_id),

  CONSTRAINT name_length          CHECK (char_length(name) > 0),
  CONSTRAINT client_id_length     CHECK (char_length(client_id) = 16),
  CONSTRAINT client_secret_length CHECK (char_length(client_secret) = 32),
  CONSTRAINT redirect_uri_length  CHECK (char_length(redirect_uri) > 0)
);
