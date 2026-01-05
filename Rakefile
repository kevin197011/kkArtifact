# frozen_string_literal: true

# Copyright (c) 2025 kk
#
# This software is released under the MIT License.
# https://opensource.org/licenses/MIT

require 'fileutils'
require 'open3'

# Rake tasks for kkArtifact project

desc 'Update agent version.json and binaries'
task :update_agent_version do
  puts 'Updating agent version.json...'
  system('ruby', 'scripts/update-agent-version.rb') || abort('Failed to update agent version')
  puts '✅ Agent version updated'
end

desc 'Build agent binaries for all platforms'
task :build_agent_all do
  puts 'Building agent binaries for all platforms...'
  system('ruby', 'scripts/build-agent-binaries.rb') || abort('Failed to build agent binaries')
  puts '✅ Agent binaries built'
end

desc 'Update agent version and build binaries'
task :build_and_update_agent => [:build_agent_all, :update_agent_version] do
  puts '✅ Agent binaries built and version updated'
end

desc 'Get current git tag version'
task :git_version do
  version = `git describe --tags --exact-match 2>/dev/null`.strip
  version = `git describe --tags 2>/dev/null`.strip if version.empty?
  version = "v#{Time.now.strftime('%Y%m%d%H%M%S')}" if version.empty?
  puts version
end

desc 'Commit changes to git'
task :git_commit do
  # Check if there are any changes to commit
  status = `git status --porcelain`.strip
  if status.empty?
    puts 'No changes to commit'
    next
  end

  # Get current git tag version for commit message
  version = `git describe --tags --exact-match 2>/dev/null`.strip
  version = `git describe --tags 2>/dev/null`.strip if version.empty?
  version = 'latest' if version.empty?

  # Stage agent-related files first (if they exist)
  agent_files = [
    'server/static/agent/version.json',
    'server/static/agent/kkartifact-agent-*'
  ]

  staged_agent_files = false
  agent_files.each do |pattern|
    if pattern.include?('*')
      # Handle glob patterns
      Dir.glob(pattern).each do |file|
        if File.exist?(file)
          system('git', 'add', file)
          staged_agent_files = true
        end
      end
    else
      if File.exist?(pattern)
        system('git', 'add', pattern)
        staged_agent_files = true
      end
    end
  end

  # Stage all other modified files (excluding bin/ and other ignored files)
  modified_files = `git status --porcelain | grep '^ M' | awk '{print $2}'`.split("\n").reject(&:empty?)
  untracked_files = `git status --porcelain | grep '^??' | awk '{print $2}'`.split("\n").reject(&:empty?)
  
  # Add modified files (excluding bin/ directory)
  modified_files.each do |file|
    next if file.start_with?('bin/')
    next if file.start_with?('server/static/agent/') # Already handled above
    system('git', 'add', file)
  end
  
  # Add untracked files (excluding bin/ and ignored patterns)
  untracked_files.each do |file|
    next if file.start_with?('bin/')
    next if file.start_with?('server/static/agent/') # Already handled above
    system('git', 'add', file)
  end

  # Check if there are staged changes
  staged = `git diff --cached --name-only`.strip
  if staged.empty?
    puts 'No staged changes to commit'
    next
  end

  # Determine commit message based on what was changed
  if staged_agent_files
    commit_message = "chore: update binaries and version.json for #{version}"
  else
    # For other changes, use a generic message or try to infer from file names
    changed_files = staged.split("\n")
    if changed_files.any? { |f| f.include?('workflow') || f.include?('.github') }
      commit_message = "ci: update GitHub Actions workflow"
    elsif changed_files.any? { |f| f.include?('Rakefile') }
      commit_message = "chore: update Rakefile"
    else
      commit_message = "chore: update files"
    end
  end

  # Auto commit without confirmation
  system('git', 'commit', '-m', commit_message) || abort('Failed to commit changes')
  puts "✅ Changes committed: #{commit_message}"
  
  # Auto-create version tag after commit (only on main branch)
  create_version_tag
end

# Helper method to create version tag
def create_version_tag
  current_branch = `git rev-parse --abbrev-ref HEAD`.strip
  return unless current_branch == 'main'
  
  # Check if HEAD already has a tag
  existing_tag = `git describe --tags --exact-match HEAD 2>/dev/null`.strip
  return if existing_tag && !existing_tag.empty?
  
  # Generate version tag based on timestamp
  version_tag = "v#{Time.now.strftime('%Y%m%d%H%M%S')}"
  
  # Create the tag
  if system('git', 'tag', version_tag)
    puts "✅ Created version tag: #{version_tag}"
  else
    puts "⚠️  Failed to create tag: #{version_tag}"
  end
end

desc 'Push changes to remote'
task :git_push do
  current_branch = `git rev-parse --abbrev-ref HEAD`.strip
  if current_branch.empty? || current_branch == 'HEAD'
    puts 'Not on a branch, skipping push'
    next
  end

  # Check if there are commits to push
  ahead = `git rev-list --count HEAD..origin/#{current_branch} 2>/dev/null`.strip
  if ahead.empty? || ahead == '0'
    # Check if there are local commits not pushed
    behind = `git rev-list --count origin/#{current_branch}..HEAD 2>/dev/null`.strip
    if behind.empty? || behind == '0'
      puts 'No commits to push'
      # Still check for tags to push
      push_tags
      next
    end
  end

  # Auto push without confirmation
  system('git', 'push', 'origin', current_branch) || abort('Failed to push changes')
  puts "✅ Changes pushed to origin/#{current_branch}"
  
  # Also push tags
  push_tags
end

# Helper method to push tags
def push_tags
  # Get local tags that are not on remote
  local_tags = `git tag -l`.split("\n")
  remote_tags = `git ls-remote --tags origin 2>/dev/null | grep -v '\^{}' | sed 's/.*\\///'`.split("\n")
  
  tags_to_push = local_tags - remote_tags
  
  if tags_to_push.any?
    puts "Pushing #{tags_to_push.size} tag(s) to remote..."
    system('git', 'push', 'origin', '--tags') || puts('⚠️  Failed to push some tags')
    tags_to_push.each { |tag| puts "  ✅ Pushed tag: #{tag}" }
  end
end

desc 'Build all components, update versions, and commit (default task)'
task :default => [:build_all]

desc 'Build all, update versions, commit, and push'
task :build_and_commit => [:build_all, :git_push] do
  puts ''
  puts '✅ All done: built, updated, committed, and pushed'
end

desc 'Build all components (server, agent binaries, and update version)'
task :build_all => [:build_server, :build_and_update_agent, :git_commit] do
  puts ''
  puts '✅ All components built, versions updated, and changes committed'
end

desc 'Build server binary'
task :build_server do
  puts 'Building server...'
  system('make build-server') || abort('Failed to build server')
  puts '✅ Server built'
end

desc 'Show all available tasks'
task :list => :help

desc 'Show help'
task :help do
  puts <<~HELP
    kkArtifact Rake Tasks
    ====================

    Main Tasks (executed in order):
      rake                          Build all, update versions, and commit changes
      rake build_all                Build server, agent binaries, update version, and commit
      rake build_and_commit         Build all, update versions, commit, and push
      rake build_server             Build server binary
      rake build_and_update_agent   Build agent binaries and update version.json

    Git Tasks:
      rake git_commit               Commit changes (with confirmation)
      rake git_push                 Push changes to remote (with confirmation)
      rake git_version              Show current git tag version

    Agent Tasks:
      rake update_agent_version     Update agent version.json (from git tag)
      rake build_agent_all          Build agent binaries for all platforms

    Utility Tasks:
      rake help                     Show this help message

    Examples:
      rake                          # Build everything, commit, and push automatically
      rake build_all                # Build and commit (no push)
      rake build_and_commit         # Build, commit, and push
      rake update_agent_version     # Only update version.json
      rake build_agent_all          # Only build agent binaries
  HELP
end

task :run do
  system 'docker compose down'
  system 'docker compose up -d --build --remove-orphans'
  system 'docker compose logs -f'
end