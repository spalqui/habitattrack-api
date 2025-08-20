"""
Property Router

This module contains all property-related endpoints for the HabitatTrack API.
Handles CRUD operations for property management.
"""

from fastapi import APIRouter, Query, status

from models.property import Property, PropertyCreate, PropertyUpdate
from services.property_service import property_service

router = APIRouter()


@router.post(
    "",
    response_model=Property,
    status_code=status.HTTP_201_CREATED,
    summary="Create a new property",
    description="""
    Create a new property with the provided details.
    ID and timestamps are auto-generated.
    """,
)
async def create_property(property_data: PropertyCreate) -> Property:
    """
    Create a new property.

    Args:
        property_data: Property details from request body

    Returns:
        Property: The created property with generated ID and timestamps

    Raises:
        ValueError: If validation fails (handled by global exception handler)
        Exception: If creation fails (handled by global exception handler)
    """
    new_property = property_service.create_property(property_data)
    return new_property


@router.get(
    "",
    response_model=list[Property],
    status_code=status.HTTP_200_OK,
    summary="Get all properties",
    description="""
    Retrieve a list of all properties with optional pagination parameters.
    """,
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
        Exception: If retrieval fails (handled by global exception handler)
    """
    properties = property_service.get_all_properties(limit=limit, offset=offset)
    return properties


@router.get(
    "/{property_id}",
    response_model=Property,
    status_code=status.HTTP_200_OK,
    summary="Get a property by ID",
    description="Retrieve a specific property by its unique identifier.",
)
async def get_property(property_id: str) -> Property:
    """
    Retrieve a specific property by ID.

    Args:
        property_id: The unique identifier of the property

    Returns:
        Property: The requested property

    Raises:
        ValueError: If property doesn't exist (handled by global exception handler)
        Exception: If retrieval fails (handled by global exception handler)
    """
    property_obj = property_service.get_property(property_id)
    if property_obj is None:
        raise ValueError(f"Property with ID '{property_id}' not found")
    return property_obj


@router.put(
    "/{property_id}",
    response_model=Property,
    status_code=status.HTTP_200_OK,
    summary="Update a property",
    description="Update an existing property with the provided details.",
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
        ValueError: If property doesn't exist or validation fails (handled by
            global exception handler)
        Exception: If update fails (handled by global exception handler)
    """
    updated_property = property_service.update_property(property_id, property_data)
    if updated_property is None:
        raise ValueError(f"Property with ID '{property_id}' not found")
    return updated_property


@router.delete(
    "/{property_id}",
    status_code=status.HTTP_204_NO_CONTENT,
    summary="Delete a property",
    description="Delete a property by its unique identifier.",
)
async def delete_property(property_id: str):
    """
    Delete a property by ID.

    Args:
        property_id: The unique identifier of the property

    Raises:
        ValueError: If property doesn't exist (handled by global exception handler)
        Exception: If deletion fails (handled by global exception handler)
    """
    deleted = property_service.delete_property(property_id)
    if not deleted:
        raise ValueError(f"Property with ID '{property_id}' not found")
    # Return 204 No Content (no response body needed)
