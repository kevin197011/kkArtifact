#!/usr/bin/env ruby
# frozen_string_literal: true

# Copyright (c) 2025 kk
#
# This software is released under the MIT License.
# https://opensource.org/licenses/MIT

require 'fileutils'
require 'open3'

# Build agent binaries for multiple platforms
class AgentBinaryBuilder
  PLATFORMS = [
    { goos: 'linux', goarch: 'amd64', suffix: '' },
    { goos: 'linux', goarch: 'arm64', suffix: '' },
    { goos: 'darwin', goarch: 'amd64', suffix: '' },
    { goos: 'darwin', goarch: 'arm64', suffix: '' },
    { goos: 'windows', goarch: 'amd64', suffix: '.exe' }
  ].freeze

  def initialize(output_dir = 'server/static/agent')
    @output_dir = output_dir
    @agent_dir = 'agent'
  end

  def build_all
    version = get_version_from_git
    puts "开始构建 agent 二进制文件..."
    puts "版本: #{version}"
    puts "输出目录: #{@output_dir}"
    puts ''

    # Create output directory
    FileUtils.mkdir_p(@output_dir)

    success_count = 0
    failed_platforms = []

    PLATFORMS.each do |platform|
      result = build_platform(platform)
      if result[:success]
        success_count += 1
        puts "✅ #{platform[:goos]}/#{platform[:goarch]} - #{result[:filename]}"
      else
        failed_platforms << platform
        puts "❌ #{platform[:goos]}/#{platform[:goarch]} - #{result[:error]}"
      end
    end

    puts ''
    puts "构建完成: #{success_count}/#{PLATFORMS.size} 个平台"
    
    if failed_platforms.any?
      puts "失败的平台:"
      failed_platforms.each do |platform|
        puts "  - #{platform[:goos]}/#{platform[:goarch]}"
      end
      exit 1
    end

    # Generate version info file
    generate_version_info
  end

  private

  def build_platform(platform)
    filename = "kkartifact-agent-#{platform[:goos]}-#{platform[:goarch]}#{platform[:suffix]}"
    output_path = File.join(@output_dir, filename)

    # Get version from git tag
    version = get_version_from_git
    build_time = Time.now.utc.iso8601
    git_commit = get_git_commit

    env = {
      'GOOS' => platform[:goos],
      'GOARCH' => platform[:goarch],
      'CGO_ENABLED' => '0'
    }

    # Build with version info injected via ldflags
    ldflags = "-s -w -X github.com/kk/kkartifact-agent/internal/cli.Version=#{version} -X github.com/kk/kkartifact-agent/internal/cli.BuildTime=#{build_time} -X github.com/kk/kkartifact-agent/internal/cli.GitCommit=#{git_commit}"
    cmd = "cd #{@agent_dir} && go build -trimpath -ldflags='#{ldflags}' -o ../#{output_path} ./main.go"

    stdout, stderr, status = Open3.capture3(env, cmd)

    if status.success?
      { success: true, filename: filename, version: version }
    else
      { success: false, error: stderr.strip }
    end
  rescue => e
    { success: false, error: e.message }
  end

  def get_version_from_git
    # Try to get exact tag first
    tag, _ = Open3.capture2('git describe --tags --exact-match 2>/dev/null')
    return tag.strip if tag && !tag.strip.empty?

    # Try to get latest tag with distance
    tag, _ = Open3.capture2('git describe --tags 2>/dev/null')
    return tag.strip if tag && !tag.strip.empty?

    # Fallback to timestamp
    Time.now.strftime('v%Y%m%d%H%M%S')
  end

  def get_git_commit
    commit, _ = Open3.capture2('git rev-parse --short HEAD 2>/dev/null')
    commit&.strip || 'unknown'
  end

  def generate_version_info
    version_file = File.join(@output_dir, 'version.json')
    
    # Get version from git tag (same logic as build_platform)
    version = get_version_from_git

    binaries = PLATFORMS.map do |platform|
      filename = "kkartifact-agent-#{platform[:goos]}-#{platform[:goarch]}#{platform[:suffix]}"
      file_path = File.join(@output_dir, filename)
      
      if File.exist?(file_path)
        size = File.size(file_path)
        {
          platform: "#{platform[:goos]}/#{platform[:goarch]}",
          filename: filename,
          size: size,
          url: "/api/v1/downloads/agent/#{filename}"
        }
      else
        nil
      end
    end.compact

    version_info = {
      version: version,
      build_time: Time.now.utc.iso8601,
      binaries: binaries
    }

    File.write(version_file, JSON.pretty_generate(version_info))
    puts "生成版本信息文件: #{version_file}"
  rescue => e
    puts "警告: 生成版本信息失败: #{e.message}"
  end
end

# Run if executed directly
if __FILE__ == $PROGRAM_NAME
  require 'json'
  
  output_dir = ARGV[0] || 'server/static/agent'
  builder = AgentBinaryBuilder.new(output_dir)
  builder.build_all
end

