CREATE TABLE IF NOT EXISTS songs (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `capsuleId` INT UNSIGNED NOT NULL,
  `userId` INT UNSIGNED NOT NULL,

  `spotifyId` VARCHAR(63) NOT NULL, -- spotify ids should be 22 long, but we'll give some buffer room
  `name` VARCHAR(255) NOT NULL,
  `artist` VARCHAR(255) NOT NULL,
  `albumName` VARCHAR(255) NOT NULL,
  `albumArtURL` VARCHAR(255),

  PRIMARY KEY (`id`),
  FOREIGN KEY (`capsuleId`) REFERENCES capsules(`id`)
  FOREIGN KEY (`userId`) REFERENCES users(`id`)
);
