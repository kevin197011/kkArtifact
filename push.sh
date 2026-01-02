./bin/kkArtifact-agent push \
  --project devops \
  --app app01 \
  --version test-$(date +%s) \
  --path ./web-ui \
  --config .kkartifact.yml







  ./bin/kkArtifact-agent push \
  --project video \
  --app srs \
  --version test-$(date +%s) \
  --path /Users/kevin/projects/work/video \
  --config .kkartifact.yml