Feature: Access network resources transparently
  In order to access network resouces behind a proxy
  As a corporate company employee
  I want to do so without configuring my environment

  Scenario Outline: Accessing web resource
    Given my system hasn't transprouter
    And a web resource <url> is not accessible
    When my system has transprouter
    And I access web resource <url>
    Then the request response body contain
    """
    You are on a proxied webserver!
    """

    Examples:
    | url                   |
    | http://web.away/lost  |
    | https://web.away/lost |

  Scenario: Connecting to remote SSH server
    Given my system hasn't transprouter
    And a SSH service on ssh.away is is not accessible
    When my system has transprouter
    And I execute "echo hello world" on ssh.away
    Then the command result is
    """
    hello world
    """
