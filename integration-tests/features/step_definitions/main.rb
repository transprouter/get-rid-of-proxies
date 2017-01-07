require 'net/ssh'

def sshexec(user, cmd)
  Net::SSH.start('localhost', user, :port => 2222) do |ssh|
    return ssh.exec!(cmd)
  end
end

Given(/^my system has transprouter$/) do
  @user = "proxied"
end

Given(/^my system hasn't transprouter$/) do
  @user = "direct"
end

When(/^I request the web resource at (https?:.+)$/) do |url|
  @http_response = sshexec(@user, "curl -sSk --max-time 5 #{url}")
end

Then(/^the HTTP reponse body contains$/) do |expected_body|
  expect(@http_response).to eq(expected_body)
end

Then(/^a HTTP timeout error occurred$/) do
  expect(@http_response).to start_with("curl: (28) Connection timed out after ")
end

When(/^I execute "([^"]*)" on (.+)$/) do |cmd, host|
  @command_output = sshexec(@user, "ssh -o ConnectTimeout=15 root@#{host} #{cmd}")
end

Then(/^the command output is$/) do |expected_output|
  expect(@command_output).to eq(expected_output)
end

Then(/^a SSH timeout error occurred$/) do
  expect(@command_output).to match(/ssh: connect to host .* port 22: Connection timed out\r\n/)
end