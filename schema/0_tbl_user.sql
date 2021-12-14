CREATE TABLE user_type
(
    id         INT NOT NULL AUTO_INCREMENT,
    name       VARCHAR(100) default '',
    createdAt  TIMESTAMP    DEFAULT CURRENT_TIMESTAMP,
    modifiedAt TIMESTAMP    DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
);

CREATE TABLE users
(
    id                INT          NOT NULL AUTO_INCREMENT,
    fullname          VARCHAR(100)          default '',
    email             VARCHAR(100) NOT NULL default '',
    password          VARCHAR(100) NOT NULL default '',
    phone             VARCHAR(20)           default '',
    type              INT                   default 0,
    cuisine           varchar(100)          default '',
    status            TINYINT(1) default 0,
    everyday          TINYINT(1) default 0,
    profileImageId    varchar(1000)         default '',
    shopLogoId        varchar(1000)         default '',
    shopBannerId      varchar(1000)         default '',
    isVeg             TINYINT(1) default 0,
    mealService       TINYINT(1) default 0,
    partyCatering     TINYINT(1) default 0,
    deliveryTakeAway  TINYINT(1) default 0,
    delivery          TINYINT(1) default 0,
    freeDelivery      TINYINT(1) default 0,
    offerType         int                   default 0,
    offer             int                   default 0,
    offerAmount       int                   default 0,
    maxDeliveryTime   int                   default 0,
    description       varchar(1000)         default '',
    location          varchar(1000)         default '',
    locLongitude      varchar(30)           default '',
    locLatitude       varchar(30)           default '',
    notificationToken varchar(1000)         default '',
    createdAt         TIMESTAMP             DEFAULT CURRENT_TIMESTAMP,
    modifiedAt        TIMESTAMP             DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
);

CREATE TABLE address_type
(
    id         INT         NOT NULL AUTO_INCREMENT,
    type       VARCHAR(20) NOT NULL DEFAULT '',
    createAt   TIMESTAMP            DEFAULT CURRENT_TIMESTAMP,
    modifiedAt TIMESTAMP            DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
);

CREATE TABLE addresses
(
    id         INT NOT NULL AUTO_INCREMENT,
    unit       VARCHAR(30) DEFAULT '',
    building   VARCHAR(30) DEFAULT '',
    street     VARCHAR(30) DEFAULT '',
    city       VARCHAR(30) DEFAULT '',
    region     VARCHAR(30) DEFAULT '',
    country    VARCHAR(30) DEFAULT '',
    zipcode    VARCHAR(30) DEFAULT '',
    createdAt  TIMESTAMP   DEFAULT CURRENT_TIMESTAMP,
    modifiedAt TIMESTAMP   DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
);

CREATE TABLE user_address
(
    id          INT NOT NULL AUTO_INCREMENT,
    userId      INT NOT NULL,
    addressId   INT NOT NULL,
    addressType INT NOT NULL,
    isBilling   BOOLEAN,
    createAt    TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    modifiedAt  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
);

CREATE TABLE sms_codes
(
    id         INT NOT NULL AUTO_INCREMENT,
    user_id    INT NOT NULL,
    code       VARCHAR(10),
    status     INT,
    createdAt  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    modifiedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
);

create TABLE open_time
(
    id         INT NOT NULL AUTO_INCREMENT,
    dayName    varchar(20) default '',
    fromOpen   varchar(20) default '',
    toOpen     varchar(20) default '',
    userId     INT         default 0,
    createdAt  TIMESTAMP   DEFAULT CURRENT_TIMESTAMP,
    modifiedAt TIMESTAMP   DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
);

create table tables
(
    id                  INT NOT NULL AUTO_INCREMENT,
    name                varchar(20)  default '',
    qrCode              varchar(100) default '',
    seats               int          default 0,
    userId              INT          default 0,
    isOpenToReservation TINYINT(1) default 0,
    createdAt           TIMESTAMP    DEFAULT CURRENT_TIMESTAMP,
    modifiedAt          TIMESTAMP    DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
);

create table reservations
(
    id           INT NOT NULL AUTO_INCREMENT,
    tableId      int           default 0,
    fromReserved varchar(100)  default '',
    toReserved   varchar(100)  default '',
    reservedById INT           default 0,
    description  varchar(1000) default '',
    createdAt    TIMESTAMP     DEFAULT CURRENT_TIMESTAMP,
    modifiedAt   TIMESTAMP     DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
);


CREATE TABLE categories
(
    id            INT          NOT NULL AUTO_INCREMENT,
    name          VARCHAR(100) NOT NULL default '',
    description   VARCHAR(350)          default '',
    status        TINYINT(1) default 0,
    categoryOrder INT                   default 0,
    userId        INT          NOT NULL default 0,
    createdAt     TIMESTAMP             DEFAULT CURRENT_TIMESTAMP,
    modifiedAt    TIMESTAMP             DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
);



CREATE TABLE addons
(
    id         INT          NOT NULL AUTO_INCREMENT,
    name       VARCHAR(100) NOT NULL default '',
    price      FLOAT        NOT NULL default 0.0,
    userId     INT          NOT NULL default 0,
    createdAt  TIMESTAMP             DEFAULT CURRENT_TIMESTAMP,
    modifiedAt TIMESTAMP             DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
);


CREATE TABLE products
(
    id           INT          NOT NULL AUTO_INCREMENT,
    name         VARCHAR(200) NOT NULL default '',
    incredients  VARCHAR(500)          default '',
    categoryId   INT          NOT NULL default 0,
    status       TINYINT(1) default 0,
    menuOrder    INT                   default 0,
    isFeatured   TINYINT(1) default 0,
    productImage varchar(1000)         default '',
    position     INT                   default 0,
    price        FLOAT        NOT NULL default 0.0,
    discount     FLOAT        NOT NULL default 0.0,
    discountType INT          NOT NULL default 0,
    userId       INT          NOT NULL default 0,
    createdAt    TIMESTAMP             DEFAULT CURRENT_TIMESTAMP,
    modifiedAt   TIMESTAMP             DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
);

CREATE TABLE discount_type
(
    id           INT NOT NULL AUTO_INCREMENT,
    discountType VARCHAR(30),
    createdAt    TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    modifiedAt   TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
);

CREATE TABLE product_addons
(
    id         INT NOT NULL AUTO_INCREMENT,
    productId  INT NOT NULL default 0,
    addonsId   INT NOT NULL default 0,
    createdAt  TIMESTAMP    DEFAULT CURRENT_TIMESTAMP,
    modifiedAt TIMESTAMP    DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
);

CREATE TABLE order_type
(
    id         INT NOT NULL AUTO_INCREMENT,
    orderType  VARCHAR(30) DEFAULT '',
    createdAt  TIMESTAMP   DEFAULT CURRENT_TIMESTAMP,
    modifiedAt TIMESTAMP   DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
);

CREATE TABLE cart
(
    id              INT          NOT NULL AUTO_INCREMENT,
    userId          INT          NOT NULL,
    orderType       INT          NOT NULL DEFAULT 0,
    deliveryAddress VARCHAR(500) NOT NULL,
    instructions    VARCHAR(650)          DEFAULT '',
    coupon          VARCHAR(20)           DEFAULT '',
    hasCoupon       TINYINT               DEFAULT 0,
    shipCost        FLOAT                 DEFAULT 0.0,
    cartTotal       FLOAT                 DEFAULT 0.0,
    cartState       varchar(100)          default 'Requested',
    createdAt       TIMESTAMP             DEFAULT CURRENT_TIMESTAMP,
    modifiedAt      TIMESTAMP             DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
);

CREATE TABLE cart_comment
(
    id         INT NOT NULL AUTO_INCREMENT,
    cartId     INT NOT NULL,
    userId     int not null,
    comment    varchar(1000) default '',
    createdAt  TIMESTAMP     DEFAULT CURRENT_TIMESTAMP,
    modifiedAt TIMESTAMP     DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
);

CREATE TABLE cart_items
(
    id         INT NOT NULL AUTO_INCREMENT,
    cartId     INT NOT NULL,
    itemId     INT NOT NULL,
    itemCount  INT NOT NULL DEFAULT 1,
    itemPrice  FLOAT        DEFAULT 0.0,
    createdAt  TIMESTAMP    DEFAULT CURRENT_TIMESTAMP,
    modifiedAt TIMESTAMP    DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
);
