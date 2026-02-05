terraform {
  backend "s3" {
    key          = "lab/foundation.tfstate"
    encrypt      = true
    use_lockfile = true
  }
}

