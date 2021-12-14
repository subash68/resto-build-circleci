insert into user_type(id, name)
values (1, 'Admin');
insert into user_type(id, name)
values (2, 'Restaurant');
insert into user_type(id, name)
values (3, 'Driver');
insert into user_type(id, name)
values (4, 'Customer');

insert into address_type(type)
values ('Home');
insert into address_type(type)
values ('Work');



/** The Pass here is: pass123# */
insert into users(fullname, email, password, type)
values ('ayoub bouroumine', 'ayoub@test.com', '$2a$10$UApBLBNW.KlF9cVO4O.aau6xt2cALxTUJXhLee0EChJKb.m1gsVaK', 2);
insert into users(fullname, email, password, type)
values ('Subash Chandra Bose', 'subash@test.com', '$2a$10$UApBLBNW.KlF9cVO4O.aau6xt2cALxTUJXhLee0EChJKb.m1gsVaK', 2);
insert into users(fullname, email, password, type)
values ('Pizza Hut', 'pizzahut@test.com', '$2a$10$UApBLBNW.KlF9cVO4O.aau6xt2cALxTUJXhLee0EChJKb.m1gsVaK', 2);
insert into users(fullname, email, password, type)
values ('Hissa al Musa', 'hissa@test.com', '$2a$10$UApBLBNW.KlF9cVO4O.aau6xt2cALxTUJXhLee0EChJKb.m1gsVaK', 2);


insert into addresses(building, street, city, region, country, zipcode)
values ('Box No. 43669', '', 'Dubai', 'Dubai', 'Emirates', '43651');


insert into user_address(userId, addressId, addressType, isBilling)
values (4, 1, 1, 1);


insert into tables(name, qrCode, seats, userId)
values ('table 1', 'qrcode-1', 6, 1);
insert into tables(name, qrCode, seats, userId)
values ('table 2', 'qrcode-2', 4, 1);
insert into tables(name, qrCode, seats, userId)
values ('table 3', 'qrcode-1', 10, 1);
insert into tables(name, qrCode, seats, userId)
values ('table 4', 'qrcode-1', 2, 1);

insert into reservations(tableId, fromReserved, toReserved, reservedById, description)
values (1, '2006-01-02T15:04:05Z07:00', '2006-01-02T16:04:05Z07:00', 1, 'Description 1 2 3');
insert into reservations(tableId, fromReserved, toReserved, reservedById, description)
values (2, '2006-01-02T15:04:05Z07:00', '2006-01-02T16:04:05Z07:00', 1, 'Description 2 3 4');
insert into reservations(tableId, fromReserved, toReserved, reservedById, description)
values (3, '2006-01-02T15:04:05Z07:00', '2006-01-02T16:04:05Z07:00', 1, 'Description 122 333dsf asd asd ');
insert into reservations(tableId, fromReserved, toReserved, reservedById, description)
values (4, '2006-01-02T15:04:05Z07:00', '2006-01-02T16:04:05Z07:00', 1, 'Description ad ksdfkn ksdkfm km');

insert into categories(name, userId)
values ('cat 1', 1);
insert into categories(name, userId)
values ('cat 2', 1);
insert into categories(name, userId)
values ('cat 3', 1);

insert into addons(name, price, userId)
values ('addons 1', 19.2, 1);
insert into addons(name, price, userId)
values ('addons 2', 29.2, 1);
insert into addons(name, price, userId)
values ('addons 3', 39.2, 1);
insert into addons(name, price, userId)
values ('addons 4', 49.2, 1);

insert into order_type(orderType)
values ('Pick up');
insert into order_type(orderType)
values ('Delivery');
insert into order_type(orderType)
values ('Dine In')


insert into discount_type(discountType)
values ('Percentage');
insert into discount_type(discountType)
values ('Amount');