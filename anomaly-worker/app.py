import asyncio
import os
import json
import logging
from nats.aio.client import Client as NATS
from app.detector import AnomalyDetector
from app.storage import Database

logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)


async def message_handler(msg, detector, db):
    """Handle incoming probe result messages"""
    try:
        data = json.loads(msg.data.decode())
        logger.info(f"Received probe result: {data}")

        # TODO: Detect anomalies using z-score/ESD
        is_anomaly = detector.detect(data)

        if is_anomaly:
            # TODO: Create incident in database
            # TODO: Publish incident event to NATS
            logger.warning(f"Anomaly detected for target {data.get('target_id')}")

    except Exception as e:
        logger.error(f"Error processing message: {e}")


async def main():
    # Connect to NATS
    nc = NATS()
    await nc.connect(os.getenv("NATS_URL", "nats://localhost:4222"))

    # Initialize detector and database
    detector = AnomalyDetector()
    db = Database(os.getenv("DATABASE_URL"))

    # Subscribe to probe results
    await nc.subscribe("probe.results", cb=lambda msg: message_handler(msg, detector, db))

    logger.info("Anomaly worker started, listening for probe results...")

    # Keep running
    while True:
        await asyncio.sleep(1)


if __name__ == "__main__":
    asyncio.run(main())
