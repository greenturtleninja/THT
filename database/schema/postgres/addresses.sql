CREATE TABLE addresses (
  addressID varchar (255) PRIMARY KEY,
  userID varchar (255) NOT NULL,
  line1 varchar(255) NOT NULL,
  line2 varchar(255),
  line3 varchar(255),
  town varchar(255),
  county varchar(255),
  postcode varchar(10) NOT NULL,
  status status_type NOT NULL DEFAULT 'active'
);

GRANT ALL PRIVILEGES ON TABLE addresses TO eagle_bank;