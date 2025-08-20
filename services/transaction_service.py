"""
Transaction Service

This module contains the business logic for transaction operations,
using Google Cloud Firestore as the database.
"""

from datetime import datetime, timezone

from models.transaction import Transaction, TransactionCreate, TransactionUpdate
from services.firestore_client import firestore_client
from services.property_service import property_service
from services.transaction_category_service import transaction_category_service


class TransactionService:
    """
    Service class for transaction operations using Firestore database.
    
    This service handles CRUD operations for transactions using Google Cloud Firestore.
    """

    def __init__(self):
        """Initialize the service with shared Firestore client."""
        # Use the shared Firestore client
        self.db = firestore_client.client
        
        # Define collection name
        self.collection_name = "transactions"

    def create_transaction(self, transaction_data: TransactionCreate) -> Transaction:
        """
        Create a new transaction with auto-generated fields.
        
        Args:
            transaction_data: The transaction data from the request
            
        Returns:
            Transaction: The created transaction with generated id and timestamp
            
        Raises:
            ValueError: If transaction data is invalid, transaction category doesn't exist, or property doesn't exist
        """
        # Validate transaction category exists
        category = transaction_category_service.get_transaction_category(transaction_data.transaction_category_id)
        if category is None:
            raise ValueError(f"Transaction category with ID '{transaction_data.transaction_category_id}' does not exist")
        
        # Validate property exists if property_id is provided
        if transaction_data.property_id is not None:
            property_obj = property_service.get_property(transaction_data.property_id)
            if property_obj is None:
                raise ValueError(f"Property with ID '{transaction_data.property_id}' does not exist")
        
        # Get current timestamp
        now = datetime.now(timezone.utc)
        
        # Prepare transaction data for Firestore
        transaction_dict = transaction_data.model_dump()
        
        # Convert datetime objects to ensure proper serialization
        if isinstance(transaction_dict.get('transaction_date'), datetime):
            transaction_dict['transaction_date'] = transaction_dict['transaction_date']
        
        # Add auto-generated fields
        transaction_dict.update({
            "date_created": now
        })
        
        # Add document to Firestore (Firestore will auto-generate the ID)
        doc_ref = self.db.collection(self.collection_name).add(transaction_dict)[1]
        
        # Get the generated ID and add it to the transaction data
        transaction_dict["id"] = doc_ref.id
        
        # Create and return Transaction instance
        new_transaction = Transaction(**transaction_dict)
        
        return new_transaction

    def get_transaction(self, transaction_id: str) -> Transaction | None:
        """
        Retrieve a transaction by ID.
        
        Args:
            transaction_id: The unique identifier of the transaction
            
        Returns:
            Transaction or None if not found
        """
        try:
            # Get document from Firestore
            doc_ref = self.db.collection(self.collection_name).document(transaction_id)
            doc = doc_ref.get()
            
            if doc.exists:
                # Convert Firestore document to Transaction object
                transaction_dict = doc.to_dict()
                transaction_dict["id"] = doc.id
                return Transaction(**transaction_dict)
            else:
                return None
        except Exception as e:
            # Log error in real application
            print(f"Error retrieving transaction {transaction_id}: {e}")
            return None

    def get_all_transactions(
        self, 
        property_id: str | None = None,
        start_date: datetime | None = None,
        end_date: datetime | None = None,
        limit: int = 100, 
        offset: int = 0
    ) -> list[Transaction]:
        """
        Retrieve all transactions with filtering and pagination support.
        
        Args:
            property_id: Optional filter by property ID
            start_date: Optional filter transactions from this date onwards
            end_date: Optional filter transactions up to this date
            limit: Maximum number of transactions to return (default: 100)
            offset: Number of transactions to skip (default: 0)
        
        Returns:
            List of transactions from Firestore
        """
        try:
            # Start with base query
            query = self.db.collection(self.collection_name)
            
            # Apply filters
            if property_id:
                query = query.where("property_id", "==", property_id)
            
            if start_date:
                query = query.where("transaction_date", ">=", start_date)
                
            if end_date:
                query = query.where("transaction_date", "<=", end_date)
            
            # Apply pagination
            query = query.offset(offset).limit(limit)
            
            # Order by transaction_date descending for better UX
            query = query.order_by("transaction_date", direction="DESCENDING")
            
            docs = query.stream()
            
            transactions = []
            for doc in docs:
                transaction_dict = doc.to_dict()
                transaction_dict["id"] = doc.id
                transactions.append(Transaction(**transaction_dict))
            
            return transactions
        except Exception as e:
            # Log error in real application
            print(f"Error retrieving transactions: {e}")
            return []

    def update_transaction(self, transaction_id: str, transaction_data: TransactionUpdate) -> Transaction | None:
        """
        Update an existing transaction.
        
        Args:
            transaction_id: The unique identifier of the transaction
            transaction_data: The updated transaction data
            
        Returns:
            Updated Transaction or None if not found
            
        Raises:
            ValueError: If transaction data is invalid, transaction category doesn't exist, or property doesn't exist
        """
        try:
            # Check if transaction exists
            doc_ref = self.db.collection(self.collection_name).document(transaction_id)
            doc = doc_ref.get()
            
            if not doc.exists:
                return None
            
            # Validate transaction category exists if being updated
            if transaction_data.transaction_category_id is not None:
                category = transaction_category_service.get_transaction_category(transaction_data.transaction_category_id)
                if category is None:
                    raise ValueError(f"Transaction category with ID '{transaction_data.transaction_category_id}' does not exist")
            
            # Validate property exists if property_id is being updated and is not None
            if transaction_data.property_id is not None:
                property_obj = property_service.get_property(transaction_data.property_id)
                if property_obj is None:
                    raise ValueError(f"Property with ID '{transaction_data.property_id}' does not exist")
            
            # Prepare update data (exclude None values)
            update_dict = {k: v for k, v in transaction_data.model_dump().items() if v is not None}
            
            # Ensure datetime objects are properly handled
            if 'transaction_date' in update_dict and isinstance(update_dict['transaction_date'], datetime):
                update_dict['transaction_date'] = update_dict['transaction_date']
            
            # Update document in Firestore
            doc_ref.update(update_dict)
            
            # Retrieve and return updated transaction
            updated_doc = doc_ref.get()
            if updated_doc.exists:
                transaction_dict = updated_doc.to_dict()
                transaction_dict["id"] = updated_doc.id
                return Transaction(**transaction_dict)
            else:
                return None
                
        except Exception as e:
            # Log error in real application
            print(f"Error updating transaction {transaction_id}: {e}")
            raise ValueError(f"Failed to update transaction: {str(e)}")

    def delete_transaction(self, transaction_id: str) -> bool:
        """
        Delete a transaction by ID.
        
        Args:
            transaction_id: The unique identifier of the transaction
            
        Returns:
            True if deleted successfully, False if not found
        """
        try:
            # Check if transaction exists
            doc_ref = self.db.collection(self.collection_name).document(transaction_id)
            doc = doc_ref.get()
            
            if not doc.exists:
                return False
            
            # Delete document from Firestore
            doc_ref.delete()
            return True
            
        except Exception as e:
            # Log error in real application
            print(f"Error deleting transaction {transaction_id}: {e}")
            return False


# Create a singleton instance for the application
transaction_service = TransactionService()