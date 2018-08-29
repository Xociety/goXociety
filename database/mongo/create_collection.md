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
    "_id" : ObjectId("5b84c7b6a6d3157bc2d4425c"),
    "category_id" : 4,
    "popular_posts" : [ 
        {
            "post_id" : NumberLong(3877),
            "user" : {
                "user_id" : NumberLong(1974),
                "username" : "della",
                "name" : "Letitia Schowalter II",
                "photo_url" : ""
            },
            "content" : "test post",
            "blob" : {
                "blob_id" : "http://storage2.1mthechildbride.com/posts/images/sample/4/3/0.jpg",
                "origin_width" : 1242,
                "origin_height" : 2004
            },
            "type" : 0,
            "like_count" : NumberLong(45),
            "dislike_count" : NumberLong(47),
            "comment_count" : NumberLong(89),
            "category_id" : 4,
            "place" : {
                "place_id" : NumberLong(500),
                "country_code" : "TWN",
                "city_id_1" : "TWN.7_1",
                "city_id_2" : "TWN.7.3_1",
                "city_id_3" : "",
                "city_id_4" : "",
                "city_id_5" : "",
                "lat" : 23.5148062139715,
                "lon" : 120.513088892444,
                "name" : "test",
                "total_check_count" : NumberLong(0)
            },
            "public" : false,
            "createtime" : 1535422675,
            "updatetime" : 1535422675
        }
    ]
}
```

## post_read

* init

```javascript
db.createCollection("post_user_read");
db.post_user_read.createIndex({ "user_id": 1, "category_id": 1, "week_timestamp": 1 });
```

* document sample

```json
{
    "_id" : ObjectId("5b84c862a6d3157bc2d44f42"),
    "category_id" : 4,
    "user_id" : NumberLong(1909),
    "week_timestamp" : 2538,
    "popular_posts" : {
        "3877" : 1535428706,
        "4083" : 1535428782,
        "4312" : 1535428782,
        "4358" : 1535428782,
        "4527" : 1535428782
    }
}
```

## 

* init

```javascript
db.createCollection("post_user_read_index");
db.post_user_read_index.createIndex({ "user_id": 1, "category_id": 1 });
```

* document sample

```json
{
    "_id" : ObjectId("5b84c7b6a6d3157bc2d44195"),
    "category_id" : 4,
    "user_id" : NumberLong(1909),
    "popular_post" : {
        "index" : 5
    }
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