create table if not exists payment
(
    id                  serial primary key,
    name                varchar,
    legal_entity_id     integer,
    legal_entity_name   varchar not null,
    amount              numeric not null,
    start_date          timestamp,
    end_date            timestamp,
    payment_type        varchar not null,
    status              varchar not null,
    bill                varchar,
    bill_payment        varchar,
    billing_at          timestamp,
    bill_payment_at     timestamp,
    confirm_payment_at  timestamp,
    created_at          timestamp not null,
    updated_at          timestamp,
    foreign key (legal_entity_id) references legal_entity (id)
);