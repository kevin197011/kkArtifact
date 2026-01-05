#!/usr/bin/env ruby
# frozen_string_literal: true
#
# Copyright (c) 2025 kk
#
# This software is released under the MIT License.
# https://opensource.org/licenses/MIT

require 'fileutils'

puts "🔧 设置本地开发环境..."
puts ""

# 创建本地存储目录
storage_path = File.expand_path("~/kkartifact-storage")
FileUtils.mkdir_p(storage_path) unless Dir.exist?(storage_path)
puts "✅ 创建存储目录: #{storage_path}"

# 创建 .env.local 文件
env_file = File.join(Dir.pwd, ".env.local")
unless File.exist?(env_file)
  env_content = <<~ENV
    # 本地开发环境变量
    SERVER_HOST=0.0.0.0
    SERVER_PORT=8080
    
    # 数据库配置（Docker Compose 服务）
    DB_HOST=localhost
    DB_PORT=5432
    DB_NAME=kkartifact
    DB_USER=kkartifact
    DB_PASSWORD=kkartifact
    DB_SSLMODE=disable
    
    # Redis 配置（Docker Compose 服务）
    REDIS_HOST=localhost
    REDIS_PORT=6379
    REDIS_PASSWORD=
    REDIS_DB=0
    
    # 存储配置（本地文件系统）
    STORAGE_TYPE=local
    STORAGE_BASE_PATH=#{storage_path}
    STORAGE_LOCAL_PATH=#{storage_path}
    
    # 日志配置
    LOG_LEVEL=info
    LOG_FORMAT=text
    
    # 版本保留
    VERSION_RETENTION_LIMIT=5
  ENV
  
  File.write(env_file, env_content)
  puts "✅ 创建环境变量文件: .env.local"
end

puts ""
puts "📝 环境变量文件: .env.local"
puts "📦 存储目录: #{storage_path}"
puts ""
puts "✅ 本地环境设置完成！"
puts ""
puts "下一步："
puts "  1. 确保依赖服务运行: docker compose up -d postgres redis"
puts "  2. 运行数据库迁移"
puts "  3. 启动 server: cd server && source ../.env.local && go run main.go"

