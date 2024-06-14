CREATE TABLE IF NOT EXISTS songs (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `userId` INT UNSIGNED NOT NULL,
  `capsuleId` INT UNSIGNED NOT NULL,

  `spotifyId` VARCHAR(255) NOT NULL, -- spotify ids should be 22 long, but we'll give some buffer room
  `name` VARCHAR(255) NOT NULL,
  `artistName` VARCHAR(255) NOT NULL,
  `albumArtURL` VARCHAR(255),

  `createdAt` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

  PRIMARY KEY (`id`),
  FOREIGN KEY (`userId`) REFERENCES users(`id`),
  FOREIGN KEY (`capsuleId`) REFERENCES capsules(`id`)
);
