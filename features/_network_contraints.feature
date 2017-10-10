Feature: Netword has constraints

  Scenario Outline: Some web resources are protected behind a proxy
    Given my system hasn't transprouter
    When I request the web resource at <url>
    Then a HTTP timeout error occurred

    Examples:
    | url                   |
    | http://web.away/lost  |
    | https://web.away/lost |

  Scenario Outline: Some web resources are accessible directly
    Given my system hasn't transprouter
    When I request the web resource at <url>
    Then the HTTP reponse body contains
    """
    You are on a direct webserver!
    """

    Examples:
    | url                   |
    | http://web.local/lost  |
    | https://web.local/lost |

  Scenario: Some SSH servers are protected behind a proxy
    Given my system hasn't transprouter
    When I execute "echo hello world" on ssh.away
    Then a SSH timeout error occurred

  Scenario Outline: Some SSH servers are accessible directly
    Given my system hasn't transprouter
    When I execute "echo -n hello world" on ssh.away
    Then the command output is
    """
    hello world
    """
