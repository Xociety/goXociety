package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"strconv"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"
)

type user struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func coerceInt(value interface{}) interface{} {
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
		return coerceInt(*value)
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
		if value < int64(math.MinInt32) || value > int64(math.MaxInt32) {
			return nil
		}
		return int(value)
	case *int64:
		return coerceInt(*value)
	case uint:
		if value > math.MaxInt32 {
			return nil
		}
		return int(value)
	case *uint:
		return coerceInt(*value)
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
		return coerceInt(*value)
	case uint64:
		if value > uint64(math.MaxInt32) {
			return nil
		}
		return int(value)
	case *uint64:
		return coerceInt(*value)
	case float32:
		if value < float32(math.MinInt32) || value > float32(math.MaxInt32) {
			return nil
		}
		return int(value)
	case *float32:
		return coerceInt(*value)
	case float64:
		if value < float64(math.MinInt32) || value > float64(math.MaxInt32) {
			return nil
		}
		return int(value)
	case *float64:
		return coerceInt(*value)
	case string:
		val, err := strconv.ParseFloat(value, 0)
		if err != nil {
			return nil
		}
		return coerceInt(val)
	case *string:
		return coerceInt(*value)
	}

	// If the value cannot be transformed into an int, return nil instead of '0'
	// to denote 'no integer found'
	return nil
}

var graphqlScalarMap = graphql.NewScalar(graphql.ScalarConfig{
	Name:        "JSON",
	Description: "custom json",
	Serialize:   func(value interface{}) interface{} { return value },
	ParseValue:  func(value interface{}) interface{} { return value },
	ParseLiteral: func(valueAST ast.Value) interface{} {
		fmt.Println(valueAST)
		return valueAST
		// return nil
	},
})
var graphqlScalarInt64 = graphql.NewScalar(graphql.ScalarConfig{
	Name:        "Int64",
	Description: "custom int64",
	Serialize:   coerceInt,
	ParseValue:  coerceInt,
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
var graphqlSampleData map[string]user
var userGraphqlType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "User",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.String,
			},
			"name": &graphql.Field{
				Type: graphql.String,
			},
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
			"createtime":  &graphql.Field{Type: graphql.Int, Description: "Unix Timestamp"},
			"updatetime":  &graphql.Field{Type: graphql.Int, Description: "Unix Timestamp"},
		},
	},
)
var postsGraphqlType = graphql.NewList(postGraphqlType)

var usersGraphqlType = graphql.NewList(userGraphqlType)
var graphqlIdsArrayType = graphql.NewList(graphql.String)
var graphqlFilterArrayType = graphql.NewList(graphql.String)
var resultGraphqlType = graphql.NewList(graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Result",
		Fields: graphql.Fields{
			"fivemin": &graphql.Field{
				Type: graphqlScalarMap,
			},
			"withLabel": &graphql.Field{
				Type: graphqlScalarInt64,
			},
			"groupID": &graphql.Field{
				Type: graphql.String,
			},
			"groupName": &graphql.Field{
				Type: graphql.String,
			},
		},
	},
))
var graphqlQueryType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"user": &graphql.Field{
				Type: userGraphqlType,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type:        graphql.String,
						Description: "search for id",
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					idQuery, isOK := p.Args["id"].(string)
					if isOK {
						return graphqlSampleData[idQuery], nil
					}
					return nil, nil
				},
				Description: "get certain user by id",
			},
			"allUsers": &graphql.Field{
				Type: usersGraphqlType,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					var users []user
					for _, v := range graphqlSampleData {
						u := user{ID: v.ID, Name: v.Name}
						users = append(users, u)
					}
					return users, nil
				},
				Description: "get all users",
			},
			"users": &graphql.Field{
				Type: usersGraphqlType,
				Args: graphql.FieldConfigArgument{
					"ids": &graphql.ArgumentConfig{
						Type:        graphqlIdsArrayType,
						Description: "search for ids",
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					ids, _ := p.Args["ids"].([]interface{})
					var users []user
					for _, v := range graphqlSampleData {
						u := user{ID: v.ID, Name: v.Name}
						for l := 0; l < len(ids); l++ {
							if ids[l] == u.ID {
								users = append(users, u)
								break
							}
						}
					}
					return users, nil
				},
				Description: "get certain user by ids",
			},
			//
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
			/*"result": &graphql.Field{
				Type: resultGraphqlType,
				Args: graphql.FieldConfigArgument{
					"Interest": &graphql.ArgumentConfig{
						Type:        graphqlFilterArrayType,
						Description: "Filter Interest",
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					interest, _ := p.Args["Interest"].([]interface{})
					filter := FilterRequest{Filter: make(map[string][]string)}
					for l := 0; l < len(interest); l++ {
						// str := fmt.Sprintf("%v", interest[l])
						str, _ := (interest[l]).(string)
						filter.Filter["Interest"] = append(filter.Filter["Interest"], str)
					}
					var result []Result
					for i := 0; i < 1; i++ {
						f := make(map[string]int64)
						f["123"] = 123
						f["456"] = 456
						r := Result{Fivemin: f, WithLabel: 123, Name: "123"}
						result = append(result, r)
					}
					return result, nil
				},
				Description: "Project result",
			},*/
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
					"file_origin": &graphql.ArgumentConfig{
						Type:        graphql.String,
						Description: "form-data key\"file_origin\"",
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					content, isOK := p.Args["content"].(string)
					if !isOK {
						return nil, nil
					}
					fileOrigin, isOK := p.Context.Value(contextKeyFileOrigin).([]byte)
					if !isOK {
						return nil, errors.New("file err")
					}
					log.Println(content, len(fileOrigin))
					post := postDB{}
					postID := postNow(post)
					return postID, nil
				},
				Description: "post",
				DeprecationReason: `please use form-data to upload file, form-data key:
					query: mutation{post(content:"45"){post_id}}
					file_origin: file, 
					. not finished yet
				`,
			},
		},
	},
)

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
