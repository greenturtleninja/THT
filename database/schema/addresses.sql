CREATE TABLE addresses (
  addressID varchar (255) PRIMARY KEY,
  userID varchar (255) NOT NULL,
  line1 varchar(255) NOT NULL,
  line2 varchar(255),
  line3 varchar(255),
  town varchar(255),
  county varchar(255),
  postcode varchar(10) NOT NULL,
  status varchar(8) CHECK (status IN ('active', 'deleted')) NOT NULL DEFAULT 'active'
);