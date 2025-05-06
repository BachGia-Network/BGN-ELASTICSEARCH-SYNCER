-- Insert Categories
INSERT INTO categories (id, name, parent_id) VALUES
    ('C101', 'Clothing', NULL),
    ('C102', 'T-Shirts', 'C101'),
    ('C103', 'Jeans', 'C101'),
    ('C104', 'Dresses', 'C101'),
    ('C105', 'Accessories', NULL),
    ('C106', 'Hats', 'C105'),
    ('C107', 'Belts', 'C105');

-- Insert Products
INSERT INTO products (id, name) VALUES
    ('P101', 'Classic White T-Shirt'),
    ('P102', 'Slim Fit Jeans'),
    ('P103', 'Summer Dress'),
    ('P104', 'Leather Belt'),
    ('P105', 'Baseball Cap');

-- Insert Product Categories (relationships)
INSERT INTO product_categories (product_id, category_id) VALUES
    ('P101', 'C102'),  -- White T-Shirt in T-Shirts
    ('P102', 'C103'),  -- Slim Fit Jeans in Jeans
    ('P103', 'C104'),  -- Summer Dress in Dresses
    ('P104', 'C107'),  -- Leather Belt in Belts
    ('P105', 'C106');  -- Baseball Cap in Hats

-- Insert Product Variants
INSERT INTO product_variants (id, product_id, name) VALUES
    ('V101', 'P101', 'Small'),
    ('V102', 'P101', 'Medium'),
    ('V103', 'P101', 'Large'),
    ('V104', 'P102', '30x32'),
    ('V105', 'P102', '32x32'),
    ('V106', 'P103', 'XS'),
    ('V107', 'P103', 'S'),
    ('V108', 'P103', 'M'),
    ('V109', 'P104', 'Small'),
    ('V110', 'P104', 'Medium'),
    ('V111', 'P105', 'One Size');

-- Insert Product Attributes
INSERT INTO product_attributes (id, product_id, name, value) VALUES
    ('A101', 'P101', 'Color', 'White'),
    ('A102', 'P101', 'Material', 'Cotton'),
    ('A103', 'P101', 'Fit', 'Regular'),
    ('A104', 'P102', 'Color', 'Blue'),
    ('A105', 'P102', 'Material', 'Denim'),
    ('A106', 'P102', 'Style', 'Slim'),
    ('A107', 'P103', 'Color', 'Floral'),
    ('A108', 'P103', 'Material', 'Cotton Blend'),
    ('A109', 'P103', 'Style', 'Casual'),
    ('A110', 'P104', 'Color', 'Brown'),
    ('A111', 'P104', 'Material', 'Genuine Leather'),
    ('A112', 'P105', 'Color', 'Black'),
    ('A113', 'P105', 'Material', 'Cotton'); 