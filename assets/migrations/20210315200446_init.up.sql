CREATE extension IF NOT EXISTS "uuid-ossp";

CREATE TABLE "webhooks" (
    "id" uuid PRIMARY KEY UNIQUE NOT NULL DEFAULT uuid_generate_v4(),
    "payload" TEXT NOT NULL,
    "payload_hash" CHAR(40) NOT NULL,
    "received_at" TIMESTAMP WITH TIME ZONE NOT NULL
);
