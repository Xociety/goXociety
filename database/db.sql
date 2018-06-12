-- check all PRIMARY KEY, FOREIGN KEY, data valid range, CRUD relation, related api, reference sequence, schema

-- alter forgign key:alter table post ADD CONSTAINT post_type_fkey FOREIGN KEY (type) REFERENCES post_type(post_type_id)

------------------------
CREATE TABLE country (
    country_id integer PRIMARY KEY NOT NULL,
    name VARCHAR(200),
    code VARCHAR(2)
);
COPY country(country_id, name, code) FROM '/var/lib/postgresql/data/country.csv' DELIMITER ';' CSV; -- HEADER;
/*
https://developers.google.com/custom-search/docs/xml_results_appendices#countryCodes
*/
------------------------
CREATE TABLE language (
    language_id integer PRIMARY KEY NOT NULL,
    display_language VARCHAR(200),
    value VARCHAR(10)
);
COPY language(language_id, display_language, value) FROM '/var/lib/postgresql/data/language.csv' DELIMITER ';' CSV;
/*
https://developers.google.com/custom-search/docs/xml_results_appendices#interfaceLanguages
*/
------------------------
CREATE TABLE actions (
    action_id integer PRIMARY KEY NOT NULL,
    name VARCHAR(15)
);
COPY actions(action_id, name) FROM '/var/lib/postgresql/data/actions.csv' DELIMITER ';' CSV;
------------------------
CREATE TABLE xuser(
    user_id BIGSERIAL    PRIMARY KEY NOT NULL,
    username        VARCHAR(50) UNIQUE NOT NULL,
    email           VARCHAR(50), -- unique and format check
    password VARCHAR(100) NOT NULL,
    name            VARCHAR(100)    NOT NULL,
    phone           VARCHAR(30) NOT NULL, -- add country calling code
    gender          integer, -- 0: not known, 1: male, 2:female
    bio             VARCHAR(50),
    credit          double precision,
    photo_url    VARCHAR(200),
    language_id        integer references language(langauge_id),
    country_id         integer references country(country_id),
    timezone        integer,
    last_ip         VARCHAR(15), -- ipv4 123.194.188.0
    -- apid
    createtime      integer, -- unix time
    updatetime      integer
);
INSERT INTO xuser 
(username, email, password, name, phone, gender, bio, credit, photo_url, language_id, country_id, timezone, last_ip, createtime, updatetime) 
VALUES ('jeff', 'jeff@gmail.com', 'salted', 'jeff', '+886-911111111', 1, 'hi', 0, '', 12, 206, 28800, '123.194.188.0', 1527496777, 1527496777),
('kyler', 'kyler@gmail.com', 'salted', 'kyler', '+886-911111111', 1, 'yo', 0, '', 12, 206, 28800, '123.194.188.0', 1527496777, 1527496777),
('robby', 'robby@gmail.com', 'salted', 'robby', '+886-911111111', 1, 'man', 0, '', 12, 206, 28800, '123.194.188.0', 1527496777, 1527496777);
-- username, phone + country code, email (lower case) logic check
/*
gender -- ISO/IEC 5218
0 = not known, 1 = male, 2 = female, 9 = not applicable.
ref: https://en.wikipedia.org/wiki/ISO/IEC_5218

language -- ISO 639-1
ref: https://en.wikipedia.org/wiki/List_of_ISO_639-1_codes

country
https://en.wikipedia.org/wiki/ISO_3166-1
*/
------------------------
CREATE TABLE follow (
    following_user_id bigint references xuser(user_id),
    follower_user_id bigint references xuser(user_id),
    valid boolean,
    createtime integer
    -- CONSTRAINT target UNIQUE (following_user_id, follower_user_id)
    -- CONSTRAINT self CHECK (followed_user_id <> follower_user_id) NOT VALID
);
-- following limit 7500 on Instagram, we might limit following 5000 people
------------------------
CREATE TABLE post_type (
    post_type_id integer PRIMARY KEY NOT NULL,
    name VARCHAR(10)
);
COPY post_type(post_type_id, name) FROM '/var/lib/postgresql/data/post_type.csv' DELIMITER ';' CSV HEADER;
-- INSERT INTO post_type (post_type_id, name) VALUES (0, 'image'), (1, 'video');
------------------------
CREATE TABLE post (
    post_id BIGSERIAL PRIMARY KEY     NOT NULL,
    user_id bigint	 references xuser(user_id),
    content VARCHAR(300), -- include @tagid, #hashtag
    blob_id VARCHAR(200), -- may be foldername
    -- notice: multiple images, videos
    type integer references post_type(post_type_id),
    point point,
    country_id integer references country(country_id),
    category_id integer,
    public boolean,
    createtime integer, -- timestamp without time zone, unix time || *index
    updatetime integer
);
--  maximum 30 hashtag
------------------------
CREATE TABLE hashtag(
   hashtag_id BIGSERIAL PRIMARY KEY,
   name VARCHAR(100) UNIQUE NOT NULL,
   count bigint
);
-- define hashtag length
------------------------
CREATE TABLE post_hashtag(
   post_id integer NOT NULL references post(post_id),
   hashtag_id integer NOT NULL references hashtag(hashtag_id)
 );
------------------------
CREATE TABLE post_tag_xuser (
    post_id bigint references post(post_id),
    user_id bigint references xuser(user_id),
    x integer, -- percentage 0-99
    y integer, -- percentage 0-99
    valid boolean, -- valid update by xuser
    createtime integer,
    updatetime integer
);
------------------------
INSERT INTO post 
(user_id, content, blob_id, type, point, country_id, category_id, public, createtime, updatetime) 
VALUES 
(4, 'hello world #happy @jeff', 'sha256 hashed id', 0, point('121.5643,25.0336'), 206, 0, true, 1527498044, 1527498044);
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
    post_id bigint references post(post_id),
    user_id bigint references xuser(user_id),
    act integer reference actions(action_id), -- 0: like, 1: dislike
    createtime integer
);
INSERT INTO post_actions 
(post_id, user_id, act, createtime) 
VALUES (33, 4, 0, 1527498711);
------------------------
CREATE TABLE comments (
    comment_id BIGSERIAL PRIMARY KEY NOT NULL,
    post_id bigint references post(post_id),
    user_id bigint references xuser(user_id),
    comment VARCHAR(300),
    createtime integer,
    updatetime integer
    -- constraint like_count (check (like_count >= 0))
);
------------------------
CREATE TABLE comment_actions (
    comment_id BIGSERIAL references comments(comment_id),
    user_id bigint references xuser(user_id),
    act integer reference actions(action_id),
    createtime integer
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
    comment_deep_id BIGSERIAL NOT NULL,
    comment_id bigint references comments(comment_id),
    user_id bigint references xuser(user_id),
    comment VARCHAR(300),
    like_count integer,
    dislike_count integer,
    createtime integer,
    updatetime integer
);
-- maximum taged user number: 5
------------------------
TABLE comments_deep_actions (
    comment_deep_id bigint references comments_deep(comment_deep_id),
    user_id bigint references xuser(user_id),
    act integer reference actions(action_id),
    createtime integer
);
------------------------
TABLE block (
    user_id BIGSERIAL references xuser(user_id),
    blocked_user_id BIGSERIAL references xuser(user_id),
    createtime integer
);
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
TABLE gender
------------------------
TABLE saved_post
------------------------
TABLE cities