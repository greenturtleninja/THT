CREATE TABLE transactions (
  transactionID varchar (255) PRIMARY KEY,
  amount BIGINT DEFAULT 0,
  currency varchar(3) CHECK (currency IN ('GBP')) NOT NULL DEFAULT 'GBP',
  type varchar(12) CHECK (type IN ('deposit', 'withdrawal')) NOT NULL DEFAULT 'deposit',
  createdTimestamp TEXT DEFAULT current_timestamp,
  reference varchar(255),
  userID varchar(255)
);
