Feature: get divisions
    In order to search divisions
    I need to do a GET request search urls and pass parameters

    Scenario: Get divisions by coordinates without the full geometry
        When I send a "GET" request to "/divisions?latitude=47.394405&longitude=0.681738" accepting "application/ld+json"
        And the response code should be 200
        And the response header "Content-Type" should be "application/ld+json; charset=utf-8"
        And the JSON should be valid according to this schema:
        """
        {
            "type": "object",
            "properties": {
                "features": {
                    "type": "array",
                    "minItems": 6,
                    "maxItems": 6,
                    "uniqueItems": true,
                    "items": {
                        "type": "object",
                        "properties": {
                            "type": { "type": "string" },
                            "geometry": { "type": "null" },
                            "properties": {
                                "type": "object",
                                "properties": {
                                    "name": {
                                        "type": "string",
                                        "enum": [
                                            "Indre-et-Loire",
                                            "Centre-Val de Loire",
                                            "France",
                                            "Tours-4",
                                            "Tours"
                                        ]
                                    }
                                }
                            }
                        },
                        "required": ["type", "geometry", "properties"],
                        "additionalProperties": false
                    }
                }
            }
        }
        """

    Scenario: Get divisions by coordinates without the full geometry and with a larger radius
        When I send a "GET" request to "/divisions?latitude=47.394405&longitude=0.681738&radius=5000" accepting "application/ld+json"
        And the response code should be 200
        And the response header "Content-Type" should be "application/ld+json; charset=utf-8"
        And the JSON should be valid according to this schema:
        """
        {
            "type": "object",
            "properties": {
                "features": {
                    "type": "array",
                    "minItems": 28,
                    "maxItems": 28,
                    "uniqueItems": true,
                    "items": {
                        "type": "object",
                        "properties": {
                            "type": { "type": "string" },
                            "geometry": { "type": "null" },
                            "properties": {
                                "type": "object",
                                "properties": {
                                    "name": {
                                        "type": "string",
                                        "enum": [
                                            "Ballan-Miré",
                                            "Centre-Val de Loire",
                                            "Chambray-lès-Tours",
                                            "Fondettes",
                                            "France",
                                            "Indre-et-Loire",
                                            "Joué-lès-Tours",
                                            "Membrolle-sur-Choisille",
                                            "Mettray",
                                            "Montlouis-sur-Loire",
                                            "Notre-Dame-d'Oé",
                                            "Parçay-Meslay",
                                            "Riche",
                                            "Rochecorbon",
                                            "Saint-Avertin",
                                            "Saint-Cyr-sur-Loire",
                                            "Saint-Pierre-des-Corps",
                                            "Tours",
                                            "Tours-1",
                                            "Tours-2",
                                            "Tours-3",
                                            "Tours-4",
                                            "Vouvray"
                                        ]
                                    }
                                }
                            }
                        },
                        "required": ["type", "geometry", "properties"],
                        "additionalProperties": false
                    }
                }
            }
        }
        """

    Scenario: Get divisions by search query
        When I send a "GET" request to "/divisions?q=Tours" accepting "application/ld+json"
        And the response code should be 200
        And the response header "Content-Type" should be "application/ld+json; charset=utf-8"
        And the JSON should be valid according to this schema:
        """
        {
            "type": "object",
            "properties": {
                "features": {
                    "type": "array",
                    "minItems": 9,
                    "maxItems": 9,
                    "uniqueItems": true,
                    "items": {
                        "type": "object",
                        "properties": {
                            "type": { "type": "string" },
                            "geometry": { "type": "null" },
                            "properties": {
                                "type": "object",
                                "properties": {
                                    "name": {
                                        "type": "string",
                                        "enum": [
                                            "Tours",
                                            "Tours-1",
                                            "Tours-2",
                                            "Tours-3",
                                            "Tours-4",
                                            "Joué-lès-Tours",
                                            "Chambray-lès-Tours"
                                        ]
                                    }
                                }
                            }
                        },
                        "required": ["type", "geometry", "properties"],
                        "additionalProperties": false
                    }
                }
            }
        }
        """

    Scenario: Get divisions by search query and type
        When I send a "GET" request to "/divisions?q=Tours&properties.administrativeName=commune" accepting "application/ld+json"
        And the response code should be 200
        And the response header "Content-Type" should be "application/ld+json; charset=utf-8"
        And the JSON should be valid according to this schema:
        """
        {
            "type": "object",
            "properties": {
                "features": {
                    "type": "array",
                    "minItems": 3,
                    "maxItems": 3,
                    "uniqueItems": true,
                    "items": {
                        "type": "object",
                        "properties": {
                            "type": { "type": "string" },
                            "geometry": { "type": "null" },
                            "properties": {
                                "type": "object",
                                "properties": {
                                    "name": {
                                        "type": "string",
                                        "enum": [
                                            "Tours",
                                            "Joué-lès-Tours",
                                            "Chambray-lès-Tours"
                                        ]
                                    }
                                }
                            }
                        },
                        "required": ["type", "geometry", "properties"],
                        "additionalProperties": false
                    }
                }
            }
        }
        """
    
    Scenario: Get only cities by coordinates without the full geometry
        When I send a "GET" request to "/divisions?latitude=47.394405&longitude=0.681738&properties.isCity=1" accepting "application/ld+json"
        And the response code should be 200
        And the response header "Content-Type" should be "application/ld+json; charset=utf-8"
        And the JSON should be valid according to this schema:
        """
        {
            "type": "object",
            "properties": {
                "features": {
                    "type": "array",
                    "minItems": 1,
                    "maxItems": 1,
                    "uniqueItems": true,
                    "items": {
                        "type": "object",
                        "properties": {
                            "type": { "type": "string" },
                            "geometry": { "type": "null" },
                            "properties": {
                                "type": "object",
                                "properties": {
                                    "name": {
                                        "type": "string",
                                        "enum": [ "Tours" ]
                                    },
                                    "isCity": {
                                        "type": "boolean",
                                        "enum": [ true ]
                                    }
                                }
                            }
                        },
                        "required": ["type", "geometry", "properties"],
                        "additionalProperties": false
                    }
                }
            }
        }
        """

    Scenario: Get divisions by coordinates with their full geometry
        When I send a "GET" request to "/divisions?latitude=47.394405&longitude=0.681738" accepting "application/geo+json"
        Then the JSON response should be a valid GeoJson Feature Collection
        And the response code should be 200
        And the response header "Content-Type" should be "application/geo+json; charset=utf-8"
        And the JSON should be valid according to this schema:
        """
        {
            "type": "object",
            "properties": {
                "features": {
                    "type": "array",
                    "minItems": 6,
                    "maxItems": 6,
                    "uniqueItems": true,
                    "items": {
                        "type": "object",
                        "properties": {
                            "type": { "type": "string" },
                            "geometry": { "type": "object" },
                            "properties": {
                                "type": "object",
                                "properties": {
                                    "name": {
                                        "type": "string",
                                        "enum": [
                                            "Indre-et-Loire",
                                            "Centre-Val de Loire",
                                            "France",
                                            "Tours-4",
                                            "Tours"
                                        ]
                                    }
                                }
                            }
                        },
                        "required": ["type", "geometry", "properties"],
                        "additionalProperties": false
                    }
                }
            }
        }
        """

    Scenario: Get divisions with pagination
        When I send a "GET" request to "/divisions?from=0&size=3" accepting "application/geo+json"
        Then the JSON response should be a valid GeoJson Feature Collection
        And the response code should be 200
        And the response header "Content-Type" should be "application/geo+json; charset=utf-8"
        And the response header "X-Total-Results" should exist
        And the JSON should be valid according to this schema:
        """
        {
            "type": "object",
            "properties": {
                "features": {
                    "type": "array",
                    "minItems": 3,
                    "maxItems": 3,
                    "uniqueItems": true
                }
            }
        }
        """

    Scenario: Get divisions with invalid latitude
        When I send a "GET" request to "/divisions?latitude=chuck&longitude=0.688802" accepting "application/geo+json"
        Then the response code should be 400
        And the error message should be "Invalid values provided for latitude"
        And the response header "Content-Type" should be "application/json; charset=utf-8"

    Scenario: Get divisions with invalid longitude
        When I send a "GET" request to "/divisions?latitude=47.390359&longitude=norris" accepting "application/geo+json"
        Then the response code should be 400
        And the error message should be "Invalid values provided for longitude"
        And the response header "Content-Type" should be "application/json; charset=utf-8"