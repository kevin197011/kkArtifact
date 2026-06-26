#!/usr/bin/env ruby
# frozen_string_literal: true
# 批量 push：遍历 root/{app}/{version}/ 并调用 kkartifact-agent push-tree

require 'optparse'

options = {
  project: ENV['KK_PROJECT'],
  dry_run: false,
  skip_existing: false,
  publish: false,
  config: '.kkartifact.yml'
}

parser = OptionParser.new do |opts|
  opts.banner = 'Usage: ruby scripts/push-tree.rb [root-path] [options]'
  opts.on('-p', '--project PROJECT', '项目名称（或在配置文件中设置 project）') { |v| options[:project] = v }
  opts.on('--config PATH', 'Agent 配置文件路径') { |v| options[:config] = v }
  opts.on('--dry-run', '仅打印将要执行的 push') { options[:dry_run] = true }
  opts.on('--skip-existing', '跳过服务端已存在的版本') { options[:skip_existing] = true }
  opts.on('--publish', '推送成功后自动发布版本') { options[:publish] = true }
end
parser.parse!

root = ARGV[0] || '/data/vcs/G02/tidb'
abort("目录不存在: #{root}") unless File.directory?(root)

cmd = ['kkartifact-agent', 'push-tree', root]
cmd += ['--project', options[:project]] if options[:project]
cmd += ['--config', options[:config]]
cmd << '--dry-run' if options[:dry_run]
cmd << '--skip-existing' if options[:skip_existing]
cmd << '--publish' if options[:publish]

puts "执行: #{cmd.join(' ')}"
success = system(*cmd)
abort('push-tree 失败') unless success
