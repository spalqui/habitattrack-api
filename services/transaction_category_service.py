"""
Transaction Category Service

This module contains the business logic for transaction category operations,
using Google Cloud Firestore as the database.
"""

from datetime import datetime, timezone
from typing import Optional

from models.transaction_category import TransactionCategory, TransactionCategoryCreate, TransactionCategoryUpdate
from services.firestore_client import firestore_client


class TransactionCategoryService:
    """
    Service class for transaction category operations using Firestore database.
    
    This service handles CRUD operations for transaction categories using Google Cloud Firestore.
    """

    def __init__(self):
        """Initialize the service with shared Firestore client."""
        # Use the shared Firestore client
        self.db = firestore_client.client
        
        # Define collection name
        self.collection_name = "transaction_categories"

    def create_transaction_category(self, category_data: TransactionCategoryCreate) -> TransactionCategory:
        """
        Create a new transaction category with auto-generated fields.
        
        Args:
            category_data: The transaction category data from the request
            
        Returns:
            TransactionCategory: The created transaction category with generated id and timestamps
            
        Raises:
            ValueError: If category data is invalid or if category name already exists
        """
        # Check if category name already exists
        existing_category = self._get_category_by_name(category_data.name)
        if existing_category is not None:
            raise ValueError(f"Transaction category with name '{category_data.name}' already exists")
        
        # Get current timestamp
        now = datetime.now(timezone.utc)
        
        # Prepare category data for Firestore
        category_dict = category_data.model_dump()
        category_dict.update({
            "date_created": now,
            "date_last_updated": now
        })
        
        # Add document to Firestore (Firestore will auto-generate the ID)
        doc_ref = self.db.collection(self.collection_name).add(category_dict)[1]
        
        # Get the generated ID and add it to the category data
        category_dict["id"] = doc_ref.id
        
        # Create and return TransactionCategory instance
        new_category = TransactionCategory(**category_dict)
        
        return new_category

    def get_transaction_category(self, category_id: str) -> Optional[TransactionCategory]:
        """
        Retrieve a transaction category by ID.
        
        Args:
            category_id: The unique identifier of the transaction category
            
        Returns:
            TransactionCategory or None if not found
        """
        try:
            # Get document from Firestore
            doc_ref = self.db.collection(self.collection_name).document(category_id)
            doc = doc_ref.get()
            
            if doc.exists:
                # Convert Firestore document to TransactionCategory object
                category_dict = doc.to_dict()
                category_dict["id"] = doc.id
                return TransactionCategory(**category_dict)
            else:
                return None
        except Exception as e:
            # Log error in real application
            print(f"Error retrieving transaction category {category_id}: {e}")
            return None

    def get_all_transaction_categories(self) -> list[TransactionCategory]:
        """
        Retrieve all transaction categories.
        
        Returns:
            List of transaction categories from Firestore
        """
        try:
            # Get documents from the transaction_categories collection
            query = self.db.collection(self.collection_name)
            docs = query.stream()
            
            categories = []
            for doc in docs:
                category_dict = doc.to_dict()
                category_dict["id"] = doc.id
                categories.append(TransactionCategory(**category_dict))
            
            return categories
        except Exception as e:
            # Log error in real application
            print(f"Error retrieving all transaction categories: {e}")
            return []

    def update_transaction_category(self, category_id: str, category_data: TransactionCategoryUpdate) -> Optional[TransactionCategory]:
        """
        Update an existing transaction category.
        
        Args:
            category_id: The unique identifier of the transaction category
            category_data: The updated transaction category data
            
        Returns:
            Updated TransactionCategory or None if not found
            
        Raises:
            ValueError: If category data is invalid or if updated name already exists
        """
        try:
            # Check if category exists
            doc_ref = self.db.collection(self.collection_name).document(category_id)
            doc = doc_ref.get()
            
            if not doc.exists:
                return None
            
            # Prepare update data (exclude None values)
            update_dict = {k: v for k, v in category_data.model_dump().items() if v is not None}
            
            # If name is being updated, check for duplicates
            if "name" in update_dict:
                existing_category = self._get_category_by_name(update_dict["name"])
                if existing_category is not None and existing_category.id != category_id:
                    raise ValueError(f"Transaction category with name '{update_dict['name']}' already exists")
            
            # Add updated timestamp
            update_dict["date_last_updated"] = datetime.now(timezone.utc)
            
            # Update document in Firestore
            doc_ref.update(update_dict)
            
            # Retrieve and return updated category
            updated_doc = doc_ref.get()
            if updated_doc.exists:
                category_dict = updated_doc.to_dict()
                category_dict["id"] = updated_doc.id
                return TransactionCategory(**category_dict)
            else:
                return None
                
        except Exception as e:
            # Log error in real application
            print(f"Error updating transaction category {category_id}: {e}")
            raise ValueError(f"Failed to update transaction category: {str(e)}")

    def delete_transaction_category(self, category_id: str) -> bool:
        """
        Delete a transaction category by ID.
        
        Args:
            category_id: The unique identifier of the transaction category
            
        Returns:
            True if deleted successfully, False if not found
            
        Raises:
            ValueError: If category is associated with existing transactions
        """
        try:
            # Check if category exists
            doc_ref = self.db.collection(self.collection_name).document(category_id)
            doc = doc_ref.get()
            
            if not doc.exists:
                return False
            
            # Check if category is referenced by any transactions
            # Note: This would require checking the transactions collection
            # For now, we'll skip this check as transactions aren't implemented yet
            # TODO: Add check for existing transactions when transaction model is implemented
            
            # Delete document from Firestore
            doc_ref.delete()
            return True
            
        except Exception as e:
            # Log error in real application
            print(f"Error deleting transaction category {category_id}: {e}")
            return False

    def _get_category_by_name(self, name: str) -> Optional[TransactionCategory]:
        """
        Helper method to get a transaction category by name.
        
        Args:
            name: The name of the transaction category
            
        Returns:
            TransactionCategory or None if not found
        """
        try:
            # Query Firestore for category with the given name
            query = self.db.collection(self.collection_name).where("name", "==", name).limit(1)
            docs = list(query.stream())
            
            if docs:
                doc = docs[0]
                category_dict = doc.to_dict()
                category_dict["id"] = doc.id
                return TransactionCategory(**category_dict)
            else:
                return None
        except Exception as e:
            # Log error in real application
            print(f"Error checking category name existence: {e}")
            return None


# Create a singleton instance for the application
transaction_category_service = TransactionCategoryService()