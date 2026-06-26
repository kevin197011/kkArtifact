#!/usr/bin/env ruby
# frozen_string_literal: true
# scripts/up.rb - 构建并启动本地开发环境（等价于 build.rb --up）

exec('ruby', File.join(__dir__, 'build.rb'), '--up', *ARGV)
