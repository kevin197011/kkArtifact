
  ./bin/kkArtifact-agent push \
  --project devops \
  --app app01 \
  --version test-$(date +%s) \
  --path ./web-ui \
  --config .kkartifact.yml

./kkartifact-agent push   --project g01   --app G01_dts   --version 885305a25cd644e9a4bcbaaed59c9543333b3701   --path /data/G01/G01_dts  --config .kkartifact.yml