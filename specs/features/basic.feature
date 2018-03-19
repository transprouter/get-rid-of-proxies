Feature: Access network resources transparently
  In order to access network resouces behind a proxy
  As a corporate company employee
  I want to do so without configuring my environment

  Scenario Outline: Accessing web resource
    Given my system has transprouter
    When I request the web resource at <url>
    Then the HTTP reponse body contains
    """
    You are on a proxied webserver!
    """

    Examples:
    | url                   |
    | http://webserver.public/lost  |
    | https://webserver.public/lost |

  Scenario: Connecting to remote SSH server
    Given my system has transprouter
    When I execute "echo -n hello world" on sshserver.public
    Then the command output is
    """
    hello world
    """