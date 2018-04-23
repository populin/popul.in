Feature: GET divisions
  In order to get one or multiple divisions
  I need to do a GET request on the right url

  Scenario: Get a known division
    When I send a "GET" request to "/divisions/fr-region-centre-val-de-loire-5727539415420288060" accepting "application/geo+json"
    Then the JSON response should be a valid GeoJson Feature
    And the response code should be 200
    And the response header "Content-Type" should be "application/geo+json; charset=utf-8"
    And the GeoJSON property "name" should be equal to "Centre-Val de Loire"
    And the JSON should be valid according to this schema:
      """
      {
        "type": "object",
        "properties": {
          "id": { "type": "string" },
          "type": { "type": "string" },
          "geometry": { "type": "object" },
          "properties": {
            "type": "object",
            "properties": {
              "name": {
                "type": "string",
                "enum": [
                  "Centre-Val de Loire"
                ]
              }
            }
          }
        },
        "required": ["id", "type", "geometry", "properties"],
        "additionalProperties": false
      }
      """

  Scenario: Get a known division
    When I send a "GET" request to "/divisions/fr-region-centre-val-de-loire-5727539415420288060" accepting "application/vnd.api+json"
    Then the response code should be 200
    And the response header "Content-Type" should be "application/vnd.api+json; charset=utf-8"
    And the JSON should be valid according to this schema:
      """
      {
        "type": "object",
        "properties": {
          "data": {
            "type": "object",
            "properties": {
              "id": {"type": "string"},
              "type": {"type": "string"},
              "attributes": {
                "type": "object",
                "properties": {
                  "properties": {
                    "name": {
                      "type": "string",
                      "enum": [
                        "Centre-Val de Loire"
                      ]
                    }
                  }
                }
              }
            },
            "required": ["type", "id", "attributes"],
            "additionalProperties": false
          }
        },
        "required": ["data"],
        "additionalProperties": false
      }
      """

  Scenario: Get an unknown division
    When I send a "GET" request to "/divisions/not-found" accepting "application/geo+json"
    Then the response code should be 404
    And the error message should be "Division not-found not found"
    And the response header "Content-Type" should be "application/vnd.api+json; charset=utf-8"
