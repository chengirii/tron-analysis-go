package config

type Server struct {
	Neo4j   Neo4j   `mapstructure:"database-neo4j" json:"database-neo4j"`
	MongoDB MongoDB `mapstructure:"database-mongodb" json:"database-mongodb"`
	Zap     Zap     `mapstructure:"zap" json:"zap"`
	Tron    Tron    `mapstructure:"tron" json:"tron"`
}
