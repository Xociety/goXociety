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
    "_id" : ObjectId("5b7a6e4da6d3157bc2b49f3a"),
    "category_id" : 10,
    "posts" : [ 
        {
            "post_id" : NumberLong(4662),
            "user" : {
                "user_id" : NumberLong(1916),
                "username" : "deron.fritsch",
                "name" : "Elza Klocko",
                "photo_url" : ""
            },
            "content" : "test post",
            "blob" : {
                "blob_id" : "http://storage2.1mthechildbride.com/posts/images/sample/10/1/0.jpg",
                "origin_width" : 1242,
                "origin_height" : 2004
            },
            "type" : 0,
            "like_count" : NumberLong(57),
            "dislike_count" : NumberLong(41),
            "comment_count" : NumberLong(86),
            "category_id" : 10,
            "place" : {
                "placeid" : NumberLong(406),
                "countrycode" : "TWN",
                "cityid1" : "TWN.7_1",
                "cityid2" : "TWN.7.6_1",
                "cityid3" : "",
                "cityid4" : "",
                "cityid5" : "",
                "lat" : 24.2904120493827,
                "lon" : 121.722583824469,
                "name" : "test",
                "totalcheckcount" : NumberLong(0)
            },
            "public" : false,
            "createtime" : 1535424486,
            "updatetime" : 1535424486
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