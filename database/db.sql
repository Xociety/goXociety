-- check all PRIMARY KEY, FOREIGN KEY, data valid range, CRUD relation, related api, reference sequence, schema
-- alter table XXX DROP CONTAINT XXX;
-- alter table post ADD CONSTAINT post_type_fkey FOREIGN KEY (type) REFERENCES post_type(post_type_id) ON DELETE CASCADE ON UPDATE CASCADE;

------------------------
CREATE TABLE country (
    country_id          integer PRIMARY KEY NOT NULL,
    country             VARCHAR(200),
    country_code        VARCHAR(2)
);
COPY country(country_id, country, country_code) FROM '/var/lib/postgresql/data/country.csv' DELIMITER ';' CSV; -- HEADER;
/*
https://developers.google.com/custom-search/docs/xml_results_appendices#countryCodes
*/
------------------------
CREATE TABLE language (
    language_id         integer PRIMARY KEY NOT NULL,
    display_language    VARCHAR(200),
    value               VARCHAR(10)
);
COPY language(language_id, display_language, value) FROM '/var/lib/postgresql/data/language.csv' DELIMITER ';' CSV;
/*
https://developers.google.com/custom-search/docs/xml_results_appendices#interfaceLanguages
*/
------------------------
CREATE TABLE actions (
    action_id           integer PRIMARY KEY NOT NULL,
    value               VARCHAR(15)
);
COPY actions(action_id, name) FROM '/var/lib/postgresql/data/actions.csv' DELIMITER ';' CSV;
------------------------
CREATE TABLE gender (
    gender_id           integer PRIMARY KEY NOT NULL,
    value               VARCHAR(15)
);
COPY gender(gender_id, value) FROM '/var/lib/postgresql/data/gender.csv' DELIMITER ';' CSV;
------------------------
CREATE TABLE post_type (
    post_type_id        integer PRIMARY KEY NOT NULL,
    value               VARCHAR(10)
);
COPY post_type(post_type_id, value) FROM '/var/lib/postgresql/data/post_type.csv' DELIMITER ';' CSV HEADER;
------------------------
CREATE TABLE xuser(
    user_id             BIGSERIAL PRIMARY KEY NOT NULL,
    username            VARCHAR(50) UNIQUE NOT NULL, -- [index]
    email               VARCHAR(50), -- [UNIQUE,index] and format check
    password            VARCHAR(100) NOT NULL,
    name                VARCHAR(100) NOT NULL,
    phone               VARCHAR(30) NOT NULL, -- add country calling code
    gender              integer, -- 0: not known, 1: male, 2:female
    bio                 VARCHAR(100),
    credit              double precision,
    photo_url           VARCHAR(200),
    language_id         integer references language(langauge_id),
    country_id          integer references country(country_id),
    timezone            integer,
    last_ip             VARCHAR(15), -- ipv4 123.194.188.0
    -- apid
    createtime          integer, -- unix time
    updatetime          integer
);
INSERT INTO xuser 
(username, email, password, name, phone, gender, bio, credit, photo_url, language_id, country_id, timezone, last_ip, createtime, updatetime) 
VALUES ('jeff', 'jeff@gmail.com', 'salted', 'jeff', '+886-911111111', 1, 'hi', 0, '', 12, 206, 28800, '123.194.188.0', 1527496777, 1527496777),
('kyler', 'kyler@gmail.com', 'salted', 'kyler', '+886-911111111', 1, 'yo', 0, '', 12, 206, 28800, '123.194.188.0', 1527496777, 1527496777),
('robby', 'robby@gmail.com', 'salted', 'robby', '+886-911111111', 1, 'man', 0, '', 12, 206, 28800, '123.194.188.0', 1527496777, 1527496777);
-- username, phone + country code, email (lower case) logic check
------------------------
CREATE TABLE follow (
    following_user_id   bigint references xuser(user_id) ON DELETE CASCADE ON UPDATE CASCADE, -- a person whom you follow
    follower_user_id    bigint references xuser(user_id) ON DELETE CASCADE ON UPDATE CASCADE, -- a person who follows you
    valid               boolean, -- update by following_user, block user usage?
    createtime          integer, -- [index] -- create and delete by follower_user
    updatetime          integer -- update by following_user
    -- CONSTRAINT target UNIQUE (following_user_id, follower_user_id)
    -- CONSTRAINT self CHECK (followed_user_id <> follower_user_id) NOT VALID
);
-- following limit 7500 on Instagram, we might limit following 5000 people
------------------------
CREATE TABLE post (
    post_id             BIGSERIAL PRIMARY KEY NOT NULL,
    user_id             bigint references xuser(user_id) ON DELETE CASCADE ON UPDATE CASCADE,
    content             VARCHAR(300), -- include @tagid, #hashtag
    blob_id             VARCHAR(200), -- may be foldername
    -- blob_count       int, -- multiple images, videos
    type                integer references post_type(post_type_id),
    like_count          bigint,
    dislike_count       bigint,
    comment_count       bigint,
    point               point,
    country_id          integer references country(country_id),
    category_id         integer, -- references
    public              boolean, -- public post in the future
    createtime          integer, -- [index] -- timestamp without time zone, unix time || *index
    updatetime          integer
    -- constraint like_count (check (like_count >= 0))
    -- constraint dislike_count (check (dislike_count >= 0))
    -- constraint comment_count (check (dislike_count >= 0))
);
--  maximum 30 hashtag
------------------------
CREATE TABLE hashtag(
   hashtag_id           BIGSERIAL PRIMARY KEY,
   value                VARCHAR(100) UNIQUE NOT NULL,
   count                bigint
);
-- define hashtag length, space check
------------------------
CREATE TABLE post_hashtag(
   hashtag_id           bigint NOT NULL references hashtag(hashtag_id) -- [primary key or index]?
   post_id              bigint NOT NULL references post(post_id) ON DELETE CASCADE ON UPDATE CASCADE,
 );
------------------------
CREATE TABLE post_tag_xuser (
    post_id             bigint references post(post_id) ON DELETE CASCADE ON UPDATE CASCADE, -- [primary key]
    user_id             bigint references xuser(user_id) ON DELETE CASCADE ON UPDATE CASCADE,
    x                   integer, -- percentage 0-99
    y                   integer, -- percentage 0-99
    valid               boolean, -- valid update by xuser
    createtime          integer, -- [index]
    updatetime          integer
);
------------------------
INSERT INTO post 
(user_id, content, blob_id, type, like_count, dislike_count, point, country_id, category_id, public, createtime, updatetime) 
VALUES 
(4, 'hello world #happy @jeff', 'sha256 hashed id', 0, 0, 0, point('121.5643,25.0336'), 206, 0, true, 1527498044, 1527498044);
-- UPDATE post SET content='hello world #happy @jeff', blob_id='sha256 hashed id' WHERE post_id = 1;
INSERT INTO hashtag 
(name) 
VALUES ('happy');
INSERT INTO post_hashtag 
(post_id, hashtag_id) 
VALUES (33, 2);
INSERT INTO post_tag_xuser 
(post_id, user_id, x, y, valid, createtime, updatetime) 
VALUES (33, 4, 0, 0, false, 1527498711, 1527498711);
------------------------
CREATE TABLE post_actions (
    post_id             bigint references post(post_id) ON DELETE CASCADE ON UPDATE CASCADE, -- [primary key]
    user_id             bigint references xuser(user_id) ON DELETE CASCADE ON UPDATE CASCADE,
    act                 integer reference actions(action_id), -- 0: like, 1: dislike
    createtime          integer -- [index]
    -- CONSTRAINT post_target UNIQUE (post_id, user_id)
);
INSERT INTO post_actions 
(post_id, user_id, act, createtime) 
VALUES (33, 4, 0, 1527498711);
------------------------
CREATE TABLE comments (
    comment_id          BIGSERIAL PRIMARY KEY NOT NULL,
    post_id             bigint references post(post_id) ON DELETE CASCADE ON UPDATE CASCADE,
    user_id             bigint references xuser(user_id) ON DELETE CASCADE ON UPDATE CASCADE,
    comment             VARCHAR(300),
    like_count          bigint,
    dislike_count       bigint,
    comment_count       bigint,
    createtime          integer, -- [index]
    updatetime          integer
    -- constraint like_count (check (like_count >= 0))
    -- constraint dislike_count (check (dislike_count >= 0))
    -- constraint comment_count (check (comment_count >= 0))
);
------------------------
CREATE TABLE comment_actions (
    comment_id          bigint references comments(comment_id), -- [PRIMARY KEY]
    user_id             bigint references xuser(user_id) ON DELETE CASCADE ON UPDATE CASCADE,
    act                 integer references actions(action_id),
    createtime          integer
);
INSERT INTO comments 
(post_id, user_id, comment, createtime, updatetime) 
VALUES (33, 4, 'yo', 1527498711, 1527498711);
-- maximum taged user number: 5
INSERT INTO comment_actions 
(comment_id, user_id, act, createtime) 
VALUES (2, 4, 0, 1527498711);

SELECT * FROM post 
FULL OUTER JOIN xuser on post.user_id = xuser.user_id 
FULL OUTER JOIN post_actions on post.post_id = post_actions.post_id 
FULL OUTER JOIN comments on post.post_id = comments.post_id 
FULL OUTER JOIN comment_actions on comments.comment_id = comment_actions.comment_id;

------------------------
TABLE comments_deep (
    comment_deep_id     PRIMARY KEY BIGSERIAL NOT NULL,
    comment_id          bigint references comments(comment_id) ON DELETE CASCADE ON UPDATE CASCADE,
    user_id             bigint references xuser(user_id) ON DELETE CASCADE ON UPDATE CASCADE,
    comment             VARCHAR(300),
    like_count          bigint,
    dislike_count       bigint,
    createtime          integer,
    updatetime          integer
    -- constraint like_count (check (like_count >= 0))
    -- constraint dislike_count (check (dislike_count >= 0))
);
-- maximum taged user number: 5
------------------------
TABLE comments_deep_actions (
    comment_deep_id     bigint references comments_deep(comment_deep_id) ON DELETE CASCADE ON UPDATE CASCADE, -- [PRIMARY KEY]
    user_id             bigint references xuser(user_id) ON DELETE CASCADE ON UPDATE CASCADE,
    act                 integer reference actions(action_id),
    createtime          integer
);
------------------------
-- TABLE block ( -- this should be implement in follow valid field
--     user_id BIGSERIAL references xuser(user_id),
--     blocked_user_id BIGSERIAL references xuser(user_id),
--     createtime integer
-- );
------------------------
TABLE report (
    user_id
    report_user_id
    content
    createtime integer
);
------------------------
TABLE category
------------------------
TABLE sub_category
------------------------
TABLE saved_post
------------------------
TABLE cities
------------------------
TABLE notice (
    user_id
    type
    post_id
    comment_id
    comment_deep_id
    blob_id
    content
);