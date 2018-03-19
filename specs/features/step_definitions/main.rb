def exec(user, cmd)
  `docker exec -itu #{@user} workstation.private #{cmd}`
end

Given(/^my system has transprouter$/) do
  @user = "transprouter"
end

Given(/^my system hasn't transprouter$/) do
  @user = "direct"
end

When(/^I request the web resource at (https?:.+)$/) do |url|
  @http_response = exec(@user, "curl -sSk --max-time 5 #{url}")
end

Then(/^the HTTP reponse body contains$/) do |expected_body|
  expect(@http_response).to eq(expected_body)
end

Then(/^a HTTP timeout error occurred$/) do
  expect(@http_response).to start_with("curl: (28) Connection timed out after ")
end

When(/^I execute "([^"]*)" on (.+)$/) do |cmd, host|
  @command_output = exec(@user, "ssh -o ConnectTimeout=15 -o StrictHostKeyChecking=no root@#{host} #{cmd}")
end

Then(/^the command output is$/) do |expected_output|
  expect(@command_output).to eq(expected_output)
end

Then(/^a SSH timeout error occurred$/) do
  expect(@command_output).to match(/ssh: connect to host .* port 22: Connection timed out\r\n/)
end