"""Main FastAPI application entry point."""

import logging
import sys
from collections.abc import AsyncGenerator
from contextlib import asynccontextmanager

from fastapi import FastAPI
from fastapi.exceptions import RequestValidationError
from fastapi.middleware.cors import CORSMiddleware
from fastapi.responses import RedirectResponse
from pydantic import ValidationError

from app.api.v1 import dashboard, health
from app.core.config import settings
from app.core.database import DatabaseManager
from app.core.exceptions import (
    AppException,
    app_exception_handler,
    generic_exception_handler,
    validation_exception_handler,
)
from app.core.rabbitmq import RabbitMQConsumerManager
from app.models import ALL_MODELS
from app.services.message_handler import create_message_handler

# Configure logging
logging.basicConfig(
    level=logging.INFO,
    format="%(asctime)s - %(name)s - %(levelname)s - %(message)s",
    handlers=[logging.StreamHandler(sys.stdout)],
)

# Set log level for specific loggers
logging.getLogger("app.services.message_handler").setLevel(logging.INFO)
logging.getLogger("app.core.rabbitmq").setLevel(logging.INFO)


@asynccontextmanager
async def lifespan(app: FastAPI) -> AsyncGenerator[None, None]:
    """
    Application lifespan manager.

    Handles startup and shutdown events.
    """
    # Startup
    print(f"🚀 Starting {settings.app_name} v{settings.app_version}")
    print(f"🌍 Environment: {settings.app_env}")
    print(f"📍 Server: http://{settings.host}:{settings.port}")
    print(f"📚 API Docs: http://{settings.host}:{settings.port}/docs")

    # Connect to MongoDB
    try:
        await DatabaseManager.connect_to_database(document_models=ALL_MODELS)
    except Exception as e:
        print(f"⚠️  Failed to connect to MongoDB: {e}")
        print("⚠️  Service will continue without database connection")

    # Start RabbitMQ Consumer
    try:
        await RabbitMQConsumerManager.start_consumer(create_message_handler)
    except Exception as e:
        print(f"⚠️  Failed to start RabbitMQ consumer: {e}")
        print("⚠️  Service will continue without RabbitMQ consumer")

    yield

    # Shutdown
    print(f"👋 Shutting down {settings.app_name}")
    await RabbitMQConsumerManager.stop_consumer()
    await DatabaseManager.close_database_connection()


# Initialize FastAPI application
app = FastAPI(
    title=settings.app_name,
    version=settings.app_version,
    description="Re:Path Analytics Service - FastAPI microservice for analytics",
    docs_url="/docs",
    redoc_url="/redoc",
    openapi_url="/openapi.json",
    lifespan=lifespan,
)

# Configure CORS
app.add_middleware(
    CORSMiddleware,
    allow_origins=settings.cors_origins,
    allow_credentials=settings.cors_allow_credentials,
    allow_methods=settings.cors_allow_methods,
    allow_headers=settings.cors_allow_headers,
)

# Register exception handlers
app.add_exception_handler(AppException, app_exception_handler)
app.add_exception_handler(RequestValidationError, validation_exception_handler)
app.add_exception_handler(ValidationError, validation_exception_handler)
app.add_exception_handler(Exception, generic_exception_handler)

# Include routers
app.include_router(health.router, prefix=settings.api_v1_prefix)
app.include_router(dashboard.router, prefix=settings.api_v1_prefix)


@app.get("/", include_in_schema=False)
async def root() -> RedirectResponse:
    """Redirect root to API documentation."""
    return RedirectResponse(url="/docs")
