# auth

```javascript
use admin
db.createUser(
  {
    user: "",
    pwd: "",
    roles: [ "root" ]
  }
)
```

# collection

`use xociety`

## post_popular

* init

```javascript
db.createCollection("post_popular_common");
db.post_popular_common.createIndex({ "category_id": 1 });
```

* document sample

```json
{
    "_id" : ObjectId("5b7a6e4da6d3157bc2b49a5e"),
    "category_id" : 4,
    "posts" : [ 
        {
            "post_id" : NumberLong(2564),
            "user" : {
                "user_id" : NumberLong(925),
                "username" : "stewart",
                "name" : "Enola Wintheiser",
                "photo_url" : ""
            },
            "content" : "test post",
            "blob" : {
                "blob_id" : "http://storage2.1mthechildbride.com/posts/images/sample/4/6/0.jpg",
                "origin_width" : 1242,
                "origin_height" : 2004
            },
            "type" : 0,
            "like_count" : NumberLong(49),
            "dislike_count" : NumberLong(44),
            "comment_count" : NumberLong(93),
            "country_id" : 0,
            "category_id" : 4,
            "public" : false,
            "createtime" : 1534826851,
            "updatetime" : 1534826851
        }, 
        {
            "post_id" : NumberLong(2319),
            "user" : {
                "user_id" : NumberLong(981),
                "username" : "gilda",
                "name" : "Helene Jerde",
                "photo_url" : ""
            },
            "content" : "test post",
            "blob" : {
                "blob_id" : "http://storage2.1mthechildbride.com/posts/images/sample/4/3/0.jpg",
                "origin_width" : 1242,
                "origin_height" : 2004
            },
            "type" : 0,
            "like_count" : NumberLong(47),
            "dislike_count" : NumberLong(46),
            "comment_count" : NumberLong(77),
            "country_id" : 0,
            "category_id" : 4,
            "public" : false,
            "createtime" : 1534826220,
            "updatetime" : 1534826220
        }
    ]
}
```

## post_read

* init

```javascript
db.createCollection("post_popular_read");
db.post_popular_read.createIndex({ "user_id": 1, "category_id": 1, "week_timestamp": 1 });
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

## 

* init

```javascript
db.createCollection("post_popular_read_index");
db.post_popular_index.createIndex({ "user_id": 1, "category_id": 1 });
```

* document sample

```json
{
    "_id" : ObjectId("5b7a8a50a6d3157bc2b58f1c"),
    "category_id" : 1,
    "user_id" : NumberLong(1),
    "index" : 0
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
        "country_code" : "TWN",
        "city_id_1" : "TWN.6_1",
        "city_id_2" : "TWN.6.1_1",
        "engtype_2" : "Special Municipality"
    }
}
```

## city_level

* data source gadm

please check python repo [xGeoCity](https://github.com/chienfuchen32/xGeoCity)

```javascript
db.createCollection("country")
db.country.createIndex({"country_code": 1})
db.createCollection("city_level_1")
db.city_level_1.createIndex({"city_id": 1})
db.createCollection("city_level_2")
db.city_level_2.createIndex({"city_id": 1})
db.createCollection("city_level_3")
db.city_level_3.createIndex({"city_id": 1})
db.createCollection("city_level_4")
db.city_level_4.createIndex({"city_id": 1})
db.createCollection("city_level_5")
db.city_level_5.createIndex({"city_id": 1})
```

* document sample
```json
{
    "_id" : ObjectId("5b7d0d11a6d3157bc2d2eed8"),
    "country_name" : "Taiwan",
    "country_code" : "TWN"
}
```
```json
{
    "_id" : ObjectId("5b7b6e94a6d3157bc2c080c6"),
    "city_id" : "TWN.6.1_1",
    "type" : "Special Municipality",
    "name" : "Taipei"
}
```