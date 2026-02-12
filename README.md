# Serverless Dispatch

## CI/CD Required Secrets

Set the following repository secrets in GitHub (`Settings` → `Secrets and variables` → `Actions`) for `.github/workflows/deploy.yml`:

- `AWS_ROLE_ARN` — IAM role ARN assumed by GitHub Actions via OIDC.
- `TF_STATE_BUCKET` — S3 bucket name used by Terraform backend state.
- `TF_DEV_TFVARS` — full contents of `terraform/infra/dev.tfvars` (multiline secret).
- `COGNITO_USER_POOL_ID` — Cognito User Pool ID passed as `TF_VAR_user_pool_id`.

## Notes

- The pipeline injects `COGNITO_USER_POOL_ID` as `TF_VAR_user_pool_id`, which Terraform uses to set `USER_POOL_ID` on backend notification Lambda.
- If `COGNITO_USER_POOL_ID` is missing, the `infra-plan` and `infra-deploy` jobs fail fast.
