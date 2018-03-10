
create table accounts (
    id integer primary key autoincrement,
    product text unique,
    funds real
);

create table history (
    unix integer primary key, 
    product text,
    low real,
    high real,
    open real,
    closing real,
    volume real
);

create table orders (
    id integer primary key autoincrement,
    exchange_id text
);
