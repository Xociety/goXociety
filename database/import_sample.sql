COPY country(country_id, country, country_code) FROM '/var/lib/postgresql/data/pgdata/xsrc/country.csv' DELIMITER ';' CSV; -- HEADER;
/*
https://developers.google.com/custom-search/docs/xml_results_appendices#countryCodes
*/
COPY language(language_id, display_language, value) FROM '/var/lib/postgresql/data/pgdata/xsrc/language.csv' DELIMITER ';' CSV;
/*
https://developers.google.com/custom-search/docs/xml_results_appendices#interfaceLanguages
*/
------------------------
COPY reaction(reaction_id, value) FROM '/var/lib/postgresql/data/pgdata/xsrc/reaction.csv' DELIMITER ';' CSV;
------------------------
COPY gender(gender_id, value) FROM '/var/lib/postgresql/data/pgdata/xsrc/gender.csv' DELIMITER ';' CSV;
------------------------
COPY post_type(post_type_id, value) FROM '/var/lib/postgresql/data/pgdata/xsrc/post_type.csv' DELIMITER ';' CSV; -- CSV HEADER;
------------------------
COPY category(category_id, category_name) FROM '/var/lib/postgresql/data/pgdata/xsrc/category.csv' DELIMITER ';' CSV; -- CSV HEADER;
------------------------
INSERT INTO xuser 
(username, email, password, name, phone, gender, bio, credit, photo_url, language_id, country_id, timezone, last_ip, createtime, updatetime) 
VALUES ('jeff', 'jeff@gmail.com', 'salted', 'jeff', '+886-911111111', 1, 'hi', 0, '', 12, 206, 28800, '123.194.188.0', 1527496777, 1527496777),
('kyler', 'kyler@gmail.com', 'salted', 'kyler', '+886-911111111', 1, 'yo', 0, '', 12, 206, 28800, '123.194.188.0', 1527496777, 1527496777),
('robby', 'robby@gmail.com', 'salted', 'robby', '+886-911111111', 1, 'man', 0, '', 12, 206, 28800, '123.194.188.0', 1527496777, 1527496777);
-- username, phone + country code, email (lower case) logic check
------------------------
INSERT INTO post 
(user_id, content, blob_id, type, like_count, dislike_count, point, country_id, category_id, public, createtime, updatetime) 
VALUES 
(4, 'hello world #happy @jeff', 'sha256 hashed id', 0, 0, 0, point('121.5643,25.0336'), 206, 0, true, 1527498044, 1527498044);
-- UPDATE post SET content='hello world #happy @jeff', blob_id='sha256 hashed id' WHERE post_id = 1;
------------------------
INSERT INTO hashtag 
(name) 
VALUES ('happy');
------------------------
INSERT INTO post_hashtag 
(post_id, hashtag_id) 
VALUES (33, 2);
------------------------
INSERT INTO post_tag_xuser 
(post_id, user_id, x, y, valid, createtime, updatetime) 
VALUES (33, 4, 0, 0, false, 1527498711, 1527498711);
------------------------
INSERT INTO post_reaction 
(post_id, user_id, reaction_id, createtime) 
VALUES (33, 4, 0, 1527498711);
------------------------
INSERT INTO comment 
(post_id, user_id, comment, createtime, updatetime) 
VALUES (33, 4, 'yo', 1527498711, 1527498711);
------------------------
INSERT INTO comment_reaction 
(comment_id, user_id, reaction_id, createtime) 
VALUES (2, 4, 0, 1527498711);
