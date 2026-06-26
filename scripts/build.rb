#!/usr/bin/env ruby
# frozen_string_literal: true
# scripts/build.rb - 本地构建（Docker 镜像 / Agent 二进制，无交互）

require 'fileutils'
require 'optparse'

ROOT = File.expand_path('..', __dir__)

options = {
  services: %w[server web-ui],
  no_cache: false,
  up: false,
  test: false,
  agent_only: false,
  agent_out: ENV.fetch('KK_AGENT_BIN', '/tmp/kkartifact-agent')
}

OptionParser.new do |opts|
  opts.banner = 'Usage: ruby scripts/build.rb [options]'
  opts.on('--server', '仅构建 server 镜像') { options[:services] = ['server'] }
  opts.on('--web-ui', '仅构建 web-ui 镜像') { options[:services] = ['web-ui'] }
  opts.on('--all', '构建 server + web-ui（默认）') { options[:services] = %w[server web-ui] }
  opts.on('--agent', '仅编译本地 agent（不构建 Docker）') { options[:agent_only] = true }
  opts.on('--agent-out PATH', 'agent 输出路径') { |v| options[:agent_out] = v }
  opts.on('--no-cache', 'Docker 构建禁用缓存') { options[:no_cache] = true }
  opts.on('--up', '构建后 docker compose up -d') { options[:up] = true }
  opts.on('--test', '构建/启动后运行集成测试') { options[:test] = true; options[:up] = true }
end.parse!

def run!(cmd, desc, chdir: ROOT, env: {})
  label = cmd.is_a?(Array) ? cmd.join(' ') : cmd
  puts ">>> #{desc}"
  puts "    #{label}"
  success = if env.empty?
              cmd.is_a?(Array) ? system(*cmd, chdir: chdir) : system(cmd, chdir: chdir)
            else
              cmd.is_a?(Array) ? system(env, *cmd, chdir: chdir) : system(env, cmd, chdir: chdir)
            end
  abort("❌ #{desc} 失败") unless success

  puts "✅ #{desc} 完成"
end

if options[:agent_only]
  run!(['go', 'build', '-o', options[:agent_out], '.'], '编译 agent', chdir: File.join(ROOT, 'agent'))
  puts "   Agent: #{options[:agent_out]}"
else
  cmd = ['docker', 'compose', 'build', *options[:services]]
  cmd << '--no-cache' if options[:no_cache]
  run!(cmd, "Docker 构建 #{options[:services].join(', ')}")
end

run!(['docker', 'compose', 'up', '-d'], '启动服务') if options[:up] && !options[:agent_only]

if options[:test]
  unless options[:agent_only]
    run!(['go', 'build', '-o', options[:agent_out], '.'], '编译测试用 agent', chdir: File.join(ROOT, 'agent'))
  end
  puts '>>> 等待服务就绪...'
  30.times do
    break if system('curl', '-sf', 'http://localhost:8080/api/v1/health', out: File::NULL, err: File::NULL)

    sleep 2
  end
  run!(
    ['ruby', File.join(ROOT, 'scripts/test_integration.rb')],
    '集成测试',
    env: {
      'KK_SERVER_URL' => ENV.fetch('KK_SERVER_URL', 'http://localhost:8080'),
      'KK_AGENT_BIN' => options[:agent_out],
      'NO_INTERACTION' => 'true'
    }
  )
end

puts
puts '🎉 完成'
puts '   Web UI:  http://localhost:3000'
puts '   API:     http://localhost:8080'
puts '   测试:    ruby scripts/test_integration.rb'
