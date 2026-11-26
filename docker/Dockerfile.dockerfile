FROM python:3.13-slim
WORKDIR /usr/local/app

# Python configuration
ENV PYTHONDONTWRITEBYTECODE=1 \
    PYTHONUNBUFFERED=1 \
    PIP_NO_CACHE_DIR=1

# System libraries needed by scientific Python packages
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
        libgomp1 \
        ca-certificates && \
    rm -rf /var/lib/apt/lists/*

# Install the application dependencies
COPY requirements.txt ./
RUN pip install --no-cache-dir -r requirements.txt

# Copy in the source code
COPY src ./src

# The model
COPY models ./models
COPY artifacts/ ./artifacts/


# Setup an app user so the container doesn't run as the root user
RUN useradd --create-home our_group_app && \
    chown -R our_group_app:our_group_app /usr/local/app

USER our_group_app

EXPOSE 8000

CMD ["uvicorn", "src.deployment.fastapi:app", "--host", "0.0.0.0", "--port", "8000"]