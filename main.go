package mgdbconnector

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//Auth copies the same pattern as mongo driver options Auth
type MongoDbConfiguration struct {
	Timeout        int64
	AppName        string
	Auth           *options.Credential
	ConnectTimeout int
	Uri            string
}

// AuthMechanism: the mechanism to use for authentication.
//Supported values include "SCRAM-SHA-256", "SCRAM-SHA-1", "MONGODB-CR", "PLAIN", "GSSAPI", "MONGODB-X509", and "MONGODB-AWS".
//This can also be set through the "authMechanism" URI option. (e.g. "authMechanism=PLAIN").
//For more information, see https://docs.mongodb.com/manual/core/authentication-mechanisms/.

func StartMongoDb(cf MongoDbConfiguration) (*mongo.Client, error) {
	timeout := 10 * time.Second
	connectionTimeout := 10 * time.Second
	appName := "mgbconnector"
	if cf.Timeout != 0 {
		timeout = time.Duration(cf.Timeout * int64(time.Second))
	}
	if cf.ConnectTimeout != 0 {
		connectionTimeout = time.Duration(cf.ConnectTimeout * int(time.Second))
	}
	if cf.AppName != "" {
		appName = cf.AppName
	}
	if cf.Uri == "" {
		return nil, errors.New("uri cannot be empty")
	}
	options := options.Client()
	options.AppName = &appName
	options.Auth = cf.Auth
	options.ConnectTimeout = &connectionTimeout
	options = options.ApplyURI(cf.Uri)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout))
	defer cancel()
	client, err := mongo.Connect(ctx, options)
	if err != nil {
		return nil, errors.New("connection couldn't be stabilished, review your infos")
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, errors.New("connection with MongoDb couldn't be stabilished, timeout")
	}
	return client, nil

}
