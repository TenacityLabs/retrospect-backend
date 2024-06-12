CREATE TABLE IF NOT EXISTS capsules (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `code` VARCHAR(10) NOT NULL,
  `createdAt` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

  -- capsules can have 1 owner, and 3 additional members
  `public` BOOLEAN NOT NULL DEFAULT FALSE, -- public refers to whether the time capsule can be joined via code
  `capsuleOwnerId` INT UNSIGNED NOT NULL,
  `capsuleMember1Id` INT UNSIGNED,
  `capsuleMember2Id` INT UNSIGNED,
  `capsuleMember3Id` INT UNSIGNED,

  -- FIXME: input actual vessel names later
  `vessel` ENUM('box', 'suitcase', 'guitar_case', 'bottle', 'shoe'),

  `dateToOpen` TIMESTAMP, -- date to open refers to the date when the time capsule can be opened
  `sealed` BOOLEAN NOT NULL DEFAULT FALSE, -- sealed refers to whether the time capsule has been sealed (and can't be opened until the specified date)

  FOREIGN KEY (`capsuleOwnerId`) REFERENCES users(`id`),
  FOREIGN KEY (`capsuleMember1Id`) REFERENCES users(`id`),
  FOREIGN KEY (`capsuleMember2Id`) REFERENCES users(`id`),
  FOREIGN KEY (`capsuleMember3Id`) REFERENCES users(`id`),

  -- TODO: finish up with other fields such as songs, photos, etc. after meeting

  PRIMARY KEY (`id`),
  UNIQUE KEY `code` (`code`)
  CONSTRAINT chk_code_length CHECK (LENGTH(code) = 10)
);
