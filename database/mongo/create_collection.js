// post_popular
db.createCollection("post_popular");
db.post_popular.createIndex({ "user_id": 1, "category_id": 1 });
// sample
/*
{ 
    "_id" : ObjectId("5b4c12bc4157b267656284c8"), 
    "category_id" : NumberInt(0), 
    "user_id" : NumberLong(4), 
    "week_timestamp" : NumberInt(2532), 
    "posts" : {
        "2858" : NumberInt(1531712188), 
        "2882" : NumberInt(1531712188), 
        "2914" : NumberInt(1531712207)
    }
}
*/

// post_read
db.createCollection("post_read");
db.post_read.createIndex({ "user_id": 1, "category_id": 1, "week_timestamp": 1 });
// sample
/*
{ 
    "_id" : ObjectId("5b46ffdf4157b2676561d91f"), 
    "category_id" : NumberInt(0), 
    "user_id" : NumberLong(4), 
    "posts" : [
        {
            "post_id" : NumberLong(2920), 
            "user" : {
                "user_id" : NumberLong(316), 
                "username" : "art", 
                "name" : "Aubrey Roob"
            }, 
            "content" : "test post", 
            "blob" : {
                "blob_id" : "http://storage.1mthechildbride.com/videos/sample/0.m3u8", 
                "origin_width" : NumberInt(1920), 
                "origin_height" : NumberInt(1080)
            }, 
            "type" : NumberInt(1), 
            "like_count" : NumberLong(3), 
            "dislike_count" : NumberLong(3), 
            "comment_count" : NumberLong(19), 
            "country_id" : NumberInt(0), 
            "category_id" : NumberInt(0), 
            "public" : false, 
            "createtime" : NumberInt(1531205242), 
            "updatetime" : NumberInt(1531205242)
        }, 
        {
            "post_id" : NumberLong(2942), 
            "user" : {
                "user_id" : NumberLong(311), 
                "username" : "eden", 
                "name" : "Mrs. Katheryn Skiles"
            }, 
            "content" : "test post", 
            "blob" : {
                "blob_id" : "http://storage.1mthechildbride.com/videos/sample/0.m3u8", 
                "origin_width" : NumberInt(1920), 
                "origin_height" : NumberInt(1080)
            }, 
            "type" : NumberInt(1), 
            "like_count" : NumberLong(5), 
            "dislike_count" : NumberLong(7), 
            "comment_count" : NumberLong(12), 
            "country_id" : NumberInt(0), 
            "category_id" : NumberInt(0), 
            "public" : false, 
            "createtime" : NumberInt(1531205257), 
            "updatetime" : NumberInt(1531205257)
        }
    ]
}
*/