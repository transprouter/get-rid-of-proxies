When(/^I execute "([^"]*)" on (.+)$/) do |cmd, host|
  begin
    @command_output = Timeout::timeout(2) {
      exec(@user, "ssh root@#{host} #{cmd}")
    }
  rescue
    @command_output = "Operation timed out"
  end
end

Then(/^the command output is$/) do |expected_output|
  expect(@command_output).to eq(expected_output)
end

Then(/^a SSH timeout error occurred$/) do
  expect(@command_output).to eq("Operation timed out")
end
