CREATE table users(
    id BIGINT NOT NULL AUTO_INCREMENT,
    google_sub VARBINARY(255) NOT NULL,
    created_at TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    PRIMARY KEY (id),
    UNIQUE KEY uq_users_google_sub (google_sub)
) ENGINE = InnoDB;