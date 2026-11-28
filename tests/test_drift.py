from src.deployment.drift_detection import load_inference_logs, load_baseline_stats, detect_drift
from src.deployment.fastapi import app
from fastapi.testclient import TestClient

artifact_dir = "artifacts"
inference_log_path = artifact_dir + "/inference_log.json"

def test_drift(drift_default_days=7):
    logs = load_inference_logs(inference_log_path)
    print(f"Logged requests: {len(logs)}")
    print(logs[-3:])

    baseline = load_baseline_stats()
    drift_results = detect_drift(logs, days=drift_default_days, baseline_stats=baseline)
    print("\nDrift check results:")
    print(drift_results)

    if drift_results["overall_drift_detected"]:
        print("\n!! Drift detected - time to retrain and redeploy")
    else:
        print("\nNo drift detected")

    # The endpoint does the same check over HTTP
    print("\nGET /drift ->")
    r = TestClient(app).get("/drift")
    assert r.status_code == 200
    assert "overall_drift_detected" in r.json()

if __name__ == "__main__":
    test_drift()