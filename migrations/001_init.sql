CREATE TABLE users (
                       id         BIGSERIAL PRIMARY KEY,
                       name       VARCHAR(255) NOT NULL,
                       email      VARCHAR(255) NOT NULL UNIQUE,
                       password   VARCHAR(255) NOT NULL,
                       role       VARCHAR(50)  NOT NULL DEFAULT 'client',
                       created_at TIMESTAMP    NOT NULL DEFAULT NOW()
);

CREATE TABLE tickets (
                         id          BIGSERIAL PRIMARY KEY,
                         title       VARCHAR(255) NOT NULL,
                         description TEXT         NOT NULL,
                         status      VARCHAR(50)  NOT NULL DEFAULT 'open',
                         priority    VARCHAR(50)  NOT NULL DEFAULT 'medium',
                         author_id   BIGINT       NOT NULL REFERENCES users(id),
                         assignee_id BIGINT       REFERENCES users(id),
                         created_at  TIMESTAMP    NOT NULL DEFAULT NOW(),
                         updated_at  TIMESTAMP    NOT NULL DEFAULT NOW()
);

CREATE TABLE comments (
                          id         BIGSERIAL PRIMARY KEY,
                          ticket_id  BIGINT    NOT NULL REFERENCES tickets(id) ON DELETE CASCADE,
                          user_id    BIGINT    NOT NULL REFERENCES users(id),
                          text       TEXT      NOT NULL,
                          created_at TIMESTAMP NOT NULL DEFAULT NOW()
);