CREATE TABLE whatsapp_messages
(
    id                VARCHAR(255) PRIMARY KEY,
    name              VARCHAR(255) NOT NULL,
    whatsapp_id        VARCHAR(255) NOT NULL,
    sender_phone_number VARCHAR(255) NOT NULL,
    timestamp         VARCHAR(255) NOT NULL,
    type              VARCHAR(255) NOT NULL,
    text              VARCHAR(255) NOT NULL
)