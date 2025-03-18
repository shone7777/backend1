package graphdb

import (
	"context"
	"fmt"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func GetHopCount(startUserID, endUserID int64) (int, error) {
	fmt.Printf("Executing GetHopCount query with UserID1: %d, UserID2: %d\n", startUserID, endUserID)
	ctx := context.Background()
	session := Neo4jDriver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(ctx)

	query := `
        MATCH (a:user {userid: $startID}), (b:user {userid: $endID})
        MATCH path = shortestPath((a)-[:CONNECTED_TO*..5]-(b))
        RETURN length(path) AS hops
    `

	result, err := session.Run(ctx, query, map[string]interface{}{
		"startID": startUserID,
		"endID":   endUserID,
	})

	if err != nil {
		return -1, fmt.Errorf("neo4j query failed: %v", err)
	}

	// FIX: Pass ctx to Next(ctx)
	if result.Next(ctx) {
		hops, _ := result.Record().Get("hops")
		return int(hops.(int64)), nil
	}

	return -1, nil
}
func ConnectUsers(user1ID, user2ID int64) error {
	ctx := context.Background()
	session := Neo4jDriver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	query := `
		MATCH (a:user {userid: $user1ID}), (b:user {userid: $user2ID})
		MERGE (a)-[:CONNECTED_TO]->(b);
	`
	_, err := session.Run(ctx, query, map[string]interface{}{
		"user1ID": user1ID,
		"user2ID": user2ID,
	})
	return err
}
func DisconnectUsers(user1ID, user2ID int64) error {
	ctx := context.Background()
	session := Neo4jDriver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	query := `
		MATCH (a:user {userid: $user1ID})-[r:CONNECTED_TO]->(b:user {userid: $user2ID})
		DELETE r;
	`
	_, err := session.Run(ctx, query, map[string]interface{}{
		"user1ID": user1ID,
		"user2ID": user2ID,
	})
	return err
}
