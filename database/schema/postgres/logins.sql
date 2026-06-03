CREATE TABLE logins (
  userID varchar (255) PRIMARY KEY,
  displayName varchar(255) NOT NULL,
  login varchar(255) NOT NULL,
  email varchar(255),
  passwordHash varchar(255) NOT NULL,
  userType user_type NOT NULL,
  createdTimestamp TIMESTAMP WITH TIME ZONE DEFAULT current_timestamp(3),
  updatedTimestamp TIMESTAMP WITH TIME ZONE DEFAULT current_timestamp(3)
);

CREATE OR REPLACE FUNCTION update_modified_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updatedTimestamp = current_timestamp(3);
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_logins_modtime
BEFORE UPDATE ON logins
FOR EACH ROW
EXECUTE FUNCTION update_modified_column();

GRANT ALL PRIVILEGES ON TABLE logins TO eagle_bank;