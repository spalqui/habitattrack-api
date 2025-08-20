"""
Transaction Category Router

This module contains all transaction category-related endpoints for the
HabitatTrack API.
Handles CRUD operations for transaction category management.
"""

from fastapi import APIRouter, status

from models.transaction_category import (
    TransactionCategory,
    TransactionCategoryCreate,
    TransactionCategoryUpdate,
)
from services.transaction_category_service import transaction_category_service

router = APIRouter()


@router.post(
    "",
    response_model=TransactionCategory,
    status_code=status.HTTP_201_CREATED,
    summary="Create a new transaction category",
    description=(
        "Create a new transaction category with the provided details. "
        "ID and timestamps are auto-generated."
    ),
)
async def create_transaction_category(
    category_data: TransactionCategoryCreate,
) -> TransactionCategory:
    """
    Create a new transaction category.

    Args:
        category_data: Transaction category details from request body

    Returns:
        TransactionCategory: The created transaction category with generated ID
        and timestamps

    Raises:
        ValueError: If validation fails or category name already exists
            (handled by global exception handler)
        Exception: If creation fails (handled by global exception handler)
    """
    new_category = transaction_category_service.create_transaction_category(
        category_data
    )
    return new_category


@router.get(
    "",
    response_model=list[TransactionCategory],
    status_code=status.HTTP_200_OK,
    summary="Get all transaction categories",
    description="Retrieve a list of all transaction categories.",
)
async def get_transaction_categories() -> list[TransactionCategory]:
    """
    Retrieve all transaction categories.

    Returns:
        list[TransactionCategory]: List of transaction categories

    Raises:
        Exception: If retrieval fails (handled by global exception handler)
    """
    categories = transaction_category_service.get_all_transaction_categories()
    return categories


@router.get(
    "/{category_id}",
    response_model=TransactionCategory,
    status_code=status.HTTP_200_OK,
    summary="Get a transaction category by ID",
    description="Retrieve a specific transaction category by its unique identifier.",
)
async def get_transaction_category(category_id: str) -> TransactionCategory:
    """
    Retrieve a specific transaction category by ID.

    Args:
        category_id: The unique identifier of the transaction category

    Returns:
        TransactionCategory: The requested transaction category

    Raises:
        ValueError: If transaction category doesn't exist (handled by global
            exception handler)
        Exception: If retrieval fails (handled by global exception handler)
    """
    category = transaction_category_service.get_transaction_category(category_id)
    if category is None:
        raise ValueError(f"Transaction category with ID '{category_id}' not found")
    return category


@router.put(
    "/{category_id}",
    response_model=TransactionCategory,
    status_code=status.HTTP_200_OK,
    summary="Update a transaction category",
    description="Update an existing transaction category with the provided details.",
)
async def update_transaction_category(
    category_id: str, category_data: TransactionCategoryUpdate
) -> TransactionCategory:
    """
    Update an existing transaction category.

    Args:
        category_id: The unique identifier of the transaction category
        category_data: Updated transaction category details from request body

    Returns:
        TransactionCategory: The updated transaction category

    Raises:
        ValueError: If transaction category doesn't exist, validation fails,
            or updated name already exists (handled by global exception handler)
        Exception: If update fails (handled by global exception handler)
    """
    updated_category = transaction_category_service.update_transaction_category(
        category_id, category_data
    )
    if updated_category is None:
        raise ValueError(f"Transaction category with ID '{category_id}' not found")
    return updated_category


@router.delete(
    "/{category_id}",
    status_code=status.HTTP_204_NO_CONTENT,
    summary="Delete a transaction category",
    description="Delete a transaction category by its unique identifier.",
)
async def delete_transaction_category(category_id: str):
    """
    Delete a transaction category by ID.

    Args:
        category_id: The unique identifier of the transaction category

    Raises:
        ValueError: If transaction category doesn't exist or is associated with
            existing transactions (handled by global exception handler)
        Exception: If deletion fails (handled by global exception handler)
    """
    deleted = transaction_category_service.delete_transaction_category(category_id)
    if not deleted:
        raise ValueError(f"Transaction category with ID '{category_id}' not found")
    # Return 204 No Content (no response body needed)
