-- check all PRIMARY KEY, FOREIGN KEY, data valid range, CRUD relation, related api, reference sequence, schema
------------------------
CREATE TABLE country (
    country_id int PRIMARY KEY     NOT NULL,
    name CHAR(200),
    code CHAR(2)
);
COPY country(country_id, name, code) FROM '/var/lib/postgresql/data/country.csv' DELIMITER ';' CSV HEADER;
/*
https://developers.google.com/custom-search/docs/xml_results_appendices#countryCodes
*/
------------------------
CREATE TABLE language (
    language_id int PRIMARY KEY     NOT NULL,
    display_language CHAR(200),
    value CHAR(10)
);
COPY language(language_id, display_language, value) FROM '/var/lib/postgresql/data/language.csv' DELIMITER ';' CSV HEADER;
/*
https://developers.google.com/custom-search/docs/xml_results_appendices#interfaceLanguages
*/
------------------------
CREATE TABLE xuser(
    user_id BIGSERIAL    PRIMARY KEY NOT NULL,
    username        CHAR(50) UNIQUE NOT NULL,
    email           CHAR(50), -- unique and format check
    password VARCHAR(100) NOT NULL,
    name            TEXT    NOT NULL,
    phone           CHAR(50) NOT NULL, -- add country calling code
    gender          INT, -- 0: not known, 1: male, 2:female
    bio             CHAR(50),
    credit          INT,
    language_id        int references language(langauge_id),
    country_id         int references country(country_id),
    timezone        INT,
    last_ip         CHAR(15), -- ipv4 123.194.188.0
    -- apid
    createtime      integer, -- unix time
    updatetime      integer
);
INSERT INTO xuser 
(username, email, password, name, phone, gender, bio, credit, language_id, country_id, timezone, last_ip, createtime, updatetime) 
VALUES ('jeff', 'jeff@gmail.com', 'salted', 'jeff', '+886-911111111', 1, 'hi', 0, 13, 207, 28800, '123.194.188.0', 1527496777, 1527496777),
('kyler', 'kyler@gmail.com', 'salted', 'kyler', '+886-911111111', 1, 'yo', 0, 13, 207, 28800, '123.194.188.0', 1527496777, 1527496777),
('robby', 'robby@gmail.com', 'salted', 'robby', '+886-911111111', 1, 'man', 0, 13, 207, 28800, '123.194.188.0', 1527496777, 1527496777);
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
CREATE TABLE post (
    post_id BIGSERIAL PRIMARY KEY     NOT NULL,
    user_id BIGSERIAL references xuser(user_id),
    content VARCHAR(300), -- include @tagid, #hashtag
    blob_id VARCHAR(100), -- may be foldername
    type int references post_type(post_type_id),
    point point,
    country_id int references country(country_id),
    category_id INT,
    public boolean,
    createtime integer, -- timestamp without time zone, unix time
    updatetime integer
);
--  maximum 30 hashtag
------------------------
CREATE TABLE post_type (
    post_type_id int PRIMARY KEY NOT NULL,
    name CHAR(10)
);
COPY country(country_id, name, code) FROM '/var/lib/postgresql/data/country.csv' DELIMITER ';' CSV HEADER;
-- INSERT INTO post_type (post_type_id, name) VALUES (1, 'jpg'), (2, 'hls');
------------------------
CREATE TABLE hashtag(
   hashtag_id BIGSERIAL PRIMARY KEY,
   name TEXT UNIQUE NOT NULL,
   count bigint
);
-- define hashtag length
------------------------
CREATE TABLE post_hashtag(
   post_id INT NOT NULL references post(post_id),
   hashtag_id INT NOT NULL references hashtag(hashtag_id)
 );
------------------------
CREATE TABLE post_tag_xuser (
    post_id BIGSERIAL references post(post_id),
    user_id SERIAL references xuser(user_id),
    x int, -- percentage 0-99
    y int, -- percentage 0-99
    valid boolean, -- valid update by xuser
    createtime integer,
    updatetime integer
);
------------------------
INSERT INTO post 
(user_id, content, blob_id, point, country_id, category_id, public, createtime, updatetime) 
VALUES 
(1, 'hello world #happy @jeff', 'sha256 hashed id', point('121.5643,25.0336'), 207, 0, true, 1527498044, 1527498044);
-- UPDATE post SET content='hello world #happy @jeff', blob_id='sha256 hashed id' WHERE post_id = 1;
INSERT INTO hashtag 
(name) 
VALUES ('happy');
INSERT INTO post_hashtag 
(post_id, hashtag_id) 
VALUES (1, 1);
INSERT INTO post_tag_xuser 
(post_id, user_id, x, y, valid, createtime, updatetime) 
VALUES (1, 1, 0, 0, false, 1527498711, 1527498711);
------------------------
CREATE TABLE post_likes (
    post_id SERIAL references post(post_id),
    user_id SERIAL references xuser(user_id),
    type int, -- 1: like, 2: dislike
    createtime integer
);
INSERT INTO post_likes 
(post_id, user_id, type, createtime) 
VALUES (1, 1, 1, 1527498711);
------------------------
CREATE TABLE comments (
    comment_id BIGSERIAL PRIMARY KEY NOT NULL,
    post_id SERIAL references post(post_id),
    user_id SERIAL references xuser(user_id),
    comment VARCHAR(300),
    createtime integer,
    updatetime integer
);
------------------------
CREATE TABLE comment_likes (
    comment_id BIGSERIAL references comments(comment_id),
    user_id SERIAL references xuser(user_id),
    type int, -- 1: like, 2: dislike
    createtime integer
);
INSERT INTO comments 
(post_id, user_id, comment, createtime, updatetime) 
VALUES (1, 1, 'yo', 1527498711, 1527498711);
-- maximum taged user number: 5
INSERT INTO comment_likes 
(comment_id, user_id, type, createtime) 
VALUES (1, 1, 1, 1527498711);

SELECT * FROM post 
JOIN xuser on post.user_id = xuser.user_id 
JOIN post_likes on post.post_id = post_likes.post_id 
JOIN comments on post.post_id = comments.post_id 
JOIN comment_likes on comments.comment_id = comment_likes.comment_id;

------------------------
CREATE TABLE comments_deep (
    comment_deep_id BIGSERIAL NOT NULL,
    comment_id BIGSERIAL references comments(comment_id),
    user_id SERIAL references xuser(user_id),
    comment VARCHAR(300),
    createtime integer,
    updatetime integer
);
-- maximum taged user number: 5
------------------------
CREATE TABLE comments_deep_likes (
    comment_deep_id BIGSERIAL references comments_deep(comment_deep_id),
    user_id SERIAL references xuser(user_id),
    type int, -- 1: like, 2: dislike
    createtime integer
);
------------------------
CREATE TABLE follow (
    following_user_id SERIAL references xuser(user_id),
    followed_user_id SERIAL references xuser(user_id),
    valid boolean,
    createtime integer
);
-- following limit 7500 on Instagram, we might limit following 5000 people
------------------------
CREATE TABLE block (
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
)
------------------------
TABLE category
------------------------
TABLE sub_category
------------------------
TABLE gender
------------------------
TABLE feeling
------------------------
