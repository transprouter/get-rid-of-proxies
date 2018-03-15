Feature: Network has constraints

  Scenario Outline: Some web resources are protected behind a proxy
    Given my system hasn't transprouter
    When I request the web resource at <url>
    Then a HTTP timeout error occurred

    Examples:
    | url                     |
    | http://srv2.public.net/ |
    | https://srv2.publc.net/ |

  Scenario Outline: Some web resources are accessible directly
    Given my system hasn't transprouter
    When I request the web resource at <url>
    Then the HTTP reponse body contains
    """
    Welcome to nginx!
    """

    Examples:
    | url                 |
    | http://srv1.local/  |
    | https://srv1.local/ |

  Scenario: Some SSH servers are protected behind a proxy
    Given my system hasn't transprouter
    When I execute "echo hello world" on srv2.public.net
    Then a SSH timeout error occurred

  Scenario: Some SSH servers are accessible directly
    Given my system hasn't transprouter
    When I execute "echo -n hello world" on srv1.local
    Then the command output is
    """
    hello world
    """
