package service

import (
	"encoding/json"
	"github.com/graphql-go/graphql"
	"kaleido/master/model/Mirror"
	"kaleido/master/model/MirrorStation"
	"net/http"
	"strconv"
)

var mirrorType = graphql.NewObject(graphql.ObjectConfig{
	Name:        "Mirror",
	Description: "一个镜像",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type:        graphql.NewNonNull(graphql.ID),
			Description: "镜像ID",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if mirror, ok := p.Source.(Mirror.Mirror); ok {
					return mirror.Id, nil
				}
				return nil, nil
			},
		},
		"name": &graphql.Field{
			Type:        graphql.String,
			Description: "镜像名称",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if mirror, ok := p.Source.(Mirror.Mirror); ok {
					return mirror.GetName()
				}
				return nil, nil
			},
		},
	},
})

var mirrorStationType = graphql.NewObject(graphql.ObjectConfig{
	Name:        "MirrorStation",
	Description: "一个镜像站",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type:        graphql.NewNonNull(graphql.ID),
			Description: "镜像站ID",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if mirrorStation, ok := p.Source.(MirrorStation.MirrorStation); ok {
					return mirrorStation.GetId(), nil
				}
				return nil, nil
			},
		},
		"name": &graphql.Field{
			Type:        graphql.String,
			Description: "镜像站名称",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if mirrorStation, ok := p.Source.(MirrorStation.MirrorStation); ok {
					return mirrorStation.GetName()
				}
				return nil, nil
			},
		},
		"mirrors": &graphql.Field{
			Type:        graphql.NewList(mirrorType),
			Description: "镜像站中的镜像",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if mirrorStation, ok := p.Source.(MirrorStation.MirrorStation); ok {
					return mirrorStation.GetMirrors()
				}
				return nil, nil
			},
		},
	},
})

func init() {
	mirrorType.AddFieldConfig("mirrorStations", &graphql.Field{
		Type:        graphql.NewList(mirrorStationType),
		Description: "镜像站",
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			if mirror, ok := p.Source.(Mirror.Mirror); ok {
				return MirrorStation.GetStations(mirror)
			}
			return nil, nil
		},
	})
}

var queryType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Query",
	Fields: graphql.Fields{
		"mirrorStations": &graphql.Field{
			Description: "查询所有镜像站",
			Type:        graphql.NewList(mirrorStationType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return MirrorStation.All()
			},
		},
		"mirrorStation": &graphql.Field{
			Description: "查询镜像站信息",
			Type:        mirrorStationType,
			Args: graphql.FieldConfigArgument{
				"ID": &graphql.ArgumentConfig{
					Description: "镜像站ID",
					Type:        graphql.NewNonNull(graphql.ID),
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				id, _ := strconv.ParseUint(p.Args["ID"].(string), 10, 64)
				return MirrorStation.Get(id)
			},
		},
		"mirrors": &graphql.Field{
			Description: "查询所有镜像",
			Type:        graphql.NewList(mirrorType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return Mirror.All()
			},
		},
		"mirror": &graphql.Field{
			Description: "查询镜像信息",
			Type:        mirrorType,
			Args: graphql.FieldConfigArgument{
				"ID": &graphql.ArgumentConfig{
					Description: "镜像ID",
					Type:        graphql.NewNonNull(graphql.ID),
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				id, _ := strconv.ParseUint(p.Args["ID"].(string), 10, 64)
				return Mirror.Mirror{Id: id}, nil
			},
		},
	},
})

var kaleidoSchema, _ = graphql.NewSchema(graphql.SchemaConfig{
	Query: queryType,
})

func StartGraphQLServer() {
	http.HandleFunc("/graphql", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		query := r.URL.Query().Get("query")
		result := graphql.Do(graphql.Params{
			Schema:        kaleidoSchema,
			RequestString: query,
		})
		json.NewEncoder(w).Encode(result)
	})
	http.ListenAndServe(":8086", nil)
}
