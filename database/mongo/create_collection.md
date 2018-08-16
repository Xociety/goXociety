# collection

`use xociety`

## post_popular

* init


```javascript
db.createCollection("post_popular");
db.post_popular.createIndex({ "user_id": 1, "category_id": 1 });
```

* document sample

```json
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
```

## post_read

* init

```javascript
db.createCollection("post_read");
db.post_read.createIndex({ "user_id": 1, "category_id": 1, "week_timestamp": 1 });
```

* document sample

```json
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
```

## city

* init

```javascript
db.createCollection("city")
db.city.createIndex({"geometry":"2dsphere"})
```

* document sample

```json
{
    "_id" : ObjectId("5b73e8fcd5dd9938e90580e1"),
    "geometry" : {
        "type" : "Polygon",
        "coordinates" : [ 
            [ 
                [ 
                    121.413475036621, 
                    25.1676387786866
                ], 
                [ 
                    121.413475036621, 
                    25.167917251587
                ], 
                [ 
                    121.413475036621, 
                    25.1676387786866
                ]
            ]
        ]
    },
    "type" : "Feature",
    "properties" : {
        "NAME_2" : "Taipei",
        "NAME_0" : "Taiwan",
        "NAME_1" : "Taipei",
        "VARNAME_2" : "Táiběi Shì",
        "NL_NAME_1" : "台北",
        "NL_NAME_2" : "台北市",
        "HASC_2" : "TW.TP.TC",
        "TYPE_2" : "Zhíxiáshì",
        "CC_2" : "",
        "GID_0" : "TWN",
        "GID_1" : "TWN.6_1",
        "GID_2" : "TWN.6.1_1",
        "ENGTYPE_2" : "Special Municipality"
    }
}
```

## city_level

```javascript
db.createCollection("city_level0")
db.city_level0.createIndex({"GID": 1})
db.createCollection("city_level1")
db.city_level1.createIndex({"GID": 1})
db.createCollection("city_level2")
db.city_level2.createIndex({"GID": 1})
db.createCollection("city_level3")
db.city_level3.createIndex({"GID": 1})
db.createCollection("city_level4")
db.city_level4.createIndex({"GID": 1})
db.createCollection("city_level5")
db.city_level5.createIndex({"GID": 1})
```