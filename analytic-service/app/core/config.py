"""Application configuration using Pydantic Settings."""

from pydantic import Field
from pydantic_settings import BaseSettings, SettingsConfigDict


class Settings(BaseSettings):
    """Application settings with environment variable support."""

    model_config = SettingsConfigDict(
        env_file=".env",
        env_file_encoding="utf-8",
        case_sensitive=False,
        extra="ignore",
    )

    # Application
    app_name: str = Field(default="analytic-service", description="Application name")
    app_version: str = Field(default="0.1.0", description="Application version")
    app_env: str = Field(default="development", description="Environment (development/production)")

    # Server
    host: str = Field(default="0.0.0.0", description="Server host")
    port: int = Field(default=8000, description="Server port")
    reload: bool = Field(default=True, description="Enable auto-reload for development")

    # API
    api_v1_prefix: str = Field(default="/api/v1", description="API v1 prefix")

    # CORS
    cors_origins: list[str] = Field(
        default=["http://localhost:3000", "http://localhost:8000"],
        description="Allowed CORS origins",
    )
    cors_allow_credentials: bool = Field(default=True, description="Allow credentials in CORS")
    cors_allow_methods: list[str] = Field(default=["*"], description="Allowed HTTP methods")
    cors_allow_headers: list[str] = Field(default=["*"], description="Allowed HTTP headers")

    # Logging
    log_level: str = Field(default="info", description="Logging level")

    # MongoDB
    mongodb_url: str = Field(
        default="mongodb://localhost:27017",
        description="MongoDB connection URL",
    )
    mongodb_database: str = Field(
        default="analytic_db",
        description="MongoDB database name",
    )
    mongodb_max_connections: int = Field(
        default=10,
        description="Maximum number of MongoDB connections in pool",
    )
    mongodb_min_connections: int = Field(
        default=1,
        description="Minimum number of MongoDB connections in pool",
    )

    # RabbitMQ
    rabbitmq_url: str = Field(
        default="amqp://guest:guest@localhost:5672/",
        description="RabbitMQ connection URL",
    )
    rabbitmq_queue: str = Field(
        default="click_events",
        description="RabbitMQ queue name for click events",
    )
    rabbitmq_exchange: str = Field(
        default="analytics",
        description="RabbitMQ exchange name",
    )
    rabbitmq_routing_key: str = Field(
        default="analytics.click",
        description="RabbitMQ routing key",
    )

    @property
    def is_production(self) -> bool:
        """Check if running in production environment."""
        return self.app_env.lower() == "production"


# Global settings instance
settings = Settings()
