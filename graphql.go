package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"time"

	"github.com/graphql-go/graphql"
)

var xuserGraphqlType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "XUser",
		Fields: graphql.Fields{
			"user_id":     &graphql.Field{Type: graphql.String},
			"username":    &graphql.Field{Type: graphql.String},
			"email":       &graphql.Field{Type: graphql.String},
			"password":    &graphql.Field{Type: graphql.String},
			"name":        &graphql.Field{Type: graphql.String},
			"phone":       &graphql.Field{Type: graphql.String},
			"gender":      &graphql.Field{Type: graphql.Int},
			"bio":         &graphql.Field{Type: graphql.String},
			"credit":      &graphql.Field{Type: graphql.Int},
			"language_id": &graphql.Field{Type: graphql.Int},
			"country_id":  &graphql.Field{Type: graphql.Int},
			"timezone":    &graphql.Field{Type: graphql.Int},
			"last_ip":     &graphql.Field{Type: graphql.String},
			"createtime":  &graphql.Field{Type: graphql.Int},
			"updatetime":  &graphql.Field{Type: graphql.Int},
		},
	},
)

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
				Type: xuserGraphqlType,
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
					user := login(email, password)
					return user, nil
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
						return nil, nil
					}
					user := getXuserByID(userID)
					return user, nil
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
						return nil, nil
					}
					user := getXuserByUsername(username)
					return user, nil
				},
				Description: "get xuser",
			},
			"posts": &graphql.Field{
				Type: postsGraphqlType,
				Args: graphql.FieldConfigArgument{
					"category_id": &graphql.ArgumentConfig{
						Type:        graphql.Int,
						Description: "category_id",
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					categoryID, isOK := p.Args["category_id"].(int)
					if !isOK {
						return nil, nil
					}
					posts := getPost(categoryID)
					return posts, nil
				},
				Description: "posts",
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
					"user_id": &graphql.ArgumentConfig{
						Type:        graphql.String,
						Description: "user_id",
					},
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
					file, isOK := p.Context.Value(contextKeyFile).(io.Reader)
					if !isOK {
						return nil, errors.New("file err")
					}
					// size check
					// log.Println(file)
					post, err := postInputCheck(p)
					if err != nil {
						return nil, err
					}
					err = untarFileAndUpload(post, file)
					log.Panicln("untar", err)
					// fmt.Println(post)
					// post.PostID = postNow(post)
					return postDB{PostID: "123"}, nil
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
	userID, isOK := p.Args["user_id"].(string)
	if !isOK {
		return post, errors.New("user_id err")
	}
	content, isOK := p.Args["content"].(string)
	if !isOK {
		return post, errors.New("content err")
	}
	blobID, isOK := p.Args["blob_id"].(string)
	if !isOK {
		return post, errors.New("blob_id err")
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
	post.UserID = userID
	post.Content = content
	post.BlobID = blobID
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
