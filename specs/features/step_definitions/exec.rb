def exec(user, cmd)
  `docker exec -itu #{@user} transprouter_priv_ws #{cmd}`
end

Given(/^my system has transprouter$/) do
  @user = "transprouter"
end

Given(/^my system hasn't transprouter$/) do
  @user = "direct"
end
