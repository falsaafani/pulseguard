import numpy as np
from scipy import stats


class AnomalyDetector:
    """Detects anomalies in latency and error rates using statistical methods"""

    def __init__(self, z_threshold=3.0):
        self.z_threshold = z_threshold
        self.history = {}  # Store recent data per target

    def detect(self, probe_result):
        """
        Detect if probe result is anomalous

        Args:
            probe_result: dict with keys: target_id, latency_ms, ok

        Returns:
            bool: True if anomaly detected
        """
        target_id = probe_result.get('target_id')
        latency = probe_result.get('latency_ms')

        # TODO: Implement z-score calculation
        # TODO: Track history per target
        # TODO: Return True if |z-score| > threshold

        return False

    def calculate_zscore(self, value, history):
        """Calculate z-score for a value given historical data"""
        if len(history) < 3:
            return 0.0

        mean = np.mean(history)
        std = np.std(history)

        if std == 0:
            return 0.0

        return abs((value - mean) / std)
