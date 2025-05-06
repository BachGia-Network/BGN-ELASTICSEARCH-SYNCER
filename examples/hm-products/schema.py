import click
import sqlalchemy as sa
from sqlalchemy.orm import DeclarativeBase, Mapped, mapped_column, relationship
from typing import Optional

from pgsync.base import pg_engine
from pgsync.utils import config_loader, get_config


class Base(DeclarativeBase):
    pass


class Product(Base):
    __tablename__ = "products"
    __table_args__ = (sa.PrimaryKeyConstraint("id"),)
    
    id: Mapped[str] = mapped_column(sa.String)
    name: Mapped[str] = mapped_column(sa.String)

    categories = relationship("Categories", secondary="product_categories")
    variants = relationship("ProductVariants", back_populates="product")
    attributes = relationship("ProductAttributes", back_populates="product")


class Categories(Base):
    __tablename__ = "categories"
    __table_args__ = (sa.PrimaryKeyConstraint("id"),)
    
    id: Mapped[str] = mapped_column(sa.String)
    name: Mapped[str] = mapped_column(sa.String)
    parent_id: Mapped[Optional[str]] = mapped_column(sa.String, nullable=True)
    
    products = relationship("Product", secondary="product_categories")


class ProductCategories(Base):
    __tablename__ = "product_categories"
    __table_args__ = (sa.PrimaryKeyConstraint("product_id", "category_id"),)
    
    product_id: Mapped[str] = mapped_column(sa.String, sa.ForeignKey("products.id"))
    category_id: Mapped[str] = mapped_column(sa.String, sa.ForeignKey("categories.id"))


class ProductVariants(Base):
    __tablename__ = "product_variants"
    __table_args__ = (sa.PrimaryKeyConstraint("id"),)
    
    id: Mapped[str] = mapped_column(sa.String)
    product_id: Mapped[str] = mapped_column(sa.String, sa.ForeignKey("products.id"))
    name: Mapped[str] = mapped_column(sa.String)
    
    product = relationship("Product", back_populates="variants")


class ProductAttributes(Base):
    __tablename__ = "product_attributes"
    __table_args__ = (sa.PrimaryKeyConstraint("id"),)
    
    id: Mapped[str] = mapped_column(sa.String)
    product_id: Mapped[str] = mapped_column(sa.String, sa.ForeignKey("products.id"))
    name: Mapped[str] = mapped_column(sa.String)
    value: Mapped[str] = mapped_column(sa.String)
    
    product = relationship("Product", back_populates="attributes")


def setup(config: str) -> None:
    for doc in config_loader(config):
        database: str = doc.get("database", doc["index"])
        with pg_engine(database) as engine:
            Base.metadata.drop_all(engine)
            Base.metadata.create_all(engine)


@click.command()
@click.option(
    "--config",
    "-c",
    help="Schema config",
    type=click.Path(exists=True),
)
def main(config: str) -> None:
    config: str = get_config(config)
    setup(config)


if __name__ == "__main__":
    main()
