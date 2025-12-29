./bin/kkArtifact-agent push \
  --project myproject \
  --app myapp4 \
  --version test-$(date +%s) \
  --path ./web-ui \
  --config .kkartifact.yml