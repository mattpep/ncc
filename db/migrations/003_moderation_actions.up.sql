CREATE TYPE mod_action AS ENUM ('flag', 'approve');
CREATE TABLE moderation_actions (
	id SERIAL UNIQUE,
	comment_id INT NOT NULL,
	action mod_action NOT NULL,
	date_time TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	actor inet,
	CONSTRAINT fk_comment FOREIGN KEY(comment_id) REFERENCES comments(id) ON DELETE CASCADE
);
