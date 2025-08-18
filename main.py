from fastapi import FastAPI, HTTPException, Query, status
from fastapi.responses import JSONResponse
from typing import List

from models.property import Property, PropertyCreate, PropertyUpdate
from services.property_service import property_service

app = FastAPI(
    title="HabitatTrack API",
    description="HabitatTrack API",
    version="0.1.0"
)


@app.get("/")
async def root():
    """Root endpoint returning welcome message."""
    return {"message": "Welcome to the HabitatTrack API"}


@app.post(
    "/properties",
    response_model=Property,
    status_code=status.HTTP_201_CREATED,
    summary="Create a new property",
    description="Create a new property with the provided details. ID and timestamps are auto-generated."
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
            detail=f"Invalid property data: {str(e)}"
        )
    except Exception as e:
        raise HTTPException(
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
            detail=f"An error occurred while creating the property: {str(e)}"
        )


@app.get(
    "/properties",
    response_model=List[Property],
    status_code=status.HTTP_200_OK,
    summary="Get all properties",
    description="Retrieve a list of all properties with optional pagination parameters."
)
async def get_properties(
    limit: int = Query(default=100, ge=1, le=1000, description="Maximum number of properties to return"),
    offset: int = Query(default=0, ge=0, description="Number of properties to skip")
) -> List[Property]:
    """
    Retrieve all properties with pagination.
    
    Args:
        limit: Maximum number of properties to return (1-1000, default: 100)
        offset: Number of properties to skip (default: 0)
        
    Returns:
        List[Property]: List of properties
        
    Raises:
        HTTPException: 500 Internal Server Error if retrieval fails
    """
    try:
        properties = property_service.get_all_properties(limit=limit, offset=offset)
        return properties
        
    except Exception as e:
        raise HTTPException(
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
            detail=f"An error occurred while retrieving properties: {str(e)}"
        )


@app.get(
    "/properties/{property_id}",
    response_model=Property,
    status_code=status.HTTP_200_OK,
    summary="Get a property by ID",
    description="Retrieve a specific property by its unique identifier."
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
                detail=f"Property with ID '{property_id}' not found"
            )
        return property_obj
        
    except HTTPException:
        # Re-raise HTTP exceptions
        raise
    except Exception as e:
        raise HTTPException(
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
            detail=f"An error occurred while retrieving the property: {str(e)}"
        )


@app.put(
    "/properties/{property_id}",
    response_model=Property,
    status_code=status.HTTP_200_OK,
    summary="Update a property",
    description="Update an existing property with the provided details."
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
                detail=f"Property with ID '{property_id}' not found"
            )
        return updated_property
        
    except HTTPException:
        # Re-raise HTTP exceptions
        raise
    except ValueError as e:
        raise HTTPException(
            status_code=status.HTTP_400_BAD_REQUEST,
            detail=f"Invalid property data: {str(e)}"
        )
    except Exception as e:
        raise HTTPException(
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
            detail=f"An error occurred while updating the property: {str(e)}"
        )


@app.delete(
    "/properties/{property_id}",
    status_code=status.HTTP_204_NO_CONTENT,
    summary="Delete a property",
    description="Delete a property by its unique identifier."
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
                detail=f"Property with ID '{property_id}' not found"
            )
        # Return 204 No Content (no response body needed)
        
    except HTTPException:
        # Re-raise HTTP exceptions
        raise
    except Exception as e:
        raise HTTPException(
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
            detail=f"An error occurred while deleting the property: {str(e)}"
        )


@app.exception_handler(ValueError)
async def value_error_handler(request, exc):
    return JSONResponse(
        status_code=status.HTTP_400_BAD_REQUEST,
        content={"detail": f"Validation error: {str(exc)}"}
    )

@app.exception_handler(Exception)
async def exception_error_handler(request, exc):
    return JSONResponse(
        status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
        content={"detail": f"Server error: {str(exc)}"}
    )