REST API developed using golang, gorilla mux web framework that aggregates foods from various suppliers and provides APIs for consumers to buy foods

There are 5 API endpoints(GET Methods) in this service

1. /buy-item/{name} - to find the item by specifying its name in the path parameter

2. /buy-item-qty/{name}&{quantity} - to find the item by specifying the name and quantity as path parameters

3. /buy-item-qty-price/{name}&{quantity}&{price} - to find the item by specifying the name, quantity and price as path parameters

4. /show-summary - to display all the cached items

5. /fast-buy-item/{name} - to fetch the item by making parallel calls to Supplier APIs
