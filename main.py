from datetime import datetime

from fastapi import FastAPI, HTTPException, Query, status
from fastapi.responses import JSONResponse

from models.property import Property, PropertyCreate, PropertyUpdate
from models.transaction import Transaction, TransactionCreate, TransactionUpdate
from models.transaction_category import (
    TransactionCategory,
    TransactionCategoryCreate,
    TransactionCategoryUpdate,
)
from services.property_service import property_service
from services.transaction_category_service import transaction_category_service
from services.transaction_service import transaction_service

app = FastAPI(
    title="HabitatTrack API",
    description="""
    HabitatTrack API runs the backend for the HabitatTrack application,
    providing endpoints to manage properties, transaction categories, and transactions.
    """,
    version="1.0.0",
)


@app.post(
    "/properties",
    response_model=Property,
    status_code=status.HTTP_201_CREATED,
    summary="Create a new property",
    description="""
    Create a new property with the provided details.
    ID and timestamps are auto-generated.
    """,
    tags=["Properties"],
)
async def create_property(property_data: PropertyCreate) -> Property:
    """
    Create a new property.

    Args:
        property_data: Property details from request body

    Returns:
        Property: The created property with generated ID and timestamps

    Raises:
        HTTPException: 400 Bad Request if validation fails
        HTTPException: 500 Internal Server Error if creation fails
    """
    try:
        new_property = property_service.create_property(property_data)
        return new_property

    except ValueError as e:
        raise HTTPException(
            status_code=status.HTTP_400_BAD_REQUEST,
            detail=f"Invalid property data: {str(e)}",
        )
    except Exception as e:
        raise HTTPException(
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
            detail=f"An error occurred while creating the property: {str(e)}",
        )


@app.get(
    "/properties",
    response_model=list[Property],
    status_code=status.HTTP_200_OK,
    summary="Get all properties",
    description="""
    Retrieve a list of all properties with optional pagination parameters.
    """,
    tags=["Properties"],
)
async def get_properties(
    limit: int = Query(
        default=100, ge=1, le=1000, description="Maximum number of properties to return"
    ),
    offset: int = Query(default=0, ge=0, description="Number of properties to skip"),
) -> list[Property]:
    """
    Retrieve all properties with pagination.

    Args:
        limit: Maximum number of properties to return (1-1000, default: 100)
        offset: Number of properties to skip (default: 0)

    Returns:
        list[Property]: List of properties

    Raises:
        HTTPException: 500 Internal Server Error if retrieval fails
    """
    try:
        properties = property_service.get_all_properties(limit=limit, offset=offset)
        return properties

    except Exception as e:
        raise HTTPException(
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
            detail=f"An error occurred while retrieving properties: {str(e)}",
        )


@app.get(
    "/properties/{property_id}",
    response_model=Property,
    status_code=status.HTTP_200_OK,
    summary="Get a property by ID",
    description="Retrieve a specific property by its unique identifier.",
    tags=["Properties"],
)
async def get_property(property_id: str) -> Property:
    """
    Retrieve a specific property by ID.

    Args:
        property_id: The unique identifier of the property

    Returns:
        Property: The requested property

    Raises:
        HTTPException: 404 Not Found if property doesn't exist
        HTTPException: 500 Internal Server Error if retrieval fails
    """
    try:
        property_obj = property_service.get_property(property_id)
        if property_obj is None:
            raise HTTPException(
                status_code=status.HTTP_404_NOT_FOUND,
                detail=f"Property with ID '{property_id}' not found",
            )
        return property_obj

    except HTTPException:
        # Re-raise HTTP exceptions
        raise
    except Exception as e:
        raise HTTPException(
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
            detail=f"An error occurred while retrieving the property: {str(e)}",
        )


@app.put(
    "/properties/{property_id}",
    response_model=Property,
    status_code=status.HTTP_200_OK,
    summary="Update a property",
    description="Update an existing property with the provided details.",
    tags=["Properties"],
)
async def update_property(property_id: str, property_data: PropertyUpdate) -> Property:
    """
    Update an existing property.

    Args:
        property_id: The unique identifier of the property
        property_data: Updated property details from request body

    Returns:
        Property: The updated property

    Raises:
        HTTPException: 404 Not Found if property doesn't exist
        HTTPException: 400 Bad Request if validation fails
        HTTPException: 500 Internal Server Error if update fails
    """
    try:
        updated_property = property_service.update_property(property_id, property_data)
        if updated_property is None:
            raise HTTPException(
                status_code=status.HTTP_404_NOT_FOUND,
                detail=f"Property with ID '{property_id}' not found",
            )
        return updated_property

    except HTTPException:
        # Re-raise HTTP exceptions
        raise
    except ValueError as e:
        raise HTTPException(
            status_code=status.HTTP_400_BAD_REQUEST,
            detail=f"Invalid property data: {str(e)}",
        )
    except Exception as e:
        raise HTTPException(
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
            detail=f"An error occurred while updating the property: {str(e)}",
        )


@app.delete(
    "/properties/{property_id}",
    status_code=status.HTTP_204_NO_CONTENT,
    summary="Delete a property",
    description="Delete a property by its unique identifier.",
    tags=["Properties"],
)
async def delete_property(property_id: str):
    """
    Delete a property by ID.

    Args:
        property_id: The unique identifier of the property

    Raises:
        HTTPException: 404 Not Found if property doesn't exist
        HTTPException: 500 Internal Server Error if deletion fails
    """
    try:
        deleted = property_service.delete_property(property_id)
        if not deleted:
            raise HTTPException(
                status_code=status.HTTP_404_NOT_FOUND,
                detail=f"Property with ID '{property_id}' not found",
            )
        # Return 204 No Content (no response body needed)

    except HTTPException:
        # Re-raise HTTP exceptions
        raise
    except Exception as e:
        raise HTTPException(
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
            detail=f"An error occurred while deleting the property: {str(e)}",
        )


# Transaction Category Endpoints


@app.post(
    "/transaction_categories",
    response_model=TransactionCategory,
    status_code=status.HTTP_201_CREATED,
    summary="Create a new transaction category",
    description=(
        "Create a new transaction category with the provided details. "
        "ID and timestamps are auto-generated."
    ),
    tags=["Transaction Categories"],
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
        HTTPException: 400 Bad Request if validation fails
        HTTPException: 409 Conflict if category name already exists
        HTTPException: 500 Internal Server Error if creation fails
    """
    try:
        new_category = transaction_category_service.create_transaction_category(
            category_data
        )
        return new_category

    except ValueError as e:
        if "already exists" in str(e):
            raise HTTPException(status_code=status.HTTP_409_CONFLICT, detail=str(e))
        else:
            raise HTTPException(
                status_code=status.HTTP_400_BAD_REQUEST,
                detail=f"Invalid transaction category data: {str(e)}",
            )
    except Exception as e:
        raise HTTPException(
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
            detail=(
                f"An error occurred while creating the transaction category: {str(e)}"
            ),
        )


@app.get(
    "/transaction_categories",
    response_model=list[TransactionCategory],
    status_code=status.HTTP_200_OK,
    summary="Get all transaction categories",
    description="Retrieve a list of all transaction categories.",
    tags=["Transaction Categories"],
)
async def get_transaction_categories() -> list[TransactionCategory]:
    """
    Retrieve all transaction categories.

    Returns:
        list[TransactionCategory]: List of transaction categories

    Raises:
        HTTPException: 500 Internal Server Error if retrieval fails
    """
    try:
        categories = transaction_category_service.get_all_transaction_categories()
        return categories

    except Exception as e:
        raise HTTPException(
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
            detail=(
                f"An error occurred while retrieving transaction categories: {str(e)}"
            ),
        )


@app.get(
    "/transaction_categories/{category_id}",
    response_model=TransactionCategory,
    status_code=status.HTTP_200_OK,
    summary="Get a transaction category by ID",
    description="Retrieve a specific transaction category by its unique identifier.",
    tags=["Transaction Categories"],
)
async def get_transaction_category(category_id: str) -> TransactionCategory:
    """
    Retrieve a specific transaction category by ID.

    Args:
        category_id: The unique identifier of the transaction category

    Returns:
        TransactionCategory: The requested transaction category

    Raises:
        HTTPException: 404 Not Found if transaction category doesn't exist
        HTTPException: 500 Internal Server Error if retrieval fails
    """
    try:
        category = transaction_category_service.get_transaction_category(category_id)
        if category is None:
            raise HTTPException(
                status_code=status.HTTP_404_NOT_FOUND,
                detail=f"Transaction category with ID '{category_id}' not found",
            )
        return category

    except HTTPException:
        # Re-raise HTTP exceptions
        raise
    except Exception as e:
        raise HTTPException(
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
            detail=(
                f"An error occurred while retrieving the transaction category: {str(e)}"
            ),
        )


@app.put(
    "/transaction_categories/{category_id}",
    response_model=TransactionCategory,
    status_code=status.HTTP_200_OK,
    summary="Update a transaction category",
    description="Update an existing transaction category with the provided details.",
    tags=["Transaction Categories"],
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
        HTTPException: 404 Not Found if transaction category doesn't exist
        HTTPException: 400 Bad Request if validation fails
        HTTPException: 409 Conflict if updated name already exists
        HTTPException: 500 Internal Server Error if update fails
    """
    try:
        updated_category = transaction_category_service.update_transaction_category(
            category_id, category_data
        )
        if updated_category is None:
            raise HTTPException(
                status_code=status.HTTP_404_NOT_FOUND,
                detail=f"Transaction category with ID '{category_id}' not found",
            )
        return updated_category

    except HTTPException:
        # Re-raise HTTP exceptions
        raise
    except ValueError as e:
        if "already exists" in str(e):
            raise HTTPException(status_code=status.HTTP_409_CONFLICT, detail=str(e))
        else:
            raise HTTPException(
                status_code=status.HTTP_400_BAD_REQUEST,
                detail=f"Invalid transaction category data: {str(e)}",
            )
    except Exception as e:
        raise HTTPException(
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
            detail=(
                f"An error occurred while updating the transaction category: {str(e)}"
            ),
        )


@app.delete(
    "/transaction_categories/{category_id}",
    status_code=status.HTTP_204_NO_CONTENT,
    summary="Delete a transaction category",
    description="Delete a transaction category by its unique identifier.",
    tags=["Transaction Categories"],
)
async def delete_transaction_category(category_id: str):
    """
    Delete a transaction category by ID.

    Args:
        category_id: The unique identifier of the transaction category

    Raises:
        HTTPException: 404 Not Found if transaction category doesn't exist
        HTTPException: 409 Conflict if category is associated with existing transactions
        HTTPException: 500 Internal Server Error if deletion fails
    """
    try:
        deleted = transaction_category_service.delete_transaction_category(category_id)
        if not deleted:
            raise HTTPException(
                status_code=status.HTTP_404_NOT_FOUND,
                detail=f"Transaction category with ID '{category_id}' not found",
            )
        # Return 204 No Content (no response body needed)

    except HTTPException:
        # Re-raise HTTP exceptions
        raise
    except ValueError as e:
        if "associated with existing transactions" in str(e):
            raise HTTPException(status_code=status.HTTP_409_CONFLICT, detail=str(e))
        else:
            raise HTTPException(status_code=status.HTTP_400_BAD_REQUEST, detail=str(e))
    except Exception as e:
        raise HTTPException(
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
            detail=(
                f"An error occurred while deleting the transaction category: {str(e)}"
            ),
        )


# Transaction Endpoints


@app.post(
    "/transactions",
    response_model=Transaction,
    status_code=status.HTTP_201_CREATED,
    summary="Create a new transaction",
    description=(
        "Create a new transaction with the provided details. "
        "ID and date_created are auto-generated."
    ),
    tags=["Transactions"],
)
async def create_transaction(transaction_data: TransactionCreate) -> Transaction:
    """
    Create a new transaction.

    Args:
        transaction_data: Transaction details from request body

    Returns:
        Transaction: The created transaction with generated ID and timestamp

    Raises:
        HTTPException: 400 Bad Request if validation fails
        HTTPException: 500 Internal Server Error if creation fails
    """
    try:
        new_transaction = transaction_service.create_transaction(transaction_data)
        return new_transaction

    except ValueError as e:
        raise HTTPException(
            status_code=status.HTTP_400_BAD_REQUEST,
            detail=f"Invalid transaction data: {str(e)}",
        )
    except Exception as e:
        raise HTTPException(
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
            detail=f"An error occurred while creating the transaction: {str(e)}",
        )


@app.get(
    "/transactions",
    response_model=list[Transaction],
    status_code=status.HTTP_200_OK,
    summary="Get all transactions",
    description=(
        "Retrieve a list of all transactions with optional filtering and "
        "pagination parameters."
    ),
    tags=["Transactions"],
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
        HTTPException: 500 Internal Server Error if retrieval fails
    """
    try:
        transactions = transaction_service.get_all_transactions(
            property_id=property_id,
            start_date=start_date,
            end_date=end_date,
            limit=limit,
            offset=offset,
        )
        return transactions

    except Exception as e:
        raise HTTPException(
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
            detail=f"An error occurred while retrieving transactions: {str(e)}",
        )


@app.get(
    "/transactions/{transaction_id}",
    response_model=Transaction,
    status_code=status.HTTP_200_OK,
    summary="Get a transaction by ID",
    description="Retrieve a specific transaction by its unique identifier.",
    tags=["Transactions"],
)
async def get_transaction(transaction_id: str) -> Transaction:
    """
    Retrieve a specific transaction by ID.

    Args:
        transaction_id: The unique identifier of the transaction

    Returns:
        Transaction: The requested transaction

    Raises:
        HTTPException: 404 Not Found if transaction doesn't exist
        HTTPException: 500 Internal Server Error if retrieval fails
    """
    try:
        transaction = transaction_service.get_transaction(transaction_id)
        if transaction is None:
            raise HTTPException(
                status_code=status.HTTP_404_NOT_FOUND,
                detail=f"Transaction with ID '{transaction_id}' not found",
            )
        return transaction

    except HTTPException:
        # Re-raise HTTP exceptions
        raise
    except Exception as e:
        raise HTTPException(
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
            detail=f"An error occurred while retrieving the transaction: {str(e)}",
        )


@app.put(
    "/transactions/{transaction_id}",
    response_model=Transaction,
    status_code=status.HTTP_200_OK,
    summary="Update a transaction",
    description="Update an existing transaction with the provided details.",
    tags=["Transactions"],
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
        HTTPException: 404 Not Found if transaction doesn't exist
        HTTPException: 400 Bad Request if validation fails
        HTTPException: 500 Internal Server Error if update fails
    """
    try:
        updated_transaction = transaction_service.update_transaction(
            transaction_id, transaction_data
        )
        if updated_transaction is None:
            raise HTTPException(
                status_code=status.HTTP_404_NOT_FOUND,
                detail=f"Transaction with ID '{transaction_id}' not found",
            )
        return updated_transaction

    except HTTPException:
        # Re-raise HTTP exceptions
        raise
    except ValueError as e:
        raise HTTPException(
            status_code=status.HTTP_400_BAD_REQUEST,
            detail=f"Invalid transaction data: {str(e)}",
        )
    except Exception as e:
        raise HTTPException(
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
            detail=f"An error occurred while updating the transaction: {str(e)}",
        )


@app.delete(
    "/transactions/{transaction_id}",
    status_code=status.HTTP_204_NO_CONTENT,
    summary="Delete a transaction",
    description="Delete a transaction by its unique identifier.",
    tags=["Transactions"],
)
async def delete_transaction(transaction_id: str):
    """
    Delete a transaction by ID.

    Args:
        transaction_id: The unique identifier of the transaction

    Raises:
        HTTPException: 404 Not Found if transaction doesn't exist
        HTTPException: 500 Internal Server Error if deletion fails
    """
    try:
        deleted = transaction_service.delete_transaction(transaction_id)
        if not deleted:
            raise HTTPException(
                status_code=status.HTTP_404_NOT_FOUND,
                detail=f"Transaction with ID '{transaction_id}' not found",
            )
        # Return 204 No Content (no response body needed)

    except HTTPException:
        # Re-raise HTTP exceptions
        raise
    except Exception as e:
        raise HTTPException(
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
            detail=f"An error occurred while deleting the transaction: {str(e)}",
        )


@app.exception_handler(ValueError)
async def value_error_handler(request, exc):
    return JSONResponse(
        status_code=status.HTTP_400_BAD_REQUEST,
        content={"detail": f"Validation error: {str(exc)}"},
    )


@app.exception_handler(Exception)
async def exception_error_handler(request, exc):
    return JSONResponse(
        status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
        content={"detail": f"Server error: {str(exc)}"},
    )
