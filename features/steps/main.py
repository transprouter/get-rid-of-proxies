from behave import given, when, then
from hamcrest import assert_that, contains_string, equal_to, starts_with

@given('my system has transprouter')
def step(context):
    context.host = context.mn.get('ws1')

@given('my system hasn\'t transprouter')
def step(context):
    context.host = context.mn.get('ws2')

@when('I request the web resource at {url}')
def step(context, url):
    context.response_body = context.host.cmd('curl -i -sSk --max-time 5 %s' % url)

@then('the HTTP reponse body contains')
def step(context):
    assert_that(context.response_body, contains_string(context.text))

@then('a HTTP timeout error occurred')
def step_impl(context):
    assert_that(context.response_body, starts_with('curl: (28) Connection timed out after'))

@when('I execute "{command}" on {host}')
def step( context, command, host ):
    context.command_output = context.host.cmd('ssh -o ConnectTimeout=15 -o StrictHostKeyChecking=no root@$(echo %s ) %s' % (host, command) )

@then('the command output is')
def step( context ):
    assert_that(context.command_output, equal_to(context.text))

@then('a SSH timeout error occurred')
def step_impl(context):
    assert_that(context.command_output, equal_to('ssh: connect to host .* port 22: Connection timed out\r\n'))
