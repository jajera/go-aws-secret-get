# go-aws-secret-get

export AWS_VAULT_FILE_PASSPHRASE="$(cat /root/.awsvaultk)"

aws-vault exec dev -- terraform -chdir=./terraform init
aws-vault exec dev -- terraform -chdir=./terraform apply --auto-approve

bash ./terraform/terraform.tmp

go mod init awsSecretGet
go mod tidy
go run cmd/awsSecretGet/main.go -secretId=secret-get-<PREFIX>
