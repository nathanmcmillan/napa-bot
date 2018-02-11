
create table trades (
    id text,
    price text,
    size text,
    product_id text,
    side text,
    stp text,
    type text,
    time_in_force text,
    post_only boolean,
    created_at text,
    fill_fees text,
    filled_size text,
    executed_value text,
    settled boolean
);
