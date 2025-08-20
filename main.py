"""
HabitatTrack API - Main Application

Clean FastAPI application entry point with modular router structure
and centralized error handling.
"""

from fastapi import FastAPI, Request, status
from fastapi.responses import JSONResponse

from routers import properties, transaction_categories, transactions

app = FastAPI(
    title="HabitatTrack API",
    description="""
    HabitatTrack API runs the backend for the HabitatTrack application,
    providing endpoints to manage properties, transaction categories, and transactions.
    """,
    version="1.0.0",
)


# Global Exception Handlers
@app.exception_handler(ValueError)
async def value_error_handler(request: Request, exc: ValueError) -> JSONResponse:
    """
    Handle ValueError exceptions globally.

    Converts ValueError to 400 Bad Request or 404 Not Found based on message content.
    """
    error_message = str(exc)

    # Check if it's a "not found" error
    if "not found" in error_message.lower():
        return JSONResponse(
            status_code=status.HTTP_404_NOT_FOUND,
            content={"detail": error_message},
        )

    # Check if it's a conflict error (already exists, associated with
    # transactions, etc.)
    if any(
        phrase in error_message.lower()
        for phrase in ["already exists", "associated with"]
    ):
        return JSONResponse(
            status_code=status.HTTP_409_CONFLICT,
            content={"detail": error_message},
        )

    # Default to bad request for other validation errors
    return JSONResponse(
        status_code=status.HTTP_400_BAD_REQUEST,
        content={"detail": f"Validation error: {error_message}"},
    )


@app.exception_handler(Exception)
async def general_exception_handler(request: Request, exc: Exception) -> JSONResponse:
    """
    Handle unexpected exceptions globally.

    Converts unhandled exceptions to 500 Internal Server Error.
    """
    return JSONResponse(
        status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
        content={"detail": f"Internal server error: {str(exc)}"},
    )


# Router Registration
app.include_router(
    properties.router,
    prefix="/properties",
    tags=["Properties"],
)

app.include_router(
    transaction_categories.router,
    prefix="/transaction_categories",
    tags=["Transaction Categories"],
)

app.include_router(
    transactions.router,
    prefix="/transactions",
    tags=["Transactions"],
)
