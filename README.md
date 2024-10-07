# Card Validation REST API

## Run using docker-compose

```shell
git clone https://github.com/leavemeal0ne/SolidgateTestTask.git

cd SolidgateTestTask

docker-compose up -d
```
## Examples

#### POST `/validate`

`POST http://localhost:8080/validate`

> Valid card of known issuer

Body:
```json
{
    "Card number": 4929164492674746,
    "Expiration month": 12,
    "Expiration year": 2055
}
```
Result:
```json
{
    "valid": true
}
```
> Card with invalid number

Body:
```json
{
    "Card number": 353011133330000011,
    "Expiration month": 12, 
    "Expiration year": 2025
} 
```
Result:
```json
{
    "valid": false,
    "error": {
        "code": "001",
        "message": "the credit card number you entered failed the Luhn Check"
    }
}
```