
create table accounts (
    id integer primary key autoincrement,
    funds real
);

create table history (
    unix integer unique, 
    product text,
    low real,
    high real,
    open real,
    closing real,
    volume real
);

create table orders (
    product text,
    price real,
    size real,
);