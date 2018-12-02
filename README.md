# AdventureWorkReview

## **Instructions**

1. *Clone the repository*
2. *Download docker images*
    ```
   docker pull mysql
   docker pull rabbitmq:3-management-alpine
   ``` 
3. *Change current folder to project*
    ``` 
   cd AdventureWorkReview
   ``` 
4. *Build docker image*
    ``` 
   docker-compose build
   ```
5. *Create and start containers*
    ``` 
   docker-compose up -d
   ```
   - Please wait several minutes for migrating database. When it is complete, the application container would start.
   - You can check whether it starts with the command below:
        ``` 
       docker ps
        ```
6. *Post a review*
    ``` 
    curl -X POST http://0.0.0.0:8888/api/reviews \ -H 'Content-Type: application/json' \
    -d '{
    "name": "Elvis Presley",
    "email": "theking@elvismansion.com", "productid": "3",
    "review": "I really love the product and will recommend!" }'
    ``` 
    
    
    
## **REST API**
``` 
- POST  '/api/reviews'          (Create a new review)
- GET   '/api/reviews/approved' (Get all approved reviews)
- GET   '/api/reviews/{id}'     (Get a review)    
``` 
