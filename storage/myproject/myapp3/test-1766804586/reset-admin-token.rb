#!/usr/bin/env ruby
# frozen_string_literal: true

# Copyright (c) 2025 kk
#
# This software is released under the MIT License.
# https://opensource.org/licenses/MIT

# Script to reset admin token by deleting the existing one
# Usage: ./scripts/reset-admin-token.sh

require 'net/http'
require 'uri'
require 'json'

# Get database connection info from environment or defaults
db_host = ENV['DB_HOST'] || 'localhost'
db_port = ENV['DB_PORT'] || '5432'
db_name = ENV['DB_NAME'] || 'kkartifact'
db_user = ENV['DB_USER'] || 'kkartifact'
db_password = ENV['DB_PASSWORD'] || 'kkartifact'

puts "🗑️  正在删除现有的 admin-initial-token..."

# Delete the token using psql
system("PGPASSWORD=#{db_password} psql -h #{db_host} -p #{db_port} -U #{db_user} -d #{db_name} -c \"DELETE FROM tokens WHERE name = 'admin-initial-token';\"") || begin
  puts "❌ 无法删除 token。请确保 PostgreSQL 客户端已安装。"
  puts ""
  puts "或者手动执行："
  puts "  docker compose exec postgres psql -U kkartifact -d kkartifact -c \"DELETE FROM tokens WHERE name = 'admin-initial-token';\""
  exit 1
end





puts "✅ Token 已删除！"
puts ""
puts "重启服务器以生成新的 token："
puts "  docker compose restart server"
puts ""
puts "然后查看新的 token："
puts "  docker compose logs server | grep -A 6 'Admin Token'"

