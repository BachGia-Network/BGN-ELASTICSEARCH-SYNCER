import click
from schema import Product, Categories, ProductCategories, ProductVariants, ProductAttributes
from sqlalchemy.orm import sessionmaker

from pgsync.base import pg_engine, subtransactions
from pgsync.helper import teardown
from pgsync.utils import config_loader, get_config


@click.command()
@click.option(
    "--config",
    "-c",
    help="Schema config",
    type=click.Path(exists=True),
)
def main(config: str) -> None:
    config: str = get_config(config)
    teardown(drop_db=False, config=config)
    doc: dict = next(config_loader(config))
    database: str = doc.get("database", doc["index"])
    with pg_engine(database) as engine:
        Session = sessionmaker(bind=engine, autoflush=True)
        session = Session()

        # Create products first
        products = [
            Product(
                id="P001",
                name="Classic White T-Shirt"
            ),
            Product(
                id="P002",
                name="Slim Fit Jeans"
            )
        ]

        # Create categories
        categories = [
            Categories(
                id="C001",
                name="Clothing",
                parent_id=None
            ),
            Categories(
                id="C002",
                name="T-Shirts",
                parent_id="C001"
            ),
            Categories(
                id="C003",
                name="Jeans",
                parent_id="C001"
            )
        ]

        # Create product variants
        variants = [
            ProductVariants(
                id="V001",
                product_id="P001",
                name="Small"
            ),
            ProductVariants(
                id="V002",
                product_id="P001",
                name="Medium"
            ),
            ProductVariants(
                id="V003",
                product_id="P002",
                name="30x32"
            )
        ]

        # Create product attributes
        attributes = [
            ProductAttributes(
                id="A001",
                product_id="P001",
                name="Color",
                value="White"
            ),
            ProductAttributes(
                id="A002",
                product_id="P001",
                name="Material",
                value="Cotton"
            ),
            ProductAttributes(
                id="A003",
                product_id="P002",
                name="Color",
                value="Blue"
            )
        ]

        # Create product-category relationships
        product_categories = [
            ProductCategories(
                product_id="P001",
                category_id="C002"
            ),
            ProductCategories(
                product_id="P002",
                category_id="C003"
            )
        ]

        with subtransactions(session):
            # First add products and categories
            session.add_all(products)
            session.add_all(categories)
            session.flush()  # Ensure products and categories are created before adding relationships
            
            # Then add variants and attributes
            session.add_all(variants)
            session.add_all(attributes)
            session.flush()  # Ensure variants and attributes are created
            
            # Finally add the relationships
            session.add_all(product_categories)


if __name__ == "__main__":
    main()
