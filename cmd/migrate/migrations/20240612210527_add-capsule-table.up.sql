CREATE TABLE IF NOT EXISTS capsules (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `code` VARCHAR(10) NOT NULL,
  `createdAt` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

  -- capsules can have 1 owner, and 3 additional members
  `public` BOOLEAN NOT NULL, -- public refers to whether the time capsule can be joined via code
  `capsuleOwnerId` INT UNSIGNED NOT NULL,
  `capsuleMember1Id` INT UNSIGNED,
  `capsuleMember2Id` INT UNSIGNED,
  `capsuleMember3Id` INT UNSIGNED,

  `vessel` ENUM('box', 'suitcase', 'guitar case', 'bottle', 'shoe', 'garbage') NOT NULL,

  `name` VARCHAR(255) NOT NULL,
  `dateToOpen` TIMESTAMP, -- date to open refers to the date when the time capsule can be opened
  `emailSent` BOOLEAN NOT NULL DEFAULT FALSE, -- whether or not an email has been sent to the capsule owner
  `sealed` ENUM('preseal', 'sealed', 'opened') DEFAULT 'preseal', -- sealed refers to whether the time capsule has been sealed (and can't be opened until the specified date)

  FOREIGN KEY (`capsuleOwnerId`) REFERENCES users(`id`),
  FOREIGN KEY (`capsuleMember1Id`) REFERENCES users(`id`),
  FOREIGN KEY (`capsuleMember2Id`) REFERENCES users(`id`),
  FOREIGN KEY (`capsuleMember3Id`) REFERENCES users(`id`),

  -- TODO: finish up with other fields such as songs, photos, etc. after meeting

  PRIMARY KEY (`id`),
  UNIQUE KEY `code` (`code`)
);
