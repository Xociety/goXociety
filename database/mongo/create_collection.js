// post_popular
db.createCollection("post_popular");
db.post_popular.createIndex({ "user_id": 1, "category_id": 1 });
// post_read
db.createCollection("post_read");
db.post_read.createIndex({ "user_id": 1, "category_id": 1, "week_timestamp": 1 });
