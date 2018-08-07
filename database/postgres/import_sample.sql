COPY country(country_id, country, country_code) FROM '/var/lib/postgresql/data/pgdata/xsrc/country.csv' DELIMITER ';' CSV; -- HEADER;
/*
data source: https://developers.google.com/custom-search/docs/xml_results_appendices#countryCodes
updated id by ./id_generator.sh
*/
------------------------
COPY language(language_id, display_language, value) FROM '/var/lib/postgresql/data/pgdata/xsrc/language.csv' DELIMITER ';' CSV;
/*
data source: https://developers.google.com/custom-search/docs/xml_results_appendices#interfaceLanguages
updated id by ./id_generator.sh
*/
------------------------
COPY reaction(reaction_id, value) FROM '/var/lib/postgresql/data/pgdata/xsrc/reaction.csv' DELIMITER ';' CSV;
------------------------
COPY gender(gender_id, value) FROM '/var/lib/postgresql/data/pgdata/xsrc/gender.csv' DELIMITER ';' CSV;
-- data source: https://en.wikipedia.org/wiki/ISO/IEC_5218
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
COPY city(category_id, category_name) FROM '/var/lib/postgresql/data/pgdata/xsrc/city.csv' DELIMITER ';' CSV; -- CSV HEADER;
/*
data source: https://dev.maxmind.com/geoip/geoip2/geolite2/, http://geolite.maxmind.com/download/geoip/database/GeoLite2-City-CSV.zip, data updated
./city_origin.csv modified data as ./city.csv
*/