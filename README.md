# ERAL Promo Library Go

ERAL Promo Library Go is a promo library for ERAL. It is written in Golang and uses MySQL as the database.

## Getting Started

To get started with ERAL Promo Library Go, you need to have the following:

- [Go 1.16 or later](https://go.dev/dl/)
- [MySQL 5.6 or later](https://dev.mysql.com/downloads/mysql/)

### Installing

To install ERAL Promo Library Go, you can use the following command:

```bash
go  get  github.com/fritz-immanuel/eral-promo-library-go
```

This will download the latest version of ERAL Promo Library Go and install it in your `$GOPATH/bin` directory.

## Running the program

### 1. Setting up the database schema

Before you can run the program, you need to set up the database schema. You can do this by running the following query in your MySQL database.

```sql
CREATE DATABASE <your_preffered_db_name>;
```

### 2. Setting up the `.env` file

The `.env` file is used to store configuration settings for the program. You can create a new `.env` file by copying the `.env.example` file and renaming it to `.env`.

Here is an example of what the `.env` file might look like:

```dotenv
{
  "SERVER_NAME": <server_name>,

  "DB_CONNECTION_STRING":"<user>:<password>@(host:port)/<db_name>?parseTime=true",
  "PORT_APPS": ":9034",

  "APP_URL": "http://localhost:9034",

  "ANDROID_APP_MINIMUM_VERSION": "1.0.0",
  "IOS_APP_MINIMUM_VERSION": "1.0.0",

  "EXTERNAL_URL": "",
  "EXTERNAL_TOKEN": "",
  "EXTERNAL_ACCESS_TOKEN": "",

  "REDIS_ADDR": "localhost:6379",
  "REDIS_TIME_OUT": "259200",
  "REDIS_DB": "0",
  "REDIS_PASSWORD": "",

  "SEND_WHATSAPP_API": "",
  "SEND_WHATSAPP_TOKEN": "",

  "TELE_BOT_TOKEN": "",
  "TELE_GROUP_ID": "",

  "FIREBASE_SERVER_KEY":<fb_server_key>,
  "FIREBASE_SENDER_ID":<fb_sender_id>,
  "FIREBASE_BUCKET_URL":<fb_bucket_url>,
  "FIREBASE_AUTH_FILE_PATH":<fb_auth_file_path>,

  "WHITELISTED_IPS": "0.0.0.0"
}
```

### 3. Congrats! You are now set-up for running.

You may now execute `go run main.go` to start the program.

## Contributing

We welcome contributions to ERAL Promo Library Go.
Please read the [CONTRIBUTING.md](CONTRIBUTING.md) file for more information.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
