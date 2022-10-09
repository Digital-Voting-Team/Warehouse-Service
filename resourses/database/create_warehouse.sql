create table INGREDIENT
(
    ID   NUMBER generated as identity
        constraint "ingredient_pk"
            primary key,
    NAME VARCHAR(50) not null
);

create table ADDRESS
(
    ID          NUMBER generated as identity
        constraint "address_pk"
            primary key,
    BUILDING    NUMBER      not null,
    STREET      VARCHAR(50) not null,
    CITY        VARCHAR(50) not null,
    DISTRICT    VARCHAR(50) not null,
    REGION      VARCHAR(50) not null,
    POSTAL_CODE VARCHAR(10) not null
);

create table WAREHOUSE
(
    ID         NUMBER generated as identity
        constraint "warehouse_pk"
            primary key,
    CAFE_ID    NUMBER not null,
    ADDRESS_ID NUMBER not null
        constraint "warehouse_address_null_fk"
            references ADDRESS,
    CAPACITY   NUMBER not null
);

create table DELIVERY
(
    ID             NUMBER not null
        constraint "delivery_pk"
            primary key,
    SOURCE_ID      NUMBER not null
        constraint "delivery_warehouse_null_fk"
            references WAREHOUSE,
    DESTINATION_ID NUMBER not null
        constraint "delivery_warehouse_id_fk"
            references WAREHOUSE,
    DELIVERY_PRICE FLOAT  not null,
    DELIVERY_DATE  DATE   not null
);

create table INGREDIENT_WAREHOUSE
(
    ID              NUMBER      not null
        constraint "ingredient_warehouse_pk"
            primary key,
    INGREDIENT_ID   NUMBER      not null
        constraint "ingredient_warehouse_ingredient_null_fk"
            references INGREDIENT,
    WAREHOUSE_ID    NUMBER      not null
        constraint "ingredient_warehouse_warehouse_null_fk"
            references WAREHOUSE,
    QUANTITY        NUMBER      not null,
    ORIGIN          VARCHAR(50) not null,
    PRICE           FLOAT       not null,
    EXPIRATION_DATE DATE        not null,
    DELIVERY_DATE   DATE        not null
);

comment on column INGREDIENT_WAREHOUSE.PRICE is 'Per item';

