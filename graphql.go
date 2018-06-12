package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"strconv"
	"time"

	"github.com/graphql-go/graphql"
)

var loginGraphqlType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "login",
		Fields: graphql.Fields{
			"token": &graphql.Field{Type: graphql.String},
		},
	},
)

var xuserGraphqlType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "XUser",
		Fields: graphql.Fields{
			"user_id":     &graphql.Field{Type: graphql.String},
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

var xuserFollowingGraphqlType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "following_user",
		Fields: graphql.Fields{
			"user_id":        &graphql.Field{Type: graphql.String},
			"username":       &graphql.Field{Type: graphql.String},
			"name":           &graphql.Field{Type: graphql.String},
			"photo_url":      &graphql.Field{Type: graphql.String},
			"following_time": &graphql.Field{Type: graphql.Int},
		},
	},
)
var xusersFollowingGraphqlType = graphql.NewList(xuserFollowingGraphqlType)
var xuserFollowerGraphqlType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "follower_user",
		Fields: graphql.Fields{
			"user_id":        &graphql.Field{Type: graphql.String},
			"username":       &graphql.Field{Type: graphql.String},
			"name":           &graphql.Field{Type: graphql.String},
			"photo_url":      &graphql.Field{Type: graphql.String},
			"following_time": &graphql.Field{Type: graphql.Int},
		},
	},
)
var xusersFollowerGraphqlType = graphql.NewList(xuserFollowerGraphqlType)

var postGraphqlType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Posts",
		Fields: graphql.Fields{
			"post_id":     &graphql.Field{Type: graphql.String},
			"user_id":     &graphql.Field{Type: graphql.String},
			"content":     &graphql.Field{Type: graphql.String},
			"blob_id":     &graphql.Field{Type: graphql.String},
			"country_id":  &graphql.Field{Type: graphql.Int},
			"category_id": &graphql.Field{Type: graphql.Int},
			"public":      &graphql.Field{Type: graphql.Boolean},
			"type":        &graphql.Field{Type: graphql.Int, Description: "0: image, 1: video"}, // config
			"createtime":  &graphql.Field{Type: graphql.Int, Description: "Unix Timestamp"},
			"updatetime":  &graphql.Field{Type: graphql.Int, Description: "Unix Timestamp"},
		},
	},
)
var postsGraphqlType = graphql.NewList(postGraphqlType)

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
				Description: "login",
			},
			"xuser_by_user_id": &graphql.Field{
				Type: xuserGraphqlType,
				Args: graphql.FieldConfigArgument{
					"user_id": &graphql.ArgumentConfig{
						Type:        graphql.String,
						Description: "user_id",
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					userID, isOK := p.Args["user_id"].(string)
					if !isOK {
						return nil, errors.New("user_id format")
					}
					user, err := getXuserByID(userID)
					return user, err
				},
				Description: "get xuser",
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
				Description: "get xuser",
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
					page, isOK := p.Args["page"].(int)
					if !isOK {
						return nil, errors.New("page format")
					}
					userToken, isOK := p.Context.Value(contextUserToken).(string)
					if !isOK {
						return nil, errors.New("user-token format")
					}
					user, err := checkSession(userToken)
					if err != nil {
						return nil, errors.New("user-token invalid")
					}
					users, err := getFollowingList(user.UserID, page)
					return users, err
				},
				Description: "get following xuser list",
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
					page, isOK := p.Args["page"].(int)
					if !isOK {
						return nil, errors.New("page format")
					}
					userToken, isOK := p.Context.Value(contextUserToken).(string)
					if !isOK {
						return nil, errors.New("user-token format")
					}
					user, err := checkSession(userToken)
					if err != nil {
						return nil, errors.New("user-token invalid")
					}
					users, err := getFollwerList(user.UserID, page)
					return users, err
				},
				Description: "get follower xuser list",
			},
			"posts_following": &graphql.Field{
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
						return nil, nil
					}
					page, isOK := p.Args["page"].(int)
					if !isOK {
						return nil, nil
					}
					posts := getFollowingUsersPosts(categoryID, page)
					return posts, nil
				},
				Description: "posts",
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
						return nil, nil
					}
					page, isOK := p.Args["page"].(int)
					if !isOK {
						return nil, nil
					}
					posts := getPostsRecent(categoryID, page)
					return posts, nil
				},
				Description: "posts_recent",
			},
		},
	})
var graphqlMutationType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			"post_now": &graphql.Field{
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
						Description: "0: image, 1: video",
					},
					"file": &graphql.ArgumentConfig{
						Type:        graphql.String,
						Description: "tar.gz[*.jpg,...]/[*.m3u8(only one), *.ts...]",
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					startTime := time.Now()
					userToken, isOK := p.Context.Value(contextUserToken).(string)
					if !isOK {
						return nil, errors.New("user-token format")
					}
					file, isOK := p.Context.Value(contextKeyFile).(io.Reader)
					if !isOK {
						return nil, errors.New("file format")
					}
					user, err := checkSession(userToken)
					if err != nil {
						return nil, errors.New("user-token invalid")
					}
					// post parameter check
					post, err := postInputCheck(p)
					post.UserID = user.UserID
					if err != nil {
						return nil, err
					}
					// size check
					err = untarFileAndUpload(post, file)
					// log.Println(err)
					if err != nil {
						return nil, err
					}
					post.PostID = postNow(post)
					log.Printf("post now total took %fs\n", time.Since(startTime).Seconds())
					return postDB{}, nil
				},
				Description: "post",
				DeprecationReason: `please use form-data to upload file, form-data key:
					query: mutation{post(...:...){post_id}}
					file: tar.gz file, 
					. not finished yet
				`,
			},
		},
	},
)

func postInputCheck(p graphql.ResolveParams) (post postDB, err error) {
	content, isOK := p.Args["content"].(string)
	if !isOK {
		return post, errors.New("content err")
	}
	countryID, isOK := p.Args["country_id"].(int)
	if !isOK {
		return post, errors.New("country_id err")
	}
	categoryID, isOK := p.Args["category_id"].(int)
	if !isOK {
		return post, errors.New("category_id err")
	}
	public, isOK := p.Args["public"].(bool)
	if !isOK {
		return post, errors.New("public err")
	}
	Type, isOK := p.Args["type"].(int)
	if !isOK {
		return post, errors.New("public err")
	}
	timestamp := int(time.Now().Unix())
	post.Content = content
	post.BlobID = post.UserID + "_" + strconv.Itoa(timestamp)
	// Point
	post.CountryID = countryID
	post.CategoryID = categoryID
	post.Public = public
	post.Type = Type
	post.Createtime = timestamp
	post.Updatetime = timestamp
	return post, nil
}

var graphqlSchema, _ = graphql.NewSchema(
	graphql.SchemaConfig{
		Query:    graphqlQueryType,
		Mutation: graphqlMutationType,
	},
)

func importJSONDataFromGraphqlFile(fileName string, result interface{}) (isOK bool) {
	isOK = true
	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Print("Error:", err)
		isOK = false
	}
	err = json.Unmarshal(content, result)
	if err != nil {
		isOK = false
		fmt.Print("Error:", err)
	}
	return
}
