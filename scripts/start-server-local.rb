#!/usr/bin/env ruby
# frozen_string_literal: true
#
# Copyright (c) 2025 kk
#
# This software is released under the MIT License.
# https://opensource.org/licenses/MIT

require 'fileutils'

puts "🚀 启动本地 Server..."
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
  puts "✅ 已加载环境变量: .env.local"
else
  puts "⚠️  未找到 .env.local，使用默认配置"
end

# 确保存储目录存在
storage_path = ENV['STORAGE_LOCAL_PATH'] || ENV['STORAGE_BASE_PATH'] || File.expand_path("~/kkartifact-storage")
FileUtils.mkdir_p(storage_path) unless Dir.exist?(storage_path)

puts "📦 存储路径: #{storage_path}"
puts "🌐 服务器地址: #{ENV['SERVER_HOST'] || '0.0.0.0'}:#{ENV['SERVER_PORT'] || '8080'}"
puts ""

# 切换到 server 目录并运行
server_dir = File.join(Dir.pwd, "server")
Dir.chdir(server_dir) do
  puts "▶️  启动 server..."
  exec("go run main.go")
end

