#!/usr/bin/env ruby
# frozen_string_literal: true

# Copyright (c) 2025 kk
#
# This software is released under the MIT License.
# https://opensource.org/licenses/MIT

require 'fileutils'
require 'json'
require 'open3'

# Quick script to update agent version.json without rebuilding binaries
class AgentVersionUpdater
  PLATFORMS = [
    { goos: 'linux', goarch: 'amd64', suffix: '' },
    { goos: 'linux', goarch: 'arm64', suffix: '' },
    { goos: 'darwin', goarch: 'amd64', suffix: '' },
    { goos: 'darwin', goarch: 'arm64', suffix: '' },
    { goos: 'windows', goarch: 'amd64', suffix: '.exe' }
  ].freeze

  def initialize(output_dir = 'server/static/agent')
    @output_dir = output_dir
  end

  def update_version
    version_file = File.join(@output_dir, 'version.json')
    
    # Get version from git tag
    version = get_version_from_git
    puts "更新版本信息到: #{version}"
    
    # Read existing version.json if it exists
    existing_data = {}
    if File.exist?(version_file)
      begin
        existing_data = JSON.parse(File.read(version_file))
        puts "当前版本: #{existing_data['version']}"
      rescue
        # Ignore parse errors
      end
    end

    # Get binaries info from existing files
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
        # If file doesn't exist, try to get from existing data
        existing_binary = existing_data['binaries']&.find { |b| b['filename'] == filename }
        existing_binary ? {
          platform: existing_binary['platform'],
          filename: filename,
          size: existing_binary['size'],
          url: existing_binary['url']
        } : nil
      end
    end.compact

    # If no binaries found, use existing ones
    if binaries.empty? && existing_data['binaries']
      binaries = existing_data['binaries']
    end

    version_info = {
      version: version,
      build_time: Time.now.utc.iso8601,
      binaries: binaries
    }

    # Ensure directory exists
    FileUtils.mkdir_p(@output_dir)
    
    # Write version.json
    File.write(version_file, JSON.pretty_generate(version_info))
    puts "✅ 已更新版本信息文件: #{version_file}"
    puts "   版本: #{version}"
    puts "   构建时间: #{version_info[:build_time]}"
    puts "   二进制文件数量: #{binaries.size}"
  end

  private

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
end

# Run if executed directly
if __FILE__ == $PROGRAM_NAME
  output_dir = ARGV[0] || 'server/static/agent'
  updater = AgentVersionUpdater.new(output_dir)
  updater.update_version
end

