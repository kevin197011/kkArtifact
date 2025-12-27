./bin/kkArtifact-agent push \
  --project myproject \
  --app myapp3 \
  --version test-$(date +%s) \
  --path ./scripts \
  --config .kkartifact.yml