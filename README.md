# Example Hearthstone API Browser Webpage Thing-o-matic 5000

## How to run locally
This would eventually be deployed as a helm chart or similar where the values are managed (plus vault for secrets)
```
export HS_API_CLIENT_ID="battle.net api client id"
export HS_API_CLIENT_SECRET="battle.net api client secret"
export HS_MAX_MANNA_COST=10
go run .
```

## Accessing
```
curl http://localhost:8080/cards/list?classes=hunter,mage&manaCost=7&rarity=epic&limit=10
```