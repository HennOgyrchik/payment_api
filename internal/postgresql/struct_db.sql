create table users (
                       id serial,
                       cash integer not null default 0);

create table transactions (
                              id serial,
                              user_id integer not null,
                              service_id integer not null,
                              order_id integer not null,
                              cost integer not null,
                              type varchar(6) not null);

create or replace function check_cash()
returns trigger as $BODY$
declare cash integer;
Begin
cash:=(select users.cash from users where id=NEW.USER_ID);
if NEW.COST>cash then
return NULL;
Elseif NEW.type='debet' then
update users set cash =users.cash-new.cost where id=new.user_id;
return NEW;
Elseif NEW.type='credit' then

update users set cash =users.cash+new.cost where id=new.user_id;
return NEW;
Else
return NULL;
end if;
end;
$BODY$
language plpgsql;

create or replace trigger check_balance before insert on transactions for each row execute function check_cash();