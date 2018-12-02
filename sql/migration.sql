ALTER TABLE adventureworks.productreview
  CHANGE ReviewDate ReviewDate TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  CHANGE ModifiedDate ModifiedDate TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  ADD CONSTRAINT fk_product_id FOREIGN KEY (ProductID) REFERENCES adventureworks.product(ProductID),
  ALTER COLUMN Rating SET DEFAULT 0,
  ADD Status varchar(255);