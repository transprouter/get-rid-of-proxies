When(/^I request the web resource at (https?:.+)$/) do |url|
  @http_response = exec(@user, "curl -sSk --max-time 2 #{url}")
end

Then(/^the HTTP reponse body contains$/) do |expected_body|
  expect(@http_response).to eq(expected_body)
end

Then(/^a HTTP timeout error occurred$/) do
  expect(@http_response).to start_with("curl: (28) Connection timed out after ")
end
