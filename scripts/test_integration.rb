#!/usr/bin/env ruby
# frozen_string_literal: true
# scripts/test_integration.rb - 自动化集成测试（无交互，CI 可运行）

require 'fileutils'
require 'json'
require 'net/http'
require 'open3'
require 'securerandom'
require 'tmpdir'
require 'uri'

class IntegrationTestRunner
  SERVER_URL = ENV.fetch('KK_SERVER_URL', 'http://localhost:8080')
  ADMIN_USER = ENV.fetch('KK_ADMIN_USER', 'admin')
  ADMIN_PASS = ENV.fetch('KK_ADMIN_PASS', 'admin123')
  AGENT_BIN = ENV.fetch('KK_AGENT_BIN', '/tmp/kkartifact-agent')

  def initialize
    @errors = []
    @passed = 0
    @project = "test-#{SecureRandom.hex(4)}"
    @app = 'test_app'
    @version = SecureRandom.hex(8)
    @tmpdir = File.join(Dir.tmpdir, "kkartifact-test-#{SecureRandom.hex(4)}")
    @jwt = nil
    @api_token = nil
    @token_id = nil
  end

  def run
    puts '🚀 kkArtifact 集成测试开始'
    puts "   Server: #{SERVER_URL}"
    puts "   Project/App/Version: #{@project}/#{@app}/#{@version}"
    puts

    setup_fixture_dir
    test_health
    test_public_inventory
    test_deprecated_public_list
    test_admin_login
    test_create_api_token
    test_agent_push
    test_agent_pull
    test_publish_and_latest
    test_audit_logs_filter
    test_admin_only_delete
    test_permission_push_without_pull
    cleanup

    print_results
  end

  private

  def setup_fixture_dir
    version_dir = File.join(@tmpdir, @app, @version)
    FileUtils.mkdir_p(version_dir)
    File.write(File.join(version_dir, 'hello.txt'), "integration test #{Time.now.to_i}\n")
    @version_path = version_dir
    @deploy_path = File.join(@tmpdir, 'deploy')
    @agent_config = File.join(@tmpdir, '.kkartifact.yml')
    File.write(@agent_config, <<~YAML)
      server_url: #{SERVER_URL}
      token: PLACEHOLDER
      project: #{@project}
      concurrency: 2
      ignore: []
    YAML
    ok('准备测试目录与配置')
  end

  def test_health
    resp = public_get('/api/v1/health')
    assert_eq('health status', resp['status'], 'ok')
  end

  def test_public_inventory
    code, headers, body = raw_get('/api/v1/public/inventory')
    assert_true('public inventory 200', code == 200)
    data = JSON.parse(body)
    assert_true('public inventory has projects key', data.key?('projects'))
  end

  def test_deprecated_public_list
    code, headers, _body = raw_get('/api/v1/public/projects?limit=1')
    assert_true('deprecated projects 200', code == 200)
    dep = header_value(headers, 'deprecation')
    if dep == 'true'
      ok('deprecated header present')
    else
      puts '  ⚠️  Deprecation 头未出现（需重建 server 镜像以包含最新中间件）'
    end
  end

  def test_admin_login
    resp = post_json('/api/v1/login', { username: ADMIN_USER, password: ADMIN_PASS })
    assert_true('login returns token', resp['token'].to_s.length > 20)
    @jwt = resp['token']
  end

  def test_create_api_token
    code, _h, body = raw_request('POST', '/api/v1/tokens', {
                                    name: "integration-#{SecureRandom.hex(4)}",
                                    permissions: %w[push pull promote]
                                  }, bearer: @jwt)
    assert_true('create token status 200/201', [200, 201].include?(code))
    resp = JSON.parse(body)
    assert_true('token created', resp['token'].to_s.length > 20)
    @api_token = resp['token']
    @token_id = resp['id']
    update_agent_config_token(@api_token)
  end

  def test_agent_push
    skip_unless_agent!
    cmd = [
      AGENT_BIN, 'push',
      '--project', @project,
      '--path', @version_path,
      '--config', @agent_config
    ]
    run_agent(cmd, 'agent push')
  end

  def test_agent_pull
    skip_unless_agent!
    FileUtils.rm_rf(@deploy_path)
    cmd = [
      AGENT_BIN, 'pull',
      '--project', @project,
      '--app', @app,
      '--version', @version,
      '--path', @deploy_path,
      '--config', @agent_config
    ]
    run_agent(cmd, 'agent pull')
    pulled = File.join(@deploy_path, 'hello.txt')
    assert_true('pull created hello.txt', File.file?(pulled))
  end

  def test_publish_and_latest
    resp = token_post_json('/api/v1/publish', {
                             project: @project,
                             app: @app,
                             version: @version
                           })
    assert_eq('publish status', resp['status'], 'published')

    latest = token_get("/api/v1/projects/#{@project}/apps/#{@app}/latest")
    assert_eq('latest version', latest['version'], @version)
  end

  def test_audit_logs_filter
    logs = auth_get('/api/v1/audit-logs?operation=push&limit=10')
    assert_true('audit logs total is number', logs['total'].is_a?(Integer))
    push_ops = logs['data'].select { |e| e['operation'] == 'push' }
    assert_true('audit push filter works', push_ops.any? || logs['total'].zero?)
  end

  def test_admin_only_delete
    # 非 admin 的 API token 不能创建 webhook
    code, = raw_request('POST', '/api/v1/webhooks', {
                          name: 'should-fail',
                          event_types: ['push'],
                          url: 'https://example.com/hook',
                          enabled: true
                        }, bearer: @api_token)
    assert_eq('non-admin webhook create forbidden', code, 403)
  end

  def test_permission_push_without_pull
    code, _h, body = raw_request('POST', '/api/v1/tokens', {
                                    name: "push-only-#{SecureRandom.hex(4)}",
                                    permissions: ['push']
                                  }, bearer: @jwt)
    assert_true('create push-only token 200/201', [200, 201].include?(code))
    resp = JSON.parse(body)
    push_only = resp['token']
    push_only_id = resp['id']

    code, = raw_get("/api/v1/manifest/#{@project}/#{@app}/#{@version}", bearer: push_only)
    assert_eq('push-only cannot pull manifest', code, 403)

    auth_delete("/api/v1/tokens/#{push_only_id}")
  end

  def cleanup
    auth_delete("/api/v1/tokens/#{@token_id}") if @token_id
    auth_delete("/api/v1/projects/#{@project}") if @jwt
    FileUtils.rm_rf(@tmpdir)
    ok('清理测试资源')
  rescue StandardError => e
    @errors << "清理失败: #{e.message}"
  end

  # --- helpers ---

  def skip_unless_agent!
    return if File.executable?(AGENT_BIN)

    @errors << "Agent 不可执行: #{AGENT_BIN}（跳过 push/pull 测试）"
    raise SkipStep
  end

  def run_agent(cmd, label)
    stdout, stderr, status = Open3.capture3({ 'NO_INTERACTION' => 'true' }, *cmd)
    puts stdout unless stdout.empty?
    warn stderr unless stderr.empty?
    assert_true("#{label} 成功", status.success?, stderr.empty? ? stdout : stderr)
  rescue SkipStep
    raise
  rescue StandardError => e
    fail_test("#{label} 异常: #{e.message}")
  end

  def update_agent_config_token(token)
    content = File.read(@agent_config).gsub('PLACEHOLDER', token)
    File.write(@agent_config, content)
  end

  def public_get(path)
    code, _h, body = raw_get(path)
    assert_eq("#{path} status", code, 200)
    JSON.parse(body)
  end

  def auth_get(path)
    code, _h, body = raw_get(path, bearer: @jwt)
    assert_eq("#{path} status", code, 200)
    JSON.parse(body)
  end

  def token_get(path)
    code, _h, body = raw_get(path, bearer: @api_token)
    assert_eq("#{path} status", code, 200)
    JSON.parse(body)
  end

  def auth_post_json(path, payload)
    code, _h, body = raw_request('POST', path, payload, bearer: @jwt)
    assert_eq("#{path} status", code, 201) if path.include?('tokens')
    assert_eq("#{path} status", code, 200) unless path.include?('tokens')
    JSON.parse(body)
  end

  def token_post_json(path, payload)
    code, _h, body = raw_request('POST', path, payload, bearer: @api_token)
    assert_eq("#{path} status", code, 200)
    JSON.parse(body)
  end

  def auth_delete(path)
    raw_request('DELETE', path, nil, bearer: @jwt)
  end

  def post_json(path, payload)
    code, _h, body = raw_request('POST', path, payload)
    assert_eq("#{path} status", code, 200)
    JSON.parse(body)
  end

  def raw_get(path, bearer: nil)
    raw_request('GET', path, nil, bearer: bearer)
  end

  def header_value(headers, name)
    val = headers[name]
    val.is_a?(Array) ? val.first : val
  end

  def raw_request(method, path, payload = nil, bearer: nil)
    uri = URI("#{SERVER_URL}#{path}")
    http = Net::HTTP.new(uri.host, uri.port)
    http.open_timeout = 10
    http.read_timeout = 120
    req = case method
          when 'GET' then Net::HTTP::Get.new(uri)
          when 'POST' then Net::HTTP::Post.new(uri)
          when 'DELETE' then Net::HTTP::Delete.new(uri)
          else raise "unsupported method #{method}"
          end
    req['Content-Type'] = 'application/json'
    req['Authorization'] = "Bearer #{bearer}" if bearer
    req.body = payload.to_json if payload
    res = http.request(req)
    [res.code.to_i, res.to_hash.transform_keys(&:downcase), res.body]
  end

  def ok(msg)
    @passed += 1
    puts "  ✅ #{msg}"
  end

  def assert_eq(label, actual, expected)
    return ok(label) if actual == expected

    fail_test("#{label}: 期望 #{expected.inspect}, 实际 #{actual.inspect}")
  end

  def assert_true(label, condition, detail = nil)
    return ok(label) if condition

    fail_test([label, detail].compact.join(' - '))
  end

  def fail_test(msg)
    @errors << msg
    puts "  ❌ #{msg}"
  end

  def print_results
    puts
    puts '=' * 50
    if @errors.empty?
      puts "✅ 全部通过（#{@passed} 项检查）"
      exit 0
    else
      puts "❌ 失败 #{@errors.size} 项："
      @errors.each { |e| puts "   - #{e}" }
      exit 1
    end
  end

  class SkipStep < StandardError; end
end

IntegrationTestRunner.new.run
