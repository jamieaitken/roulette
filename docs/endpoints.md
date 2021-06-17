# Table
A Table represents the roulette table where Bets are placed, and an outcome is decided.

## Create
Create a table.
```http request
POST http://localhost:8080/v1/tables
```

## List
Retrieve all created tables.
```http request
GET http://localhost:8080/v1/tables
```

## Get
Fetch a specific table.
```http request
GET http://localhost:8080/v1/tables/{table}
```

## Spin
Accept no more bets and generate the Outcome.
```http request
PUT http://localhost:8080/v1/tables/{table}/spin
```

## Settle
Move all bets to settled and find any winners.
```http request
PUT http://localhost:8080/v1/tables/{table}/settle
```

# Bet
A Bet represents an individuals stake for a given Table.

## Create
Create a bet for the given table; the required body can be found below.
```http request
POST http://localhost:8080/v1/tables/{table}/bet
```

```json
{
  "stake": {
    "amount": 400,
    "currency": "GBP"
  },
  "selectedSpaces": [1, 14, 17]
}
```

## Get
Fetch a specific bet.
```http request
GET http://localhost:8080/v1/bets/{bet}
```
