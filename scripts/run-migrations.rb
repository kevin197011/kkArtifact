#!/usr/bin/env ruby
# frozen_string_literal: true
#
# Copyright (c) 2025 kk
#
# This software is released under the MIT License.
# https://opensource.org/licenses/MIT

require 'fileutils'

puts "🔄 运行数据库迁移..."
puts ""

# 加载环境变量
env_file = File.join(Dir.pwd, ".env.local")
if File.exist?(env_file)
  File.readlines(env_file).each do |line|
    line = line.strip
    next if line.empty? || line.start_with?('#')
    key, value = line.split('=', 2)
    ENV[key] = value if key && value
  end
end

db_host = ENV['DB_HOST'] || 'localhost'
db_port = ENV['DB_PORT'] || '5432'
db_name = ENV['DB_NAME'] || 'kkartifact'
db_user = ENV['DB_USER'] || 'kkartifact'
db_password = ENV['DB_PASSWORD'] || 'kkartifact'
db_sslmode = ENV['DB_SSLMODE'] || 'disable'

db_url = "postgres://#{db_user}:#{db_password}@#{db_host}:#{db_port}/#{db_name}?sslmode=#{db_sslmode}"
migrations_path = File.join(Dir.pwd, "server/migrations/migrations")

puts "📦 数据库: #{db_host}:#{db_port}/#{db_name}"
puts "📁 迁移路径: #{migrations_path}"
puts ""

# 检查 migrate 工具
unless system("which migrate > /dev/null 2>&1")
  puts "❌ 未找到 migrate 工具"
  puts ""
  puts "请安装 golang-migrate:"
  puts "  brew install golang-migrate"
  puts "  或: go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest"
  exit 1
end

# 运行迁移
puts "▶️  执行迁移..."
system("migrate -path #{migrations_path} -database \"#{db_url}\" up")

if $?.success?
  puts ""
  puts "✅ 数据库迁移完成！"
else
  puts ""
  puts "❌ 数据库迁移失败"
  exit 1
end

