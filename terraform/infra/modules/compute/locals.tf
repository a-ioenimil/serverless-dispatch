locals {
  binary_name = "bootstrap"
  # Build key: CGO disabled, OS linux, Arch arm64. Output must be 'bootstrap' for provided.al2023
  build_command = "GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -tags lambda.norpc -o bootstrap main.go"
}
