#!/usr/bin/env ruby
# frozen_string_literal: true
#
# Copyright (c) 2025 kk
#
# This software is released under the MIT License.
# https://opensource.org/licenses/MIT

require 'fileutils'

# 加载环境变量
env_file = File.join(__dir__, "..", ".env.local")
if File.exist?(env_file)
  File.readlines(env_file).each do |line|
    line = line.strip
    next if line.empty? || line.start_with?('#')
    key, value = line.split('=', 2)
    ENV[key] = value.gsub(/^~/, Dir.home) if key && value
  end
end

# 确保存储目录存在
storage_path = (ENV['STORAGE_LOCAL_PATH'] || ENV['STORAGE_BASE_PATH'] || File.expand_path("~/kkartifact-storage")).gsub(/^~/, Dir.home)
FileUtils.mkdir_p(storage_path) unless Dir.exist?(storage_path)
ENV['STORAGE_LOCAL_PATH'] = storage_path
ENV['STORAGE_BASE_PATH'] = storage_path

puts "🚀 启动 kkArtifact Server (本地模式)"
puts ""
puts "📦 存储路径: #{storage_path}"
puts "🌐 服务器: #{ENV['SERVER_HOST'] || '0.0.0.0'}:#{ENV['SERVER_PORT'] || '8080'}"
puts ""

# 切换到 server 目录
Dir.chdir(File.join(__dir__, "..", "server")) do
  exec("go run main.go")
end

