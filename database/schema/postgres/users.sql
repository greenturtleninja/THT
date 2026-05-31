CREATE TABLE users (
  userID varchar (255) PRIMARY KEY,
  name varchar(255) NOT NULL,
  email varchar(255),
  phoneNumber varchar(15) DEFAULT '',
  createdTimestamp TEXT DEFAULT current_timestamp,
  updatedTimestamp TEXT,
  status status_type NOT NULL DEFAULT 'active'
);

GRANT ALL PRIVILEGES ON TABLE users TO eagle_bank;