"""
Property Service

This module contains the business logic for property operations,
using Google Cloud Firestore as the database.
"""

from datetime import datetime, timezone

from models.property import Property, PropertyCreate, PropertyUpdate
from services.firestore_client import firestore_client


class PropertyService:
    """
    Service class for property operations using Firestore database.

    This service handles CRUD operations for properties using Google Cloud Firestore.
    """

    def __init__(self):
        """Initialize the service with shared Firestore client."""
        # Use the shared Firestore client
        self.db = firestore_client.client

        # Define collection name
        self.collection_name = "properties"

    def create_property(self, property_data: PropertyCreate) -> Property:
        """
        Create a new property with auto-generated fields.

        Args:
            property_data: The property data from the request

        Returns:
            Property: The created property with generated id and timestamps

        Raises:
            ValueError: If property data is invalid
        """
        # Get current timestamp
        now = datetime.now(timezone.utc)

        # Prepare property data for Firestore
        property_dict = property_data.model_dump()
        property_dict.update({"date_created": now, "date_last_updated": now})

        # Add document to Firestore (Firestore will auto-generate the ID)
        doc_ref = self.db.collection(self.collection_name).add(property_dict)[1]

        # Get the generated ID and add it to the property data
        property_dict["id"] = doc_ref.id

        # Create and return Property instance
        new_property = Property(**property_dict)

        return new_property

    def get_property(self, property_id: str) -> Property | None:
        """
        Retrieve a property by ID.

        Args:
            property_id: The unique identifier of the property

        Returns:
            Property or None if not found
        """
        try:
            # Get document from Firestore
            doc_ref = self.db.collection(self.collection_name).document(property_id)
            doc = doc_ref.get()

            if doc.exists:
                # Convert Firestore document to Property object
                property_dict = doc.to_dict()
                property_dict["id"] = doc.id
                return Property(**property_dict)
            else:
                return None
        except Exception as e:
            # Log error in real application
            print(f"Error retrieving property {property_id}: {e}")
            return None

    def get_all_properties(self, limit: int = 100, offset: int = 0) -> list[Property]:
        """
        Retrieve all properties with pagination support.

        Args:
            limit: Maximum number of properties to return (default: 100)
            offset: Number of properties to skip (default: 0)

        Returns:
            List of properties from Firestore
        """
        try:
            # Get documents from the properties collection with pagination
            query = self.db.collection(self.collection_name).offset(offset).limit(limit)
            docs = query.stream()

            properties = []
            for doc in docs:
                property_dict = doc.to_dict()
                property_dict["id"] = doc.id
                properties.append(Property(**property_dict))

            return properties
        except Exception as e:
            # Log error in real application
            print(f"Error retrieving all properties: {e}")
            return []

    def update_property(
        self, property_id: str, property_data: PropertyUpdate
    ) -> Property | None:
        """
        Update an existing property.

        Args:
            property_id: The unique identifier of the property
            property_data: The updated property data

        Returns:
            Updated Property or None if not found

        Raises:
            ValueError: If property data is invalid
        """
        try:
            # Check if property exists
            doc_ref = self.db.collection(self.collection_name).document(property_id)
            doc = doc_ref.get()

            if not doc.exists:
                return None

            # Prepare update data (exclude None values)
            update_dict = {
                k: v for k, v in property_data.model_dump().items() if v is not None
            }

            # Add updated timestamp
            update_dict["date_last_updated"] = datetime.now(timezone.utc)

            # Update document in Firestore
            doc_ref.update(update_dict)

            # Retrieve and return updated property
            updated_doc = doc_ref.get()
            if updated_doc.exists:
                property_dict = updated_doc.to_dict()
                property_dict["id"] = updated_doc.id
                return Property(**property_dict)
            else:
                return None

        except Exception as e:
            # Log error in real application
            print(f"Error updating property {property_id}: {e}")
            raise ValueError(f"Failed to update property: {str(e)}")

    def delete_property(self, property_id: str) -> bool:
        """
        Delete a property by ID.

        Args:
            property_id: The unique identifier of the property

        Returns:
            True if deleted successfully, False if not found
        """
        try:
            # Check if property exists
            doc_ref = self.db.collection(self.collection_name).document(property_id)
            doc = doc_ref.get()

            if not doc.exists:
                return False

            # Delete document from Firestore
            doc_ref.delete()
            return True

        except Exception as e:
            # Log error in real application
            print(f"Error deleting property {property_id}: {e}")
            return False


# Create a singleton instance for the application
property_service = PropertyService()
