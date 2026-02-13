moved {
  from = aws_amplify_app.my_app
  to   = module.amplify.aws_amplify_app.this
}
