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
db.createCollection("post_common");
db.post_common.createIndex({ "category_id": 1 });
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
                "name" : "test"
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

## popular_post_user_read_index

* init

```javascript
db.createCollection("popular_post_user_read_index");
db.popular_post_user_read_index.createIndex({ "user_id": 1 });
```

* sample data

```json
{
    "_id" : ObjectId("5b8d0e106ed3c50a23fd2e7d"),
    "user_id" : "3172",
    "city_sup_popular_post_index" : {
        "PRY_10_6" : 0,
        "ETH_8_7_7" : 0,
        "ITA_13_3_19" : 0,
        "PHL_36_21_91" : 0,
        "FRA_10_12_1_2" : 0,
        "EST_16_11_46" : 0,
        "KEN_26_6_2" : 0
    },
    "common_popular_post_index" : {
        "0" : 0,
        "7" : 0,
        "1" : 0,
        "10" : 0,
        "2" : 0,
        "4" : 0,
        "13" : 0,
        "8" : 0,
        "11" : 0,
        "15" : 0,
        "12" : 0,
        "6" : 0,
        "5" : 0,
        "14" : 0,
        "3" : 0,
        "9" : 0
    },
    "country_sup_popular_post_index" : {
        "BVT" : 0,
        "CCK" : 0,
        "MWI" : 0,
        "JOR" : 0,
        "LSO" : 0,
        "PAN" : 0
    }
}
```


## city_geometry

* init

```javascript
db.createCollection("city_geometry")
db.city_geometry.createIndex({"geometry":"2dsphere"})
```

* data source gadmpost_popular_common

please check python repo [xGeoCity](https://github.com/chienfuchen32/xGeoCity/releases/tag/v0.0.2)

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

## city info

* data source gadm

please check python repo [xGeoCity](https://github.com/chienfuchen32/xGeoCity/releases/tag/v0.0.3)

```javascript
db.createCollection("city")
db.city.createIndex({"level": 1, "country_code": 1, "city_id_1": 1, "city_id_2": 1, "city_id_3": 1, "city_id_4": 1, "city_id_5": 1})
```

* document sample
```json
{
    "_id" : ObjectId("5b891bc56ed3c50a23ebbe14"),
    "name" : "",
    "country_code" : "AFG",
    "level" : "2",
    "country_name" : "Afghanistan",
    "type" : "District",
    "city_id_5" : "",
    "city_id_4" : "",
    "city_id_3" : "",
    "city_id_2" : "AFG_1_11",
    "city_id_1" : "AFG_1"
}
```