create table users (
                       id integer unique not null,
                       cash integer default 0 check (cash>=0));

create table transactions (
                              id serial,
                              user_id integer not null,
                              service_id integer,
                              order_id integer ,
                              cost integer not null,
                              type varchar(20) not null);

create unique index order_type_idx on transactions (order_id, type);

CREATE OR REPLACE function check_cash()
RETURNS TRIGGER AS $BODY$
DECLARE cash integer;
BEGIN
CASE NEW.type
    WHEN 'replenishment' THEN
        update users set cash =users.cash+new.cost where id=new.user_id;
        return NEW;
    WHEN 'revenue' THEN
        update users set cash =users.cash-new.cost where id=-1;
        return NEW;
    WHEN 'buy' THEN
        cash:=(select users.cash from users where id=NEW.USER_ID);
        IF NEW.COST>cash THEN
            return NULL;
        ELSE
            update users set cash =users.cash-new.cost where id=new.user_id;
            update users set cash=users.cash+new.cost where id=-1;
            return NEW;
        END IF;
    ELSE return NULL;
END CASE;
END;
$BODY$
language plpgsql;

create or replace trigger check_balance before insert on transactions for each row execute function check_cash();

