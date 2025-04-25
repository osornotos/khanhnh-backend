> khanhnh's home assignment

## Deployment

Prerequisites: Install Go, Docker and Docker Compose

1. Start PostgreSQL and Redis
```sh
docker-compose --up
```

2. Build docker image

```sh
docker build -t khanhnh-backend:1.0.0 .
```

3. Migrate database and run application

```sh
docker build -t khanhnh-backend:1.0.0 .
docker run khanhnh-backend:1.0.0 migrate
docker run -d -p 8080:8080 khanhnh-backend:1.0.0 api
```

4. Run test
```sh
I'm using mockery to generate mock at tests folder (use mock.sh script to generate mock)
I've added a command to run the test cases in the Dockerfile (RUN go test ./... -v) -> I comment it out because it takes a lot of time to run the test cases and I don't want to run it every time I build the image.
but I think it should be done locally (local machine), or in CI/CD pipeline (real world application).
```

5. Import Postman collection using json file in docs folder

6. Get jwt token
```sh
Use users/login endpoint to get jwt token
```

7. Generate test data
```sh
Use products/generate-test-data endpoint to generate test data
It will run in background and run in about 3 minutes
You can check the progress in the container log or rerun the API to know if the data generation is in progress
("Another test data generation process is running" in log)
```

## Explanation
```
Basically, This application is a simple REST API that allows you to manage products, users, categories and visualize statistics as well.
- Using cobra to create a command line application that can run migrate and api commands separately.
- Using docker to pack the application to run anywhere that install docker. 
- Using onion architecture to separate the application into layers. It's easy to test and maintain when the codebase becomes large but hard to understand for new comers at the beginning.
At first, I think about using simple codebase but it's not a good idea when the codebase becomes large.
The application is a simple REST API that allows you to create, read, update, and delete users. 
- Using PostgreSQL as the database, written in Go, Fiber as HTTP Server, Go Bun as ORM, Redis as cache/distributed lock, JWT for authentication. 
- I've changed database schema a little bit: 
+ product-category is 1-1 so I've added category_id to products table. Removed relation table product_category.
+ created an index for product_id column in products table to get dashboard faster.
+ some indexes for products table to search faster.

- I don't have much time to implement other features like CRUD user/wishlist/category, so I just implemented product features and dashboard.
+ In the future, login feature will be implemented using argon2 to hash. To implement role based access control, after login, I will generate correct jwt based on their role. 
After that I will change authorization to check whether that user is admin or not. Currently, I just check if the user is logged in or not based on jwt expiry time.
```

```
Real-time Data Processing (Focus on real-time problem solving):
- I've implemented dashboard feature that shows the number of products in each category already. 
+ I created a index for category_id column in products table to get dashboard faster then cache the result in redis (in a single key). Whenever a new product is created, I will delete the cache. 
+ It can be improved by store 4 keys for each category and update the cache but not delete it.
```

```
Log Aggregation Optimization (Focus on eﬃcient log handling):
- I'm using zerolog for logging. The configuration defined in setupLogger function in root.go. It overrided the global logger from zerolog
We can integrate with 3rd party monitoring tools like Sentry (without an separate adapter and call to Sentry each failure). I've written a function to send log to sentry in case of failure
```

```
Product Search Optimization (Focus on data search performance):
1. To generate test data, I'm using a worker pool and send dummy datas to worker, after exceeded batch size, I will insert them to database. 
- I must set a small batch size to avoid timeout error when inserting to postgresql (1000)
- But I think it can be improved by using a task queue like go machinery background worker to process the data in the background to achive fault tolerance.
- I also think about using COPY query to get the data from a CSV file and insert them to products table. MAYBE it can be faster than using batch insert.

2. I've created these indexes for products table. For description and name columns, I'm using GIN index with description/name columns to search faster cause of ILIKE (case insensitive).
CREATE EXTENSION IF NOT EXISTS pg_trgm;
CREATE INDEX products_description_idx ON products USING GIN (description gin_trgm_ops);
CREATE INDEX products_name_idx ON products USING GIN (name gin_trgm_ops);
CREATE INDEX products_created_at_idx ON products (created_at);
CREATE INDEX products_category_id_idx ON products (category_id);
CREATE INDEX products_status_idx ON products (status);
3. With 10 million of records, my get products endpoint takes about 2s to get the data, of course with pagination (I've set default page size to 20).
```
