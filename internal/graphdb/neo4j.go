package graphdb

import (
	"context"
	"fmt"
	"log"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

var Neo4jDriver neo4j.DriverWithContext

func InitializeNeo4j() error {
	var err error
	Neo4jDriver, err = neo4j.NewDriverWithContext(
		NEO4J_URI,
		neo4j.BasicAuth(NEO4J_USERNAME, NEO4J_PASSWORD, ""),
	)
	if err != nil {
		return fmt.Errorf("failed to create Neo4j driver: %w", err)
	}

	err = Neo4jDriver.VerifyConnectivity(context.Background())
	if err != nil {
		return fmt.Errorf("failed to connect to Neo4j: %w", err)
	}

	log.Println("✅ Connected to Neo4j!")
	return nil
}

func CloseNeo4j() {
	if Neo4jDriver != nil {
		Neo4jDriver.Close(context.Background())
		log.Println("🛑 Neo4j connection closed.")
	}
}
