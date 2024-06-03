create table if not exists orders
(
	id serial primary key, 
    order_uid varchar,
    track_number varchar,
    entry varchar,
    locale varchar,
    internal_signature varchar,
    customer_id varchar,
    delivery_service varchar,
    shardkey varchar,
    sm_id integer,
    date_created timestamp,
    oof_shard varchar
	);

create table if not exists delivery
(
    order_id integer,
    name varchar,
    phone varchar,
    zip varchar,
    city varchar,
    address varchar,
    region varchar,
    email varchar,
    
    constraint fk_delivery_order_id foreign key (order_id) references orders (id)
);

create table if not exists payment
(
    order_id integer,
	transaction varchar,
    request_id varchar,
    currency varchar,
    provider varchar,
    amount integer,
    payment_dt integer,
    bank varchar,
    delivery_cost integer,
    goods_total integer,
    custom_fee integer,
    
    constraint fk_payment_order_id foreign key (order_id) references orders (id)
);

create table if not exists items
(
	order_id integer,
    chrt_id integer,
    track_number varchar,
    price integer,
    rid varchar,
    name varchar,
    sale integer,
    size varchar,
    total_price integer,
    nm_id integer,
    brand varchar,
    status integer,

    constraint fk_items_order_id foreign key (order_id) references orders (id)
    );