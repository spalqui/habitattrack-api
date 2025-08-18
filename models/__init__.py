"""
HabitatTrack API Data Models

This package contains Pydantic models for data validation and serialization
in the HabitatTrack API, designed to work with Google Cloud Firestore.
"""

from .property import Property, PropertyCreate, PropertyUpdate

__all__ = [
    "Property",
    "PropertyCreate", 
    "PropertyUpdate"
]