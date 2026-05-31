CREATE TABLE users_accounts (
  accountID varchar (255) NOT NULL,
  userID varchar (255) NOT NULL,
  FOREIGN KEY(accountID) REFERENCES accounts(accountID),
  FOREIGN KEY(userID) REFERENCES users(userID)
);

GRANT ALL PRIVILEGES ON TABLE users_accounts TO eagle_bank;