package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

type SysCred struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func main() {
	secretId := flag.String("secretId", "", "The ID of the secret to retrieve")
	flag.Parse()

	if *secretId == "" {
		log.Fatalf("secretId is required")
	}

	ctx := context.Background()

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatalf("unable to load AWS config: %v", err)
	}

	svc := secretsmanager.NewFromConfig(cfg)

	secretValue, err := getSecretValue(ctx, svc, *secretId)
	if err != nil {
		log.Fatalf("unable to get secret value: %v", err)
	}

	var sysCred SysCred
	err = json.Unmarshal([]byte(secretValue), &sysCred)
	if err != nil {
		log.Fatalf("unable to parse secret value: %v", err)
	}

	fmt.Printf("Username: %s\n", sysCred.Username)
	fmt.Printf("Password: %s\n", sysCred.Password)
}

func getSecretValue(ctx context.Context, svc *secretsmanager.Client, secretId string) (string, error) {
	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretId),
	}

	result, err := svc.GetSecretValue(ctx, input)
	if err != nil {
		if err.Error() == "ResourceNotFoundException" {
			return "", fmt.Errorf("secret %s not found", secretId)
		}
		return "", err
	}

	if result.SecretString != nil {
		return *result.SecretString, nil
	}

	return "", fmt.Errorf("secret value is nil")
}
