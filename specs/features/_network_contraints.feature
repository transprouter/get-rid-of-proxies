Feature: Netword has constraints

  Scenario Outline: Some web resources are protected behind a proxy
    Given my system hasn't transprouter
    When I request the web resource at <url>
    Then a HTTP timeout error occurred

    Examples:
    | url                     |
    | http://web.public/lost  |
    | https://web.public/lost |

  Scenario Outline: Some web resources are accessible directly
    Given my system hasn't transprouter
    When I request the web resource at <url>
    Then the HTTP reponse body contains
    """
    private
    """

    Examples:
    | url                      |
    | http://web.private/lost  |
    | https://web.private/lost |

  Scenario: Some SSH servers are protected behind a proxy
    Given my system hasn't transprouter
    When I execute "echo hello world" on ssh.public
    Then a SSH timeout error occurred

  Scenario: Some SSH servers are accessible directly
    Given my system hasn't transprouter
    When I execute "echo -n hello world" on ssh.private
    Then the command output is
    """
    hello world
    """
