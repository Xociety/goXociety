-- check all PRIMARY KEY, FOREIGN KEY, data valid range, Default, CRUD relation, related api, reference sequence, schema
-- alter table XXX DROP CONTAINT XXX;
-- alter table post ADD CONSTAINT post_type_fkey FOREIGN KEY (type) REFERENCES post_type(post_type_id) ON DELETE CASCADE ON UPDATE CASCADE;
------------------------
CREATE TABLE country (
    country_id          integer PRIMARY KEY NOT NULL,
    country             VARCHAR(200),
    country_code        VARCHAR(2)
);
------------------------
CREATE TABLE language (
    language_id         integer PRIMARY KEY NOT NULL,
    display_language    VARCHAR(200),
    value               VARCHAR(10)
);
------------------------
CREATE TABLE reaction (
    reaction_id         integer PRIMARY KEY NOT NULL,
    value               VARCHAR(15)
);
------------------------
CREATE TABLE gender (
    gender_id           integer PRIMARY KEY NOT NULL,
    value               VARCHAR(15)
);
------------------------
CREATE TABLE post_type (
    post_type_id        integer PRIMARY KEY NOT NULL,
    value               VARCHAR(10)
);
------------------------
CREATE TABLE category (
    category_id     integer PRIMARY KEY NOT NULL,
    category_name   VARCHAR(30)
);
------------------------
TABLE city
------------------------
CREATE TABLE xuser(
    user_id             BIGSERIAL PRIMARY KEY NOT NULL,
    username            VARCHAR(50) UNIQUE NOT NULL, -- [index]
    email               VARCHAR(50) UNIQUE, -- format check
    password            VARCHAR(100) NOT NULL, -- bcrypt format check
    name                VARCHAR(100) NOT NULL,
    phone               VARCHAR(30) NOT NULL, -- add country calling code, UNIQUE?
    gender              integer references gender(gender_id),
    bio                 VARCHAR(100),
    credit              double precision,
    photo_url           VARCHAR(200),
    language_id         integer references language(language_id),
    country_id          integer references country(country_id),
    timezone            integer,
    last_ip             VARCHAR(15), -- ipv4 123.194.188.0
    -- apid
    createtime          integer, -- unix time
    updatetime          integer
);
CREATE INDEX xuser_createtime ON xuser USING btree (createtime);
------------------------
CREATE TABLE follow (
    following_user_id   bigint references xuser(user_id) ON DELETE CASCADE ON UPDATE CASCADE, -- a person whom you follow
    follower_user_id    bigint references xuser(user_id) ON DELETE CASCADE ON UPDATE CASCADE, -- a person who follows you
    valid               boolean, -- update by following_user, block user usage?
    createtime          integer, -- create and delete by follower_user
    updatetime          integer, -- update by following_user
    CONSTRAINT follow_user_target UNIQUE (following_user_id, follower_user_id),
    CONSTRAINT follow_user_self_check CHECK (following_user_id <> follower_user_id) NOT VALID
);
CREATE INDEX follow_createtime ON follow USING btree (createtime);
-- following limit 7500 on Instagram, we might limit following 5000 people
------------------------
CREATE TABLE place (
    place_id    BIGSERIAL PRIMARY KEY NOT NULL,
    -- city_id
    position    point,
    name        VARCHAR(50),
    check_count bigint
    -- CONSTRAINT position_point_name_unique UNIQUE (position, name) [ERROR:  data type point has no default operator class for access method "btree", HINT:  You must specify an operator class for the index or define a default operator class for the data type.]
);
------------------------
CREATE TABLE post (
    post_id             BIGSERIAL PRIMARY KEY NOT NULL,
    user_id             bigint references xuser(user_id) ON DELETE CASCADE ON UPDATE CASCADE,
    content             VARCHAR(300), -- include @tagid, #hashtag
    blob_id             VARCHAR(200), -- may be foldername
    origin_width        int,
    origin_height       int,
    -- blob_count       int, -- multiple images, videos, preview image.. multiple size
    type                integer references post_type(post_type_id),
    like_count          bigint,
    dislike_count       bigint,
    comment_count       bigint,
    place_id            bigint references place(place_id),
    country_id          integer references country(country_id),
    category_id         integer, -- references
    public              boolean, -- public post in the future
    createtime          integer, -- timestamp without time zone, unix time
    updatetime          integer
    -- constraint like_count (check (like_count >= 0))
    -- constraint dislike_count (check (dislike_count >= 0))
    -- constraint comment_count (check (dislike_count >= 0))
);
--  maximum 30 hashtag
CREATE INDEX post_createtime ON post USING btree (createtime);
------------------------
CREATE TABLE hashtag(
    hashtag_id  BIGSERIAL PRIMARY KEY,
    value       VARCHAR(100) UNIQUE NOT NULL, -- lower case [index]
    count       bigint
);
CREATE INDEX hashtag_count ON hashtag USING btree (count);
-- define hashtag length, space check
------------------------
CREATE TABLE post_hashtag(
    hashtag_id           bigint NOT NULL references hashtag(hashtag_id),
    post_id              bigint NOT NULL references post(post_id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT post_hashtag_hashtag_post_unique UNIQUE (hashtag_id, post_id)
 );
------------------------
CREATE TABLE post_tag_xuser (
    post_id             bigint references post(post_id) ON DELETE CASCADE ON UPDATE CASCADE, -- [primary key?]
    user_id             bigint references xuser(user_id) ON DELETE CASCADE ON UPDATE CASCADE, -- [index?]
    x                   integer, -- percentage 0-99
    y                   integer, -- percentage 0-99
    valid               boolean, -- valid update by xuser
    createtime          integer,
    updatetime          integer
);
CREATE INDEX post_tag_xuser_createtime ON post_tag_xuser USING btree (createtime);
------------------------
CREATE TABLE post_reaction (
    post_id             bigint references post(post_id) ON DELETE CASCADE ON UPDATE CASCADE, -- [primary key]
    user_id             bigint references xuser(user_id) ON DELETE CASCADE ON UPDATE CASCADE,
    reaction_id         integer references reaction(reaction_id), -- 0: like, 1: dislike
    createtime          integer,
    CONSTRAINT post_reaction_post_user_unique UNIQUE (post_id, user_id)
);
CREATE INDEX post_reaction_createtime ON post_reaction USING btree (createtime);
------------------------
CREATE TABLE comment (
    comment_id          BIGSERIAL PRIMARY KEY NOT NULL,
    post_id             bigint references post(post_id) ON DELETE CASCADE ON UPDATE CASCADE,
    user_id             bigint references xuser(user_id) ON DELETE CASCADE ON UPDATE CASCADE,
    comment             VARCHAR(300),
    like_count          bigint,
    dislike_count       bigint,
    reply_count         bigint,
    createtime          integer,
    updatetime          integer
    -- constraint like_count (check (like_count >= 0))
    -- constraint dislike_count (check (dislike_count >= 0))
    -- constraint reply_count (check (reply_count >= 0))
);
CREATE INDEX comment_createtime ON comment USING btree (createtime);
------------------------
CREATE TABLE comment_reaction (
    comment_id          bigint references comment(comment_id), -- [PRIMARY KEY]
    user_id             bigint references xuser(user_id) ON DELETE CASCADE ON UPDATE CASCADE,
    reaction_id         integer references reaction(reaction_id),
    createtime          integer,
    CONSTRAINT comment_reaction_comment_user_unique UNIQUE (comment_id, user_id)
);
CREATE INDEX comment_reaction_createtime ON comment_reaction USING btree (createtime);
------------------------
CREATE TABLE reply (
    reply_id            BIGSERIAL PRIMARY KEY NOT NULL,
    comment_id          bigint references comment(comment_id) ON DELETE CASCADE ON UPDATE CASCADE,
    user_id             bigint references xuser(user_id) ON DELETE CASCADE ON UPDATE CASCADE,
    reply               VARCHAR(300),
    like_count          bigint,
    dislike_count       bigint,
    createtime          integer,
    updatetime          integer
    -- constraint like_count (check (like_count >= 0))
    -- constraint dislike_count (check (dislike_count >= 0))
);
CREATE INDEX reply_createtime ON reply USING btree (createtime);
------------------------
CREATE TABLE reply_reaction (
    reply_id            bigint references reply(reply_id) ON DELETE CASCADE ON UPDATE CASCADE, -- [PRIMARY KEY]
    user_id             bigint references xuser(user_id) ON DELETE CASCADE ON UPDATE CASCADE,
    reaction_id         integer references reaction(reaction_id),
    createtime          integer,
    CONSTRAINT reply_reaction_reply_user_unique UNIQUE (reply_id, user_id)
);
CREATE INDEX reply_reaction_createtime ON reply_reaction USING btree (createtime);
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
TABLE saved_post(
    user_id
    post_id
    createtime
);
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