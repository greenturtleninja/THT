CREATE TABLE accounts (
  accountID varchar (255) PRIMARY KEY,
  name varchar (255) NOT NULL,
  accountNumber varchar(8) NOT NULL,
  sortCode varchar(8) NOT NULL,
  accountType varchar(12) CHECK (accountType IN ('personal', 'business')) NOT NULL DEFAULT 'personal',
  balance BIGINT DEFAULT 0,
  currency varchar(3) CHECK (currency IN ('GBP')) NOT NULL DEFAULT 'GBP',
  createdTimestamp TEXT DEFAULT current_timestamp,
  updatedTimestamp TEXT,
  status varchar(8) CHECK (status IN ('active', 'deleted')) NOT NULL DEFAULT 'active'
);