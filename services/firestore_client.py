"""
Firestore Client Configuration

This module provides a shared Firestore client instance for the application.
Following the singleton pattern to avoid creating multiple connections.
"""

import os

from dotenv import load_dotenv
from google.cloud import firestore

# Load environment variables
load_dotenv()

class FirestoreClient:
    """Singleton Firestore client"""
    
    _instance = None
    _client = None
    
    def __new__(cls):
        if cls._instance is None:
            cls._instance = super().__new__(cls)
        return cls._instance
    
    @property
    def client(self) -> firestore.Client:
        if self._client is None:
            project_id = os.getenv("GOOGLE_CLOUD_PROJECT")
            if not project_id:
                raise ValueError("GOOGLE_CLOUD_PROJECT environment variable is required")

            database_name = os.getenv("FIRESTORE_DATABASE_NAME")

            self._client = firestore.Client(project=project_id, database=database_name)
        
        return self._client

# Create shared instance
firestore_client = FirestoreClient()