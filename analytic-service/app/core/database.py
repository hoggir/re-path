"""MongoDB database connection manager using Beanie ODM."""

from typing import List, Optional

from beanie import init_beanie
from motor.motor_asyncio import AsyncIOMotorClient, AsyncIOMotorDatabase
from pydantic import BaseModel

from app.core.config import settings


class DatabaseManager:
    """MongoDB database connection manager."""

    client: Optional[AsyncIOMotorClient] = None
    database: Optional[AsyncIOMotorDatabase] = None

    @classmethod
    async def connect_to_database(cls, document_models: List[type]) -> None:
        """
        Connect to MongoDB and initialize Beanie ODM.

        Args:
            document_models: List of Beanie document model classes
        """
        print(f"📦 Connecting to MongoDB at {settings.mongodb_url}...")

        # Create MongoDB client
        cls.client = AsyncIOMotorClient(
            settings.mongodb_url,
            maxPoolSize=settings.mongodb_max_connections,
            minPoolSize=settings.mongodb_min_connections,
        )

        # Get database
        cls.database = cls.client[settings.mongodb_database]

        # Initialize Beanie with document models
        await init_beanie(
            database=cls.database,
            document_models=document_models,
        )

        print(f"✅ Connected to MongoDB database: {settings.mongodb_database}")

    @classmethod
    async def close_database_connection(cls) -> None:
        """Close MongoDB connection."""
        if cls.client:
            print("📦 Closing MongoDB connection...")
            cls.client.close()
            cls.client = None
            cls.database = None
            print("✅ MongoDB connection closed")

    @classmethod
    def get_database(cls) -> Optional[AsyncIOMotorDatabase]:
        """Get MongoDB database instance."""
        return cls.database

    @classmethod
    async def ping(cls) -> bool:
        """
        Ping MongoDB to check connection.

        Returns:
            bool: True if connection is alive, False otherwise
        """
        if not cls.client:
            return False

        try:
            await cls.client.admin.command("ping")
            return True
        except Exception as e:
            print(f"❌ MongoDB ping failed: {e}")
            return False


class DatabaseHealthCheck(BaseModel):
    """Database health check response."""

    connected: bool
    database: Optional[str] = None
    ping: bool = False


async def get_database_health() -> DatabaseHealthCheck:
    """
    Get database health status.

    Returns:
        DatabaseHealthCheck: Database health information
    """
    if DatabaseManager.client is None or DatabaseManager.database is None:
        return DatabaseHealthCheck(connected=False)

    ping_result = await DatabaseManager.ping()

    return DatabaseHealthCheck(
        connected=True,
        database=DatabaseManager.database.name,
        ping=ping_result,
    )


# Dependency for getting database
async def get_database() -> AsyncIOMotorDatabase:
    """
    FastAPI dependency to get database instance.

    Returns:
        AsyncIOMotorDatabase: MongoDB database instance

    Raises:
        RuntimeError: If database is not initialized
    """
    if DatabaseManager.database is None:
        raise RuntimeError("Database is not initialized")
    return DatabaseManager.database
