Feature: Access network resources transparently
  In order to access network resouces behind a proxy
  As a corporate company employee
  I want to do so without configuring my environment

  Scenario Outline: Accessing web resource
    Given my system has transprouter
    When I request the web resource at <url>
    Then the HTTP reponse body contains
    """
    public webserver
    """

    Examples:
    | url                     |
    | http://web.public/lost  |
    | https://web.public/lost |

  Scenario: Connecting to remote SSH server
    Given my system has transprouter
    When I execute "echo -n hello world" on ssh.public
    Then the command output is
    """
    hello world
    """
