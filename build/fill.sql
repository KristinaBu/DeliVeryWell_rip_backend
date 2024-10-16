-- Вставляем данные
INSERT INTO users (login, password, is_admin) VALUES
('пользователь1', 'пароль1', false),
('админ', 'пароль_админ', true);

INSERT INTO delivery_items (image, title, price, description) VALUES
('http://127.0.0.1:9000/images/1.png', 'Посылки', 1075, 'Доставка посылок весом менее 1 кг в Новокузнецке'),
('http://127.0.0.1:9000/images/2.png', 'Спортивные товары', 430, 'Доставка спортивных товаров, включая аксессуары (кольца, гантели) и спортивное оборудование (швейцарские мячи)'),
('http://127.0.0.1:9000/images/3.png', 'Пицца', 149, 'Доставка пиццы из ресторанов ДоДо, Maestrello, FoodBand за 149 рублей. При покупке 4 штук - доставка 59 рублей'),
('http://127.0.0.1:9000/images/4.png', 'Цветы', 799, 'Доставка цветов, букетов, упаковочных материалов в Москве'),
('http://127.0.0.1:9000/images/5.png', 'Суши', 200, 'Доставка суши в Москве за 200 рублей. При покупке товара на сумму более 500 рублей - доставка 100 рублей (Суши Мастер)');


-- Создание функции для проверки черновиковitem_r
CREATE FUNCTION check_draft_request() RETURNS trigger AS $$
BEGIN
    IF (SELECT COUNT(*) FROM delivery_requests WHERE user_id = NEW.user_id AND status = 'черновик') > 0 THEN
        RAISE EXCEPTION 'User already has a draft request';
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Создание триггера для проверки статуса черновиков
CREATE TRIGGER draft_request_trigger
    BEFORE INSERT OR UPDATE ON delivery_requests
    FOR EACH ROW
    WHEN (NEW.status = 'черновик')
EXECUTE FUNCTION check_draft_request();

ALTER TABLE delivery_items
    ALTER COLUMN image SET DEFAULT '';

ALTER TABLE delivery_items
    ALTER COLUMN image TYPE varchar(255),
    ALTER COLUMN image SET DEFAULT '',
    ALTER COLUMN image SET NOT NULL,
    ALTER COLUMN title TYPE varchar(255),
    ALTER COLUMN title SET NOT NULL,
    ALTER COLUMN price TYPE integer,
    ALTER COLUMN price SET NOT NULL,
    ALTER COLUMN description TYPE text,
    ALTER COLUMN description SET NOT NULL,
    ALTER COLUMN is_delete SET DEFAULT false;

