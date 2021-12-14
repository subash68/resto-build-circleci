DELIMITER //
CREATE PROCEDURE sp_GET_PRODUCT_PRICE(IN PRODUCT_ID INT)
BEGIN
    -- DECLARE PROD_NAME                VARCHAR(200)    DEFAULT '';
    -- DECLARE PRICE               FLOAT           DEFAULT 0.0;
    -- DECLARE DISCOUNT            FLOAT           DEFAULT 0.0;
    -- DECLARE DISCOUNT_TYPE       INT             DEFAULT 1;
    -- DECLARE DISCOUNTED_PRICE    FLOAT           DEFAULT 0.0;

    SELECT 
        name,
        price,
        discount,
        discountType
    -- INTO
    --     PROD_NAME,
    --     PRICE,
    --     DISCOUNT,
    --     DISCOUNT_TYPE
    FROM products
    WHERE id = PRODUCT_ID;

    -- IF DISCOUNT_TYPE = 1 THEN
    --     DISCOUNTED_PRICE = PRICE - (PRICE * (DISCOUNT / 100));
    -- ELSEIF DISCOUNT_TYPE = 2 THEN
    --     DISCOUNTED_PRICE = PRICE - DISCOUNT;
    -- END IF;

    -- SELECT NAME, PRICE, DISCOUNTED_PRICE, DISCOUNT;
END //
DELIMITER ;