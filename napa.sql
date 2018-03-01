
create table accounts (id integer primary key autoincrement, funds real);
create table book (product text, unix integer, buy integer, price real, size real, complete integer);
create table btc_usd (unix integer unique, low real, high real, open real, closing real, volume real);
create table eth_usd (unix integer unique, low real, high real, open real, closing real, volume real);
create table ltc_usd (unix integer unique, low real, high real, open real, closing real, volume real);
