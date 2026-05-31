CREATE TABLE accounts (
  accountID varchar (255) PRIMARY KEY,
  name varchar (255) NOT NULL,
  accountNumber varchar(8) NOT NULL,
  sortCode varchar(8) NOT NULL,
  accountType account_type NOT NULL DEFAULT 'personal',
  balance BIGINT DEFAULT 0,
  currency currency_type NOT NULL DEFAULT 'GBP',
  createdTimestamp TEXT DEFAULT current_timestamp,
  updatedTimestamp TEXT,
  status status_type NOT NULL DEFAULT 'active'
);

GRANT ALL PRIVILEGES ON TABLE accounts TO eagle_bank;