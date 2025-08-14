ALTER TABLE deliveries ADD CONSTRAINT uq_deliveries_order_uid UNIQUE (order_uid);
ALTER TABLE payments   ADD CONSTRAINT uq_payments_order_uid   UNIQUE (order_uid);


