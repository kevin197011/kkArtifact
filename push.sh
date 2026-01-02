
  ./bin/kkArtifact-agent push \
  --project devops \
  --app app01 \
  --version test-$(date +%s) \
  --path ./web-ui \
  --config .kkartifact.yml

