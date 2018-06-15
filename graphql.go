package main

import (
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

// graphql type

// common
var deleteStatusGraphqlType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "delete_status",
		Fields: graphql.Fields{
			"rows_affected": &graphql.Field{Type: graphql.Int},
		},
	},
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
var xuserGraphqlType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "XUser",
		Fields: graphql.Fields{
			"user_id":     &graphql.Field{Type: int64GraphqlScalar},
			"username":    &graphql.Field{Type: graphql.String},
			"email":       &graphql.Field{Type: graphql.String},
			"name":        &graphql.Field{Type: graphql.String},
			"phone":       &graphql.Field{Type: graphql.String},
			"gender":      &graphql.Field{Type: graphql.Int},
			"bio":         &graphql.Field{Type: graphql.String},
			"credit":      &graphql.Field{Type: graphql.Int},
			"photo_url":   &graphql.Field{Type: graphql.String},
			"language_id": &graphql.Field{Type: graphql.Int},
			"country_id":  &graphql.Field{Type: graphql.Int},
			"timezone":    &graphql.Field{Type: graphql.Int},
			"last_ip":     &graphql.Field{Type: graphql.String},
			"createtime":  &graphql.Field{Type: graphql.Int},
			"updatetime":  &graphql.Field{Type: graphql.Int},
		},
	},
)
var xusersGraphqlType = graphql.NewList(xuserGraphqlType)

// following
var xuserFollowingGraphqlType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "following_user",
		Fields: graphql.Fields{
			"user_id":        &graphql.Field{Type: int64GraphqlScalar},
			"username":       &graphql.Field{Type: graphql.String},
			"name":           &graphql.Field{Type: graphql.String},
			"photo_url":      &graphql.Field{Type: graphql.String},
			"following_time": &graphql.Field{Type: graphql.Int},
		},
	},
)
var xusersFollowingGraphqlType = graphql.NewList(xuserFollowingGraphqlType)

// follower
var xuserFollowerGraphqlType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "follower_user",
		Fields: graphql.Fields{
			"user_id":        &graphql.Field{Type: int64GraphqlScalar},
			"username":       &graphql.Field{Type: graphql.String},
			"name":           &graphql.Field{Type: graphql.String},
			"photo_url":      &graphql.Field{Type: graphql.String},
			"following_time": &graphql.Field{Type: graphql.Int},
		},
	},
)
var xusersFollowerGraphqlType = graphql.NewList(xuserFollowerGraphqlType)

var followingGraphqlType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "following",
		Fields: graphql.Fields{
			"user_id":        &graphql.Field{Type: graphql.String},
			"following_time": &graphql.Field{Type: graphql.Int},
		},
	},
)

// post
var postGraphqlType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "post",
		Fields: graphql.Fields{
			"post_id":       &graphql.Field{Type: int64GraphqlScalar},
			"user_id":       &graphql.Field{Type: int64GraphqlScalar},
			"username":      &graphql.Field{Type: graphql.String},
			"name":          &graphql.Field{Type: graphql.String},
			"content":       &graphql.Field{Type: graphql.String},
			"blob_id":       &graphql.Field{Type: graphql.String},
			"type":          &graphql.Field{Type: graphql.Int, Description: "media type"},
			"like_count":    &graphql.Field{Type: int64GraphqlScalar},
			"dislike_count": &graphql.Field{Type: int64GraphqlScalar},
			"comment_count": &graphql.Field{Type: int64GraphqlScalar},
			"country_id":    &graphql.Field{Type: graphql.Int},
			"category_id":   &graphql.Field{Type: graphql.Int},
			"createtime":    &graphql.Field{Type: graphql.Int, Description: "Unix Timestamp"},
			"updatetime":    &graphql.Field{Type: graphql.Int, Description: "Unix Timestamp"},
		},
	},
)
var postsGraphqlType = graphql.NewList(postGraphqlType)

// comments
var commentGraphqlType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "comment",
		Fields: graphql.Fields{
			"comment_id":    &graphql.Field{Type: int64GraphqlScalar},
			"post_id":       &graphql.Field{Type: int64GraphqlScalar},
			"user_id":       &graphql.Field{Type: int64GraphqlScalar},
			"username":      &graphql.Field{Type: graphql.String},
			"name":          &graphql.Field{Type: graphql.String},
			"comment":       &graphql.Field{Type: graphql.String},
			"like_count":    &graphql.Field{Type: int64GraphqlScalar},
			"dislike_count": &graphql.Field{Type: int64GraphqlScalar},
			"comment_count": &graphql.Field{Type: int64GraphqlScalar},
			"createtime":    &graphql.Field{Type: graphql.Int, Description: "Unix Timestamp"},
			"updatetime":    &graphql.Field{Type: graphql.Int, Description: "Unix Timestamp"},
		},
	},
)
var commentsGraphqlType = graphql.NewList(commentGraphqlType)

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
func parsePost(p graphql.ResolveParams, userID int64) (post postDB, err error) {
	content, isOK := p.Args["content"].(string)
	if !isOK {
		return post, errors.New("content format")
	}
	countryID, isOK := p.Args["country_id"].(int)
	if !isOK {
		return post, errors.New("country_id format")
	}
	categoryID, isOK := p.Args["category_id"].(int)
	if !isOK {
		return post, errors.New("category_id format")
	}
	public, isOK := p.Args["public"].(bool)
	if !isOK {
		return post, errors.New("public format")
	}
	Type, isOK := p.Args["type"].(int)
	if !isOK {
		return post, errors.New("type format")
	}
	timestamp := getNowUnixTimestamp()
	post.Content = content

	post.BlobID = strconv.FormatInt(userID, 10) + "_" + strconv.Itoa(timestamp)
	// Point
	post.Type = Type
	post.CommentCount = 0
	post.LikeCount = 0
	post.DislikeCount = 0
	post.CountryID = countryID
	post.CategoryID = categoryID
	post.Public = public
	post.Createtime = timestamp
	post.Updatetime = timestamp
	return post, nil
}
func parseComment(p graphql.ResolveParams) (c commentAPI, err error) {
	postID, isOK := p.Args["post_id"].(int64)
	if !isOK {
		return c, errors.New("post_id format")
	}
	comment, isOK := p.Args["comment"].(string)
	if !isOK {
		return c, errors.New("comment format")
	}
	timestamp := getNowUnixTimestamp()
	c.PostID = postID
	c.Comment = comment
	c.LikeCount = 0
	c.DislikeCount = 0
	c.CommentCount = 0
	c.Createtime = timestamp
	c.Updatetime = timestamp
	return c, nil
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
			"xuser_by_user_id": &graphql.Field{
				Type: xuserGraphqlType,
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
					user, err := getXuserByID(userID)
					return user, err
				},
				Description: "",
			},
			"xuser_by_username": &graphql.Field{
				Type: xuserGraphqlType,
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
					user, err := getXuserByUsername(username)
					return user, err
				},
				Description: "",
			},
			"following_list": &graphql.Field{
				Type: xusersFollowingGraphqlType,
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
					users, err := getFollowingList(user.UserID, page)
					return users, err
				},
				Description: "",
			},
			"follower_list": &graphql.Field{
				Type: xusersFollowerGraphqlType,
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
					users, err := getFollwerList(user.UserID, page)
					return users, err
				},
				Description: "",
			},
			"is_following": &graphql.Field{
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
					isFollowing, err := checkIfFollowing(userID, user.UserID)
					return isFollowing, err
				},
				Description: "",
			},
			"posts_following": &graphql.Field{
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
					posts := getFollowingUsersPosts(user.UserID, page)
					return posts, nil
				},
				Description: "",
			},
			"posts_recent": &graphql.Field{
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
					posts := getPostsRecent(categoryID, page)
					return posts, nil
				},
				Description: "",
			},
			"comments": &graphql.Field{
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
					comments, err := getCommentsPost(postID, page)
					return comments, err
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
				Type: followingGraphqlType,
				Args: graphql.FieldConfigArgument{
					"user_id": &graphql.ArgumentConfig{
						Type:        int64GraphqlScalar,
						Description: "user_id",
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					userToken, isOK := p.Context.Value(contextUserToken).(string)
					if !isOK {
						return nil, errors.New("user-token format")
					}
					user, err := checkSession(userToken)
					if err != nil {
						return nil, errors.New("user-token invalid")
					}
					userID, isOK := p.Args["user_id"].(int64)
					if !isOK {
						return nil, errors.New("user_id format")
					}
					userFollowing, err := follow(userID, user.UserID)
					return userFollowing, err
				},
				Description: "",
			},
			"unfollow": &graphql.Field{
				Type: deleteStatusGraphqlType,
				Args: graphql.FieldConfigArgument{
					"user_id": &graphql.ArgumentConfig{
						Type:        int64GraphqlScalar,
						Description: "user_id",
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					userToken, isOK := p.Context.Value(contextUserToken).(string)
					if !isOK {
						return nil, errors.New("user-token format")
					}
					user, err := checkSession(userToken)
					if err != nil {
						return nil, errors.New("user-token invalid")
					}
					userID, isOK := p.Args["user_id"].(int64)
					if !isOK {
						return nil, errors.New("user_id format")
					}
					ds, err := unfollow(userID, user.UserID)
					return ds, err
				},
				Description: "",
			},
			"post": &graphql.Field{
				Type: postGraphqlType,
				Args: graphql.FieldConfigArgument{
					"content": &graphql.ArgumentConfig{
						Type:        graphql.String,
						Description: "content",
					},
					"blob_id": &graphql.ArgumentConfig{
						Type:        graphql.String,
						Description: "blob_id",
					},
					"country_id": &graphql.ArgumentConfig{
						Type:        graphql.Int,
						Description: "country_id",
					},
					"category_id": &graphql.ArgumentConfig{
						Type:        graphql.Int,
						Description: "category_id",
					},
					"public": &graphql.ArgumentConfig{
						Type:        graphql.Boolean,
						Description: "public",
					},
					"type": &graphql.ArgumentConfig{
						Type:        graphql.Int,
						Description: "media type",
					},
					"file": &graphql.ArgumentConfig{
						Type:        graphql.String,
						Description: "tar.gz[*.jpg,...]/[*.m3u8(only one), *.ts...]",
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
					post, err := parsePost(p, user.UserID)
					if err != nil {
						return nil, err
					}
					// size check
					err = untarFileAndUpload(post, file)
					if err != nil {
						return nil, err
					}
					post.PostID = postNow(post)
					log.Printf("post now total took %fs\n", time.Since(startTime).Seconds())
					return post, nil
				},
				Description: "",
				DeprecationReason: `please use form-data to upload file, form-data key:
						query: mutation{post(...:...){post_id}}
						file: tar.gz file,
						. not finished yet
					`,
			},
			"comment": &graphql.Field{
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
					comment, err := parseComment(p)
					comment.UserID = user.UserID
					// block check
					comment.CommentID, err = commentNow(comment)
					return comment, err
				},
				Description: "",
			},
		},
	},
)
