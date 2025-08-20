"""
Transaction Router

This module contains all transaction-related endpoints for the HabitatTrack API.
Handles CRUD operations for transaction management with filtering capabilities.
"""

from datetime import datetime

from fastapi import APIRouter, Query, status

from models.transaction import Transaction, TransactionCreate, TransactionUpdate
from services.transaction_service import transaction_service

router = APIRouter()


@router.post(
    "",
    response_model=Transaction,
    status_code=status.HTTP_201_CREATED,
    summary="Create a new transaction",
    description=(
        "Create a new transaction with the provided details. "
        "ID and date_created are auto-generated."
    ),
)
async def create_transaction(transaction_data: TransactionCreate) -> Transaction:
    """
    Create a new transaction.

    Args:
        transaction_data: Transaction details from request body

    Returns:
        Transaction: The created transaction with generated ID and timestamp

    Raises:
        ValueError: If validation fails (handled by global exception handler)
        Exception: If creation fails (handled by global exception handler)
    """
    new_transaction = transaction_service.create_transaction(transaction_data)
    return new_transaction


@router.get(
    "",
    response_model=list[Transaction],
    status_code=status.HTTP_200_OK,
    summary="Get all transactions",
    description=(
        "Retrieve a list of all transactions with optional filtering and "
        "pagination parameters."
    ),
)
async def get_transactions(
    property_id: str | None = Query(
        None, description="Filter transactions by property ID"
    ),
    start_date: datetime | None = Query(
        None, description="Filter transactions from this date onwards"
    ),
    end_date: datetime | None = Query(
        None, description="Filter transactions up to this date"
    ),
    limit: int = Query(
        default=100,
        ge=1,
        le=1000,
        description="Maximum number of transactions to return",
    ),
    offset: int = Query(default=0, ge=0, description="Number of transactions to skip"),
) -> list[Transaction]:
    """
    Retrieve all transactions with filtering and pagination.

    Args:
        property_id: Optional filter by property ID
        start_date: Optional filter transactions from this date onwards
        end_date: Optional filter transactions up to this date
        limit: Maximum number of transactions to return (1-1000, default: 100)
        offset: Number of transactions to skip (default: 0)

    Returns:
        list[Transaction]: List of transactions

    Raises:
        Exception: If retrieval fails (handled by global exception handler)
    """
    transactions = transaction_service.get_all_transactions(
        property_id=property_id,
        start_date=start_date,
        end_date=end_date,
        limit=limit,
        offset=offset,
    )
    return transactions


@router.get(
    "/{transaction_id}",
    response_model=Transaction,
    status_code=status.HTTP_200_OK,
    summary="Get a transaction by ID",
    description="Retrieve a specific transaction by its unique identifier.",
)
async def get_transaction(transaction_id: str) -> Transaction:
    """
    Retrieve a specific transaction by ID.

    Args:
        transaction_id: The unique identifier of the transaction

    Returns:
        Transaction: The requested transaction

    Raises:
        ValueError: If transaction doesn't exist (handled by global exception handler)
        Exception: If retrieval fails (handled by global exception handler)
    """
    transaction = transaction_service.get_transaction(transaction_id)
    if transaction is None:
        raise ValueError(f"Transaction with ID '{transaction_id}' not found")
    return transaction


@router.put(
    "/{transaction_id}",
    response_model=Transaction,
    status_code=status.HTTP_200_OK,
    summary="Update a transaction",
    description="Update an existing transaction with the provided details.",
)
async def update_transaction(
    transaction_id: str, transaction_data: TransactionUpdate
) -> Transaction:
    """
    Update an existing transaction.

    Args:
        transaction_id: The unique identifier of the transaction
        transaction_data: Updated transaction details from request body

    Returns:
        Transaction: The updated transaction

    Raises:
        ValueError: If transaction doesn't exist or
        validation fails (handled by global exception handler)
        Exception: If update fails (handled by global exception handler)
    """
    updated_transaction = transaction_service.update_transaction(
        transaction_id, transaction_data
    )
    if updated_transaction is None:
        raise ValueError(f"Transaction with ID '{transaction_id}' not found")
    return updated_transaction


@router.delete(
    "/{transaction_id}",
    status_code=status.HTTP_204_NO_CONTENT,
    summary="Delete a transaction",
    description="Delete a transaction by its unique identifier.",
)
async def delete_transaction(transaction_id: str):
    """
    Delete a transaction by ID.

    Args:
        transaction_id: The unique identifier of the transaction

    Raises:
        ValueError: If transaction doesn't exist (handled by global exception handler)
        Exception: If deletion fails (handled by global exception handler)
    """
    deleted = transaction_service.delete_transaction(transaction_id)
    if not deleted:
        raise ValueError(f"Transaction with ID '{transaction_id}' not found")
    # Return 204 No Content (no response body needed)
