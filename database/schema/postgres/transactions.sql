CREATE TABLE transactions (
  transactionID varchar (255) PRIMARY KEY,
  amount BIGINT DEFAULT 0,
  currency currency_type NOT NULL DEFAULT 'GBP',
  accountType account_type NOT NULL DEFAULT 'personal',
  createdTimestamp TEXT DEFAULT current_timestamp,
  reference varchar(255),
  userID varchar(255)
);

GRANT ALL PRIVILEGES ON TABLE transactions TO eagle_bank;