# Ben-Jerry

## Building an API for Ben &amp; Jerry's fans

### Objects handled:  JSON array of below structure

```
{
  "name": "Vanilla Toffee Bar Crunch",
  "image_closed": "/files/live/sites/systemsite/files/flavors/products/us/pint/open-closed-pints/vanilla-toffee-landing.png",
  "image_open": "/files/live/sites/systemsite/files/flavors/products/us/pint/open-closed-pints/vanilla-toffee-landing-open.png",
  "description": "Vanilla Ice Cream with Fudge-Covered Toffee Pieces",
  "story": "Vanilla What Bar Crunch? We gave this flavor a new name to go with the new toffee bars we’re using as part of our commitment to source Fairtrade Certified and non-GMO ingredients. We love it and know you will too!",
  "sourcing_values": [
    "Non-GMO",
    "Cage-Free Eggs",
    "Fairtrade",
    "Responsibly Sourced Packaging",
    "Caring Dairy"
  ],
  "ingredients": [
    "cream",
    "skim milk",
    "liquid sugar",
    "water",
    "sugar",
    "coconut oil",
    "egg yolks",
    "butter",
    "vanilla extract",
    "almonds",
    "cocoa (processed with alkali)",
    "milk",
    "soy lecithin",
    "cocoa",
    "natural flavor",
    "salt",
    "vegetable oil
    "guar gum",
    "carrageenan"
  ],
  "allergy_info": "may contain wheat, peanuts and other tree nuts",
  "dietary_certifications": "Kosher",
  "productId": "646"
}
```
## Tasks Implemented

- Implemented a REST API with CRUD functionality with PosgressSQL as DB 
- Authenicated using JWT token 
- Documentation done using Swagger ( // Swagger json is not created properly at the time of writing )
- Extensive Testing is done using Post Main, Go Test and manual testing

## Architecural Decision

### Table Definition 


Users |  Products  |
-------------|------|
id|  id (PK) |
email| name |
password| image_open |
|| image_close |
|| description |
|| story |
|| allergy_info |
|| dietary_certifications |



ingredientsindex | sourcingvalueindex |
-----------------|--------------------|
id (PK) | id (PK) |
value (text) | value (text) |




ingredients | sourcing_values |
-------------|-----------------|
id | id |
product_id | product_id | 
value_id | value_id | 
PRIMARY KEY (product_id, value_id) |PRIMARY KEY (product_id, value_id) |
FOREIGN KEY (product_id) REFERENCES products (id) |FOREIGN KEY (product_id) REFERENCES products (id) |
FOREIGN KEY (value_id) REFERENCES ingredientsindex (id) |FOREIGN KEY (value_id) REFERENCES sourcingvalueindex (id)|






## API Definition
List of APIs supported

### GET ALL Products

[GET] /products 

Response : JSON Array

### GET Product by id

[GET] /products/{:id}

Request Parameter : id

Response : JSON element with same structure as example



### CREATE Product by id
[POST] /products/{:id}

Request Parameter: id

Request Body :  JSON ( with same structure as example) # productID field not supported # new ID generated



Response :  JSON { "productID" : value }
value - ID of newly created element


### UPDATE Product by id
[PUT] /products/{:id}

Request Parameter: id

Request Body :  JSON ( with same structure as example) 

Response:  0 - No change  , 1 - Element updated


### DELETE Product by id

[DELETE] /products/{:id}

Request Parameter: id

Response:  0 - No change  , 1 - Element updated


