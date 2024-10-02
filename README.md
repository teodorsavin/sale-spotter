# SaleSpotter

This is a brief description of the project.

## Setup

To set up the project, follow these steps:

1. Clone the repository: `git clone https://github.com/teodorsavin/sale-spotter.git`
2. Change to the project directory: `cd sale-spotter`
3. Install the dependencies: `make setup`

## Usage

To run the project locally using Docker Compose, execute the following command:

```bash
make run
```

This will build and run the project containers locally.

Access the PhpMyAdmin database at `http://0.0.0.0:9090/`
Before requesting /api/products import the `./database-sample/ah_bonus.sql` file in the database just to have the initial tables and structure.

## Cleanup

If you want to start over fresh, you can use the following command:

```bash
make nuke
```

This command will clean up all the project dependencies and artifacts.

## Interact with the API

The api has the following endpoints:
- `POST /login` - Get the Bearer token
  - parameters: `username`
  - curl command:
    - `curl -X POST http://0.0.0.0:8080/login -H "Content-Type: application/json" -d '{"username": "teodor"}'`
- `GET /api/brands` - Get Brands
  - curl command:
    - `curl http://0.0.0.0:8080/api/brands -H "Authorization: Bearer #TOKEN#" | jq`
- `GET /api/products` - Get Products
  - curl command:
    - `curl http://0.0.0.0:8080/api/products -H "Authorization: Bearer #TOKEN#" | jq`
    - This just saves the products in the DB. It doesn't return anything.
