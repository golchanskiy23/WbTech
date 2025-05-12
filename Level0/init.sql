CREATE TABLE IF NOT EXISTS orders (
    order_uid TEXT PRIMARY KEY,
    track_number TEXT NOT NULL,
    entry TEXT,
    locale TEXT,
    internal_signature TEXT,
    customer_id TEXT,
    delivery_service TEXT,
    shard_key TEXT,
    sm_id INTEGER,
    data_created TEXT,
    oof_shard TEXT
);

CREATE TABLE IF NOT EXISTS deliveries (
    order_uid TEXT PRIMARY KEY REFERENCES orders(order_uid),
    name TEXT,
    phone TEXT,
    zip TEXT,
    city TEXT,
    address TEXT,
    region TEXT,
    email TEXT
);

CREATE TABLE IF NOT EXISTS payments (
    order_uid TEXT PRIMARY KEY REFERENCES orders(order_uid),
    transaction TEXT,
    request_id TEXT,
    currency TEXT,
    provider TEXT,
    amount INTEGER,
    payment_dt BIGINT,
    bank TEXT,
    delivery_cost INTEGER,
    goods_total INTEGER,
    custom_fee INTEGER
);

CREATE TABLE IF NOT EXISTS items (
    id SERIAL PRIMARY KEY,
    order_uid TEXT REFERENCES orders(order_uid),
    chrt_id BIGINT,
    track_number TEXT,
    price INTEGER,
    rid TEXT,
    name TEXT,
    sale INTEGER,
    size TEXT,
    total_price INTEGER,
    nm_id BIGINT,
    brand TEXT,
    status INTEGER
);

TRUNCATE TABLE orders, deliveries, payments, items CASCADE;

INSERT INTO orders (
    order_uid, track_number, entry, locale, internal_signature,
    customer_id, delivery_service, shard_key, sm_id, data_created, oof_shard
) VALUES (
             'b563feb7b2b84b6test', 'WBILMTESTTRACK', 'WBIL', 'en', 'distributed',
             'test', 'meest', '9', 99, '2021-11-26T06:22:19Z', '1'
         );

INSERT INTO deliveries (
    order_uid, name, phone, zip, city, address, region, email
) VALUES (
             'b563feb7b2b84b6test', 'Test Testov', '+9720000000', '2639809',
             'Kiryat Mozkin', 'Ploshad Mira 15', 'Kraiot', 'test@gmail.com'
         );

INSERT INTO payments (
    order_uid, transaction, request_id, currency, provider,
    amount, payment_dt, bank, delivery_cost, goods_total, custom_fee
) VALUES (
             'b563feb7b2b84b6test', 'b563feb7b2b84b6test', '32356bdhe', 'USD', 'wbpay',
             1817, 1637907727, 'alpha', 1500, 317, 0
         );

INSERT INTO items (
    order_uid, chrt_id, track_number, price, rid,
    name, sale, size, total_price, nm_id, brand, status
) VALUES (
             'b563feb7b2b84b6test', 9934930, 'WBILMTESTTRACK', 453, 'ab4219087a764ae0btest',
             'Mascaras', 30, '0', 317, 2389212, 'Vivienne Sabo', 202
         );
