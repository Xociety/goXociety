# root

use admin
db.createUser(
  {
    user: "",
    pwd: "",
    roles: [ "root" ]
  }
)

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

* data source gadm

please check python repo [xGeoCity](https://github.com/chienfuchen32/xGeoCity)

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
        "name_2" : "Taipei",
        "name_0" : "Taiwan",
        "name_1" : "Taipei",
        "varname_2" : "Táiběi Shì",
        "nl_name_1" : "台北",
        "nl_name_2" : "台北市",
        "hasc_2" : "TW.TP.TC",
        "type_2" : "Zhíxiáshì",
        "cc_2" : "",
        "gid_0" : "TWN",
        "gid_1" : "TWN.6_1",
        "gid_2" : "TWN.6.1_1",
        "engtype_2" : "Special Municipality"
    }
}
```

## city_level

* data source gadm

please check python repo [xGeoCity](https://github.com/chienfuchen32/xGeoCity)

```javascript
db.createCollection("city_level0")
db.city_level0.createIndex({"gid": 1})
db.createCollection("city_level1")
db.city_level1.createIndex({"gid": 1})
db.createCollection("city_level2")
db.city_level2.createIndex({"gid": 1})
db.createCollection("city_level3")
db.city_level3.createIndex({"gid": 1})
db.createCollection("city_level4")
db.city_level4.createIndex({"gid": 1})
db.createCollection("city_level5")
db.city_level5.createIndex({"gid": 1})
```