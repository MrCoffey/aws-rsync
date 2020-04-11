# s3sync

s3sync is a simple script used to migrate s3 objects from one bucket to another. It also updates the keys in the database.

Arguments are provided using the flags described in the CLI and can be consulted running `go run main.go --help`.

## How to Run

```bash
go build

go run main.go \
    --origin-bucket=${LEGACY_BUCKET_NAME} \
    --destination-bucket=${NEW_BUCKET_NAME} \
    --database-url=${DATABASE_URL} \
    --s3-secret-key=${SECRET_KEY} \
    --s3-access-key-id=${ACCESS_KEY_ID} \
    --s3-region=${REGION} \
    --s3-endpoint=${S3_ENDPOINT} \
    --test-mode=true
```

**Note:** In case you don't have an schema in the database it can be created by using the flag `—test-mode=true`

### It is also possible to run the whole project using docker or docker-compose

```bash
docker run -it --rm \
	-e SECRET_KEY=wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY \
	-e ACCESS_KEY_ID=AKIAIOSFODNN7EXAMPLE \
	-e REGION=us-east-1 \
	-e DATABASE_URL="test:root123@/dbname" \
	-e LEGACY_BUCKET_NAME=legacybucket \
	-e NEW_BUCKET_NAME=newbucket \
	-e S3_ENDPOINT=minio:9000 \
	coffey0container/s3sync
```

### Running the project with docker-compose

1.) Start all the services

```bash
docker-compose up -d
```

2.) Run the binary

```bash
docker-compose run s3sync
```

3.) You can login into MariaDB

```
docker-compose run mariadb bash -c "mysql -u test -p dbname"
```

4.) Following is the db schema:

    MySQL [dbname]> DESCRIBE objects;
    +------------+--------------+------+-----+---------+----------------+
    | Field      | Type         | Null | Key | Default | Extra          |
    +------------+--------------+------+-----+---------+----------------+
    | id         | int unsigned | NO   | PRI | NULL    | auto_increment |
    | created_at | datetime     | YES  |     | NULL    |                |
    | updated_at | datetime     | YES  |     | NULL    |                |
    | deleted_at | datetime     | YES  | MUL | NULL    |                |
    | path       | varchar(255) | YES  |     | NULL    |                |
    | bucket     | varchar(255) | YES  |     | NULL    |                |
    +------------+--------------+------+-----+---------+----------------+

5.) You'll have access to the Minio dashboard in http://0.0.0.0:9000

### Permissions Needed to run Using AWS S3

We will need an IAM policy providing full access to both buckets.

    {
        "Version": "2012-10-17",
        "Statement": [
            {
                "Effect": "Allow",
                "Action": "s3:*",
                "Resource": [
                    "arn:aws:s3:::legacybucket/*",
                    "arn:aws:s3:::newbucket/*"
                ]
            }
        ]
    }

### Permissions Needed to run Using Google Cloud Storage

The configuration with GCP is quite different and requires a different approach since the authentication happens using a service account, this service account is created in the GCP GUI and most of the times needs to be stored as a Json file, then in your code the gcloud SDK will require the file.

The service account should have the role **Storage Object Admin** so it can have full control of objects in the bucket (but not control of the bucket itself).

### Permissions in the database

The provided user should be owner of the table or at least to have a `GRANT` for `SELECT` and `UPDATE` operations in the target tables.
