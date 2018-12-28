package main

import (
	"github.com/graphql-go/graphql"
	"github.com/ivahaev/go-logger"
)

func executeQuery(query string, schema graphql.Schema) *graphql.Result {
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
	})
	if len(result.Errors) > 0 {
		logger.Error("errors: %v", result.Errors)
	}
	return result
}

var schema, _ = graphql.NewSchema(
	graphql.SchemaConfig{
		Query:    queryType,
		Mutation: mutationType,
	},
)

var queryType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			/* Получение продукта по ID
			   /api/v1/images?query={image(id:1){url,catalog}}
			*/
			"image": &graphql.Field{
				Type:        imageType,
				Description: "Get images by id",
				Args: graphql.FieldConfigArgument{
					"ID": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					id, ok := p.Args["ID"].(string)
					if ok {
						imgs, _ := GetDB(id, false)
						return imgs, nil
					}
					return nil, nil
				},
			},
			/*
			   /api/v1/images?query={images(Catalog:"city"){url,catalog}}
			*/
			"images": &graphql.Field{
				Type:        graphql.NewList(imageType),
				Description: "Get images by catalog",
				Args: graphql.FieldConfigArgument{
					"Catalog": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					catalog, ok := p.Args["Catalog"].(string)
					if ok {
						if catalog == "All" {
							catalogs, err := GetCatalogDB(false)
							return catalogs, err
						} else {
							imgs, err := GetImageCatalogDB(catalog, false)
							return imgs, err
						}

					}
					return nil, nil
				},
			},
			/* Получение списка продуктов
			   /api/v1/images?query={list(count:10){url,catalog}}
			*/
			"list": &graphql.Field{
				Type:        graphql.NewList(imageType),
				Description: "Get Images list",
				Args: graphql.FieldConfigArgument{
					"count": &graphql.ArgumentConfig{
						Type: graphql.Int,
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					count, countOk := params.Args["count"].(int)
					if !countOk {
						imgs, _ := GetsDB(false)
						return imgs, nil
					} else {
						imgs, _ := GetsCountDB(false, count)
						return imgs, nil
					}
				},
			},
		},
	})

var mutationType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Mutation",
	Fields: graphql.Fields{
		/*
		  /api/v1/images?query=mutation+_{create(id:asd123,Url:adad){url,catalog}}
		*/
		"create": &graphql.Field{
			Type:        imageType,
			Description: "Create new image",
			Args: graphql.FieldConfigArgument{
				"ID": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"Url": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"Catalog": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"OriginPath": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"IconPath": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				img := Image{
					ID:      uuID(),
					Url:     params.Args["Url"].(string),
					Catalog: params.Args["Catalog"].(string),
				}
				AddDB(img)
				return img, nil
			},
		},

		/*
		   /api/v1/images?query=mutation+_{update(Id:1,Url:195){url,catalog}}
		*/
		"update": &graphql.Field{
			Type:        imageType,
			Description: "Update product by id",
			Args: graphql.FieldConfigArgument{
				"ID": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"Url": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"Catalog": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"OriginPath": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"IconPath": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				id, _ := params.Args["ID"].(string)
				url, _ := params.Args["Url"].(string)
				catalog, _ := params.Args["Catalog"].(string)
				originPath, _ := params.Args["OriginPath"].(string)
				iconPath, _ := params.Args["IconPath"].(string)
				img := Image{
					ID:         id,
					Url:        url,
					Catalog:    catalog,
					OriginPath: originPath,
					IconPath:   iconPath,
				}
				UpdateDB(id, img)
				return img, nil
			},
		},

		/*
		   /api/v1/images?query=mutation+_{delete(id:1){url,catalog}}
		   /api/v1/images?query=mutation+_{delete(Url:asdas12){url,catalog}}
		   /api/v1/images?query=mutation+_{delete(All:true){url,catalog}}
		*/
		"delete": &graphql.Field{
			Type:        imageType,
			Description: "Delete product by id",
			Args: graphql.FieldConfigArgument{
				"ID": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"Url": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"Catalog": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"All": &graphql.ArgumentConfig{
					Type: graphql.Boolean,
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				id, idOK := params.Args["ID"].(string)
				url, urlOK := params.Args["Url"].(string)
				catalog, catalogOK := params.Args["Catalog"].(string)
				all, allOK := params.Args["All"].(bool)
				//img := Image{}
				if idOK {
					DeleteDB(id)
					return nil, nil
				}
				if urlOK {
					DeleteUrlDB(url)
					return nil, nil
				}
				if catalogOK {
					DeleteCatalogDB(catalog)
					return nil, nil
				}
				if allOK {
					if all {
						DeleteAllDB()
						return nil, nil
					}
				}
				return nil, nil
			},
		},
	},
})

var imageType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Image",
		Fields: graphql.Fields{
			"ID": &graphql.Field{
				Type: graphql.String,
			},
			"Url": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
			},
			"Catalog": &graphql.Field{
				Type: graphql.String,
			},
			"OriginPath": &graphql.Field{
				Type: graphql.String,
			},
			"IconPath": &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)
