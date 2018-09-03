package main

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"math"
	"strconv"
	"time"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"
)

// scalar
func coerceInt64(value interface{}) interface{} {
	switch value := value.(type) {
	case bool:
		if value == true {
			return 1
		}
		return 0
	case int:
		if value < int(math.MinInt32) || value > int(math.MaxInt32) {
			return nil
		}
		return value
	case *int:
		return coerceInt64(*value)
	case int8:
		return int(value)
	case *int8:
		return int(*value)
	case int16:
		return int(value)
	case *int16:
		return int(*value)
	case int32:
		return int(value)
	case *int32:
		return int(*value)
	case int64:
		if value < int64(math.MinInt64) || value > int64(math.MaxInt64) {
			return nil
		}
		return int64(value)
	case *int64:
		return coerceInt64(*value)
	case uint:
		if value > math.MaxInt32 {
			return nil
		}
		return int(value)
	case *uint:
		return coerceInt64(*value)
	case uint8:
		return int(value)
	case *uint8:
		return int(*value)
	case uint16:
		return int(value)
	case *uint16:
		return int(*value)
	case uint32:
		if value > uint32(math.MaxInt32) {
			return nil
		}
		return int(value)
	case *uint32:
		return coerceInt64(*value)
	case uint64:
		if value > uint64(math.MaxInt64) {
			return nil
		}
		return int(value)
	case *uint64:
		return coerceInt64(*value)
	case float32:
		if value < float32(math.MinInt32) || value > float32(math.MaxInt32) {
			return nil
		}
		return int(value)
	case *float32:
		return coerceInt64(*value)
	case float64:
		if value < float64(math.MinInt32) || value > float64(math.MaxInt32) {
			return nil
		}
		return int(value)
	case *float64:
		return coerceInt64(*value)
	case string:
		val, err := strconv.ParseFloat(value, 0)
		if err != nil {
			return nil
		}
		return coerceInt64(val)
	case *string:
		return coerceInt64(*value)
	}

	// If the value cannot be transformed into an int, return nil instead of '0'
	// to denote 'no integer found'
	return nil
}

func coerceJSON(value interface{}) interface{} {
	switch value := value.(type) {
	case string:
		val := string(value)
		return coerceJSON(val)
	case *string:
		return coerceJSON(*value)
	}
	return nil
}

var int64GraphqlScalar = graphql.NewScalar(graphql.ScalarConfig{
	Name:        "Int64",
	Description: "int64",
	Serialize:   coerceInt64,
	ParseValue:  coerceInt64,
	ParseLiteral: func(valueAST ast.Value) interface{} {
		switch valueAST := valueAST.(type) {
		case *ast.IntValue:
			if intValue, err := strconv.ParseInt(valueAST.Value, 10, 64); err == nil {
				return intValue
			}
		}
		return nil
	},
})

var tagOnPostGraphqlScalar = graphql.NewScalar(graphql.ScalarConfig{
	Name:        "post_tag_str",
	Description: `string input sample: post tags => "{\\"user_id\\":1,\\"x\\":40,\\"y\\":50}"`,
	Serialize:   coerceJSON,
	ParseValue:  coerceJSON,
	ParseLiteral: func(valueAST ast.Value) interface{} {
		switch valueAST := valueAST.(type) {
		case *ast.StringValue:
			objectValue := tagOnPostSetAPI{}
			if err := json.Unmarshal([]byte(valueAST.Value), &objectValue); err == nil {
				return objectValue
			}
		}
		return nil
	},
})
var tagsOnPostGraphqlScalar = graphql.NewScalar(graphql.ScalarConfig{
	Name:        "post_tags_str",
	Description: `string input sample: post tags => "[{\\"user_id\\":1,\\"x\\":40,\\"y\\":50}]", tag no user => "[]"`,
	Serialize:   coerceJSON,
	ParseValue:  coerceJSON,
	ParseLiteral: func(valueAST ast.Value) interface{} {
		switch valueAST := valueAST.(type) {
		case *ast.StringValue:
			objectValue := []tagOnPostSetAPI{}
			if err := json.Unmarshal([]byte(valueAST.Value), &objectValue); err == nil {
				return objectValue
			}
		}
		return nil
	},
})

// graphql type

// common
var updateStatusGraphqlType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "update_status",
		Fields: graphql.Fields{
			"rows_affected": &graphql.Field{Type: graphql.Int},
		},
	},
)

var countryBasicGraphqlType = graphql.NewList(
	graphql.NewObject(
		graphql.ObjectConfig{
			Name: "country_basic",
			Fields: graphql.Fields{
				"country_code": &graphql.Field{Type: graphql.String},
				"country_name": &graphql.Field{Type: graphql.String},
			},
		},
	),
)
var languageGraphqlType = graphql.NewList(
	graphql.NewObject(
		graphql.ObjectConfig{
			Name: "language",
			Fields: graphql.Fields{
				"display_language": &graphql.Field{Type: graphql.String},
				"value":            &graphql.Field{Type: graphql.String},
			},
		},
	),
)
var genderGraphqlType = graphql.NewList(
	graphql.NewObject(
		graphql.ObjectConfig{
			Name: "gender",
			Fields: graphql.Fields{
				"gender_id": &graphql.Field{Type: graphql.String},
				"value":     &graphql.Field{Type: graphql.String},
			},
		},
	),
)
var reactionGraphqlType = graphql.NewList(
	graphql.NewObject(
		graphql.ObjectConfig{
			Name: "reaction",
			Fields: graphql.Fields{
				"reaction_id": &graphql.Field{Type: graphql.Int},
				"value":       &graphql.Field{Type: graphql.String},
			},
		},
	),
)
var postTypeGraphqlType = graphql.NewList(
	graphql.NewObject(
		graphql.ObjectConfig{
			Name: "post_type",
			Fields: graphql.Fields{
				"post_type_id": &graphql.Field{Type: graphql.Int},
				"value":        &graphql.Field{Type: graphql.String},
				"file_format":  &graphql.Field{Type: graphql.NewList(graphql.String)},
			},
		},
	),
)
var categoryGraphqlType = graphql.NewList(
	graphql.NewObject(
		graphql.ObjectConfig{
			Name: "category",
			Fields: graphql.Fields{
				"category_id":   &graphql.Field{Type: graphql.Int},
				"category_name": &graphql.Field{Type: graphql.String},
			},
		},
	),
)

// login
var loginGraphqlType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "login",
		Fields: graphql.Fields{
			"token": &graphql.Field{Type: graphql.String},
		},
	},
)

// user
var userGraphqlType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "User",
		Fields: graphql.Fields{
			"user_id":      &graphql.Field{Type: int64GraphqlScalar},
			"username":     &graphql.Field{Type: graphql.String},
			"email":        &graphql.Field{Type: graphql.String},
			"name":         &graphql.Field{Type: graphql.String},
			"phone":        &graphql.Field{Type: graphql.String},
			"gender":       &graphql.Field{Type: graphql.Int},
			"bio":          &graphql.Field{Type: graphql.String},
			"credit":       &graphql.Field{Type: graphql.Int},
			"photo_url":    &graphql.Field{Type: graphql.String},
			"language_id":  &graphql.Field{Type: graphql.Int},
			"country_code": &graphql.Field{Type: graphql.String},
			"timezone":     &graphql.Field{Type: graphql.Int},
			"last_ip":      &graphql.Field{Type: graphql.String},
			"createtime":   &graphql.Field{Type: graphql.Int},
			"updatetime":   &graphql.Field{Type: graphql.Int},
		},
	},
)

// following
var userBasicGraphqlType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "user_basic",
		Fields: graphql.Fields{
			"user_id":   &graphql.Field{Type: int64GraphqlScalar},
			"username":  &graphql.Field{Type: graphql.String},
			"name":      &graphql.Field{Type: graphql.String},
			"photo_url": &graphql.Field{Type: graphql.String},
		},
	},
)
var usersFollowingGraphqlType = graphql.NewList(
	graphql.NewObject(
		graphql.ObjectConfig{
			Name: "user_following",
			Fields: graphql.Fields{
				"user_id":        &graphql.Field{Type: int64GraphqlScalar},
				"username":       &graphql.Field{Type: graphql.String},
				"name":           &graphql.Field{Type: graphql.String},
				"photo_url":      &graphql.Field{Type: graphql.String},
				"following_time": &graphql.Field{Type: graphql.Int},
			},
		},
	),
)

// follower
var usersFollowerGraphqlType = graphql.NewList(
	graphql.NewObject(
		graphql.ObjectConfig{
			Name: "user_follower",
			Fields: graphql.Fields{
				"user_id":        &graphql.Field{Type: int64GraphqlScalar},
				"username":       &graphql.Field{Type: graphql.String},
				"name":           &graphql.Field{Type: graphql.String},
				"photo_url":      &graphql.Field{Type: graphql.String},
				"following_time": &graphql.Field{Type: graphql.Int},
			},
		},
	),
)

// city
var cityGeometryGraphqlType = graphql.NewList(
	graphql.NewObject(
		graphql.ObjectConfig{
			Name: "city_geometry",
			Fields: graphql.Fields{
				"properties": &graphql.Field{Type: cityGeometryPropertiesGraphqlType},
			},
		},
	),
)
var cityGeometryPropertiesGraphqlType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "city_geometry_properties",
		Fields: graphql.Fields{
			"country_code": &graphql.Field{Type: graphql.String},
			"city_id_1":    &graphql.Field{Type: graphql.String},
			"city_id_2":    &graphql.Field{Type: graphql.String},
			"city_id_3":    &graphql.Field{Type: graphql.String},
			"city_id_4":    &graphql.Field{Type: graphql.String},
			"city_id_5":    &graphql.Field{Type: graphql.String},
			"country_name": &graphql.Field{Type: graphql.String},
			"name_1":       &graphql.Field{Type: graphql.String},
			"name_2":       &graphql.Field{Type: graphql.String},
			"name_3":       &graphql.Field{Type: graphql.String},
			"name_4":       &graphql.Field{Type: graphql.String},
			"name_5":       &graphql.Field{Type: graphql.String},
			"type_1":       &graphql.Field{Type: graphql.String},
			"type_2":       &graphql.Field{Type: graphql.String},
			"type_3":       &graphql.Field{Type: graphql.String},
			"type_4":       &graphql.Field{Type: graphql.String},
			"type_5":       &graphql.Field{Type: graphql.String},
		},
	},
)
var cityGraphqlType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "city2",
		Fields: graphql.Fields{
			"level":             &graphql.Field{Type: graphql.String},
			"country_code":      &graphql.Field{Type: graphql.String},
			"country_name":      &graphql.Field{Type: graphql.String},
			"city_id_1":         &graphql.Field{Type: graphql.String},
			"city_id_2":         &graphql.Field{Type: graphql.String},
			"city_id_3":         &graphql.Field{Type: graphql.String},
			"city_id_4":         &graphql.Field{Type: graphql.String},
			"city_id_5":         &graphql.Field{Type: graphql.String},
			"name":              &graphql.Field{Type: graphql.String},
			"type":              &graphql.Field{Type: graphql.String},
			"post_count":        &graphql.Field{Type: int64GraphqlScalar},
			"sup_popular_posts": &graphql.Field{Type: postsGraphqlType},
		},
	},
)

// blob
var blobGraphqlType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "blob",
		Fields: graphql.Fields{
			"blob_id":       &graphql.Field{Type: graphql.String},
			"origin_width":  &graphql.Field{Type: graphql.Int},
			"origin_height": &graphql.Field{Type: graphql.Int},
		},
	},
)

// place
var placesLookupGraphql = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "places_lookup",
		Fields: graphql.Fields{
			"place":           &graphql.Field{Type: placesGraphqlType},
			"next_page_token": &graphql.Field{Type: graphql.String},
		},
		Description: "if place_id equals 0, it means that it's not existed in our db",
	},
)
var placeGraphqlType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "place",
		Fields: graphql.Fields{
			"place_id":     &graphql.Field{Type: int64GraphqlScalar},
			"country_code": &graphql.Field{Type: graphql.String},
			"city_id_1":    &graphql.Field{Type: graphql.String},
			"city_id_2":    &graphql.Field{Type: graphql.String},
			"city_id_3":    &graphql.Field{Type: graphql.String},
			"city_id_4":    &graphql.Field{Type: graphql.String},
			"city_id_5":    &graphql.Field{Type: graphql.String},
			"lat":          &graphql.Field{Type: graphql.Float},
			"lon":          &graphql.Field{Type: graphql.Float},
			"name":         &graphql.Field{Type: graphql.String},
		},
	},
)
var placesGraphqlType = graphql.NewList(placeGraphqlType)

// post
var postGraphqlType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "post",
		Fields: graphql.Fields{
			"post_id":       &graphql.Field{Type: int64GraphqlScalar},
			"user":          &graphql.Field{Type: userBasicGraphqlType},
			"content":       &graphql.Field{Type: graphql.String},
			"blob":          &graphql.Field{Type: blobGraphqlType},
			"type":          &graphql.Field{Type: graphql.Int, Description: "post type"},
			"like_count":    &graphql.Field{Type: int64GraphqlScalar},
			"dislike_count": &graphql.Field{Type: int64GraphqlScalar},
			"comment_count": &graphql.Field{Type: int64GraphqlScalar},
			"place":         &graphql.Field{Type: placeGraphqlType},
			"category_id":   &graphql.Field{Type: graphql.Int},
			"createtime":    &graphql.Field{Type: graphql.Int, Description: "Unix Timestamp"},
			"updatetime":    &graphql.Field{Type: graphql.Int, Description: "Unix Timestamp"},
		},
	},
)
var postWithoutPlaceGraphqlType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "post_without_place",
		Fields: graphql.Fields{
			"post_id":       &graphql.Field{Type: int64GraphqlScalar},
			"user":          &graphql.Field{Type: userBasicGraphqlType},
			"content":       &graphql.Field{Type: graphql.String},
			"blob":          &graphql.Field{Type: blobGraphqlType},
			"type":          &graphql.Field{Type: graphql.Int, Description: "post type"},
			"like_count":    &graphql.Field{Type: int64GraphqlScalar},
			"dislike_count": &graphql.Field{Type: int64GraphqlScalar},
			"comment_count": &graphql.Field{Type: int64GraphqlScalar},
			"category_id":   &graphql.Field{Type: graphql.Int},
			"createtime":    &graphql.Field{Type: graphql.Int, Description: "Unix Timestamp"},
			"updatetime":    &graphql.Field{Type: graphql.Int, Description: "Unix Timestamp"},
		},
	},
)
var postsWithoutPlaceGraphqlType = graphql.NewList(postWithoutPlaceGraphqlType)
var postsGraphqlType = graphql.NewList(postGraphqlType)

// hashtags
var hashtagGraphqlType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "hashtag",
		Fields: graphql.Fields{
			"hashtag_id": &graphql.Field{Type: int64GraphqlScalar},
			"value":      &graphql.Field{Type: graphql.String},
			"count":      &graphql.Field{Type: int64GraphqlScalar},
		},
	},
)
var hashtagsGraphqlType = graphql.NewList(hashtagGraphqlType)

// tags
var tagGraphqlType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "post_tags",
		Fields: graphql.Fields{
			"post_id":    &graphql.Field{Type: int64GraphqlScalar},
			"user":       &graphql.Field{Type: userBasicGraphqlType},
			"x":          &graphql.Field{Type: graphql.Int},
			"y":          &graphql.Field{Type: graphql.Int},
			"valid":      &graphql.Field{Type: graphql.Boolean},
			"createtime": &graphql.Field{Type: graphql.Int},
			"updatetime": &graphql.Field{Type: graphql.Int},
		},
	},
)
var tagsGraphqlType = graphql.NewList(tagGraphqlType)

// comment
var commentGraphqlType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "comment",
		Fields: graphql.Fields{
			"comment_id":    &graphql.Field{Type: int64GraphqlScalar},
			"post_id":       &graphql.Field{Type: int64GraphqlScalar},
			"user":          &graphql.Field{Type: userBasicGraphqlType},
			"comment":       &graphql.Field{Type: graphql.String},
			"like_count":    &graphql.Field{Type: int64GraphqlScalar},
			"dislike_count": &graphql.Field{Type: int64GraphqlScalar},
			"reply_count":   &graphql.Field{Type: int64GraphqlScalar},
			"createtime":    &graphql.Field{Type: graphql.Int, Description: "Unix Timestamp"},
			"updatetime":    &graphql.Field{Type: graphql.Int, Description: "Unix Timestamp"},
		},
	},
)
var commentsGraphqlType = graphql.NewList(commentGraphqlType)

// comment
var replyGraphqlType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "reply",
		Fields: graphql.Fields{
			"reply_id":      &graphql.Field{Type: int64GraphqlScalar},
			"comment_id":    &graphql.Field{Type: int64GraphqlScalar},
			"user":          &graphql.Field{Type: userBasicGraphqlType},
			"reply":         &graphql.Field{Type: graphql.String},
			"like_count":    &graphql.Field{Type: int64GraphqlScalar},
			"dislike_count": &graphql.Field{Type: int64GraphqlScalar},
			"createtime":    &graphql.Field{Type: graphql.Int, Description: "Unix Timestamp"},
			"updatetime":    &graphql.Field{Type: graphql.Int, Description: "Unix Timestamp"},
		},
	},
)
var repliesGraphqlType = graphql.NewList(replyGraphqlType)

// reaction
var postReactionGraphqlType = graphql.NewList(
	graphql.NewObject(
		graphql.ObjectConfig{
			Name: "post_reaction",
			Fields: graphql.Fields{
				"post_id":     &graphql.Field{Type: int64GraphqlScalar},
				"user":        &graphql.Field{Type: userBasicGraphqlType},
				"reaction_id": &graphql.Field{Type: graphql.Int, Description: "reaction"},
				"createtime":  &graphql.Field{Type: graphql.Int, Description: "Unix Timestamp"},
			},
		},
	),
)
var commentReactionGraphqlType = graphql.NewList(
	graphql.NewObject(
		graphql.ObjectConfig{
			Name: "comment_reaction",
			Fields: graphql.Fields{
				"comment_id":  &graphql.Field{Type: int64GraphqlScalar},
				"user":        &graphql.Field{Type: userBasicGraphqlType},
				"reaction_id": &graphql.Field{Type: graphql.Int},
				"createtime":  &graphql.Field{Type: graphql.Int, Description: "Unix Timestamp"},
			},
		},
	),
)
var replyReactionGraphqlType = graphql.NewList(
	graphql.NewObject(
		graphql.ObjectConfig{
			Name: "reply_reaction",
			Fields: graphql.Fields{
				"reply_id":    &graphql.Field{Type: int64GraphqlScalar},
				"user":        &graphql.Field{Type: userBasicGraphqlType},
				"reaction_id": &graphql.Field{Type: graphql.Int},
				"createtime":  &graphql.Field{Type: graphql.Int, Description: "Unix Timestamp"},
			},
		},
	),
)

// funcs
func parseAuth(p graphql.ResolveParams) (user xuserAPI, err error) {
	userToken, isOK := p.Context.Value(contextUserToken).(string)
	if !isOK {
		return user, errors.New("user-token format")
	}
	user, err = checkSession(userToken)
	if err != nil {
		return user, errors.New("user-token invalid")
	}
	return user, nil
}
func parsePost(p graphql.ResolveParams, postID, userID int64, postType, originWidth, originHeight int) (post postAPI, err error) {
	content, isOK := p.Args["content"].(string)
	if !isOK {
		return post, errors.New("content format")
	}
	categoryID, isOK := p.Args["category_id"].(int)
	if !isOK {
		return post, errors.New("category_id format")
	}
	placeID, isOK := p.Args["place_id"].(int64)
	if !isOK {
		return post, errors.New("place_id format")
	}
	timestamp := getNowUnixTimestamp()
	post.PostID = postID
	post.User.UserID = userID
	post.Content = content
	post.Blob.BlobID = strconv.FormatInt(userID, 10) + "_" + strconv.Itoa(timestamp)
	post.Blob.OriginWidth = originWidth
	post.Blob.OriginHeight = originHeight
	post.Type = postType
	post.LikeCount = 0
	post.DislikeCount = 0
	post.CommentCount = 0
	post.Place.PlaceID = placeID
	post.CategoryID = categoryID
	post.Public = true
	post.Createtime = timestamp
	post.Updatetime = timestamp
	return post, nil
}
func checkPost(post postAPI) error {
	if post.Place.PlaceID == 0 && post.CategoryID == categorySup {
		return errors.New(categoryMapID2Name[categorySup] + " must with place")
	}
	return nil
}
func parseComment(p graphql.ResolveParams, commentID, postID, userID int64) (c commentAPI, err error) {
	comment, isOK := p.Args["comment"].(string)
	if !isOK {
		return c, errors.New("comment format")
	}
	timestamp := getNowUnixTimestamp()
	c.CommentID = commentID
	c.PostID = postID
	c.User.UserID = userID
	c.Comment = comment
	c.LikeCount = 0
	c.DislikeCount = 0
	c.ReplyCount = 0
	c.Createtime = timestamp
	c.Updatetime = timestamp
	return c, nil
}
func parseReply(p graphql.ResolveParams, replyID, commentID, userID int64) (r replyAPI, err error) {
	reply, isOK := p.Args["reply"].(string)
	if !isOK {
		return r, errors.New("reply format")
	}
	timestamp := getNowUnixTimestamp()
	r.ReplyID = replyID
	r.CommentID = commentID
	r.User.UserID = userID
	r.Reply = reply
	r.LikeCount = 0
	r.DislikeCount = 0
	r.Createtime = timestamp
	r.Updatetime = timestamp
	return r, nil
}

// schema
var graphqlSchema, _ = graphql.NewSchema(
	graphql.SchemaConfig{
		Query:    graphqlQueryType,
		Mutation: graphqlMutationType,
	},
)

var graphqlQueryType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"country": &graphql.Field{
				Type: countryBasicGraphqlType,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					ct, err := getCountries()
					return ct, err
				},
				Description: "",
			},
			"language": &graphql.Field{
				Type: languageGraphqlType,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return languageConfigAPI, nil
				},
				Description: "",
			},
			"gender": &graphql.Field{
				Type: genderGraphqlType,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return genderConfigAPI, nil
				},
				Description: "",
			},
			"reaction": &graphql.Field{
				Type: reactionGraphqlType,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return reactionConfigAPI, nil
				},
				Description: "",
			},
			"post_type": &graphql.Field{
				Type: postTypeGraphqlType,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return postTypeConfigAPI, nil
				},
				Description: "",
			},
			"category": &graphql.Field{
				Type: categoryGraphqlType,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return categoryConfigAPI, nil
				},
				Description: "",
			},
			"login": &graphql.Field{
				Type: loginGraphqlType,
				Args: graphql.FieldConfigArgument{
					"email": &graphql.ArgumentConfig{
						Type:        graphql.String,
						Description: "email",
					},
					"password": &graphql.ArgumentConfig{
						Type:        graphql.String,
						Description: "password",
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					email, isOK := p.Args["email"].(string)
					if !isOK {
						return nil, nil
					}
					password, isOK := p.Args["password"].(string)
					if !isOK {
						return nil, nil
					}
					lc, err := login(email, password)
					if err != nil {
						return lc, err
					}
					return lc, nil
				},
				Description: "",
			},
			"user_by_user_id": &graphql.Field{
				Type: userGraphqlType,
				Args: graphql.FieldConfigArgument{
					"user_id": &graphql.ArgumentConfig{
						Type:        int64GraphqlScalar,
						Description: "user_id",
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					userID, isOK := p.Args["user_id"].(int64)
					if !isOK {
						return nil, errors.New("user_id format")
					}
					user, err := getUserByUserID(userID)
					return user, err
				},
				Description: "",
			},
			"user_by_username": &graphql.Field{
				Type: userGraphqlType,
				Args: graphql.FieldConfigArgument{
					"username": &graphql.ArgumentConfig{
						Type:        graphql.String,
						Description: "username",
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					username, isOK := p.Args["username"].(string)
					if !isOK {
						return nil, errors.New("username format")
					}
					user, err := getUserByUsername(username)
					return user, err
				},
				Description: "",
			},
			"users_by_following": &graphql.Field{
				Type: usersFollowingGraphqlType,
				Args: graphql.FieldConfigArgument{
					"page": &graphql.ArgumentConfig{
						Type:        graphql.Int,
						Description: "page",
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					user, err := parseAuth(p)
					if err != nil {
						return nil, err
					}
					page, isOK := p.Args["page"].(int)
					if !isOK {
						return nil, errors.New("page format")
					}
					users, err := getUsersByFollowing(user.UserID, page)
					return users, err
				},
				Description: "",
			},
			"users_by_follower": &graphql.Field{
				Type: usersFollowerGraphqlType,
				Args: graphql.FieldConfigArgument{
					"page": &graphql.ArgumentConfig{
						Type:        graphql.Int,
						Description: "page",
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					user, err := parseAuth(p)
					if err != nil {
						return nil, err
					}
					page, isOK := p.Args["page"].(int)
					if !isOK {
						return nil, errors.New("page format")
					}
					users, err := getUsersByFollower(user.UserID, page)
					return users, err
				},
				Description: "",
			},
			"user_is_following": &graphql.Field{
				Type: graphql.Boolean,
				Args: graphql.FieldConfigArgument{
					"user_id": &graphql.ArgumentConfig{
						Type:        int64GraphqlScalar,
						Description: "user_id",
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					user, err := parseAuth(p)
					if err != nil {
						return nil, err
					}
					userID, isOK := p.Args["user_id"].(int64)
					if !isOK {
						return nil, errors.New("user_id format")
					}
					isFollowing, err := checkUserIfFollowing(userID, user.UserID)
					return isFollowing, err
				},
				Description: "",
			},
			"taged_users_by_post": &graphql.Field{
				Type: tagsGraphqlType,
				Args: graphql.FieldConfigArgument{
					"post_id": &graphql.ArgumentConfig{
						Type:        int64GraphqlScalar,
						Description: "post_id",
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					postID, isOK := p.Args["post_id"].(int64)
					if !isOK {
						return nil, errors.New("post_id format")
					}
					// block check
					tags, err := getAllTagsByPost(postID)
					return tags, err
				},
				Description: "get all taged user in post",
			},
			"city_by_location": &graphql.Field{
				Type: cityGeometryGraphqlType,
				Args: graphql.FieldConfigArgument{
					"lat": &graphql.ArgumentConfig{
						Type:        graphql.Float,
						Description: "latitude",
					},
					"lon": &graphql.ArgumentConfig{
						Type:        graphql.Float,
						Description: "longitude",
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					lat, isOK := p.Args["lat"].(float64)
					if !isOK {
						return nil, errors.New("lat format")
					}
					lon, isOK := p.Args["lon"].(float64)
					if !isOK {
						return nil, errors.New("lon format")
					}
					cities, err := getCityByLocation(lat, lon)
					return cities, err
				},
				Description: "",
			},
			"country_post_count": &graphql.Field{
				Type: cityGraphqlType,
				Args: graphql.FieldConfigArgument{
					"country_code": &graphql.ArgumentConfig{
						Type:        graphql.String,
						Description: "country_code",
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					countryCode, isOK := p.Args["country_code"].(string)
					if !isOK {
						return nil, errors.New("country_code format")
					}
					country, err := getCountry(countryCode)
					return country, err
				},
				Description: "",
			},
			"city_post_count": &graphql.Field{
				Type: cityGraphqlType,
				Args: graphql.FieldConfigArgument{
					"level": &graphql.ArgumentConfig{
						Type:        graphql.Int,
						Description: "level",
					},
					"city_id": &graphql.ArgumentConfig{
						Type:        graphql.String,
						Description: "city_id",
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					level, isOK := p.Args["level"].(int)
					if !isOK {
						return nil, errors.New("level format")
					}
					if level < cityLevelRangeFirst || level > cityLevelRangeLast {
						return nil, errors.New("level format")
					}
					cityID, isOK := p.Args["city_id"].(string)
					if !isOK {
						return nil, errors.New("city_id id format")
					}
					city, err := getCity(strconv.Itoa(level), cityID)
					return city, err
				},
				Description: "",
			},
			"place_by_location": &graphql.Field{
				Type: placesLookupGraphql,
				Args: graphql.FieldConfigArgument{
					"lat": &graphql.ArgumentConfig{
						Type:        graphql.Float,
						Description: "latitude",
					},
					"lon": &graphql.ArgumentConfig{
						Type:        graphql.Float,
						Description: "longitude",
					},
					"name": &graphql.ArgumentConfig{
						Type:        graphql.String,
						Description: "enter\"\", if you have nothing special to look for",
					},
					"page_token": &graphql.ArgumentConfig{
						Type:        graphql.String,
						Description: "enter\"\" for first page, iter this from response",
					},
					// place type field provide by google place api?
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					lat, isOK := p.Args["lat"].(float64)
					if !isOK {
						return nil, errors.New("lat format")
					}
					lon, isOK := p.Args["lon"].(float64)
					if !isOK {
						return nil, errors.New("lon format")
					}
					name, isOK := p.Args["name"].(string)
					if !isOK {
						return nil, errors.New("name format")
					}
					pageToken, isOK := p.Args["page_token"].(string)
					if !isOK {
						return nil, errors.New("page_token format")
					}
					// place generated by google place api
					placesGCP, nextPageToken, err := getPlaceByLocationGCP(lat, lon, name, pageToken)
					if err != nil {
						return nil, err
					}
					placesDB := []placeAPI{}
					if len(placesGCP) != 0 {
						// check place name, lat, lon in db
						// limited number by gcp
						placesDB, err = getPlacesByPlacesGCP(placesGCP)
						if err != nil {
							return nil, err
						}
					}
					// combine places then return
					places := []placeAPI{}
					for i := 0; i < len(placesGCP); i++ {
						place := placesGCP[i]
						for j := 0; j < len(placesDB); j++ {
							if placesGCP[i].Name == placesDB[j].Name &&
								placesGCP[i].Lat == placesDB[j].Lat &&
								placesGCP[i].Lon == placesDB[j].Lon { // geo near check?
								place = placesDB[j]
								break
							}
						}
						places = append(places, place)
					}
					return placesLookupAPI{places, nextPageToken}, err
				},
				Description: "",
			},
			"place_by_name": &graphql.Field{
				Type: placesLookupGraphql,
				Args: graphql.FieldConfigArgument{
					"name": &graphql.ArgumentConfig{
						Type:        graphql.String,
						Description: "",
					},
					"page_token": &graphql.ArgumentConfig{
						Type:        graphql.String,
						Description: "enter\"\" for first page, iter this from response",
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					name, isOK := p.Args["name"].(string)
					if !isOK {
						return nil, errors.New("name format")
					}
					pageToken, isOK := p.Args["page_token"].(string)
					if !isOK {
						return nil, errors.New("page_token format")
					}
					placesGCP, nextPageToken, err := getPlaceByNameGCP(name, pageToken)
					if err != nil {
						return nil, err
					}
					placesDB := []placeAPI{}
					if len(placesGCP) != 0 {
						// check place name, lat, lon in db
						// limited number by gcp
						placesDB, err = getPlacesByPlacesGCP(placesGCP)
						if err != nil {
							return nil, err
						}
					}
					// combine places then return
					places := []placeAPI{}
					for i := 0; i < len(placesGCP); i++ {
						place := placesGCP[i]
						for j := 0; j < len(placesDB); j++ {
							if placesGCP[i].Name == placesDB[j].Name &&
								placesGCP[i].Lat == placesDB[j].Lat &&
								placesGCP[i].Lon == placesDB[j].Lon { // geo near check?
								place = placesDB[j]
								break
							}
						}
						places = append(places, place)
					}
					return placesLookupAPI{places, nextPageToken}, err
				},
				Description: "",
			},
			// post_detail
			"posts_by_recent": &graphql.Field{
				Type: postsGraphqlType,
				Args: graphql.FieldConfigArgument{
					"category_id": &graphql.ArgumentConfig{
						Type:        graphql.Int,
						Description: "category_id",
					},
					"page": &graphql.ArgumentConfig{
						Type:        graphql.Int,
						Description: "page",
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					categoryID, isOK := p.Args["category_id"].(int)
					if !isOK {
						return nil, errors.New("category format")
					}
					page, isOK := p.Args["page"].(int)
					if !isOK {
						return nil, errors.New("page format")
					}
					posts, err := getPostsByRecent(categoryID, page)
					return posts, err
				},
				Description: "",
			},
			"posts_by_recent_with_country": &graphql.Field{
				Type: postsGraphqlType,
				Args: graphql.FieldConfigArgument{
					"category_id": &graphql.ArgumentConfig{
						Type:        graphql.Int,
						Description: "category_id",
					},
					"country_code": &graphql.ArgumentConfig{
						Type:        graphql.String,
						Description: "country_code",
					},
					"page": &graphql.ArgumentConfig{
						Type:        graphql.Int,
						Description: "page",
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					categoryID, isOK := p.Args["category_id"].(int)
					if !isOK {
						return nil, errors.New("category format")
					}
					countryCode, isOK := p.Args["country_code"].(string)
					if !isOK {
						return nil, errors.New("country_code format")
					}
					page, isOK := p.Args["page"].(int)
					if !isOK {
						return nil, errors.New("page format")
					}
					// block check
					posts, err := getPostsByRecentWithCountry(countryCode, categoryID, page)
					return posts, err
				},
				Description: "",
			},
			"posts_by_recent_with_city": &graphql.Field{
				Type: postsGraphqlType,
				Args: graphql.FieldConfigArgument{
					"category_id": &graphql.ArgumentConfig{
						Type:        graphql.Int,
						Description: "category_id",
					},
					"level": &graphql.ArgumentConfig{
						Type:        graphql.Int,
						Description: "level",
					},
					"city_id": &graphql.ArgumentConfig{
						Type:        graphql.String,
						Description: "city_id",
					},
					"page": &graphql.ArgumentConfig{
						Type:        graphql.Int,
						Description: "page",
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					categoryID, isOK := p.Args["category_id"].(int)
					if !isOK {
						return nil, errors.New("category format")
					}
					level, isOK := p.Args["level"].(int)
					if !isOK {
						return nil, errors.New("level format")
					}
					if level < cityLevelRangeFirst || level > cityLevelRangeLast {
						return nil, errors.New("level format")
					}
					cityID, isOK := p.Args["city_id"].(string)
					if !isOK {
						return nil, errors.New("city_id format")
					}
					page, isOK := p.Args["page"].(int)
					if !isOK {
						return nil, errors.New("page format")
					}
					// block check
					posts, err := getPostsByRecentWithCity(strconv.Itoa(level), cityID, categoryID, page)
					return posts, err
				},
				Description: "",
			},
			"posts_by_following_users": &graphql.Field{
				Type: postsGraphqlType,
				Args: graphql.FieldConfigArgument{
					"page": &graphql.ArgumentConfig{
						Type:        graphql.Int,
						Description: "page",
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					user, err := parseAuth(p)
					if err != nil {
						return nil, err
					}
					page, isOK := p.Args["page"].(int)
					if !isOK {
						return nil, errors.New("page format")
					}
					// block check
					posts, err := getPostsByFollowingUsers(user.UserID, page)
					return posts, err
				},
				Description: "",
			},
			"posts_by_following_users_with_country": &graphql.Field{
				Type: postsGraphqlType,
				Args: graphql.FieldConfigArgument{
					"country_code": &graphql.ArgumentConfig{
						Type:        graphql.String,
						Description: "country_code",
					},
					"category_id": &graphql.ArgumentConfig{
						Type:        graphql.Int,
						Description: "category_id",
					},
					"page": &graphql.ArgumentConfig{
						Type:        graphql.Int,
						Description: "page",
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					user, err := parseAuth(p)
					if err != nil {
						return nil, err
					}
					countryCode, isOK := p.Args["country_code"].(string)
					if !isOK {
						return nil, errors.New("country_code format")
					}
					categoryID, isOK := p.Args["category_id"].(int)
					if !isOK {
						return nil, errors.New("category format")
					}
					page, isOK := p.Args["page"].(int)
					if !isOK {
						return nil, errors.New("page format")
					}
					// block check
					posts, err := getPostsByFollowingUsersWithCountry(user.UserID, countryCode, categoryID, page)
					return posts, err
				},
				Description: "",
			},
			"posts_by_following_users_with_city": &graphql.Field{
				Type: postsGraphqlType,
				Args: graphql.FieldConfigArgument{
					"level": &graphql.ArgumentConfig{
						Type:        graphql.Int,
						Description: "level",
					},
					"city_id": &graphql.ArgumentConfig{
						Type:        graphql.String,
						Description: "city_id",
					},
					"category_id": &graphql.ArgumentConfig{
						Type:        graphql.Int,
						Description: "category_id",
					},
					"page": &graphql.ArgumentConfig{
						Type:        graphql.Int,
						Description: "page",
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					user, err := parseAuth(p)
					if err != nil {
						return nil, err
					}
					level, isOK := p.Args["level"].(int)
					if !isOK {
						return nil, errors.New("level format")
					}
					if level < cityLevelRangeFirst || level > cityLevelRangeLast {
						return nil, errors.New("level format")
					}
					cityID, isOK := p.Args["city_id"].(string)
					if !isOK {
						return nil, errors.New("city_id format")
					}
					categoryID, isOK := p.Args["category_id"].(int)
					if !isOK {
						return nil, errors.New("category format")
					}
					page, isOK := p.Args["page"].(int)
					if !isOK {
						return nil, errors.New("page format")
					}
					// block check
					posts, err := getPostsByFollowingUsersWithCity(user.UserID, strconv.Itoa(level), cityID, categoryID, page)
					return posts, err
				},
				Description: "",
			},
			"posts_by_user": &graphql.Field{
				Type: postsGraphqlType,
				Args: graphql.FieldConfigArgument{
					"user_id": &graphql.ArgumentConfig{
						Type:        int64GraphqlScalar,
						Description: "user_id",
					},
					"page": &graphql.ArgumentConfig{
						Type:        graphql.Int,
						Description: "page",
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					userID, isOK := p.Args["user_id"].(int64)
					if !isOK {
						return nil, errors.New("user_id format")
					}
					page, isOK := p.Args["page"].(int)
					if !isOK {
						return nil, errors.New("page format")
					}
					// block check
					posts, err := getPostsByUser(userID, page)
					return posts, err
				},
				Description: "",
			},
			"posts_by_popular": &graphql.Field{
				Type: postsGraphqlType,
				Args: graphql.FieldConfigArgument{
					"category_id": &graphql.ArgumentConfig{
						Type:        graphql.Int,
						Description: "category_id",
					},
					"page": &graphql.ArgumentConfig{
						Type:        graphql.Int,
						Description: "page",
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					startTime := time.Now()
					user, err := parseAuth(p)
					if err != nil {
						return nil, err
					}
					categoryID, isOK := p.Args["category_id"].(int)
					if !isOK {
						return nil, errors.New("category format")
					}
					page, isOK := p.Args["page"].(int)
					if !isOK {
						return nil, errors.New("page format")
					}
					// block check
					posts, err := getPostsByPopular(user.UserID, categoryID, page)
					log.Printf("posts_by_popular total took %fs\n", time.Since(startTime).Seconds())
					return posts, err
				},
				Description: "",
			},
			"sup_posts_by_popular_with_country": &graphql.Field{
				Type: postsGraphqlType,
				Args: graphql.FieldConfigArgument{
					"country_code": &graphql.ArgumentConfig{
						Type:        graphql.String,
						Description: "country_code",
					},
					"page": &graphql.ArgumentConfig{
						Type:        graphql.Int,
						Description: "page",
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					startTime := time.Now()
					user, err := parseAuth(p)
					if err != nil {
						return nil, err
					}
					countryCode, isOK := p.Args["country_code"].(string)
					if !isOK {
						return nil, errors.New("country_code format")
					}
					page, isOK := p.Args["page"].(int)
					if !isOK {
						return nil, errors.New("page format")
					}
					// block check
					posts, err := getSupPostsByPopularWithCountry(countryCode, user.UserID, page)
					log.Printf("sup_posts_by_popular_with_country total took %fs\n", time.Since(startTime).Seconds())
					return posts, err
				},
				Description: "",
			},
			"sup_posts_by_popular_with_city": &graphql.Field{
				Type: postsGraphqlType,
				Args: graphql.FieldConfigArgument{
					"level": &graphql.ArgumentConfig{
						Type:        graphql.Int,
						Description: "level",
					},
					"city_id": &graphql.ArgumentConfig{
						Type:        graphql.String,
						Description: "city_id",
					},
					"page": &graphql.ArgumentConfig{
						Type:        graphql.Int,
						Description: "page",
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					startTime := time.Now()
					user, err := parseAuth(p)
					if err != nil {
						return nil, err
					}
					level, isOK := p.Args["level"].(int)
					if !isOK {
						return nil, errors.New("level format")
					}
					if level < cityLevelRangeFirst || level > cityLevelRangeLast {
						return nil, errors.New("level format")
					}
					cityID, isOK := p.Args["city_id"].(string)
					if !isOK {
						return nil, errors.New("city_id format")
					}
					page, isOK := p.Args["page"].(int)
					if !isOK {
						return nil, errors.New("page format")
					}
					// block check
					posts, err := getSupPostsByPopularWithCity(strconv.Itoa(level), cityID, user.UserID, page)
					log.Printf("sup_posts_by_popular_with_country total took %fs\n", time.Since(startTime).Seconds())
					return posts, err
				},
				Description: "",
			},
			"posts_by_hashtag": &graphql.Field{
				Type: postsWithoutPlaceGraphqlType,
				Args: graphql.FieldConfigArgument{
					"hashtag_id": &graphql.ArgumentConfig{
						Type:        int64GraphqlScalar,
						Description: "hashtag_id",
					},
					"page": &graphql.ArgumentConfig{
						Type:        graphql.Int,
						Description: "page",
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					hashtagID, isOK := p.Args["hashtag_id"].(int64)
					if !isOK {
						return nil, errors.New("hashtag_id format")
					}
					page, isOK := p.Args["page"].(int)
					if !isOK {
						return nil, errors.New("page format")
					}
					// block check
					posts, err := getPostsByHashtag(hashtagID, page)
					return posts, err
				},
				Description: "",
			},
			"posts_by_tag": &graphql.Field{
				Type: postsWithoutPlaceGraphqlType,
				Args: graphql.FieldConfigArgument{
					"user_id": &graphql.ArgumentConfig{
						Type:        int64GraphqlScalar,
						Description: "user_id",
					},
					"page": &graphql.ArgumentConfig{
						Type:        graphql.Int,
						Description: "page",
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					userID, isOK := p.Args["user_id"].(int64)
					if !isOK {
						return nil, errors.New("user_id format")
					}
					page, isOK := p.Args["page"].(int)
					if !isOK {
						return nil, errors.New("page format")
					}
					// block check
					posts, err := getPostsByTag(userID, page)
					return posts, err
				},
				Description: "",
			},
			"posts_by_place": &graphql.Field{
				Type: postsGraphqlType,
				Args: graphql.FieldConfigArgument{
					"place_id": &graphql.ArgumentConfig{
						Type:        int64GraphqlScalar,
						Description: "user_id",
					},
					"page": &graphql.ArgumentConfig{
						Type:        graphql.Int,
						Description: "page",
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					placeID, isOK := p.Args["place_id"].(int64)
					if !isOK {
						return nil, errors.New("place_id format")
					}
					page, isOK := p.Args["page"].(int)
					if !isOK {
						return nil, errors.New("page format")
					}
					// block check
					posts, err := getPostsByPlaceID(placeID, page)
					return posts, err
				},
				Description: "",
			},
			"hashtags": &graphql.Field{
				Type: hashtagsGraphqlType,
				Args: graphql.FieldConfigArgument{
					"value": &graphql.ArgumentConfig{
						Type:        graphql.String,
						Description: "value",
					},
					"page": &graphql.ArgumentConfig{
						Type:        graphql.Int,
						Description: "page",
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					value, isOK := p.Args["value"].(string)
					if !isOK {
						return nil, errors.New("value format")
					}
					page, isOK := p.Args["page"].(int)
					if !isOK {
						return nil, errors.New("page format")
					}
					hashtags, err := getHashtags(value, page)
					return hashtags, err
				},
				Description: "",
			},
			"reactions_on_post": &graphql.Field{
				Type: postReactionGraphqlType,
				Args: graphql.FieldConfigArgument{
					"post_id": &graphql.ArgumentConfig{
						Type:        int64GraphqlScalar,
						Description: "post_id",
					},
					"page": &graphql.ArgumentConfig{
						Type:        graphql.Int,
						Description: "page",
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					postID, isOK := p.Args["post_id"].(int64)
					if !isOK {
						return nil, errors.New("post_id format")
					}
					page, isOK := p.Args["page"].(int)
					if !isOK {
						return nil, errors.New("page format")
					}
					reactionsPost, err := getReactionsOnPost(postID, page)
					return reactionsPost, err
				},
				Description: "",
			},
			"reactions_on_comment": &graphql.Field{
				Type: commentReactionGraphqlType,
				Args: graphql.FieldConfigArgument{
					"comment_id": &graphql.ArgumentConfig{
						Type:        int64GraphqlScalar,
						Description: "comment_id",
					},
					"page": &graphql.ArgumentConfig{
						Type:        graphql.Int,
						Description: "page",
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					commentID, isOK := p.Args["comment_id"].(int64)
					if !isOK {
						return nil, errors.New("comment_id format")
					}
					page, isOK := p.Args["page"].(int)
					if !isOK {
						return nil, errors.New("page format")
					}
					reactionsComment, err := getReactionsOnComment(commentID, page)
					return reactionsComment, err
				},
				Description: "",
			},
			"reactions_on_reply": &graphql.Field{
				Type: replyReactionGraphqlType,
				Args: graphql.FieldConfigArgument{
					"reply_id": &graphql.ArgumentConfig{
						Type:        int64GraphqlScalar,
						Description: "reply_id",
					},
					"page": &graphql.ArgumentConfig{
						Type:        graphql.Int,
						Description: "page",
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					replyID, isOK := p.Args["reply_id"].(int64)
					if !isOK {
						return nil, errors.New("reply_id format")
					}
					page, isOK := p.Args["page"].(int)
					if !isOK {
						return nil, errors.New("page format")
					}
					reactionsReply, err := getReactionsOnReply(replyID, page)
					return reactionsReply, err
				},
				Description: "",
			},
			"comments_on_post": &graphql.Field{
				Type: commentsGraphqlType,
				Args: graphql.FieldConfigArgument{
					"post_id": &graphql.ArgumentConfig{
						Type:        int64GraphqlScalar,
						Description: "post_id",
					},
					"page": &graphql.ArgumentConfig{
						Type:        graphql.Int,
						Description: "page",
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					postID, isOK := p.Args["post_id"].(int64)
					if !isOK {
						return nil, errors.New("post_id format")
					}
					page, isOK := p.Args["page"].(int)
					if !isOK {
						return nil, errors.New("page format")
					}
					comments, err := getCommentsOnPost(postID, page)
					return comments, err
				},
				Description: "",
			},
			"replies_on_comment": &graphql.Field{
				Type: repliesGraphqlType,
				Args: graphql.FieldConfigArgument{
					"comment_id": &graphql.ArgumentConfig{
						Type:        int64GraphqlScalar,
						Description: "comment_id",
					},
					"page": &graphql.ArgumentConfig{
						Type:        graphql.Int,
						Description: "page",
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					commentID, isOK := p.Args["comment_id"].(int64)
					if !isOK {
						return nil, errors.New("post_id format")
					}
					page, isOK := p.Args["page"].(int)
					if !isOK {
						return nil, errors.New("page format")
					}
					replies, err := getRepliesOnComment(commentID, page)
					return replies, err
				},
				Description: "",
			},
		},
	})
var graphqlMutationType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			"follow": &graphql.Field{
				Type: updateStatusGraphqlType,
				Args: graphql.FieldConfigArgument{
					"user_id": &graphql.ArgumentConfig{
						Type:        int64GraphqlScalar,
						Description: "user_id",
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					user, err := parseAuth(p)
					if err != nil {
						return nil, err
					}
					userID, isOK := p.Args["user_id"].(int64)
					if !isOK {
						return nil, errors.New("user_id format")
					}
					// follow count limit?
					us, err := follow(userID, user.UserID)
					return us, err
				},
				Description: "",
			},
			"unfollow": &graphql.Field{
				Type: updateStatusGraphqlType,
				Args: graphql.FieldConfigArgument{
					"user_id": &graphql.ArgumentConfig{
						Type:        int64GraphqlScalar,
						Description: "user_id",
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					user, err := parseAuth(p)
					if err != nil {
						return nil, err
					}
					userID, isOK := p.Args["user_id"].(int64)
					if !isOK {
						return nil, errors.New("user_id format")
					}
					us, err := unfollow(userID, user.UserID)
					return us, err
				},
				Description: "",
			},
			"place_insert": &graphql.Field{
				Type: placeGraphqlType,
				Args: graphql.FieldConfigArgument{
					"lat": &graphql.ArgumentConfig{
						Type:        graphql.Float,
						Description: "latitude",
					},
					"lon": &graphql.ArgumentConfig{
						Type:        graphql.Float,
						Description: "longitude",
					},
					"name": &graphql.ArgumentConfig{
						Type:        graphql.String,
						Description: "",
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					lat, isOK := p.Args["lat"].(float64)
					if !isOK {
						return nil, errors.New("lat format")
					}
					lon, isOK := p.Args["lon"].(float64)
					if !isOK {
						return nil, errors.New("lon format")
					}
					name, isOK := p.Args["name"].(string)
					if !isOK {
						return nil, errors.New("name format")
					}
					// format check
					// source check
					place := placeAPI{Lat: lat, Lon: lon, Name: name}
					// city lookup
					cities, err := getCityByLocation(lat, lon)
					if err != nil {
						return nil, err
					}
					if len(cities) >= 1 { // which is closest city sortby mongo
						place.CountryCode = cities[0].Properties.CountryCode
						place.CityID1 = cities[0].Properties.CityID1
						place.CityID2 = cities[0].Properties.CityID2
						place.CityID3 = cities[0].Properties.CityID3
						place.CityID4 = cities[0].Properties.CityID4
						place.CityID5 = cities[0].Properties.CityID5
					}
					// place_insert
					place.PlaceID, err = placeInsert(place)
					return place, err
				},
				Description: "please enter data which query place provide",
			},
			"post_insert": &graphql.Field{
				Type: postGraphqlType,
				Args: graphql.FieldConfigArgument{
					"content": &graphql.ArgumentConfig{
						Type:        graphql.String,
						Description: "content",
					},
					"category_id": &graphql.ArgumentConfig{
						Type:        graphql.Int,
						Description: "category_id",
					},
					"place_id": &graphql.ArgumentConfig{
						Type:        int64GraphqlScalar,
						Description: "place_id",
					},
					"type": &graphql.ArgumentConfig{
						Type:        graphql.Int,
						Description: "post type",
					},
					"file": &graphql.ArgumentConfig{
						Type:        graphql.String,
						Description: "tar.gz[*.jpg,...]/[*.m3u8(only one), *.ts...]",
					},
					"origin_width": &graphql.ArgumentConfig{
						Type:        graphql.Int,
						Description: "origin width",
					},
					"origin_height": &graphql.ArgumentConfig{
						Type:        graphql.Int,
						Description: "origin height",
					},
					"tags": &graphql.ArgumentConfig{
						Type:        tagsOnPostGraphqlScalar,
						Description: "tags",
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					startTime := time.Now()
					user, err := parseAuth(p)
					if err != nil {
						return nil, err
					}
					file, isOK := p.Context.Value(contextKeyFile).(io.Reader)
					if !isOK {
						return nil, errors.New("file format")
					}
					// post parameter check
					postType, isOK := p.Args["type"].(int)
					if !isOK || postTypeMapID2Type[postType] == "" {
						return nil, errors.New("type format")
					}
					originWidth, isOK := p.Args["origin_width"].(int)
					if !isOK {
						return nil, errors.New("origin_width format")
					}
					originHeight, isOK := p.Args["origin_height"].(int)
					if !isOK {
						return nil, errors.New("origin_height format")
					}
					tags, isOK := p.Args["tags"].([]tagOnPostSetAPI)
					if !isOK {
						return nil, errors.New("tags format")
					}
					post, err := parsePost(p, 0, user.UserID, postType, originWidth, originHeight)
					if err != nil {
						return post, err
					}
					if err = checkPost(post); err != nil {
						return post, err
					}
					// file size check
					err = untarFileAndUpload(post, file, true)
					if err != nil {
						return nil, err
					}
					post.PostID, err = postInsert(post)
					if err != nil {
						return post, err
					}
					// hashtag check[max, rule]
					hashtags, _ := checkMention(post.Content)
					if len(hashtags) > 0 {
						hashtagsID, err := hashtagInsert(hashtags)
						if err != nil {
							return post, err
						}
						if post.PostID != 0 {
							_, err = hashtagOnPostSet(post.PostID, hashtagsID)
						}
					}
					// tag check[max] notification
					_, err = tagsOnPostSet(post.PostID, tags)
					log.Printf("post now total took %fs\n", time.Since(startTime).Seconds())
					return post, err
				},
				Description: "place_id:0 => no place",
				DeprecationReason: `please use form-data to upload file, form-data key:
						query: mutation{post(...:...){post_id}}
						file: tar.gz file,
						. not finished yet. check ~/upload.html
					`,
			},
			"post_update": &graphql.Field{
				Type: updateStatusGraphqlType,
				Args: graphql.FieldConfigArgument{
					"post_id": &graphql.ArgumentConfig{
						Type:        int64GraphqlScalar,
						Description: "post_id",
					},
					"content": &graphql.ArgumentConfig{
						Type:        graphql.String,
						Description: "content",
					},
					"is_update_content": &graphql.ArgumentConfig{
						Type:        graphql.Boolean,
						Description: "is_update_content",
					},
					"category_id": &graphql.ArgumentConfig{
						Type:        graphql.Int,
						Description: "category_id",
					},
					"place_id": &graphql.ArgumentConfig{
						Type:        int64GraphqlScalar,
						Description: "place_id",
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					user, err := parseAuth(p)
					if err != nil {
						return nil, err
					}
					postID, isOK := p.Args["post_id"].(int64)
					if !isOK {
						return nil, errors.New("post_id format")
					}
					isUpdateContent, isOK := p.Args["is_update_content"].(bool)
					if !isOK {
						return nil, errors.New("is_update_content format")
					}
					// post parameter check
					post, err := parsePost(p, postID, user.UserID, 0, 0, 0)
					if err != nil {
						return post, err
					}
					if err = checkPost(post); err != nil {
						return post, err
					}
					us, err := postUpdate(post)
					if err != nil {
						return us, err
					}
					if isUpdateContent { // prevent unnecessary process
						// hashtag check[max, rule]
						hashtags, _ := checkMention(post.Content)
						if len(hashtags) > 0 {
							hashtagsID, err := hashtagInsert(hashtags)
							if err != nil {
								return post, err
							}
							if post.PostID != 0 {
								_, err = hashtagOnPostSet(post.PostID, hashtagsID)
							}
						}
					}
					return us, err
				},
				Description: "",
			},
			"post_tag_conform": &graphql.Field{
				Type: updateStatusGraphqlType,
				Args: graphql.FieldConfigArgument{
					"post_id": &graphql.ArgumentConfig{
						Type:        int64GraphqlScalar,
						Description: "post_id",
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					user, err := parseAuth(p)
					if err != nil {
						return nil, err
					}
					postID, isOK := p.Args["post_id"].(int64)
					if !isOK {
						return nil, errors.New("post_id format")
					}
					us, err := tagOnPostConfirm(postID, user.UserID)
					return us, err
				},
				Description: "user can confirm tag in post",
			},
			"post_tag_delete": &graphql.Field{
				Type: updateStatusGraphqlType,
				Args: graphql.FieldConfigArgument{
					"post_id": &graphql.ArgumentConfig{
						Type:        int64GraphqlScalar,
						Description: "post_id",
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					user, err := parseAuth(p)
					if err != nil {
						return nil, err
					}
					postID, isOK := p.Args["post_id"].(int64)
					if !isOK {
						return nil, errors.New("post_id format")
					}
					us, err := tagOnPostDelete(postID, user.UserID)
					return us, err
				},
				Description: "user can remove tag in post",
			},
			"post_tag_update": &graphql.Field{
				Type: updateStatusGraphqlType,
				Args: graphql.FieldConfigArgument{
					"post_id": &graphql.ArgumentConfig{
						Type:        int64GraphqlScalar,
						Description: "post_id",
					},
					"tag": &graphql.ArgumentConfig{
						Type:        tagOnPostGraphqlScalar,
						Description: "tag",
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					postID, isOK := p.Args["post_id"].(int64)
					if !isOK {
						return nil, errors.New("post_id format")
					}
					tag, isOK := p.Args["tag"].(tagOnPostSetAPI)
					if !isOK {
						return nil, errors.New("tag format")
					}
					// tag check[max] notification
					us, err := tagOnPostUpdate(postID, tag)
					return us, err
				},
				Description: "",
			},
			"post_delete": &graphql.Field{
				Type: updateStatusGraphqlType,
				Args: graphql.FieldConfigArgument{
					"post_id": &graphql.ArgumentConfig{
						Type:        int64GraphqlScalar,
						Description: "post_id",
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					user, err := parseAuth(p)
					if err != nil {
						return nil, err
					}
					postID, isOK := p.Args["post_id"].(int64)
					if !isOK {
						return nil, errors.New("post_id format")
					}
					post := postAPI{PostID: postID, User: userBasicAPI{UserID: user.UserID}}
					us, err := postDelete(post)
					return us, err
				},
				Description: "",
			},
			"post_popular_read": &graphql.Field{
				Type: postsGraphqlType,
				Args: graphql.FieldConfigArgument{
					"category_id": &graphql.ArgumentConfig{
						Type:        graphql.Int,
						Description: "category_id",
					},
					"index_read": &graphql.ArgumentConfig{
						Type:        graphql.Int,
						Description: "index_read, [0 - ∞)",
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					startTime := time.Now()
					user, err := parseAuth(p)
					if err != nil {
						return nil, err
					}
					categoryID, isOK := p.Args["category_id"].(int)
					if !isOK {
						return nil, errors.New("category format")
					}
					indexRead, isOK := p.Args["index_read"].(int)
					if !isOK {
						return nil, errors.New("index_read format")
					}
					if indexRead < 0 {
						return nil, errors.New("index_read format")
					}
					posts, err := postPopularRead(categoryID, indexRead, user.UserID)
					log.Printf("post_popular_read total took %fs\n", time.Since(startTime).Seconds())
					return posts, err
				},
				Description: "",
			},
			"sup_post_popular_read_with_country": &graphql.Field{
				Type: postsGraphqlType,
				Args: graphql.FieldConfigArgument{
					"country_code": &graphql.ArgumentConfig{
						Type:        graphql.String,
						Description: "country_code",
					},
					"index_read": &graphql.ArgumentConfig{
						Type:        graphql.Int,
						Description: "index_read, [0 - ∞)",
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					startTime := time.Now()
					user, err := parseAuth(p)
					if err != nil {
						return nil, err
					}
					countryCode, isOK := p.Args["country_code"].(string)
					if !isOK {
						return nil, errors.New("country_code format")
					}
					indexRead, isOK := p.Args["index_read"].(int)
					if !isOK {
						return nil, errors.New("index_read format")
					}
					if indexRead < 0 {
						return nil, errors.New("index_read format")
					}
					posts, err := supPostPopularReadCountry(countryCode, indexRead, user.UserID)
					log.Printf("sup_post_popular_read_with_country total took %fs\n", time.Since(startTime).Seconds())
					return posts, err
				},
				Description: "",
			},
			"sup_post_popular_read_with_city": &graphql.Field{
				Type: postsGraphqlType,
				Args: graphql.FieldConfigArgument{
					"level": &graphql.ArgumentConfig{
						Type:        graphql.Int,
						Description: "level",
					},
					"city_id": &graphql.ArgumentConfig{
						Type:        graphql.String,
						Description: "city_id",
					},
					"index_read": &graphql.ArgumentConfig{
						Type:        graphql.Int,
						Description: "index_read, [0 - ∞)",
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					startTime := time.Now()
					user, err := parseAuth(p)
					if err != nil {
						return nil, err
					}
					level, isOK := p.Args["level"].(int)
					if !isOK {
						return nil, errors.New("level format")
					}
					if level < cityLevelRangeFirst || level > cityLevelRangeLast {
						return nil, errors.New("level format")
					}
					cityID, isOK := p.Args["city_id"].(string)
					if !isOK {
						return nil, errors.New("city_id format")
					}
					indexRead, isOK := p.Args["index_read"].(int)
					if !isOK {
						return nil, errors.New("index_read format")
					}
					if indexRead < 0 {
						return nil, errors.New("index_read format")
					}
					posts, err := supPostPopularReadCity(strconv.Itoa(level), cityID, indexRead, user.UserID)
					log.Printf("sup_post_popular_read_with_city total took %fs\n", time.Since(startTime).Seconds())
					return posts, err
				},
				Description: "",
			},
			"reaction_on_post": &graphql.Field{
				Type: updateStatusGraphqlType,
				Args: graphql.FieldConfigArgument{
					"post_id": &graphql.ArgumentConfig{
						Type:        int64GraphqlScalar,
						Description: "post_id",
					},
					"reaction_id": &graphql.ArgumentConfig{
						Type:        graphql.Int,
						Description: "reaction_id",
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					user, err := parseAuth(p)
					if err != nil {
						return nil, err
					}
					postID, isOK := p.Args["post_id"].(int64)
					if !isOK {
						return nil, errors.New("post_id format")
					}
					reactionID, isOK := p.Args["reaction_id"].(int)
					if !isOK {
						return nil, errors.New("reaction_id format")
					}
					if reactionsMapID2Description[reactionID] == "" {
						return nil, errors.New("reaction_id format")
					}
					reactionOnPost := reactionOnPostAPI{
						PostID:     postID,
						User:       userBasicAPI{UserID: user.UserID},
						ReactionID: reactionID,
						Createtime: getNowUnixTimestamp(),
					}
					us, err := reactionOnPostSet(reactionOnPost)
					return us, err
				},
				Description: "",
			},
			"reaction_on_post_delete": &graphql.Field{
				Type: updateStatusGraphqlType,
				Args: graphql.FieldConfigArgument{
					"post_id": &graphql.ArgumentConfig{
						Type:        int64GraphqlScalar,
						Description: "post_id",
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					user, err := parseAuth(p)
					if err != nil {
						return nil, err
					}
					postID, isOK := p.Args["post_id"].(int64)
					if !isOK {
						return nil, errors.New("post_id format")
					}
					reactionOnPost := reactionOnPostAPI{
						PostID: postID,
						User:   userBasicAPI{UserID: user.UserID},
					}
					us, err := reactionOnPostDelete(reactionOnPost)
					return us, err
				},
				Description: "",
			},
			"reaction_on_comment": &graphql.Field{
				Type: updateStatusGraphqlType,
				Args: graphql.FieldConfigArgument{
					"comment_id": &graphql.ArgumentConfig{
						Type:        int64GraphqlScalar,
						Description: "comment_id",
					},
					"reaction_id": &graphql.ArgumentConfig{
						Type:        graphql.Int,
						Description: "reaction_id",
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					user, err := parseAuth(p)
					if err != nil {
						return nil, err
					}
					commentID, isOK := p.Args["comment_id"].(int64)
					if !isOK {
						return nil, errors.New("comment_id format")
					}
					reactionID, isOK := p.Args["reaction_id"].(int)
					if !isOK {
						return nil, errors.New("reaction_id format")
					}
					if reactionsMapID2Description[reactionID] == "" {
						return nil, errors.New("reaction_id format")
					}
					reactionOnComment := reactionOnCommentAPI{
						CommentID:  commentID,
						User:       userBasicAPI{UserID: user.UserID},
						ReactionID: reactionID,
						Createtime: getNowUnixTimestamp(),
					}
					us, err := reactionOnCommentSet(reactionOnComment)
					return us, err
				},
				Description: "",
			},
			"reaction_on_comment_delete": &graphql.Field{
				Type: updateStatusGraphqlType,
				Args: graphql.FieldConfigArgument{
					"comment_id": &graphql.ArgumentConfig{
						Type:        int64GraphqlScalar,
						Description: "comment_id",
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					user, err := parseAuth(p)
					if err != nil {
						return nil, err
					}
					commentID, isOK := p.Args["comment_id"].(int64)
					if !isOK {
						return nil, errors.New("comment_id format")
					}
					reactionOnComment := reactionOnCommentAPI{
						CommentID: commentID,
						User:      userBasicAPI{UserID: user.UserID},
					}
					us, err := reactionOnCommentDelete(reactionOnComment)
					return us, err
				},
				Description: "",
			},
			"reaction_on_reply": &graphql.Field{
				Type: updateStatusGraphqlType,
				Args: graphql.FieldConfigArgument{
					"reply_id": &graphql.ArgumentConfig{
						Type:        int64GraphqlScalar,
						Description: "reply_id",
					},
					"reaction_id": &graphql.ArgumentConfig{
						Type:        graphql.Int,
						Description: "reaction_id",
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					user, err := parseAuth(p)
					if err != nil {
						return nil, err
					}
					replyID, isOK := p.Args["reply_id"].(int64)
					if !isOK {
						return nil, errors.New("reply_id format")
					}
					reactionID, isOK := p.Args["reaction_id"].(int)
					if !isOK {
						return nil, errors.New("reaction_id format")
					}
					if reactionsMapID2Description[reactionID] == "" {
						return nil, errors.New("reaction_id format")
					}
					reactionOnReply := reactionOnReplyAPI{
						ReplyID:    replyID,
						User:       userBasicAPI{UserID: user.UserID},
						ReactionID: reactionID,
						Createtime: getNowUnixTimestamp(),
					}
					us, err := reactionOnReplySet(reactionOnReply)
					return us, err
				},
				Description: "",
			},
			"reaction_on_reply_delete": &graphql.Field{
				Type: updateStatusGraphqlType,
				Args: graphql.FieldConfigArgument{
					"reply_id": &graphql.ArgumentConfig{
						Type:        int64GraphqlScalar,
						Description: "reply_id",
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					user, err := parseAuth(p)
					if err != nil {
						return nil, err
					}
					replyID, isOK := p.Args["reply_id"].(int64)
					if !isOK {
						return nil, errors.New("reply_id format")
					}
					reactionOnReply := reactionOnReplyAPI{
						ReplyID: replyID,
						User:    userBasicAPI{UserID: user.UserID},
					}
					us, err := reactionOnReplyDelete(reactionOnReply)
					return us, err
				},
				Description: "",
			},
			"comment_on_post_insert": &graphql.Field{
				Type: commentGraphqlType,
				Args: graphql.FieldConfigArgument{
					"post_id": &graphql.ArgumentConfig{
						Type:        int64GraphqlScalar,
						Description: "post_id",
					},
					"comment": &graphql.ArgumentConfig{
						Type:        graphql.String,
						Description: "comment",
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					user, err := parseAuth(p)
					if err != nil {
						return nil, err
					}
					postID, isOK := p.Args["post_id"].(int64)
					if !isOK {
						return nil, err
					}
					comment, err := parseComment(p, 0, postID, user.UserID)
					if err != nil {
						return comment, err
					}
					comment.CommentID, err = commentOnPostInsert(comment)
					return comment, err
				},
				Description: "",
			},
			"comment_on_post_update": &graphql.Field{
				Type: updateStatusGraphqlType,
				Args: graphql.FieldConfigArgument{
					"comment_id": &graphql.ArgumentConfig{
						Type:        int64GraphqlScalar,
						Description: "comment_id",
					},
					"comment": &graphql.ArgumentConfig{
						Type:        graphql.String,
						Description: "comment",
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					user, err := parseAuth(p)
					if err != nil {
						return nil, err
					}
					commentID, isOK := p.Args["comment_id"].(int64)
					if !isOK {
						return nil, errors.New("comment_id format")
					}
					comment, err := parseComment(p, commentID, 0, user.UserID)
					if err != nil {
						return comment, err
					}
					us, err := commentOnPostUpdate(comment)
					return us, err
				},
				Description: "",
			},
			"comment_on_post_delete": &graphql.Field{
				Type: updateStatusGraphqlType,
				Args: graphql.FieldConfigArgument{
					"comment_id": &graphql.ArgumentConfig{
						Type:        int64GraphqlScalar,
						Description: "comment_id",
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					user, err := parseAuth(p)
					if err != nil {
						return nil, err
					}
					commentID, isOK := p.Args["comment_id"].(int64)
					if !isOK {
						return nil, errors.New("comment_id format")
					}
					comment := commentAPI{
						CommentID: commentID,
						User:      userBasicAPI{UserID: user.UserID},
					}
					us, err := commentOnPostDelete(comment)
					return us, err
				},
				Description: "",
			},
			"reply_on_comment_insert": &graphql.Field{
				Type: replyGraphqlType,
				Args: graphql.FieldConfigArgument{
					"comment_id": &graphql.ArgumentConfig{
						Type:        int64GraphqlScalar,
						Description: "comment_id",
					},
					"reply": &graphql.ArgumentConfig{
						Type:        graphql.String,
						Description: "reply",
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					user, err := parseAuth(p)
					if err != nil {
						return nil, err
					}
					commentID, isOK := p.Args["comment_id"].(int64)
					if !isOK {
						return nil, err
					}
					reply, err := parseReply(p, 0, commentID, user.UserID)
					if err != nil {
						return reply, err
					}
					reply.ReplyID, err = replyOnCommentInsert(reply)
					return reply, err
				},
				Description: "",
			},
			"reply_on_comment_update": &graphql.Field{
				Type: updateStatusGraphqlType,
				Args: graphql.FieldConfigArgument{
					"reply_id": &graphql.ArgumentConfig{
						Type:        int64GraphqlScalar,
						Description: "reply_id",
					},
					"reply": &graphql.ArgumentConfig{
						Type:        graphql.String,
						Description: "reply",
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					user, err := parseAuth(p)
					if err != nil {
						return nil, err
					}
					replyID, isOK := p.Args["reply_id"].(int64)
					if !isOK {
						return nil, errors.New("reply_id format")
					}
					reply, err := parseReply(p, replyID, 0, user.UserID)
					if err != nil {
						return reply, err
					}
					us, err := replyOnCommentUpdate(reply)
					return us, err
				},
				Description: "",
			},
			"reply_on_comment_delete": &graphql.Field{
				Type: updateStatusGraphqlType,
				Args: graphql.FieldConfigArgument{
					"reply_id": &graphql.ArgumentConfig{
						Type:        int64GraphqlScalar,
						Description: "reply_id",
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					user, err := parseAuth(p)
					if err != nil {
						return nil, err
					}
					replyID, isOK := p.Args["reply_id"].(int64)
					if !isOK {
						return nil, errors.New("reply_id format")
					}
					reply := replyAPI{
						ReplyID: replyID,
						User:    userBasicAPI{UserID: user.UserID},
					}
					us, err := replyOnCommentDelete(reply)
					return us, err
				},
				Description: "",
			},
		},
	},
)
