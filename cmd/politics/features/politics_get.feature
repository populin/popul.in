Feature: GET politics
  In order to get one or multiple politics
  I need to do a GET request on the right url

  Scenario: Get a known division
    When I send a "GET" request to "/politics/chuck-norris" accepting "application/vnd.api+json"
    And the response code should be 200
    And the response header "Content-Type" should be "application/vnd.api+json; charset=utf-8"
