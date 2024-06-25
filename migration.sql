-- Disable foreign key checks
SET FOREIGN_KEY_CHECKS = 0;

# Dump of table swipes
# ------------------------------------------------------------
DROP TABLE IF EXISTS `swipes`;

# Dump of table swipe_matches
# ------------------------------------------------------------
DROP TABLE IF EXISTS `swipe_matches`;

# Dump of table users
# ------------------------------------------------------------
DROP TABLE IF EXISTS `users`;

# Create table users
# ------------------------------------------------------------
CREATE TABLE `users`
(
    `id`         char(36) COLLATE utf8mb4_unicode_ci     NOT NULL,
    `name`       varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
    `email`      varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
    `password`   varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
    `age`        int                                     NOT NULL,
    `gender`     enum('male','female') COLLATE utf8mb4_unicode_ci NOT NULL,
    `longitude`  decimal(10, 7)                          NOT NULL,
    `latitude`   decimal(10, 7)                          NOT NULL,
    `created_at` timestamp NULL DEFAULT NULL,
    `updated_at` timestamp NULL DEFAULT NULL,
    `deleted_at` timestamp NULL DEFAULT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `users_email_unique` (`email`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

# Create table swipe_matches
# ------------------------------------------------------------
CREATE TABLE `swipe_matches`
(
    `id`         char(36) COLLATE utf8mb4_unicode_ci NOT NULL,
    `match_one`  char(36) COLLATE utf8mb4_unicode_ci NOT NULL,
    `match_two`  char(36) COLLATE utf8mb4_unicode_ci NOT NULL,
    `created_at` timestamp NULL DEFAULT NULL,
    `updated_at` timestamp NULL DEFAULT NULL,
    PRIMARY KEY (`id`),
    KEY          `swipe_matches_match_one_foreign` (`match_one`),
    KEY          `swipe_matches_match_two_foreign` (`match_two`),
    CONSTRAINT `swipe_matches_match_one_foreign` FOREIGN KEY (`match_one`) REFERENCES `users` (`id`),
    CONSTRAINT `swipe_matches_match_two_foreign` FOREIGN KEY (`match_two`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

# Create table swipes
    # ------------------------------------------------------------
CREATE TABLE `swipes`
(
    `id`         char(36) COLLATE utf8mb4_unicode_ci NOT NULL,
    `user_id`    char(36) COLLATE utf8mb4_unicode_ci NOT NULL,
    `swiped_on`  char(36) COLLATE utf8mb4_unicode_ci NOT NULL,
    `preference` enum('yes','no') COLLATE utf8mb4_unicode_ci NOT NULL,
    `created_at` timestamp NULL DEFAULT NULL,
    `updated_at` timestamp NULL DEFAULT NULL,
    PRIMARY KEY (`id`),
    KEY          `swipes_user_id_foreign` (`user_id`),
    KEY          `swipes_swiped_on_foreign` (`swiped_on`),
    CONSTRAINT `swipes_swiped_on_foreign` FOREIGN KEY (`swiped_on`) REFERENCES `users` (`id`),
    CONSTRAINT `swipes_user_id_foreign` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Re-enable foreign key checks
SET FOREIGN_KEY_CHECKS = 1;
