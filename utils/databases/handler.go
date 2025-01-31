package databases

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"handbook-scraper/utils/log"
)

// StorageType represents the different storage strategies we support
type StorageType string

const (
	Timetable StorageType = "timetable" // Direct MongoDB storage
	Handbook  StorageType = "handbook"  // Redis-cached MongoDB storage
	Cache     StorageType = "cache"     // Pure Redis storage
)

var (
	dbHandler *DatabaseHandler
	dbOnce    sync.Once
)

// DatabaseHandler provides a unified interface for different storage strategies
type DatabaseHandler struct {
	redisClient *redis.Client
	mongoClient *mongo.Client
	mongoDB     *mongo.Database
}

// GetDatabaseHandler returns the singleton instance of DatabaseHandler
func GetDatabaseHandler() *DatabaseHandler {
	dbOnce.Do(func() {
		dbHandler = newDatabaseHandler()
	})
	return dbHandler
}

// newDatabaseHandler creates a new database handler with environment variables
func newDatabaseHandler() *DatabaseHandler {
	// Get configuration from environment variables
	mongoURI := os.Getenv("MONGO_URI")
	mongoDB := os.Getenv("MONGO_DB")
	redisURL := os.Getenv("REDIS_URL")

	var redisAddr, redisPass string
	var redisDB int
	if redisURL == "" {
		redisAddr = os.Getenv("REDIS_ADDR")
		redisPass = os.Getenv("REDIS_PASSWORD")

		var err error
		redisDB, err = strconv.Atoi(os.Getenv("REDIS_DB"))
		if err != nil {
			log.Fatalf("Invalid REDIS_DB value: %v", err)
		}
	}

	// Initialize MongoDB
	mongoClient, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	var redisClient *redis.Client
	if redisURL != "" {
		opts, err := redis.ParseURL(redisURL)
		if err != nil {
			log.Fatalf("Failed to parse Redis URL: %v", err)
		}
		redisClient = redis.NewClient(opts)
	} else {
		// Initialize Redis
		redisClient = redis.NewClient(&redis.Options{
			Addr:     redisAddr,
			Password: redisPass,
			DB:       redisDB,
		})
	}

	// Verify connections
	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	if err := mongoClient.Ping(context.Background(), nil); err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}

	log.Successf("Successfully connected to databases for the first time!")

	return &DatabaseHandler{
		redisClient: redisClient,
		mongoClient: mongoClient,
		mongoDB:     mongoClient.Database(mongoDB),
	}
}

// GetMongoClient returns the underlying MongoDB client for direct access
func (h *DatabaseHandler) GetMongoClient() *mongo.Client {
	return h.mongoClient
}

// GetMongoDatabase returns the underlying MongoDB database for direct access
func (h *DatabaseHandler) GetMongoDatabase() *mongo.Database {
	return h.mongoDB
}

// Close closes all database connections
func (h *DatabaseHandler) Close() error {
	var errs []error

	if err := h.redisClient.Close(); err != nil {
		errs = append(errs, fmt.Errorf("redis close error: %w", err))
	}

	if err := h.mongoClient.Disconnect(context.Background()); err != nil {
		errs = append(errs, fmt.Errorf("mongodb close error: %w", err))
	}

	if len(errs) > 0 {
		return fmt.Errorf("errors closing databases: %v", errs)
	}
	return nil
}

// Store stores data using the specified storage strategy
func (h *DatabaseHandler) Store(storageType StorageType, key string, data interface{}, ttl time.Duration) error {
	switch storageType {
	case Timetable:
		return h.storeMongo("timetable", key, data)
	case Handbook:
		if err := h.storeRedis(key, data, ttl); err != nil {
			return fmt.Errorf("failed to store in Redis cache: %w", err)
		}
		return h.storeMongo("handbook", key, data)
	case Cache:
		return h.storeRedis(key, data, ttl)
	default:
		return fmt.Errorf("unsupported storage type: %s", storageType)
	}
}

// storeMongo stores data in MongoDB
func (h *DatabaseHandler) storeMongo(collection string, key string, data interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Convert data to BSON
	bsonData, err := toBSON(data)
	if err != nil {
		return fmt.Errorf("failed to convert data to BSON: %w", err)
	}

	// Upsert document
	_, err = h.mongoDB.Collection(collection).UpdateOne(
		ctx,
		bson.M{"_id": key},
		bson.M{"$set": bsonData},
		options.Update().SetUpsert(true),
	)
	return err
}

// storeRedis stores data in Redis
func (h *DatabaseHandler) storeRedis(key string, data interface{}, ttl time.Duration) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return h.redisClient.Set(ctx, key, jsonData, ttl).Err()
}

// toBSON converts data to BSON format
func toBSON(data interface{}) (bson.M, error) {
	var bsonData bson.M
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	err = bson.UnmarshalExtJSON(jsonData, true, &bsonData)
	return bsonData, err
}

// Retrieve retrieves data using the specified storage strategy
func (h *DatabaseHandler) Retrieve(storageType StorageType, key string, result interface{}) error {
	switch storageType {
	case Timetable:
		return h.retrieveMongo("timetable", key, result)
	case Handbook:
		// Try Redis first
		if err := h.retrieveRedis(key, result); err == nil {
			return nil
		}
		// Fallback to MongoDB
		if err := h.retrieveMongo("handbook", key, result); err != nil {
			return err
		}
		// Cache the result back in Redis
		return h.storeRedis(key, result, 24*time.Hour)
	case Cache:
		return h.retrieveRedis(key, result)
	default:
		return fmt.Errorf("unsupported storage type: %s", storageType)
	}
}

// retrieveMongo retrieves data from MongoDB
func (h *DatabaseHandler) retrieveMongo(collection string, key string, result interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var doc bson.M
	err := h.mongoDB.Collection(collection).FindOne(ctx, bson.M{"_id": key}).Decode(&doc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return fmt.Errorf("document not found")
		}
		return fmt.Errorf("failed to retrieve document: %w", err)
	}

	jsonData, err := json.Marshal(doc)
	if err != nil {
		return fmt.Errorf("failed to marshal document: %w", err)
	}
	return json.Unmarshal(jsonData, result)
}

// retrieveRedis retrieves data from Redis
func (h *DatabaseHandler) retrieveRedis(key string, result interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	data, err := h.redisClient.Get(ctx, key).Result()
	if err != nil {
		return fmt.Errorf("failed to retrieve from Redis: %w", err)
	}
	return json.Unmarshal([]byte(data), result)
}

// Delete removes data using the specified storage strategy
func (h *DatabaseHandler) Delete(storageType StorageType, key string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	switch storageType {
	case Timetable:
		_, err := h.mongoDB.Collection("timetable").DeleteOne(ctx, bson.M{"_id": key})
		return err
	case Handbook:
		if err := h.redisClient.Del(ctx, key).Err(); err != nil {
			return err
		}
		_, err := h.mongoDB.Collection("handbook").DeleteOne(ctx, bson.M{"_id": key})
		return err
	case Cache:
		return h.redisClient.Del(ctx, key).Err()
	default:
		return fmt.Errorf("unsupported storage type: %s", storageType)
	}
}

// Exists checks if a key exists using the specified storage strategy
func (h *DatabaseHandler) Exists(storageType StorageType, key string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	switch storageType {
	case Timetable:
		count, err := h.mongoDB.Collection("timetable").CountDocuments(ctx, bson.M{"_id": key})
		return count > 0, err
	case Handbook:
		// Check Redis first
		exists, err := h.redisClient.Exists(ctx, key).Result()
		if err != nil || exists > 0 {
			return exists > 0, err
		}
		// Check MongoDB
		count, err := h.mongoDB.Collection("handbook").CountDocuments(ctx, bson.M{"_id": key})
		return count > 0, err
	case Cache:
		exists, err := h.redisClient.Exists(ctx, key).Result()
		return exists > 0, err
	default:
		return false, fmt.Errorf("unsupported storage type: %s", storageType)
	}
}

// ListKeys returns all keys matching a pattern using the specified storage strategy
func (h *DatabaseHandler) ListKeys(storageType StorageType, pattern string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	switch storageType {
	case Timetable:
		return h.listMongoKeys("timetable", pattern, ctx)
	case Handbook:
		return h.listMongoKeys("handbook", pattern, ctx)
	case Cache:
		return h.redisClient.Keys(ctx, pattern).Result()
	default:
		return nil, fmt.Errorf("unsupported storage type: %s", storageType)
	}
}

// listMongoKeys is a helper function to list keys from MongoDB
func (h *DatabaseHandler) listMongoKeys(collection string, pattern string, ctx context.Context) ([]string, error) {
	filter := bson.M{"_id": bson.M{"$regex": pattern}}
	cursor, err := h.mongoDB.Collection(collection).Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var keys []string
	for cursor.Next(ctx) {
		var result bson.M
		if err := cursor.Decode(&result); err != nil {
			return nil, err
		}
		if id, ok := result["_id"].(string); ok {
			keys = append(keys, id)
		}
	}
	return keys, nil
}

// Flush clears data using the specified storage strategy
func (h *DatabaseHandler) Flush(storageType StorageType) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	switch storageType {
	case Timetable:
		_, err := h.mongoDB.Collection("timetable").DeleteMany(ctx, bson.M{})
		return err
	case Handbook:
		if err := h.redisClient.FlushDB(ctx).Err(); err != nil {
			return err
		}
		_, err := h.mongoDB.Collection("handbook").DeleteMany(ctx, bson.M{})
		return err
	case Cache:
		return h.redisClient.FlushDB(ctx).Err()
	default:
		return fmt.Errorf("unsupported storage type: %s", storageType)
	}
}
